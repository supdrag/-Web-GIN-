package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"supxming.com/my_project/WeCircle_Distributed/Circle_root/config"
	"supxming.com/my_project/WeCircle_Distributed/Circle_root/logs"
)

func Tools_url(use string, lock_log *sync.Mutex, Secret []byte, sql *gorm.DB) func(ctx *gin.Context) {
	if use == "get" {
		return func(ctx *gin.Context) {
			var res = "Welcome to my web<url>!\n" +
				"Here you can enjoy your shopping!\n"
			ctx.String(http.StatusOK, res)
		}
	} else if use == "log_update" {
		return func(ctx *gin.Context) {
			ctx.Next()
			logs.Log_update(ctx, lock_log, Secret, sql)
		}
	} else if use == "verify" {
		return func(ctx *gin.Context) {
			var user config.Token_data
			tk_data, err := Token_check(ctx, Secret, &user)
			if tk_data == nil {
				f_b := "Token check failed!\n" +
					"Please move to the login website !\n" +
					"err:" + err.Error() + "\n"
				fmt.Println(f_b)
				ctx.String(http.StatusUnauthorized, f_b)
				ctx.Abort()
				return
			}
			fmt.Println("<Token passed>\n" +
				"USER:" + strconv.Itoa(user.User_ID) + "\n")
			var s string
			Session_func(ctx, &user, &s)
			ctx.Next()
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
			panic("Unknown middleware in url...\n")
		}
	}
}
