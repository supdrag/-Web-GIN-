package config

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// 交流圈相关========================================

type Circle_ship struct {
	ID        int       `gorm:"column:ID;primaryKey;autoIncrement" json:"id"`
	User_ID   int       `gorm:"column:USER_ID;not null" json:"user_id"`
	Circle_ID int       `gorm:"column:CIRCLE_ID;not null" json:"group_id"` //圈子ID
	Role      int       `gorm:"column:ROLE;not null" json:"role"`          //用户身份(1-master,2-管理员,3-普通用户)
	Status    int       `gorm:"column:STATUS;not null" json:"status"`      //当前状态：0-禁言 1-正常 2-拉黑 3-待确认 4-解散
	Time      time.Time `gorm:"column:TIME;not null" json:"time"`          //加入时间
}

type Circle struct { //注意，主键用primaryKey，自增用autoIncrement，一定要声明
	Circle_ID int       `gorm:"column:CIRCLE_ID;primaryKey;autoIncrement" json:"circle_id"`
	Profile   string    `gorm:"column:PROFILE;not null" json:"profile"`
	Num       int       `gorm:"column:NUM;not null" json:"num"`         //群聊人数
	Lmt_num   int       `gorm:"column:LMT_NUM;not null" json:"lmt_num"` //人数限制
	Status    int       `gorm:"column:STATUS;not null" json:"status"`   //群聊状态 0-解散 1-存在
	Time      time.Time `gorm:"column:TIME;not null" json:"time"`       //创建时间
}

// 交流圈消息
type Circle_message struct {
	Message_ID int64     `gorm:"column:MESSAGE_ID;primaryKey;autoIncrement" json:"message_id"`
	User_ID    int       `gorm:"column:USER_ID;not null" json:"user_id"`
	User_name  string    `gorm:"column:USER_NAME;not null" json:"user_name"`
	Circle_ID  int       `gorm:"column:CIRCLE_ID;not null" json:"circle_id"`
	Content    string    `gorm:"column:CONTENT;type:varchar(512);not null" json:"content"`
	Time       time.Time `gorm:"column:TIME;not null" json:"time"`
}
type Connect struct {
	Conn    *websocket.Conn //连接实例
	User_id int
	Status  int //连接状态。0-禁言 1-正常 2-拉黑 3-待确认 4-解散
}

type Connect_manage struct {
	sync.RWMutex //GO的map并发读写不安全，需要加读写锁
	//这里锁没有给变量名，即匿名嵌入。可以直接 结构体.Lock()，就能锁了。
	Sum_num    int //总连接数量
	Num_online int //在线数量
	Connects   map[int](map[int]*Connect)
	// 连接池实例，是嵌套的map，
	// 第一层map之下是各个群聊，第二层则是群聊中具体的用户
}
