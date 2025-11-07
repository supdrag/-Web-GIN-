package webskt

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"supxming.com/my_project/WeCircle_Distributed/Circle_root/config"
	"supxming.com/my_project/WeCircle_Distributed/Circle_root/controllers"
)

func Solo_status_get(status int) string {
	if status == 0 {
		return "User Had Been Canceled"
	} else if status == 1 {
		return "Normal"
	} else if status == 2 {
		return "User Had Been Banned"
	} else {
		return "Wrong Status"
	}
}

func Solo_msg_bd(user_id int,
	user_name string,
	friend_id int,
	content string,
	time time.Time) config.Solo_message {
	return config.Solo_message{
		User_ID:   user_id,
		User_name: user_name,
		Friend_ID: friend_id,
		Content:   content,
		Time:      time,
	}
}

func Solo_message_send(msg_get config.Solo_message) []byte {
	return ([]byte)(fmt.Sprintf("<%s-(%d)>:%s \n------[%s]\n",
		msg_get.User_name, msg_get.User_ID, msg_get.Content,
		msg_get.Time.Format("2006-01-02 15:04:05")))
}

// 连接管理初始化
func Solo_connects_init(manage *config.Solo_connect_manage, sql *gorm.DB) {
	manage.Lock() //分配空间后立马锁住

	defer func() { //这里是个细节，最好是defer，保证退出函数时解锁，防止死锁
		(*manage).Unlock()
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()
	var (
		users []config.User //群关系
	)
	fmt.Println("===<<Solo Connects Init Start>>===")

	manage.Connects = make(map[int]*config.Solo_connect)
	//查用户
	rst_users := sql.Table("users").
		Select("ID,STATUS").
		Find(&users)
	if len(users) == 0 {
		fmt.Println("<No Users>")
		return
	}
	if rst_users.Error != nil && errors.Is(rst_users.Error, gorm.ErrRecordNotFound) {
		panic("<Solo Ships Search Error>")
		return
	}
	if rst_users.RowsAffected == 0 {
		fmt.Println("\n<Solo Ships Not Found>" + rst_users.Error.Error())
	}
	//初始化链接
	for i, user := range users {
		manage.Connects[user.ID] = new(config.Solo_connect)
		manage.Connects[user.ID].Conn = nil
		manage.Connects[user.ID].User_id = user.ID
		manage.Connects[user.ID].Status = user.Status
		manage.Connects[user.ID].Status_friend = 1
		fmt.Println(fmt.Sprintf(
			"->(%d)-Solo connect User<%d>-Status[%s] Create Success",
			i+1, user.ID, Solo_status_get(user.Status)))
	}
	fmt.Println("===<<Solo Connects Init Success>>===\n")
}

func Solo_subscrb(ctx *context.Context,
	manage *config.Solo_connect_manage,
	user_id int,
	cn string,
	chan_name string) {
	sub := config.RDB.Subscribe(*ctx, chan_name)
	ch := sub.Channel()
	for rm := range ch {
		// 只写给**当前这个私聊连接**
		manage.RLock()
		if c := manage.Connects[user_id]; c != nil && c.Conn != nil && c.Status == 1 {
			c.Lock()
			_ = c.Conn.WriteMessage(websocket.TextMessage, []byte(rm.Payload))
			c.Unlock()
		}
		manage.RUnlock()
	}
	fmt.Printf("<Redis sub exit> %s\n", cn)
}

// 收听消息
func Solo_read_worker(
	manage *config.Solo_connect_manage,
	user_id int,
	friend_id int,
	user_name string,
	sql *gorm.DB) {
	var (
		status        = manage.Connects[user_id].Status
		status_friend = manage.Connects[user_id].Status_friend
		cn            = fmt.Sprintf("User(%d)->Friend(%d) Status[%s]",
			user_id, friend_id, Solo_status_get(status))
		cc_msg         config.Solo_message
		share          = true
		ctx            context.Context
		cancel         context.CancelFunc
		chan_name      = fmt.Sprintf("User_%d", user_id)
		friend_connect *config.Solo_connect
		connect        *config.Solo_connect
	)
	ctx, cancel = context.WithCancel(context.TODO())
	go Solo_subscrb(&ctx, manage, user_id, cn, chan_name)
	defer func() {
		cancel()
		manage.Lock() //一定要在读之前就上锁，因为这里读写都有，保证读和写的状态一样
		if manage.Connects[user_id] != nil {
			err := manage.Connects[user_id].Conn.Close()
			manage.Connects[user_id].Conn = nil
			manage.Unlock()
			if err != nil {
				fmt.Println("<Solo Connect close fail>" + cn)
			} else {
				fmt.Println("<Solo Connect close success>" + cn)
			}
		} else {
			fmt.Println("<Solo Connect close success>" + cn)
			manage.Unlock()
		}
	}()

	if manage.Connects[user_id].Status == 1 && status_friend != 2 {
		fmt.Println("<USER Online>" + cn)
		_ = config.RDB.Publish(context.TODO(),
			chan_name, string(Solo_message_send(Solo_msg_bd(user_id, user_name, friend_id,
				"<[Online]>", time.Now())))).Err()
	}
	for {
		share = true
		// 获取到用户当前的连接状态
		manage.RLock()
		if manage.Connects[user_id] == nil {
			fmt.Println("<USER Has Been Banned>" + cn)
			manage.RUnlock()
			break
		}
		status_friend = manage.Connects[user_id].Status_friend
		connect = manage.Connects[user_id]
		status = manage.Connects[user_id].Status
		friend_connect = manage.Connects[friend_id]
		manage.RUnlock()

		_, msg, err := connect.Conn.ReadMessage()
		if err != nil || string(msg) == "[EXIT]" {
			break
		}
		//判断当前状态
		if status == 2 {
			connect.Lock()
			err = connect.Conn.WriteMessage(websocket.TextMessage, ([]byte)("<You Had Blocked Friend>\n"))
			connect.Unlock()
			if err != nil {
				fmt.Println("<Feedback Message Send Fail> Banned" + cn)
			} else {
				fmt.Println("<Feedback Message Send Success>" + cn)
			}
			share = false
		} else if status_friend == 0 {
			connect.Lock()
			err = connect.Conn.WriteMessage(websocket.TextMessage, ([]byte)("<Had Been Muted By Friend>\n"))
			connect.Unlock()
			if err != nil {
				fmt.Println("<Feedback Message Send Fail> Muted" + cn)
			} else {
				fmt.Println("<Feedback Message Send Success>" + cn)
			}
			share = false
		} else if status_friend == 2 {
			connect.Lock()
			err = connect.Conn.WriteMessage(websocket.TextMessage, ([]byte)("<Had Been Blocked By Friend>\n"))
			connect.Unlock()
			if err != nil {
				fmt.Println("<Feedback Message Send Fail> Blocked" + cn)
			} else {
				fmt.Println("<Feedback Message Send Success>" + cn)
			}
			share = false
		} else if friend_connect == nil {
			connect.Lock()
			err = connect.Conn.WriteMessage(websocket.TextMessage, ([]byte)("<Friend Status Wrong>\n"))
			connect.Unlock()
			if err != nil {
				fmt.Println("<Feedback Message Send Fail> Wrong Status" + cn)
			} else {
				fmt.Println("<Feedback Message Send Success>" + cn)
			}
			share = false
		} else if friend_connect.Status == 0 {
			connect.Lock()
			err = connect.Conn.WriteMessage(websocket.TextMessage, ([]byte)("<Friend Canceled>\n"))
			connect.Unlock()
			if err != nil {
				fmt.Println("<Feedback Message Send Fail> Friend Canceled" + cn)
			} else {
				fmt.Println("<Feedback Message Send Success>" + cn)
			}
			share = false
		} else if friend_connect.Status == 2 {
			connect.Lock()
			err = connect.Conn.WriteMessage(websocket.TextMessage, ([]byte)("<Friend Has Been Banned>\n"))
			connect.Unlock()
			if err != nil {
				fmt.Println("<Feedback Message Send Fail> Friend Has Been Banned" + cn)
			} else {
				fmt.Println("<Feedback Message Send Success>" + cn)
			}
			share = false
		}

		if !share {
			continue
		}
		// 刷到数据库
		cc_msg = Solo_msg_bd(user_id, user_name, friend_id, string(msg), time.Now())

		err_crt := sql.Table("solo_messages").
			Omit("MESSAGE_ID").
			Create(&cc_msg).Error
		if err_crt != nil {
			fmt.Println("<Solo Message Update Error>" +
				cn + "\n" +
				err.Error() + "\n" +
				"<MESSAGE>" + string(msg) + "\n" +
				"<TIME>" + controllers.Time_now() + "\n")
		}
		fmt.Println("<Solo Message Saved Success>" + cn)
		//广播给所有用户，把数据塞到公共的chan队列里面就行
		_ = config.RDB.Publish(context.TODO(),
			chan_name, string(Solo_message_send(cc_msg))).Err()
		fmt.Println("<MESSAGE SHARE SUCCESS>" +
			"(" + cc_msg.Content + ")" +
			"--Time:" + controllers.Time_now())
	}

	_ = config.RDB.Publish(context.TODO(),
		"solo_broadcast",
		fmt.Sprintf("%d|%d|%s", cc_msg.User_ID, cc_msg.Friend_ID, string(Solo_message_send(cc_msg)))).Err()

	if manage.Connects[user_id] == nil {
		_ = config.RDB.Publish(context.TODO(),
			chan_name, string(Solo_message_send(Solo_msg_bd(user_id, user_name, friend_id,
				">[Has been Banned]<", time.Now())))).Err()
	} else if manage.Connects[user_id].Status == 1 && manage.Connects[user_id].Status_friend == 1 {
		_ = config.RDB.Publish(context.TODO(),
			chan_name, string(Solo_message_send(Solo_msg_bd(user_id, user_name, friend_id,
				">[Exit]<", time.Now())))).Err()
	}
}

func Circle_connect_num() {

}

func Solo_connect_num() {

}
