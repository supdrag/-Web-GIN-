package dao

//dao用于定义数据访问层操作
import (
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"supxming.com/my_project/WeCircle_Project/Circle_root/config"
)
/*
var mysql_dsn = config.Mysql_config{
	Host:      "127.0.0.1",
	Port:      3306,
	User:      "xiaoming",
	Database:  "web_db",
	Password:  "173743lll",
	Charset:   "utf8mb4",
	ParseTime: "true",
}

var rds_cfg = config.Redis_config{
	Connect_pool_size: 10,
	Net_type:          "tcp",
	Net_url:           "127.0.0.1:6379",
	User_name:         "",
	Passwd:            "173743lcm",
	Redis_code:        "0",
	Secret:            "XiaoMing2024!!@#1234567890", //密钥
	Max_age:           86400,
	Path:              "/",
	Domain:            "",
	Secure:            false,
	Http_only:         true,
	Same_site:         http.SameSiteLaxMode,
	Cookie_name:       "Session_xm",
}
*/
var (
	CfgMysql config.Mysql_config // 可外部写
	CfgRedis config.Redis_config    // 可外部写
)

func Mysql_dsn(c *config.Mysql_config) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%s&loc=Local",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
		c.Charset,
		c.ParseTime,
	)
}

func Sql_connect(sql_name string) *gorm.DB {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("err:", err)
		}
	}()
	fmt.Println(CfgMysql)
	if sql_name == "mysql" {
		dsn := Mysql_dsn(&CfgMysql)
		dlt := mysql.Open(dsn)
		sql, err := gorm.Open(dlt, &gorm.Config{})
		if err != nil || sql == nil {
			fmt.Println("err:", err)
			panic("mysql open error...")
		}
		return sql
	} else {
		panic("Unknown database type...")
	}
}

func Sql_Close(sql *gorm.DB) {
	defer func() {
		err := recover()
		fmt.Println("err:", err)
	}()

	sql_DB, _ := sql.DB()
	err := sql_DB.Close()
	if err != nil {
		panic("Database close error...")
	}
}

func Mysql_tools(tool string) func(db *gorm.DB) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("err:", recover())
		}
	}()

	switch tool {
	case "db_create":
		fmt.Println("create database...")
		return func(db *gorm.DB) {
		}
	case "table_create":
		fmt.Println("create table...")
		return func(db *gorm.DB) {
		}
	case "row_create":
		fmt.Println("create row...")
		return func(db *gorm.DB) {
		}
	default:
		panic("Unknown tool...")
	}
}

// 将redis配置为Session的存储后端，注册到gin的中间件。
func Setup_session(r *gin.Engine) {
	defer func() {
		fmt.Println("redis connect error:", recover())
	}()
	// 将redis和gin建立连接
	store, err := redis.NewStoreWithDB(
		CfgRedis.Connect_pool_size,
		CfgRedis.Net_type,
		CfgRedis.Net_url,
		CfgRedis.User_name,
		CfgRedis.Passwd,
		CfgRedis.Redis_code,
		[]byte(CfgRedis.Secret))
	//展示连接验证信息
	fmt.Printf("Redis config now: <net:%s> <addr:%s> <pwd:%s>\n",
		CfgRedis.Net_type, CfgRedis.Net_url, CfgRedis.Passwd)
	if err != nil {
		panic(fmt.Errorf("err:%w", err))
	}
	// redis存储后端的具体配置
	store.Options(sessions.Options{
		Path:     CfgRedis.Path,
		Domain:   CfgRedis.Domain,
		MaxAge:   CfgRedis.Max_age,
		Secure:   CfgRedis.Secure,
		HttpOnly: CfgRedis.Http_only,
		SameSite: CfgRedis.Same_site})
	// 将redis作为session存储后端注册到中间件。
	r.Use(sessions.Sessions(CfgRedis.Cookie_name, store))
	fmt.Println(fmt.Sprintf("Redis set up!<name:%s> <domain:%s>", CfgRedis.Cookie_name, CfgRedis.Domain))
}
