package webskt

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"supxming.com/my_project/WeCircle_Distributed/Circle_root/config"
	"supxming.com/my_project/WeCircle_Distributed/Circle_root/controllers"
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

// 订阅拉消息实例
func Circle_subscrb(circle_id int,
	ctx context.Context,
	chan_name string,
	manage *config.Connect_manage) {
	sub := config.RDB.Subscribe(ctx, chan_name)
	ch := sub.Channel()
	for rm := range ch {
		// 只写给本机连在该群且在线的客户端
		manage.RLock()
		for _, c := range manage.Connects[circle_id] {
			if c != nil && c.Conn != nil && c.Status == 1 {
				c.Lock()
				_ = c.Conn.WriteMessage(websocket.TextMessage, []byte(rm.Payload))
				c.Unlock()
			}
		}
		manage.RUnlock()
	}
	fmt.Printf("<Circle Sub(%d) Started> \n", circle_id)
}

// 连接管理初始化
func Connects_init(pool *config.Connect_manage,
	sql *gorm.DB,
	cancels *config.Circle_cancels) {
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
		ctx      context.Context
		cancel   context.CancelFunc
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
		Select("USER_ID,CIRCLE_ID,STATUS").
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
		ctx, cancel = context.WithCancel(context.TODO())
		cancels.Cancels[circle.Circle_ID] = cancel
		go Circle_subscrb(circle.Circle_ID, ctx, fmt.Sprintf("Circle_%d", circle.Circle_ID), pool)
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

// 收听消息
func Read_worker(manage *config.Connect_manage,
	user_id int,
	circle_id int,
	sql *gorm.DB,
	user_name string) {
	var (
		cn        = fmt.Sprintf("User(%d)->Circle[%d]", user_id, circle_id)
		cc_msg    config.Circle_message
		share     = true
		status    int
		chan_name = fmt.Sprintf("Circle_%d", circle_id)
		connect   *config.Connect
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
		_ = config.RDB.Publish(context.TODO(),
			chan_name, Message_send(Circle_msg_bd(time.Now(),
				user_id, circle_id, user_name, "<[Has Joined the Circle]>"))).Err()
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
		connect = manage.Connects[circle_id][user_id]
		status = manage.Connects[circle_id][user_id].Status
		manage.RUnlock()

		_, msg, err := connect.Conn.ReadMessage()
		if err != nil || string(msg) == "[EXIT]" {
			break
		}
		if status == 0 {
			connect.Lock()
			err = connect.Conn.WriteMessage(websocket.TextMessage, ([]byte)("<Been Muted in Circle>\n"))
			connect.Unlock()
			if err != nil {
				fmt.Println("<Feedback Message Send Success> Muted" + cn)
			} else {
				fmt.Println("<Feedback Message Send Success>" + cn)
			}
			share = false
		} else if status == 2 {
			connect.Lock()
			err = connect.Conn.WriteMessage(websocket.TextMessage, ([]byte)("<Been Blocked in Circle>\n"))
			connect.Unlock()
			if err != nil {
				fmt.Println("<Feedback Message Send Success> Blocked" + cn)
			} else {
				fmt.Println("<Feedback Message Send Success>" + cn)
			}
			share = false
		} else if status == 4 {
			connect.Lock()
			err = connect.Conn.WriteMessage(websocket.TextMessage, ([]byte)("<Circle Dissolved>\n"))
			connect.Unlock()
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
		// 落库后立即广播到 Redis
		//Publish的第三个参数是要发布的文字内容，可以根据writer改改格式。
		_ = config.RDB.Publish(context.TODO(),
			chan_name, string(Message_send(cc_msg))).Err()
		fmt.Println("<MESSAGE SHARE SUCCESS>" +
			"(" + cc_msg.Content + ")" +
			"--Time:" + controllers.Time_now())
	}

	if manage.Connects[circle_id][user_id] == nil {
		_ = config.RDB.Publish(context.TODO(),
			chan_name, string(Message_send(Circle_msg_bd(time.Now(), user_id, circle_id, user_name, ">[Has been Kicked]<")))).Err()
	} else if manage.Connects[circle_id][user_id].Status == 0 || manage.Connects[circle_id][user_id].Status == 1 {
		_ = config.RDB.Publish(context.TODO(),
			chan_name, string(Message_send(Circle_msg_bd(time.Now(), user_id, circle_id, user_name, ">[EXIT]<")))).Err()
	}
}
