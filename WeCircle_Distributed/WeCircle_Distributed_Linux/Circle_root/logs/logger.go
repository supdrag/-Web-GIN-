package logs

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"supxming.com/my_project/WeCircle_Distributed/Circle_root/config"
)

func user_id(ctx *gin.Context, Secret []byte) int {
	//从token中获取到账号
	token_data := ctx.GetHeader("Authorization")
	if token_data == "" {
		return -1
	}
	claims := &config.Token_data{}
	token, err := jwt.ParseWithClaims(token_data,
		claims,
		func(t *jwt.Token) (any, error) {
			return Secret, nil
		})
	if err != nil || !token.Valid {
		return -1
	}
	return claims.User_ID
}

func WriteLogToFile(filename, content string) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Can not open the log: %v\n", err)
		return
	}
	defer file.Close() //保证函数结束时关闭文件夹

	_, err = file.WriteString(content)
	if err != nil {
		fmt.Printf("Log update error: %v\n", err)
		return
	}
	fmt.Println("log updated successfully\n")
}

// 从token中得到用户账号
func Get_UAC(ctx *gin.Context, Secret []byte) string {
	//从token中获取到账号
	token_data := ctx.GetHeader("Authorization")
	if token_data == "" {
		return "<Bad Request>"
	}
	claims := &config.Token_data{}
	token, err := jwt.ParseWithClaims(token_data,
		claims,
		func(t *jwt.Token) (any, error) {
			return Secret, nil
		})
	if err != nil || !token.Valid {
		return "<Bad request>"
	}
	return claims.Account
}
func Log_update(ctx *gin.Context, lock *sync.Mutex, Secret []byte, sql *gorm.DB) {
	lock.Lock()
	start := time.Now()
	ctx.Next() // 暂时跳过记录函数，去执行核心中间件，
	// 目的是为了得到实际访问的结果，然后将结果也录入日志
	// 暂时没有写录入访问效果的代码。
	// 日志内容

	var (
		v_log config.Visit_log
		u_id  = user_id(ctx, Secret)
	)

	log_data := fmt.Sprintf(
		"[%s] %s %s %s %s %d %v\n",
		time.Now().Format("2006-01-02 15:04:05"),
		ctx.ClientIP(),
		Get_UAC(ctx, Secret),
		ctx.Request.Method,
		ctx.Request.URL.Path,
		ctx.Writer.Status(),
		time.Since(start),
	)
	// 写入日志文件
	path_now, err := os.Getwd()
	if err != nil {
		fmt.Println("Path error!")
		return
	}
	path_logs := filepath.Join(path_now,
		"Circle_root", "logs", "access.log")
	if u_id != -1 {
		v_log.User_id = u_id
		v_log.Path = ctx.Request.URL.Path
		v_log.Type = ctx.Request.Method
		v_log.Time = time.Now()

		err_err := sql.Table("visit_logs").
			Omit("CODE").
			Create(&v_log).Error
		if err_err != nil {
			fmt.Println("<Log Saved Error>", err_err)
		} else {
			fmt.Println("<Log Saved Successfully>")
		}
	}
	WriteLogToFile(path_logs, log_data)
	lock.Unlock()
}
