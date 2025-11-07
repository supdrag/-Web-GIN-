package config

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// 关于用户表格的数据行结构
type User struct {
	ID       int       `gorm:"column:ID;primary_key" json:"id"`
	Name     string    `gorm:"column:NAME;type:varchar(30);not null" json:"name"`
	Account  string    `gorm:"column:ACCOUNT;type:varchar(20);unique_index" json:"account"`
	Passwd   string    `gorm:"column:PASSWD;type:varchar(30);not null" json:"passwd"`
	CTF      string    `gorm:"column:CTF;type:varchar(10);not null" json:"ctf"`
	Phone    string    `gorm:"column:PHONE;type:varchar(20);unique;not null" json:"phone"`
	CreateAt time.Time `gorm:"column:CreateAt;not null" json:"create_at"`
	UpdateAt time.Time `gorm:"column:UpdateAt;not null" json:"update_at"`
	Status   int       `gorm:"column:STATUS;type:int;not null" json:"status"`
}

//==User==
//ID: 用户唯一标识
//Name: 用户姓名
//QQ: 用户QQ
//CTF: 用户身份
//CreatedAt: 记录创建时间
//UpdatedAt: 记录更新时间
//Status: 0-注销 1-正常 2-封号

// 用户界面相关========================================
// 用户每次访问所需要记录的数据
type Visit_log struct { //用户访问历史
	Code    int64     `gorm:"column:CODE;primaryKey;bigint;autoIncrement" json:"code"`
	User_id int       `gorm:"column:USER_ID;not null" json:"user_id"`
	Path    string    `gorm:"column:PATH;type:varchar(150);not null" json:"path"`
	Type    string    `gorm:"column:TYPE;type:varchar(10);not null" json:"type"`
	Time    time.Time `gorm:"column:TIME;not null" json:"time"`
}

type Friend_ship struct { //好友关系
	ID           int64     `gorm:"column:ID;primary_key;autoIncrement" json:"id"`
	Small_ID     int       `gorm:"column:SMALL_ID;not null" json:"small_id"`
	Big_ID       int       `gorm:"column:BIG_ID;not null" json:"big_id"`
	Status_small int       `gorm:"column:STATUS_SMALL;not null" json:"status_small"` //0-屏蔽 1-正常 2-拉黑 3-待确认
	Status_big   int       `gorm:"column:STATUS_BIG;not null" json:"status_big"`     //都是自己对对方的处理
	Time         time.Time `gorm:"column:TIME;not null" json:"time"`
}

type User_profile struct { //用户主页
	User_id    int       `gorm:"column:USER_ID;primary_key" json:"user_id"`
	Signature  string    `gorm:"column:SIGNATURE;type:varchar(255);not null" json:"signature"`
	Popularity int       `gorm:"column:POPULARITY;not null" json:"popularity"`
	Age        int       `gorm:"column:AGE;not null" json:"age"`
	Gender     bool      `gorm:"column:GENDER;type:bool;not null" json:"gender"`
	Location   string    `gorm:"column:LOCATION;type:varchar(50);not null" json:"location"`
	Job        string    `gorm:"column:JOB;type:varchar(30);not null" json:"job"`
	Time       time.Time `gorm:"column:TIME;not null" json:"time"`
}

type Moment struct { //动态
	Moment_ID   int64     `gorm:"column:MOMENT_ID;type:bigint;primary_key;autoIncrement" json:"moment_id"`
	User_ID     int       `gorm:"column:USER_ID;not null" json:"user_id"`
	Content     string    `gorm:"column:CONTENT;type:varchar(400);not null" json:"content"`
	Time        time.Time `gorm:"column:TIME;not null" json:"time"`
	Comment_num int       `gorm:"column:COMMENT_NUM" json:"comment_num"`
	Likes       int       `gorm:"column:LIKES" json:"likes"`
	Status      int       `gorm:"column:STATUS" json:"status"` //0-私密 1-公开
}

type Goods_recommend struct { //用户推荐商品
	ID        int64     `gorm:"column:ID;primary_key;autoIncrement" json:"id"`
	User_id   int       `gorm:"column:USER_ID;not null" json:"user_id"`
	Goods_url string    `gorm:"column:GOODS_URL;type:varchar(255);not null" json:"goods_url"`
	Reason    string    `gorm:"column:REASON;type:varchar(255);not null" json:"reason"`
	Score     int       `gorm:"column:SCORE" json:"score"` //0-100
	Time      time.Time `gorm:"column:TIME;not null" json:"time"`
}

type Game_recommend struct { //用户推荐游戏
	ID       int64     `gorm:"column:ID;primary_key;autoIncrement" json:"id"`
	User_id  int       `gorm:"column:USER_ID;not null" json:"user_id"`
	Game_url string    `gorm:"column:GAME_URL;varchar(255);not null" json:"game_url"`
	Reason   string    `gorm:"column:REASON;type:varchar(255);not null" json:"reason"`
	Score    int       `gorm:"column:SCORE" json:"score"` //0-100
	Time     time.Time `gorm:"column:TIME;not null" json:"time"`
}

type Collect struct { //收藏
	ID      int64     `gorm:"column:ID;primary_key;autoIncrement" json:"id"`
	User_id int       `gorm:"column:USER_ID;not null" json:"user_id"`
	Data    string    `gorm:"column:DATA;type:varchar(255);not null" json:"data"`
	Time    time.Time `gorm:"column:TIME;not null" json:"time"`
}

type Playing struct { //最近在玩
	User_id   int    `gorm:"column:USER_ID;primaryKey;autoIncrement;not null" json:"user_id"`
	Game_name string `gorm:"column:GAME_NAME;not null" json:"game_name"`
}

type Moment_Comment struct { //推荐商品
	Comment_ID int64     `gorm:"column:COMMENT_ID;type:bigint;primary_key;autoIncrement" json:"comment_id"`
	Moment_ID  int64     `gorm:"column:MOMENT_ID;type:bigint;not null" json:"moment_id"`
	User_ID    int       `gorm:"column:USER_ID;not null" json:"user_id"`
	Time       time.Time `gorm:"column:TIME;not null" json:"time"`
	Content    string    `gorm:"column:CONTENT;type:varchar(255);not null" json:"content"`
}

type Goods_Comment struct { //商品评论
	Comment_ID int64     `gorm:"column:COMMENT_ID;type:bigint;primary_key;autoIncrement" json:"comment_id"`
	Goods_Code string    `gorm:"column:GOODS_CODE;type:varchar(50);not null" json:"moment_id"`
	User_ID    int       `gorm:"column:USER_ID;not null" json:"user_id"`
	Time       time.Time `gorm:"column:TIME;not null" json:"time"`
	Content    string    `gorm:"column:CONTENT;type:varchar(255);not null" json:"content"`
	Score      int       `gorm:"column:SCORE;not null" json:"score"` //between 1 and 5
}

type Solo_message struct { //用户消息
	Message_ID int64     `gorm:"column:MESSAGE_ID;primaryKey;autoIncrement" json:"message_id"`
	User_ID    int       `gorm:"column:USER_ID;not null" json:"user_id"`
	User_name  string    `gorm:"column:USER_NAME;not null" json:"user_name"`
	Friend_ID  int       `gorm:"column:FRIEND_ID;not null" json:"circle_id"`
	Content    string    `gorm:"column:CONTENT;type:varchar(512);not null" json:"content"`
	Time       time.Time `gorm:"column:TIME;not null" json:"time"`
}

type Solo_connect struct { //用户连接
	sync.Mutex
	Conn          *websocket.Conn
	User_id       int
	Status        int //0-注销 1-正常 2-封号
	Status_friend int //0-被屏蔽 1-正常 2-被拉黑 是朋友对自己的状态
}

type Solo_connect_manage struct { //用户连接管理
	sync.RWMutex
	Connects map[int](*Solo_connect) //索引是用户id
}
