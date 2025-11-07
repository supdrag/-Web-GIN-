package webskt

import (
	"errors"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"supxming.com/my_project/WeCircle_Project/Circle_root/config"
	"supxming.com/my_project/WeCircle_Project/Circle_root/controllers"
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

// 广播消息
func Solo_write_worker(manage *config.Solo_connect_manage, news_space <-chan config.Solo_message) {
	//注意chan做形参不允许取指针，默认原地操作，不管通道有多少个[]byte，统一写为<-[]byte
	fmt.Println("<Solo News Worker STARTED>")
	var tar_connect = new(config.Solo_connect)
	for {
		msg_get := <-news_space
		send_data := Solo_message_send(msg_get)
		manage.RLock()
		tar_connect.Conn = manage.Connects[msg_get.Friend_ID].Conn
		tar_connect.Status = manage.Connects[msg_get.Friend_ID].Status
		manage.RUnlock()
		tar_connect.User_id = msg_get.Friend_ID

		if tar_connect != nil &&
			tar_connect.Conn != nil &&
			tar_connect.Status == 1 &&
			tar_connect.User_id != msg_get.User_ID {
			err := tar_connect.Conn.WriteMessage(websocket.TextMessage, send_data)
			if err != nil {
				if tar_connect.Conn == nil {
					continue
				} else {
					fmt.Println("<Broadcast fail>" + err.Error() + "\n" +
						"<Message>" + string(send_data))
				}
			}
		}
	}
}

// 收听消息
func Solo_read_worker(
	manage *config.Solo_connect_manage,
	user_id int,
	friend_id int,
	user_name string,
	news_space chan<- config.Solo_message,
	sql *gorm.DB) {
	var (
		status        = manage.Connects[user_id].Status
		status_friend = manage.Connects[user_id].Status_friend
		cn            = fmt.Sprintf("User(%d)->Friend(%d) Status[%s]",
			user_id, friend_id, Solo_status_get(status))
		cc_msg         config.Solo_message
		share          = true
		friend_connect *config.Solo_connect
		connect        *config.Solo_connect
	)
	defer func() {
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
		news_space <- Solo_msg_bd(user_id, user_name, friend_id, "<[Online]>", time.Now())
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
			err = connect.Conn.WriteMessage(websocket.TextMessage, ([]byte)("<You Had Blocked Friend>\n"))
			if err != nil {
				fmt.Println("<Feedback Message Send Fail> Banned" + cn)
			} else {
				fmt.Println("<Feedback Message Send Success>" + cn)
			}
			share = false
		} else if status_friend == 0 {
			err = connect.Conn.WriteMessage(websocket.TextMessage, ([]byte)("<Had Been Muted By Friend>\n"))
			if err != nil {
				fmt.Println("<Feedback Message Send Fail> Muted" + cn)
			} else {
				fmt.Println("<Feedback Message Send Success>" + cn)
			}
			share = false
		} else if status_friend == 2 {
			err = connect.Conn.WriteMessage(websocket.TextMessage, ([]byte)("<Had Been Blocked By Friend>\n"))
			if err != nil {
				fmt.Println("<Feedback Message Send Fail> Blocked" + cn)
			} else {
				fmt.Println("<Feedback Message Send Success>" + cn)
			}
			share = false
		} else if friend_connect == nil {
			err = connect.Conn.WriteMessage(websocket.TextMessage, ([]byte)("<Friend Status Wrong>\n"))
			if err != nil {
				fmt.Println("<Feedback Message Send Fail> Wrong Status" + cn)
			} else {
				fmt.Println("<Feedback Message Send Success>" + cn)
			}
			share = false
		} else if friend_connect.Status == 0 {
			err = connect.Conn.WriteMessage(websocket.TextMessage, ([]byte)("<Friend Canceled>\n"))
			if err != nil {
				fmt.Println("<Feedback Message Send Fail> Friend Canceled" + cn)
			} else {
				fmt.Println("<Feedback Message Send Success>" + cn)
			}
			share = false
		} else if friend_connect.Status == 2 {
			err = connect.Conn.WriteMessage(websocket.TextMessage, ([]byte)("<Friend Has Been Banned>\n"))
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
		news_space <- cc_msg
		fmt.Println("<MESSAGE SHARE SUCCESS>" +
			"(" + cc_msg.Content + ")" +
			"--Time:" + controllers.Time_now())
	}

	if manage.Connects[user_id] == nil {
		news_space <- Solo_msg_bd(
			user_id, user_name, friend_id,
			">[Has been Banned]<", time.Now())
	} else if manage.Connects[user_id].Status == 1 && manage.Connects[user_id].Status_friend == 1 {
		news_space <- Solo_msg_bd(user_id, user_name, friend_id,
			">[Exit]<", time.Now())
	}
}

func Circle_connect_num() {

}

func Solo_connect_num() {

}
