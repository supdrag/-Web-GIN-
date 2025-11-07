package main

import (
	"context"
	"flag"
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
	"supxming.com/my_project/WeCircle_Distributed/Circle_root/Router"
	"supxming.com/my_project/WeCircle_Distributed/Circle_root/config"
	"supxming.com/my_project/WeCircle_Distributed/Circle_root/dao"
	"supxming.com/my_project/WeCircle_Distributed/Circle_root/webskt"
)

var (
	log_lock sync.Mutex
	Mysql    *gorm.DB
	upgrader = &websocket.Upgrader{
		ReadBufferSize:  1024,                                       //发送字节限制
		WriteBufferSize: 1024,                                       //接收字节限制
		CheckOrigin:     func(r *http.Request) bool { return true }, // 开发期放行
	}
	//连接管理,这里一定要这样构造出实例才行，光有类型声明则是空指针，无法访问。
	solo_cncts_manage = &config.Solo_connect_manage{}
	Cncts_manage      = &config.Connect_manage{}
	circle_cancels    = &config.Circle_cancels{}
)

func main() {
	// 1. 注册 flag，默认值 = 旧硬编码值
	dbHost := flag.String("db-host", config.DBHost, "MySQL host")
	dbPort := flag.Int("db-port", config.DBPort, "MySQL port")
	dbUser := flag.String("db-user", config.DBUser, "MySQL user")
	dbPass := flag.String("db-pass", "", "MySQL password (or use DB_PASS env)")
	dbName := flag.String("db-name", config.DBName, "MySQL database")

	redisAddr := flag.String("redis-addr", config.RedisAddr, "Redis address")
	redisPass := flag.String("redis-pass", "", "Redis password (or use REDIS_PASS env)")

	httpPort := flag.String("port", config.HTTPPort, "HTTP listen port")

	flag.Parse()

	mysqlPass := *dbPass
	if mysqlPass == "" {
		mysqlPass = os.Getenv("DB_PASS")
	}
	if mysqlPass == "" {
		mysqlPass = config.DBPass // 旧默认值
	}

	redisPwd := *redisPass
	if redisPwd == "" {
		redisPwd = os.Getenv("REDIS_PASS")
	}
	if redisPwd == "" {
		redisPwd = config.RedisPass // 旧默认值
	}

	config.DBHost = *dbHost
	config.DBPort = *dbPort
	config.DBUser = *dbUser
	config.DBPass = mysqlPass
	config.DBName = *dbName
	config.RedisAddr = *redisAddr
	config.RedisPass = redisPwd
	config.HTTPPort = *httpPort

	config.InitRedis()
	//fmt.Println(fmt.Sprintf("<db-user:%s>\n<db-pass:%s>", config.DBUser, config.DBPass))
	Mysql = dao.Sql_connect("mysql") //数据库连接实例
	rt := gin.Default()
	//rt.Static("/url/chat/chatroom/contact", "./static/chat_room")
	dao.Setup_session(rt) //一定要先注册session中间件，否则其他路由无法调用，会报错
	circle_cancels.Cancels = make(map[int]context.CancelFunc)
	//初始化用户连接管理
	webskt.Connects_init(Cncts_manage, Mysql, circle_cancels)
	webskt.Solo_connects_init(solo_cncts_manage, Mysql)

	//路由注册
	r := Router.Router_register(rt, Mysql, &log_lock,
		solo_cncts_manage,
		Cncts_manage, upgrader, circle_cancels)
	srv := &http.Server{Addr: ":" + config.HTTPPort, Handler: r}
	//Handler指定处理http的服务器变量。
	//之前用的是r.run，此处改为http.Server{}，生成一个服务器对象，
	//利用其接口运行后续代码，比r.run更加灵活，r.run没有这么多接口。
	go dao.Circle_clear(Mysql, circle_cancels, Cncts_manage)
	go srv.ListenAndServe()
	//调用srv的监听端口，go修饰是为了启动另一个goroutine，
	//相当于多线程(不一定是操作系统线程)，使得主程序和监听执行变得并行。
	fmt.Println("Server Began at:", config.HTTPPort)

	quit := make(chan os.Signal, 1) //创建容量为1的信号接收，专门接收Ctrl-C、kill等指令
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	//指定将os.Interrupt(终止信号)、syscall.SIGTERM(停机信号)作为接收对象，quit则作为接收容器。
	//Os.Interrupt：由终端直接发动的Ctrl+C(编译器中按下停止就是这个)
	//Syscall.SIGTERM：shell发动，kill 类指令
	//这两个信号最后都会被进程捕获、响应。Notify可以加多个信号参数，根据需求来。
	<-quit // 阻塞直到接收到 Ctrl-C / kill -15

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//ctx：携带五秒倒计时和取消信号通道，五秒到期就停止ctx.Done通道，监听这个通道的代码就会停止。
	defer cancel() //退出之前把计时器停掉，调用cancel可以停止计时器。

	_ = srv.Shutdown(ctx)
	// 阻塞主进程，进入关闭流程，不接受新的连接，等待已有请求处理，
	//直到时间到了第五秒，将关闭的结果返回出来，如果结果为nil，说明关闭成功。
	//context.Canceled：被提前取消。  context.DeadlineExceeded：超时。  其他：系统级错误。
	dao.Sql_Close(Mysql) // 关闭数据库连接池（你的 Close 内部调用 sqlDB.Close()）
	fmt.Println("Program over！")
}
