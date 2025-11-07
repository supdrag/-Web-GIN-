package webskt

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"supxming.com/my_project/WeCircle_Project/Circle_root/config"
	"supxming.com/my_project/WeCircle_Project/Circle_root/controllers"
)

func Message_send(msg_get config.Circle_message) []byte {
	return ([]byte)(fmt.Sprintf("[%s]\n<%s-(%d)>:%s \n",
		msg_get.Time.Format("2006-01-02 15:04:05"),
		msg_get.User_name, msg_get.User_ID, msg_get.Content))
}

func Circle_msg_bd(time time.Time,
	u_id int,
	c_id int,
	u_name string,
	ctt string) config.Circle_message {
	return config.Circle_message{
		User_ID:   u_id,
		Time:      time,
		User_name: u_name,
		Content:   ctt,
		Circle_ID: c_id,
	}
}

// 连接管理初始化
func Connects_init(pool *config.Connect_manage,
	sql *gorm.DB) {
	//构建连接系统实例，先锁住
	(*pool).Lock()
	defer func() { //这里是个细节，最好是defer，保证退出函数时解锁，防止死锁
		(*pool).Unlock()
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()
	var (
		sum      int64
		connects []config.Circle_ship //群关系
		circles  []config.Circle      //群聊
	)
	(*pool).Connects = make(map[int]map[int]*config.Connect)
	fmt.Println("\n===<Connects Init Start>===")
	(*pool).Num_online = 0
	sql.Table("circle_ships").
		Where("STATUS != ?", "3").
		Count(&sum)
	(*pool).Sum_num = int(sum)
	if sum == 0 {
		fmt.Println("<No Circle_ships>")
		return
	}
	//找到当前有的群聊
	err_c := sql.Table("circles").
		Find(&circles).Error
	if err_c != nil {
		fmt.Println("<Circle data get err>", err_c)
		panic("<Init fail>")
		return
	}
	fmt.Println("<Circle data get success>")
	//找到所有群关系
	err_cs := sql.Table("circle_ships").
		Where("STATUS != ?", "3").
		Find(&connects).Error
	if err_cs != nil {
		fmt.Println("<Connects data get err>", err_cs)
		panic("<Init fail>")
		return
	}
	fmt.Println("<Connects data get success>\n" +
		"CNN_NUM:" + strconv.Itoa((*pool).Sum_num))

	//构建各个CIRCLE,外层循环
	for _, circle := range circles {
		(*pool).Connects[circle.Circle_ID] = make(map[int]*config.Connect)
	}
	//构建各个user的连接，内层循环
	for i, connect := range connects {
		(*pool).Connects[connect.Circle_ID][connect.User_ID] = new(config.Connect)
		(*pool).Connects[connect.Circle_ID][connect.User_ID].Status = connect.Status
		(*pool).Connects[connect.Circle_ID][connect.User_ID].Conn = nil
		(*pool).Connects[connect.Circle_ID][connect.User_ID].User_id = connect.User_ID
		fmt.Println(fmt.Sprintf(
			"->(%d)-CONNNECT ENVIRONMENT:User<%d>-Circle[%d] BUILTED!",
			i, connect.User_ID, connect.Circle_ID))
	}
	fmt.Println("===<<Connects init success>>===\n")
}

// 广播消息
func Write_worker(manage *config.Connect_manage,
	news_space <-chan config.Circle_message) {
	//注意chan做形参不允许取指针，默认原地操作，不管通道有多少个[]byte，统一写为<-[]byte
	fmt.Println("<BROADCAST Worker STARTED>")
	for {
		msg_get := <-news_space
		send_data := Message_send(msg_get)
		for _, connect := range manage.Connects[msg_get.Circle_ID] {
			if connect != nil && connect.Conn != nil && connect.Status != 2 && connect.Status != 4 && connect.User_id != msg_get.User_ID {
				err := connect.Conn.WriteMessage(websocket.TextMessage, send_data)
				if err != nil {
					if connect.Conn == nil {
						continue
					} else {
						fmt.Println("<Broadcast fail>" + err.Error() + "\n" +
							"<Message>" + string(send_data))
						fmt.Println("\n<Disconnected>" +
							fmt.Sprintf("user(%d)->circle(%d)\n",
								msg_get.User_ID, msg_get.Circle_ID))
					}
				}
			}
		}
	}
}

// 收听消息
func Read_worker(manage *config.Connect_manage,
	user_id int,
	circle_id int,
	sql *gorm.DB,
	user_name string,
	news_space chan<- config.Circle_message) {
	var (
		cn     = fmt.Sprintf("User(%d)->Circle[%d]", user_id, circle_id)
		cc_msg config.Circle_message
		share  = true
		status int
		conn   *websocket.Conn
	)
	defer func() {
		manage.Lock() //一定要在读之前就上锁，因为这里读写都有，保证读和写的状态一样
		if manage.Connects[circle_id][user_id] != nil {
			err := manage.Connects[circle_id][user_id].Conn.Close()
			manage.Connects[circle_id][user_id].Conn = nil
			manage.Unlock()
			if err != nil {
				fmt.Println("<Connect close fail>" + cn)
			} else {
				fmt.Println("<Connect close success>" + cn)
			}
		} else {
			fmt.Println("<Connect close success>" + cn)
			manage.Unlock()
		}
	}()

	if manage.Connects[circle_id][user_id].Status != 2 && manage.Connects[circle_id][user_id].Status != 4 {
		fmt.Println("<USER JOIN>" + cn)
		news_space <- Circle_msg_bd(time.Now(), user_id, circle_id, user_name, "<[Has Joined the Circle]>")
	}
	for {
		share = true
		// 读取用户消息
		manage.RLock()
		if manage.Connects[circle_id][user_id] == nil {
			fmt.Println("<USER Has Been KICKED>" + cn)
			manage.RUnlock()
			break
		}
		conn = manage.Connects[circle_id][user_id].Conn
		status = manage.Connects[circle_id][user_id].Status
		manage.RUnlock()

		_, msg, err := conn.ReadMessage()
		if err != nil || string(msg) == "[EXIT]" {
			break
		}
		if status == 0 {
			err = conn.WriteMessage(websocket.TextMessage, ([]byte)("<Been Muted in Circle>\n"))
			if err != nil {
				fmt.Println("<Feedback Message Send Success> Muted" + cn)
			} else {
				fmt.Println("<Feedback Message Send Success>" + cn)
			}
			share = false
		} else if status == 2 {
			err = conn.WriteMessage(websocket.TextMessage, ([]byte)("<Been Blocked in Circle>\n"))
			if err != nil {
				fmt.Println("<Feedback Message Send Success> Blocked" + cn)
			} else {
				fmt.Println("<Feedback Message Send Success>" + cn)
			}
			share = false
		} else if status == 4 {
			err = conn.WriteMessage(websocket.TextMessage, ([]byte)("<Circle Dissolved>\n"))
			if err != nil {
				fmt.Println("<Feedback Message Send Success>" + cn)
			} else {
				fmt.Println("<Feedback Message Send Success>" + cn)
			}
			share = false
		}
		if !share {
			continue
		}

		// 刷到数据库
		cc_msg = Circle_msg_bd(time.Now(), user_id, circle_id, user_name, string(msg))
		//这个地方有个细节，将消息的ID改为mysql自增，这样就不用在进程又花时间去定ID
		//使得并发性又提高，然后create的时候，用Omit直接忽略掉MESSAGE_ID
		//因为gorm会直接把默认值传入mysql，导致即使自增，还是为1，与原主键冲突
		err_crt := sql.Table("circle_messages").
			Omit("MESSAGE_ID").
			Create(&cc_msg).Error
		if err_crt != nil {
			fmt.Println("<SQL Create Error>" +
				cn + "\n" +
				err.Error() + "\n" +
				"<MESSAGE>" + string(msg) + "\n" +
				"<TIME>" + controllers.Time_now() + "\n")
		}
		//广播给所有用户，把数据塞到公共的chan队列里面就行
		news_space <- cc_msg
		fmt.Println("<MESSAGE SHARE SUCCESS>" +
			"(" + cc_msg.Content + ")" +
			"--Time:" + controllers.Time_now())
	}
	if manage.Connects[circle_id][user_id] == nil {
		news_space <- Circle_msg_bd(time.Now(), user_id, circle_id, user_name, ">[Has been Kicked]<")
	} else if manage.Connects[circle_id][user_id].Status == 0 || manage.Connects[circle_id][user_id].Status == 1 {
		news_space <- Circle_msg_bd(time.Now(), user_id, circle_id, user_name, ">[EXIT]<")
	}
}
