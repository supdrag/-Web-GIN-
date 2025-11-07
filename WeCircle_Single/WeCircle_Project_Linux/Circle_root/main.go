package main

import (
	"flag"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"supxming.com/my_project/WeCircle_Project/Circle_root/Router"
	"supxming.com/my_project/WeCircle_Project/Circle_root/config"
	"supxming.com/my_project/WeCircle_Project/Circle_root/dao"
	"supxming.com/my_project/WeCircle_Project/Circle_root/webskt"
)

var (
	dbHost     = flag.String("db-host", "127.0.0.1", "MySQL host")
	dbPort     = flag.Int("db-port", 3306, "MySQL port")
	dbUser     = flag.String("db-user", "xiaoming", "MySQL user")
	dbPass     = flag.String("db-pass", "", "MySQL password (or use DB_PASS env)")
	dbName     = flag.String("db-name", "web_db", "MySQL database")
	redisAddr  = flag.String("redis-addr", "127.0.0.1:6379", "Redis address")
	redisPass  = flag.String("redis-pass", "", "Redis password (or use REDIS_PASS env)")
	httpPort   = flag.String("port", "9999", "HTTP listen port")
	chanCircle = flag.Int("chan-circle", 10240, "Circle message channel buffer size")
	chanSolo   = flag.Int("chan-solo", 10240, "Solo message channel buffer size")
	paraCircle = flag.Int("para-circle", 10, "Circle write worker count")
	paraSolo   = flag.Int("para-solo", 10, "Solo write worker count")
)

var (
	log_lock      sync.Mutex
	Mysql         *gorm.DB
	port          string
	chan_circle_len, chan_solo_len int
	chan_circle_parallerl_num, chan_solo_parallerl_num int
	upgrader      = &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	solo_cncts_manage = &config.Solo_connect_manage{}
	solo_news_spaces  []chan config.Solo_message
	Cncts_manage      = &config.Connect_manage{}
	news_spaces       []chan config.Circle_message
)

func main() {
    flag.Parse()

	if *dbPass == "" {
		*dbPass = os.Getenv("DB_PASS")
		if *dbPass == "" {
			*dbPass = "173743lll"
		}
	}
	if *redisPass == "" {
		*redisPass = os.Getenv("REDIS_PASS")
		if *redisPass == "" {
			*redisPass = "173743lcm"
		}
	}

	
	dao.CfgMysql = config.Mysql_config{
		Host: *dbHost, 
		Port: *dbPort, 
		User: *dbUser,
		Password: *dbPass, 
		Database: *dbName,
		Charset: "utf8mb4", 
		ParseTime: "true",
	}
	dao.CfgRedis = config.Redis_config{
		Connect_pool_size: 10, 
		Net_type: "tcp",
		Net_url: *redisAddr, 
		Passwd: *redisPass,
		Redis_code: "0", 
		Secret: "XiaoMing2024!!@#1234567890",
		Max_age: 86400, 
		Path: "/", 
		Http_only: true,
		Same_site: http.SameSiteLaxMode, 
		Cookie_name: "Session_xm",
	}

	
	port = *httpPort
	chan_circle_len, chan_solo_len = *chanCircle, *chanSolo
	chan_circle_parallerl_num, chan_solo_parallerl_num = *paraCircle, *paraSolo

    fmt.Println(dao.CfgMysql)
	Mysql = dao.Sql_connect("mysql")
	rt := gin.Default()
	dao.Setup_session(rt)
	webskt.Connects_init(Cncts_manage, Mysql)
	webskt.Solo_connects_init(solo_cncts_manage, Mysql)

	fmt.Printf("%d CIRCLE CHANNELS WILL BE STARTED AS EXPECTED\n", chan_circle_parallerl_num)
	fmt.Printf("%d SOLO CHANNELS WILL BE STARTED AS EXPECTED\n", chan_solo_parallerl_num)

	solo_news_spaces = make([]chan config.Solo_message, chan_solo_parallerl_num)
	news_spaces = make([]chan config.Circle_message, chan_circle_parallerl_num)
	for i := 0; i < chan_solo_parallerl_num; i++ {
		solo_news_spaces[i] = make(chan config.Solo_message, chan_solo_len)
		go webskt.Solo_write_worker(solo_cncts_manage, solo_news_spaces[i]) //开启群聊广播
	}
	for i := 0; i < chan_circle_parallerl_num; i++ {
		news_spaces[i] = make(chan config.Circle_message, chan_circle_len)
		go webskt.Write_worker(Cncts_manage, news_spaces[i]) //开启私聊信息传输
	}

	r := Router.Router_register(rt, Mysql, &log_lock,
		solo_cncts_manage, solo_news_spaces,
		Cncts_manage, news_spaces, upgrader)

	srv := &http.Server{Addr: ":" + port, Handler: r}
	go srv.ListenAndServe()
	fmt.Println("Server Began at:", port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	dao.Sql_Close(Mysql)
	fmt.Println("Program over！")
}