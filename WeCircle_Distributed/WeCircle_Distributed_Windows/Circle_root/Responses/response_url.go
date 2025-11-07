package responses

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"supxming.com/my_project/WeCircle_Distributed/Circle_root/config"
	"supxming.com/my_project/WeCircle_Distributed/Circle_root/controllers"
)

func Response_url(typ string, sql *gorm.DB, Secret []byte, solo_manage *config.Solo_connect_manage) func(ctx *gin.Context) {
	if typ == "get" {
		return func(ctx *gin.Context) {
			var fdbk = "\n=====================\n" +
				"<SITE>url\n"
			fdbk += "hello world!\n" +
				"If you want to enter the web,\n" +
				"Please log in first.\n"
			normal_tip(&fdbk, ctx,
				"Visit Success\n")
			ctx.Abort()
			return
		}
	} else if typ == "login" {
		return func(ctx *gin.Context) {
			var (
				user       config.Token_data
				session_xm string
				user_ctf   config.User //用户发来的
				user_exist config.User //数据库里找的
				ctf_result string
				pass       = true
				fdbk       = "\n=====================\n" +
					"<SITE>url/login\n" + "==验证结果如下==\n"
			)
			cl, _ := controllers.Token_check(ctx, Secret, &user)
			if cl != nil {
				controllers.Session_func(ctx, &user, &session_xm)
				fdbk += "Welcome!\nYour token is valid!\n" + session_xm
				normal_tip(&fdbk, ctx, "Token Pass")
				ctx.Abort()
				return
			}
			_ = ctx.ShouldBindJSON(&user_ctf)
			fmt.Println("get:\n", user_ctf)
			ac_err := sql.Table("users").
				Select("*").
				Where("ACCOUNT = ?", user_ctf.Account).
				First(&user_exist).Error
			if ac_err != nil || user_exist.Account == "" {
				ctf_result = "账号不存在!\n"
				pass = false
			} else {
				if user_ctf.Passwd != user_exist.Passwd {
					ctf_result = "密码错误\n"
					pass = false
				} else {
					if user_exist.Status == 0 {
						ctf_result = "账号已注销!\n"
						pass = false
					} else if user_exist.Status == 2 {
						ctf_result = "账号已被封禁!\n"
						pass = false
					} else {
						user.User_ID = user_ctf.ID
						user.Account = user_ctf.Account
						controllers.Session_func(ctx, &user, &session_xm)
						token_data, _ := controllers.Token_get(user_exist.ID, user_exist.Account, Secret)
						ctf_result = "Account:" + user_ctf.Account + " <Status:pass>\n" +
							"Passwd:" + "Be Hidden And Protected" + " <Status:pass>\n" +
							"<Your token>:" + token_data +
							"\n<Your Session_xm>:" + session_xm +
							"\nToken got! Valid last<24h>.\n"
					}
				}
			}
			fdbk += ctf_result + "==============\n"
			if pass {
				normal_tip(&fdbk, ctx, "Login Success")
			} else {
				wrong_tip(&fdbk, ctx, "Login Failed")
			}

		}
	} else if typ == "register" { //注册：名字，账号，密码，电话
		return func(ctx *gin.Context) {
			var (
				user    config.User
				profile config.User_profile
				fdbk    = "--------------------\n" +
					"[Register Feedback]\n"
			)
			//信息接收
			_ = ctx.ShouldBindJSON(&user)
			fmt.Println("User data:\n" +
				user.Name + "\n" +
				user.Account + "\n" +
				user.Passwd + "\n" +
				user.Phone)
			//信息检查
			app_pass, app_fdbk := controllers.APP_check(user.Account, user.Passwd, user.Phone, user.Name)
			if !app_pass {
				fdbk += app_fdbk + "--------------------\n"
				ctx.String(http.StatusBadRequest, app_fdbk)
				return
			}
			fdbk += app_fdbk
			//创建赋值准备
			user.CTF = "common"
			user.CreateAt = time.Now()
			user.UpdateAt = time.Now()
			//创建
			err := sql.Table("users").
				Omit("ID").
				Create(&user).Error
			if err != nil {
				if errors.Is(err, gorm.ErrDuplicatedKey) {
					fdbk += "Account Existed!\n"
				} else {
					fdbk += "Create User Error!\n"
				}
				fmt.Println(err)
				fdbk += "<Create Fail>\n" + "--------------------\n"
				ctx.String(http.StatusBadRequest, fdbk)
				ctx.Abort()
				return
			}
			profile.User_id = user.ID
			profile.Time = time.Now()
			rst_prf := sql.Table("user_profiles").
				Create(&profile)
			if rst_prf.Error != nil {
				wrong_tip(&fdbk, ctx, "Profile Update Error")
				ctx.Abort()
				return
			}
			//创建连接实例
			solo_manage.Lock()
			solo_manage.Connects[user.ID] = new(config.Solo_connect)
			solo_manage.Connects[user.ID].User_id = user.ID
			solo_manage.Connects[user.ID].Status = user.Status
			solo_manage.Connects[user.ID].Status_friend = 1
			solo_manage.Unlock()

			fdbk += "<Create Success!>\n" +
				"Your Account:" + user.Account + "\n" +
				"Your ID:" + strconv.Itoa(user.ID) + "\n" +
				"Welcome!\n" +
				"--------------------\n"
			ctx.String(http.StatusOK, fdbk)
		}
	} else if typ == "rank" {
		return func(ctx *gin.Context) {
			var fdbk = "<SITE>:url/rank\n" +
				"<BY>:GET\n" +
				"<STATUS>:Success\n" +
				"<INTRODUCTION>\n" +
				"This is the ranking system—hurry up and find the gaming gear that helps you climb ranks!\n" +
				"(Next, you can check the gaming <gear rankings>, look up detailed <device info>, <vote> for gaming gear, or learn about the <voting rules>.)\n" +
				"goods info->url/rank/goods action(check)\n" +
				"vote->url/rank/goods action(vote)\n" +
				"rule-url/rank/rule\n" +
				"tally->url/rank/tally\n" +
				"Before that, head to url/login to log in first!\n"
			ctx.String(http.StatusOK, fdbk)
		}
	} else if typ == "living" {
		return func(ctx *gin.Context) {
			var fdbk = "<SITE>:url/living\n" +
				"<BY>:GET\n" +
				"<STATUS>:Success\n" +
				"<INTRODUCTION>\n" +
				"API is being built" +
				"open live->url/living/livopen\n" +
				"living check->url/living/livcheck\n" +
				"join living->url/livjoin\n" +
				"Before that, head to url/login to log in first!\n"
			ctx.String(http.StatusOK, fdbk)
		}
	} else if typ == "user" {
		return func(ctx *gin.Context) {
			var fdbk = "<SITE>:url/user\n" +
				"<BY>:GET\n" +
				"<STATUS>:Success\n" +
				"<INTRODUCTION>\n" +
				"Step into the chat—where the lobby never sleeps and the meta is born!\n" +
				"(Next you can <join/chat> in rooms, <create><search><invite><check logs>.)\n" +
				"goods info->url/rank/goods action(check)\n" +
				"profile->url/rank/goods action(vote)\n" +
				"friends list-url/user/friends\n" +
				"user cancel ->url/user/profile\n" +
				"user history->url/user/history\n" +
				"contact friend->url/friends/contact\n" +
				"manage friend list->url/friends/manage\n" +
				"friend applies->url/friends/applies\n" +
				"friend interaction->url/user/friends/interaction\n" +
				"Before that, head to url/login to log in first!\n"
			ctx.String(http.StatusOK, fdbk)
		}
	} else if typ == "chat" {
		return func(ctx *gin.Context) {
			var fdbk = "<SITE>:url/chat\n" +
				"<BY>:GET\n" +
				"<STATUS>:Success\n" +
				"<INTRODUCTION>\n" +
				"This is the ranking system—hurry up and find the gaming gear that helps you climb ranks!\n" +
				"(Next, you can check the gaming <gear rankings>, look up detailed <device info>, <vote> for gaming gear, or learn about the <voting rules>.)\n" +
				"chat room->url/chat/chatroom\n" +
				"chat room check->url/chat/chatroom/ctcheck\n" +
				"contact in room->url/chat/chatroom/contact\n" +
				"room manage->url/chat/chatroom/manage\n" +
				"room invite->url/chat/chatroom/invitations\n" +
				"Before that, head to url/login to log in first!\n"
			ctx.String(http.StatusOK, fdbk)
		}
	} else {
		return func(ctx *gin.Context) {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println(err)
					ctx.String(http.StatusBadRequest,
						"There is something wrong in your request!"+
							" \n<Err>:"+err.(string)+"\n")
				}
			}()
			panic("Unsupported type of response: " + typ)
		}
	}
}
