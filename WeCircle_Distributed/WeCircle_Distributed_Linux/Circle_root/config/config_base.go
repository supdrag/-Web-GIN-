package config

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

var (
	DBHost    = "127.0.0.1"
	DBPort    = 3306
	DBUser    = "xiaoming"
	DBPass    = "173743lll"
	DBName    = "web_db"
	RedisAddr = "127.0.0.1:6379"
	RedisPass = "173743lcm"
	HTTPPort  = "9999"
)

type JsonStruct struct {
	Code   int    `json:"code"`   //注意，必须要以反引号包裹
	Msg    string `json:"msg"`    //请求的文字说明 空接口，可以是任意类型
	Type   string `json:"tp"`     //请求携带的具体数据 空接口，可以是任意类型
	Status string `json:"status"` //访问状态
} //此处请求体的字段首字母必须大写，这样才可以导出。

type Token_data struct {
	Account string `json:"account"`
	User_ID int    `json:"user_id"`
	jwt.RegisteredClaims
}

type Search struct {
	Target   string `json:"target"`
	Username string `json:"username"`
	Token    string `json:"token"`
}

// mysql访问配置文件
type Mysql_config struct {
	Host      string //服务器地址，本地是127.0.0.1
	Port      int    //端口号
	User      string //访问用户名
	Database  string //数据库
	Password  string //用户密码
	Charset   string //编码,一般用utf8mb4
	ParseTime string //是否将mysql的时间类型变量解析为time.Time
}

type Redis_config struct {
	Connect_pool_size int           //定义所创建的连接池的大小
	Net_type          string        //网络类型，一般用tcp与redis通信
	Net_url           string        //服务器地址,本地的话是127.0.0.1:6379，端口一般是6379
	User_name         string        //用户名
	Passwd            string        //redis的密码，这里是173743lcm
	Secret            string        //Session的加密密钥,动态生成
	Redis_code        string        //redis数据库编号
	Max_age           int           //Session的有效期
	Path              string        //Session ID的有效路径，即能够向服务器发送ID的路径
	Domain            string        //有效域名，有效的定义同上
	Secure            bool          //True表示只用https传输，false则表示也能用http
	Http_only         bool          //true表示只能由浏览器将cookie发出，js无法读取，用于保护Session。
	Same_site         http.SameSite //跨站请求时，是否发送cookie
	Cookie_name       string        //对应Session的cookie的名字
}

var RDB *redis.Client

func InitRedis() {
	pwd := RedisPass
	if pwd == "" { // 环境变量兜底
		pwd = os.Getenv("REDIS_PASS")
	}
	RDB = redis.NewClient(&redis.Options{
		Addr:     RedisAddr,
		Password: pwd,
		DB:       0,
	})
	if err := RDB.Ping(context.Background()).Err(); err != nil {
		panic(fmt.Sprintf("redis connect fail: %v", err))
	}
}
