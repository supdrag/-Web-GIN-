package controllers

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"supxming.com/my_project/WeCircle_Project/Circle_root/config"
)

func Chan_ID_get(num int64) int { //根据时间的秒数取余，获取对应通道号
	var result = int(time.Now().Unix() % num)
	fmt.Println(fmt.Sprintf("CHANNEL[%d] Started!", result))
	return result
}

// 比较时间
func Time_less(t1_str, t2_str string) bool {
	const layout = "2006-01-02 15:04:05"
	if t1_str == "*" || t2_str == "*" {
		return true
	}
	full1 := t1_str + " 00:00:00"
	full2 := t2_str + " 00:00:00"
	t1, _ := time.ParseInLocation(layout, full1, time.Local)
	t2, _ := time.ParseInLocation(layout, full2, time.Local)
	return t1.Before(t2)
}

// 检查日期格式 2025-10-10
func Time_check(s string) bool {
	const layout = "2006-01-02"
	_, err := time.ParseInLocation(layout, s, time.Local)
	return err == nil
}

// 找用户
func User_select(user_id int, user *config.User, sql *gorm.DB) error {
	if user_id <= 0 {
		return errors.New("ID Error")
	}
	fmt.Println("user id:", user_id)
	rst := sql.Table("users").
		Select("*").
		Where("ID = ?", user_id).
		Find(user)
	if rst.Error != nil && errors.Is(rst.Error, gorm.ErrRecordNotFound) {
		return errors.New("User Find Error")
	}
	if rst.RowsAffected == 0 {
		return errors.New("User Not Found")
	}
	return nil
}

// 颁发token
func Token_get(uid int, account string, Secret []byte) (string, error) {
	claims := config.Token_data{
		User_ID: uid,
		Account: account,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(Secret)
}

// 验证token
func Token_check(ctx *gin.Context, Secret []byte, user *config.Token_data) (*config.Token_data, error) {
	//返回两个值，第一个是token，第二个是错误提示
	token_data := ctx.GetHeader("Authorization") //从请求头中读取token
	if token_data == "" {
		return nil, errors.New("Authorization missing.")
	}
	//使用ParseWithClaims，对token变量内部填入数据。
	claims := &config.Token_data{}
	token, err := jwt.ParseWithClaims(token_data,
		claims,
		func(t *jwt.Token) (any, error) {
			return Secret, nil
		}) //此处的token变量中装有解析的所有结果，通过对它调用获取相应数据
	if err != nil || !token.Valid {
		fmt.Println("Token analysis failed！")
		return nil, errors.New("Token is invalid!")
	}

	//获取指向token的claims数据的指针，让变量claims直接指向它，节省空间
	user.Account = claims.Account
	user.User_ID = claims.User_ID // 确保 Token_data 里有 User_ID 字段

	fmt.Println("Token analysis success!")
	return claims, nil
	//ok的值反映了Token中Claims对应的数据类型与类型参数是否一致。
	//解析错误，则token加密内容结构体是空
}

func Session_func(ctx *gin.Context, user *config.Token_data, session_xm *string) {
	//从当前请求提取session id，并利用id找到redis中对应的session
	//如果session过期，返回值也仍然是正常session，但是内部所有数据都是nil
	session := sessions.Default(ctx)
	//在对应session中找到user_id字段，如果是nil，则说明session过期了。
	user_id := session.Get("user_id")
	if user_id == nil {
		session.Set("user_id", user.User_ID)
		session.Set("user_account", user.Account)
		_ = session.Save()
		fmt.Println("未找到当前用户的有效Session。\n正在重新创建...")
		cookies := ctx.Writer.Header().Values("Set-Cookie")
		for _, c := range cookies {
			if strings.HasPrefix(c, "Session_xm=") {
				// 提取 Session_xm=xxx; 中的值
				sessionValue := strings.Split(strings.Split(c, ";")[0], "=")[1]
				*session_xm = sessionValue
				fmt.Println("<Session_xm>:", sessionValue)
				break
			}
		}
		fmt.Println("Session rebuild success!")
	} else {
		_ = session.Save()
		*session_xm = "Your Session_xm has been refreshed!\n"
		fmt.Println("已为当前用户刷新Session有效时间！")
	}
}

func Is_number(str string) bool {
	for _, r := range str {
		if r > '9' || r < '0' {
			return false
		}
	}
	return true
}

func APP_check(account string, passwd string, phone string, name string) (bool, string) {
	var pass = false
	var fdbk = ""
	//判断是不是全为空格
	for _, r := range name {
		if !unicode.IsSpace(r) {
			pass = true
			break
		}
	}
	if name == "" {
		pass = false
		fdbk = "Your Name Can't be Empty!"
	}
	if len(account) > 20 {
		pass = false
		fdbk += "Your Account is too long!\n"
	} else if len(account) < 6 {
		pass = false
		fdbk += "Your Account is too short!\n"
	}

	if len(passwd) > 30 {
		pass = false
		fdbk += "Your Password is too long!\n"
	} else if len(passwd) < 6 {
		pass = false
		fdbk += "Your Password is too short!\n"
	}

	if len(phone) != 11 {
		pass = false
		fdbk += "Wrong length of your phone!\n"
	}
	if !Is_number(phone) {
		pass = false
		fdbk += "Wrong character in your phone!\n"
	}

	if pass {
		fdbk += "<APP Success>\n"
	} else {
		fdbk += "<APP Fail>\n"
	}

	return pass, fdbk
}

// 从token中得到用户ID
func Get_UID(ctx *gin.Context, Secret []byte) int {
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

// 返回当前时间文本
func Time_now() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
