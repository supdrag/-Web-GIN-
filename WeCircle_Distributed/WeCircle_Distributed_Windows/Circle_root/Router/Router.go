package Router

//用来存放所有的路由、接口，返回的是服务器引擎。
import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	responses "supxming.com/my_project/WeCircle_Distributed/Circle_root/Responses"
	"supxming.com/my_project/WeCircle_Distributed/Circle_root/config"
	"supxming.com/my_project/WeCircle_Distributed/Circle_root/controllers"
)

// 为全局定义一个密钥
var Secret = []byte("SuperXiaoMing200406150037lcmlovelovelove")

// 为路由注册各个中间件
func Router_register(r *gin.Engine, sql *gorm.DB, log_lock *sync.Mutex,
	solo_cncts_manage *config.Solo_connect_manage,
	circle_cncts_manage *config.Connect_manage,
	upgrader *websocket.Upgrader,
	cancels *config.Circle_cancels) *gin.Engine {
	//添加跨域支持（可选，解决Talend加载问题）
	//添加 Swagger 路由
	r.Use(controllers.Tools_url("log_update", log_lock, Secret, sql))

	r.GET("/url", controllers.Tools_url("get", log_lock, Secret, sql), responses.Response_url("get", sql, Secret, solo_cncts_manage))
	url := r.Group("/url")
	{
		//注册
		url.POST("/register", responses.Response_url("register", sql, Secret, solo_cncts_manage))
		//登录
		url.POST("/login", responses.Response_url("login", sql, Secret, solo_cncts_manage))
		//用户页面
		url.GET("/user", responses.Response_url("user", sql, Secret, solo_cncts_manage))
		//聊天室页面
		url.GET("/chat", responses.Response_url("chat", sql, Secret, solo_cncts_manage))
		//排行榜页面
		url.GET("/rank", responses.Response_url("rank", sql, Secret, solo_cncts_manage))
		//直播功能
		url.GET("living", responses.Response_url("living", sql, Secret, solo_cncts_manage))
	}

	//====================用户管理路由组=====================
	url_user := url.Group("/user") //user是路径封装后的引擎，路径从/user开始
	{                              //以下均表示在user对应路径下开启路由接口。
		url_user.Use(controllers.Tools_url("verify", log_lock, Secret, sql))

		//个人信息总览
		url_user.GET("/profile", responses.Response_user("profile", sql, Secret, solo_cncts_manage))
		//好友列表
		url_user.GET("/friends", responses.Response_user("friends", sql, Secret, solo_cncts_manage))
		//用户注销页面
		url_user.DELETE("/cancel", responses.Response_user("cancel", sql, Secret, solo_cncts_manage))
		//用户浏览历史
		url_user.GET("/history", responses.Response_user("history", sql, Secret, solo_cncts_manage))
	}
	//用户好友
	url_user_friends := url_user.Group("/friends")
	{
		url_user_friends.Use(controllers.Tools_url("verify", log_lock, Secret, sql))

		//联系好友
		url_user_friends.GET("/contact", responses.Response_user_friends("contact", sql, Secret, solo_cncts_manage, upgrader))
		//好友管理
		url_user_friends.GET("/manage", responses.Response_user_friends("manage", sql, Secret, solo_cncts_manage, upgrader))
		//好友请求相关
		url_user_friends.GET("/applies", responses.Response_user_friends("applies", sql, Secret, solo_cncts_manage, upgrader))
		//非实时互动
		url_user_friends.GET("/interaction", responses.Response_user_friends("interaction", sql, Secret, solo_cncts_manage, upgrader))
	}
	//====================聊天室路由组====================
	url_chat := url.Group("/chat")
	{
		url_chat.Use(controllers.Tools_url("verify", log_lock, Secret, sql))

		//交流圈
		url_chat.GET("/chatroom")
		//直播圈子
		url_chat.GET("/living")
	}
	//聊天室路由
	url_chat_chatroom := url_chat.Group("/chatroom")
	{
		//聊天室查找
		url_chat_chatroom.GET("/ctcheck", responses.Response_chat("url_chat_chatroom_ctcheck", sql, Secret, circle_cncts_manage, upgrader, cancels))
		//参与交流
		url_chat_chatroom.GET("/contact", responses.Response_chat("url_chat_chatroom_contact", sql, Secret, circle_cncts_manage, upgrader, cancels))
		//交流圈管理
		url_chat_chatroom.GET("/manage", responses.Response_chat("url_chat_chatroom_manage", sql, Secret, circle_cncts_manage, upgrader, cancels))
		//交流圈邀请
		url_chat_chatroom.GET("/invitations", responses.Response_chat("url_chat_chatroom_invitations", sql, Secret, circle_cncts_manage, upgrader, cancels))
	}
	//====================直播路由====================
	url_living := url_chat.Group("/url")
	{
		//开启直播
		url_living.PUT("/livopen")
		//查找直播/直播广场
		url_living.GET("/livcheck")
		//加入直播
		url_living.GET("/livjoin")
	}

	//====================商品投票排行榜总接口====================
	url_rank := url.Group("/rank")
	{
		url_rank.Use(controllers.Tools_url("verify", log_lock, Secret, sql))

		//投票端口
		url_rank.GET("/goods", responses.Response_rank("goods", sql, Secret))
		//投票统计总榜
		url_rank.GET("/tally", responses.Response_rank("tally", sql, Secret))
		//排行榜规则页面
		url_rank.GET("/rule", responses.Response_rank("rule", sql, Secret))
		//商品上传
		url_rank.GET("/public", responses.Response_rank("public", sql, Secret))
	}

	return r
}
