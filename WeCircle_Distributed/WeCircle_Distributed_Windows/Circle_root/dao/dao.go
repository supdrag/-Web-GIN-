package dao

//dao用于定义数据访问层操作
import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"supxming.com/my_project/WeCircle_Distributed/Circle_root/config"
)

func Mysql_cfg() *config.Mysql_config {
	return &config.Mysql_config{
		Host:      config.DBHost,
		Port:      config.DBPort,
		User:      config.DBUser,
		Database:  config.DBName,
		Password:  config.DBPass,
		Charset:   "utf8mb4",
		ParseTime: "true",
	}
}

func Redis_cfg() *config.Redis_config {
	return &config.Redis_config{
		Connect_pool_size: 10,
		Net_type:          "tcp",
		Net_url:           config.RedisAddr,
		User_name:         "",
		Passwd:            config.RedisPass,
		Redis_code:        "0",
		Secret:            "XiaoMing2024!!@#1234567890",
		Max_age:           86400,
		Path:              "/",
		Domain:            "",
		Secure:            false,
		Http_only:         true,
		Same_site:         http.SameSiteLaxMode,
		Cookie_name:       "Session_xm",
	}
}

var mysql_dsn *config.Mysql_config

var redis_dsn *config.Redis_config

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
	mysql_dsn = Mysql_cfg()
	redis_dsn = Redis_cfg()
	//fmt.Println(fmt.Sprintf("<db-user:%s>\n<db-pass:%s>", mysql_dsn.User, mysql_dsn.Password))
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("err:", err)
		}
	}()
	if sql_name == "mysql" {
		dsn := Mysql_dsn(mysql_dsn)
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

// 将redis配置为Session的存储后端，注册到gin的中间件。
func Setup_session(r *gin.Engine) {
	defer func() {
		fmt.Println("redis connect error:", recover())
	}()
	// 将redis和gin建立连接
	store, err := redis.NewStoreWithDB(
		redis_dsn.Connect_pool_size,
		redis_dsn.Net_type,
		redis_dsn.Net_url,
		redis_dsn.User_name,
		redis_dsn.Passwd,
		redis_dsn.Redis_code,
		[]byte(redis_dsn.Secret))
	//展示连接验证信息
	fmt.Printf("Redis config now: <net:%s> <addr:%s> <pwd:%s>\n",
		redis_dsn.Net_type, redis_dsn.Net_url, redis_dsn.Passwd)
	if err != nil {
		panic(fmt.Errorf("err:%w", err))
	}
	// redis存储后端的具体配置
	store.Options(sessions.Options{
		Path:     redis_dsn.Path,
		Domain:   redis_dsn.Domain,
		MaxAge:   redis_dsn.Max_age,
		Secure:   redis_dsn.Secure,
		HttpOnly: redis_dsn.Http_only,
		SameSite: redis_dsn.Same_site})
	// 将redis作为session存储后端注册到中间件。
	r.Use(sessions.Sessions(redis_dsn.Cookie_name, store))
	fmt.Println(fmt.Sprintf("Redis set up!<name:%s> <domain:%s>", redis_dsn.Cookie_name, redis_dsn.Domain))
}

func Circle_clear(sql *gorm.DB,
	cancels *config.Circle_cancels,
	manage *config.Connect_manage) {
	var (
		circles   []*config.Circle
		clear_num int
	)
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()

	for true {
		clear_num = 0
		err := sql.Table("circles").
			Select("CIRCLE_ID,NUM").
			Find(&circles).Error
		if err != nil {
			fmt.Println("<Circle Clear Select Error>", err)
			panic("<Circle Clear Select Error>" + err.Error())
		}

		for _, circle := range circles {
			if circle.Num == 0 {
				err = sql.Table("circles").
					Where("CIRCLE_ID = ?", circle.Circle_ID).
					Delete(&circle).Error
				if err != nil {
					fmt.Println("<Circle Clear Delete Error>", err)
					panic("<Circle Clear Delete Error>" + err.Error())
				}
				manage.Lock()
				manage.Connects[circle.Circle_ID] = nil
				manage.Unlock()
				cancels.Lock()
				if cancels.Cancels[circle.Circle_ID] != nil {
					cancels.Cancels[circle.Circle_ID]()
				}
				cancels.Unlock()
				clear_num++
			}
		}
		fmt.Println(fmt.Sprintf("<Circle Clear (%d) Circle>", clear_num))
		time.Sleep(60 * time.Second)
	}
}
