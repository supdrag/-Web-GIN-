package responses

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"supxming.com/my_project/WeCircle_Project/Circle_root/config"
	"supxming.com/my_project/WeCircle_Project/Circle_root/controllers"
	"supxming.com/my_project/WeCircle_Project/Circle_root/webskt"
)

func status_get(n int) string {
	if n == 0 {
		return "禁言"
	} else if n == 1 {
		return "正常"
	} else if n == 2 {
		return "拉黑"
	} else if n == 4 {
		return "解散"
	} else {
		return "Wrong Status"
	}
}

func role_get(role int) string {
	if role == 1 {
		return "Circle Master"
	} else if role == 2 {
		return "Circle Manager"
	} else if role == 3 {
		return "Common Member"
	} else {
		return "Wrong Role"
	}
}

func wrong_action(ctx *gin.Context, path string) {
	fd := "\n=====================\n" +
		"<SITE>url/chat/chatroom/" + path + "\n" +
		"<STATUS>Refused\n" +
		"<REASON>Wrong Action\n" +
		"=====================\n"
	ctx.String(http.StatusBadRequest, fd)
}

func wrong_tip(fdbk *string, ctx *gin.Context, tip string) {
	fmt.Println("<" + tip + ">")
	*fdbk += "<STATUS>Refused\n" +
		"<REASON>" + tip + "\n" +
		"=====================\n"
	ctx.String(http.StatusBadRequest, *fdbk)
}

func normal_tip(fdbk *string, ctx *gin.Context, tip string) {
	fmt.Println("<" + tip + ">")
	*fdbk += "<STATUS>Success\n" +
		"<DONE>" + tip + "\n" +
		"=====================\n"
	ctx.String(http.StatusOK, *fdbk)
}

// 路由部分=====
func Response_chat(
	typ string,
	sql *gorm.DB,
	Secret []byte,
	cncts_manage *config.Connect_manage,
	upgrader *websocket.Upgrader,
	news_spaces []chan config.Circle_message) func(ctx *gin.Context) {
	if typ == "url_chat_chatroom_ctcheck" {
		return func(ctx *gin.Context) {
			chatroom_ctcheck(ctx, sql, Secret)
		}
	} else if typ == "url_chat_chatroom_contact" {
		return func(ctx *gin.Context) {
			chatroom_contact(ctx, sql, cncts_manage, Secret, upgrader, news_spaces)
		}
	} else if typ == "url_chat_chatroom_manage" {
		//圈子管理
		//Action=join & Circle_id =... 申请加入圈子
		//Action=create & Limit_num=... & Profile=... 创建圈子
		//Action=dissolve & Circle_id =... 解散圈子
		//Action=check  查看已有圈子
		//Action=invite & User_id=... & Circle_id=... 邀请其他用户加入圈子
		//Action=exit & Circle_id=... 退出圈子
		return func(ctx *gin.Context) {
			action := ctx.Query("Action")
			if action == "join" {
				chatroom_manage_join(ctx, sql, Secret)
			} else if action == "create" {
				chatroom_manage_create(ctx, sql, cncts_manage, Secret)
			} else if action == "dissolve" {
				chatroom_manage_dissolve(ctx, sql, cncts_manage, Secret)
			} else if action == "check" {
				chatroom_manage_check(ctx, sql, Secret)
			} else if action == "exit" {
				chatroom_manage_exit(ctx, sql, cncts_manage, Secret, news_spaces)
			} else if action == "kick" {
				chatroom_manage_kick(ctx, sql, cncts_manage, Secret, news_spaces)
			} else if action == "role" {
				chatroom_manage_role(ctx, sql, Secret, news_spaces)
			} else {
				wrong_action(ctx, "manage")
				ctx.Abort()
				return
			}
		}
	} else if typ == "url_chat_chatroom_invitations" {
		return func(ctx *gin.Context) {
			action := ctx.Query("Action")
			if action == "create" {
				chatroom_invitations_create(ctx, sql, Secret)
			} else if action == "check" {
				chatroom_invitations_check(ctx, sql, Secret)
			} else if action == "solve" {
				chatroom_invitations_solve(ctx, sql, cncts_manage, Secret)
			} else {
				wrong_action(ctx, "invitations")
				ctx.Abort()
				return
			}

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

func chatroom_ctcheck(ctx *gin.Context,
	sql *gorm.DB,
	Secret []byte,
) {
	var (
		user_id        = controllers.Get_UID(ctx, Secret)
		circles        []config.Circle
		cid_in         []int
		circle_ship_in []config.Circle_ship
		fdbk           = "\n=====================\n" +
			"<SITE>url/chat/chatroom/ctcheck\n"
		rag = ctx.Query("Range")
	)
	rag_int, err := strconv.Atoi(rag)
	//获取用户搜索范围
	if (rag != "all" && err != nil) || rag == "" {
		wrong_tip(&fdbk, ctx, "Wrong Range")
		ctx.Abort()
		return
	}
	//先找用户当前已经存在于什么交流圈
	err_in := sql.Table("circle_ships").
		Select("CIRCLE_ID").
		Where("USER_ID = ? AND (STATUS = 1 OR STATUS = 0)", user_id).
		Find(&circle_ship_in).Error
	if err_in != nil {
		wrong_tip(&fdbk, ctx, "Circle Search Error")
		ctx.Abort()
		return
	} //创建已加入交流圈的切片
	for _, cc_ship := range circle_ship_in {
		cid_in = append(cid_in, cc_ship.Circle_ID)
	}

	if rag == "all" {
		//找出所有交流圈
		err_all := sql.Table("circles").
			Select("CIRCLE_ID,PROFILE,NUM,LMT_NUM").
			Where("STATUS = 1").
			Find(&circles).Error
		if err_all != nil {
			wrong_tip(&fdbk, ctx, "Circle Search Error")
			ctx.Abort()
			return
		}
	} else {
		//对应交流圈
		err_all := sql.Table("circles").
			Select("CIRCLE_ID,PROFILE,NUM,LMT_NUM").
			Where("STATUS = 1 and CIRCLE_ID = ?", rag_int).
			Find(&circles).Error
		if err_all != nil {
			wrong_tip(&fdbk, ctx, "Circle Search Error")
			ctx.Abort()
			return
		}
	}
	if len(circles) == 0 {
		fdbk += "-------\n(No Circles Meet the Range)\n-------\n"
	} else {
		for i, circle := range circles {
			fdbk += fmt.Sprintf("----\\\n<<%d>>------------------\n[CIRCLE_ID]:%d\n[PROFILE]:%s\n[MEMBERS]:%d/%d\n",
				i+1, circle.Circle_ID,
				circle.Profile, circle.Num, circle.Lmt_num)
			if slices.Contains(cid_in, circle.Circle_ID) {
				fdbk += "(Already Join)\n"
			}
			fdbk += "^^^\n"
		}
	}
	normal_tip(&fdbk, ctx, "Check Success")
}

// contact-接入交流圈频道
func chatroom_contact(ctx *gin.Context,
	sql *gorm.DB,
	cncts_manage *config.Connect_manage,
	Secret []byte,
	upgrader *websocket.Upgrader,
	news_spaces []chan config.Circle_message) {
	var (
		user_id     = controllers.Get_UID(ctx, Secret)
		circle_id   int
		messages    []config.Circle_message
		circle      config.Circle
		circle_ship config.Circle_ship
		chan_id     = controllers.Chan_ID_get(int64(len(news_spaces)))
		fdbk        = "\n=====================\n" +
			"<SITE>url/chat/chatroom/contact\n"
		data_extra = ""
	)
	circle_id, err := strconv.Atoi(ctx.Query("Circle_id"))
	if err != nil {
		fmt.Println("Circle_ID Error")
		ctx.Abort()
		return
	}
	//找出群聊
	rst := sql.Table("circles").
		Select("PROFILE").
		Where("CIRCLE_ID = ?", circle_id).
		First(&circle)
	if rst.Error != nil {
		fmt.Println("Circle Search Error")
		ctx.Abort()
		return
	}
	if rst.RowsAffected == 0 {
		fmt.Println("Circle Not Found")
		ctx.Abort()
		return
	}
	//找出群聊相关资料
	rst = sql.Table("circle_ships").
		Select("*").
		Where("CIRCLE_ID = ? and USER_ID = ?", circle_id, user_id).
		First(&circle_ship)
	if rst.Error != nil {
		fmt.Println("Circle Search Error")
		ctx.Abort()
		return
	}
	if rst.RowsAffected == 0 {
		fmt.Println("Circle Not Found")
		ctx.Abort()
		return
	}
	//先判断连接实例有没有初始化
	if cncts_manage.Connects[circle_id][user_id] == nil {
		fmt.Println("Connect Init Error")
		ctx.Abort()
		return
	}
	//建立连接
	cncts_manage.Lock()
	cncts_manage.Connects[circle_id][user_id].User_id = user_id
	cncts_manage.Connects[circle_id][user_id].Status = circle_ship.Status
	cncts_manage.Connects[circle_id][user_id].Conn, err = upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	cncts_manage.Unlock()

	if err != nil {
		fmt.Println("Connect Error" + err.Error())
		ctx.Abort()
		return
	}

	//找出历史记录，嵌套查询，先找降序五十条，再升序排列
	err = sql.Table("(?) as sub", (sql.Table("circle_messages").
		Select("*").
		Order("TIME desc").
		Where("CIRCLE_ID = ?", circle_id).
		Limit(50))).Order("TIME ASC").Find(&messages).Error
	if err != nil {
		fmt.Println("SQL Connect Error")
		ctx.Abort()
		return
	}

	fdbk += "<STATUS>Success\n" +
		"<CONNECT>" +
		fmt.Sprintf("user[%d]->circle[%d]\n", user_id, circle_id) +
		"<PROFILE>" + circle.Profile + "\n" +
		"---------------------\n"
	cncts_manage.RLock()
	if cncts_manage.Connects[circle_id][user_id].Status != 2 {
		if len(messages) == 0 {
			fdbk += "(No News in This Circle)\n"
		} else {
			for _, msg := range messages {
				fdbk += (string)(webskt.Message_send(msg))
			}
		}
	} else {
		fdbk += "(You Are Blocked in This Circle and Can't See News Here.)\n"
		data_extra = "(Blocked)"
	}
	cncts_manage.RUnlock()

	err = cncts_manage.Connects[circle_id][user_id].Conn.WriteMessage(websocket.TextMessage, ([]byte)(fdbk))
	if err != nil {
		fmt.Println("<Websocket Welcome-Write Error>")
	}
	fmt.Println("<Connection Success>" + fmt.Sprintf("user[%d]->circle[%d]%s\n", user_id, circle_id, data_extra))
	var user config.User
	sql.Table("users").
		Where("ID = ?", user_id).
		First(&user)
	fmt.Println(fmt.Sprintf("CIRCLE USING CHANNEL(%d)", chan_id))
	go webskt.Read_worker(cncts_manage, user_id, circle_id, sql, user.Name, news_spaces[chan_id])
}

// manage-交流圈申请加入
// Action=join & Circle_id =... 申请加入圈子
func chatroom_manage_join(ctx *gin.Context,
	sql *gorm.DB,
	Secret []byte) {
	var (
		circle_id, err_ccid = strconv.Atoi(ctx.Query("Circle_id"))
		user_id             = controllers.Get_UID(ctx, Secret)
		circle              config.Circle
		circle_ship         config.Circle_ship
		fdbk                = "=====================\n" +
			"<SITE>url/chat/chatroom/manage\n" +
			"<ACTION>Join\n"
	)
	//检查user_id和circle_id
	if err_ccid != nil {
		wrong_tip(&fdbk, ctx, "Wrong Circle_id")
		ctx.Abort()
		return
	}
	//检查circle人数限制
	err_cc := sql.Table("circles").
		Where("circle_id = ?", circle_id).
		First(&circle).Error
	if err_cc != nil { //没有circle
		wrong_tip(&fdbk, ctx, "Circle Error")
		ctx.Abort()
		return
	}
	if circle.Num == circle.Lmt_num { //群聊人数已满
		wrong_tip(&fdbk, ctx, "Circle is Full")
		ctx.Abort()
		return
	}
	//检查circle_ship的状态
	rst_ccsp := sql.Table("circle_ships").
		Where("circle_id = ? and user_id = ?", circle_id, user_id).
		First(&circle_ship)
	exist := rst_ccsp.RowsAffected
	if exist != 0 {
		if circle_ship.Status == 2 { //被拉黑
			wrong_tip(&fdbk, ctx, "Blocked in This Circle")
			ctx.Abort()
			return
		} else if circle_ship.Status == 1 { //已经在群聊了
			wrong_tip(&fdbk, ctx, "Already in This Circle")
			ctx.Abort()
			return
		} else if circle_ship.Status == 4 { //群聊已解散
			wrong_tip(&fdbk, ctx, "Circle Dissolved")
			ctx.Abort()
			return
		} else if circle_ship.Status == 3 { //已经申请过
			wrong_tip(&fdbk, ctx, "Apply Had been Created")
			ctx.Abort()
			return
		}

	}
	//创建circle_ship
	circle_ship.User_ID = user_id
	circle_ship.Circle_ID = circle_id
	circle_ship.Time = time.Now()
	circle_ship.Status = 3
	circle_ship.Role = 3
	err := sql.Table("circle_ships").
		Omit("ID").
		Create(&circle_ship).Error
	if err != nil {
		wrong_tip(&fdbk, ctx, "Create Circle_ship Error")
		ctx.Abort()
		return
	}
	fd := fmt.Sprintf("Application {user[%d]->circle[%d]} Sent\n", user_id, circle_id)
	normal_tip(&fdbk, ctx, fd)
}

// manage-创建圈子
// Action=create & Limit_num=... & Profile=... 创建圈子
func chatroom_manage_create(ctx *gin.Context,
	sql *gorm.DB,
	cncts_manage *config.Connect_manage,
	Secret []byte) {
	var (
		limit_num_str = ctx.Query("Limit_num")
		limit_num     = 0
		profile       = ctx.Query("Profile")
		circle        config.Circle
		circle_ship   config.Circle_ship
		user_id       = controllers.Get_UID(ctx, Secret)
		fdbk          = "======================\n" +
			"<SITE>url/chat/chatroom/manage\n" +
			"<ACTION>Create\n"
	)
	//检查基本信息
	limit_num, err := strconv.Atoi(limit_num_str)
	if err != nil {
		wrong_tip(&fdbk, ctx, "<Wrong Limit_num>")
		ctx.Abort()
		return
	}
	if profile == "" {
		wrong_tip(&fdbk, ctx, "<Abnormal Profile>")
		ctx.Abort()
		return
	}
	//创建圈子
	circle.Profile = profile
	circle.Lmt_num = limit_num
	circle.Time = time.Now()
	circle.Num = 1
	circle.Status = 1
	//主键已在数据库定为自增，因此保持写时主键为0值即可自动分配
	//不要用Omit，因为这样写入的行数据才能完整回传到circle，
	//否则circle接收到的值还是0值
	err_crt := sql.Table("circles").
		Create(&circle).Error
	if err_crt != nil {
		wrong_tip(&fdbk, ctx, "<Circle Create Error>")
		ctx.Abort()
		return
	}
	circle_ship.User_ID = user_id
	circle_ship.Circle_ID = circle.Circle_ID
	circle_ship.Time = time.Now()
	circle_ship.Role = 1
	circle_ship.Status = 1
	err_ccs := sql.Table("circle_ships").
		Omit("ID").
		Create(&circle_ship).Error
	if err_ccs != nil {
		wrong_tip(&fdbk, ctx, "<Circle_ship Create Error>")
		ctx.Abort()
		return
	}
	fd := fmt.Sprintf("[master(%d)-circle(%d)]created", user_id, circle.Circle_ID)
	//需要读的时候用读锁，需要写的时候用写锁，都需要的时候直接写锁Lock
	cncts_manage.Lock()
	if cncts_manage.Connects[circle.Circle_ID] == nil {
		cncts_manage.Connects[circle.Circle_ID] = make(map[int]*config.Connect)
	}
	if cncts_manage.Connects[circle.Circle_ID][user_id] == nil {
		cncts_manage.Connects[circle.Circle_ID][user_id] = new(config.Connect)
		cncts_manage.Connects[circle.Circle_ID][user_id].Conn = nil
		cncts_manage.Connects[circle.Circle_ID][user_id].Status = 1
		cncts_manage.Connects[circle.Circle_ID][user_id].User_id = user_id
	}
	cncts_manage.Unlock()
	normal_tip(&fdbk, ctx, fd)
}

// manage-解散圈子
// Action=dissolve & Circle_id =... 解散圈子
func chatroom_manage_dissolve(ctx *gin.Context,
	sql *gorm.DB,
	cncts_manage *config.Connect_manage,
	Secret []byte) {
	var (
		circle_id_str = ctx.Query("Circle_id")
		circle_id     int
		circle_ship   config.Circle_ship
		user_id       = controllers.Get_UID(ctx, Secret)
		fdbk          = "======================\n" +
			"<SITE>url/chat/chatroom/manage\n" +
			"<ACTION>Dissolve\n"
	)
	//判断circle_id格式
	circle_id, err := strconv.Atoi(circle_id_str)
	if err != nil {
		wrong_tip(&fdbk, ctx, "<Wrong Circle_id>")
		ctx.Abort()
		return
	}
	//找找用户的群关系
	e := sql.Table("circle_ships").
		Select("*").
		Where("USER_ID = ? and CIRCLE_ID = ? and STATUS != ?", user_id, circle_id, "0").
		First(&circle_ship).Error
	if e != nil && !errors.Is(e, gorm.ErrRecordNotFound) {
		wrong_tip(&fdbk, ctx, "<Circle_ship Find Error >")
		ctx.Abort()
		return
	}
	if circle_ship.ID == 0 {
		wrong_tip(&fdbk, ctx, "<Circle_ship Not Found>")
		ctx.Abort()
		return
	}
	//看是不是群主
	if circle_ship.Role != 1 {
		wrong_tip(&fdbk, ctx, "<Have No Privilege>")
		ctx.Abort()
		return
	}
	//看看群聊当前的状态
	if circle_ship.Status == 4 {
		wrong_tip(&fdbk, ctx, "<Circle Already Dissolved>")
		ctx.Abort()
		return
	}
	//将状态更新为解散
	cc_err := sql.Table("circles").
		Where("CIRCLE_ID = ?", circle_id).
		Update("STATUS", 0).Error
	if cc_err != nil {
		wrong_tip(&fdbk, ctx, "<Circle Dissolve Error>")
		ctx.Abort()
		return
	}
	//更新所有ship为解散
	cs_err := sql.Table("circle_ships").
		Where("CIRCLE_ID = ?", circle_id).
		Update("STATUS", 4).Error
	if cs_err != nil {
		e := sql.Table("circles").
			Where("CIRCLE_ID = ?", circle_id).
			Update("STATUS", 1)
		if e == nil {
			fmt.Println("(roll back success)")
			fdbk += "(roll back success)\n"
		} else {
			fmt.Println("(roll back error)")
			fdbk += "(roll back error)\n"
		}
		wrong_tip(&fdbk, ctx, "<Circle Dissolve Error>")
		ctx.Abort()
		return
	}
	cncts_manage.Lock()
	for _, cnnt := range cncts_manage.Connects[circle_id] {
		cnnt.Status = 4
	}
	cncts_manage.Unlock()
	don := fmt.Sprintf("circle(%d) dissolved by master(%d)", circle_id, user_id)
	normal_tip(&fdbk, ctx, don)
}

// manage-查找已有圈子
// Action=check & Range=all/Circle_id & Detail=true/false
func chatroom_manage_check(ctx *gin.Context,
	sql *gorm.DB,
	Secret []byte) {
	var (
		circle_ships   []config.Circle_ship
		member_ships   []config.Circle_ship
		circle         config.Circle
		user_id        = controllers.Get_UID(ctx, Secret)
		rag_str        = ctx.Query("Range")
		rag_int, er_rg = strconv.Atoi(rag_str)
		detail         = ctx.Query("Detail")
		detail_bool    bool
		fdbk           = "======================\n" +
			"<SITE>url/chat/chatroom/manage\n" +
			"<ACTION>Check\n"
	)
	if detail != "" {
		if detail == "true" {
			detail_bool = true
		} else if detail == "false" {
			detail_bool = false
		} else {
			wrong_tip(&fdbk, ctx, "<Detail Error>")
			ctx.Abort()
			return
		}
	}
	if rag_str == "all" {
		err := sql.Table("circle_ships").
			Select("CIRCLE_ID,STATUS").
			Where("USER_ID = ? AND STATUS != ?", user_id, 3).
			Find(&circle_ships).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			wrong_tip(&fdbk, ctx, "<Ship Find Error>")
			ctx.Abort()
			return
		}
		if len(circle_ships) == 0 {
			normal_tip(&fdbk, ctx, "Checked, But You Have No Circle")
			ctx.Abort()
			return
		}
	} else {
		if er_rg != nil {
			wrong_tip(&fdbk, ctx, "<Range Error>")
			ctx.Abort()
			return
		}
		err := sql.Table("circle_ships").
			Select("CIRCLE_ID,STATUS").
			Where("USER_ID = ? AND STATUS != ? AND CIRCLE_ID = ?", user_id, 3, rag_int).
			Find(&circle_ships).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			wrong_tip(&fdbk, ctx, "<Ship Find Error>")
			ctx.Abort()
			return
		}
		if len(circle_ships) == 0 {
			normal_tip(&fdbk, ctx, "Checked, But You Have No Circle")
			ctx.Abort()
			return
		}

	}
	fdbk += "--------------------------\n"
	for i, circle_ship := range circle_ships {
		//查当前圈子
		er := sql.Table("circles").
			Select("NUM,LMT_NUM,TIME").
			Where("CIRCLE_ID = ? and STATUS=1", circle_ship.Circle_ID).
			First(&circle).Error
		if er != nil {
			wrong_tip(&fdbk, ctx, "<Circle Find Error>")
			ctx.Abort()
			return
		}
		fdbk += fmt.Sprintf("|-Circle(%d)  ID<%d> Status[%s] Member(%d/%d) --Time{%s}\n",
			i+1, circle_ship.Circle_ID, status_get(circle_ship.Status),
			circle.Num, circle.Lmt_num, circle.Time.Format("2006-01-02 15:04:05"))
		if detail_bool {
			//找出成员详情
			err_ship := sql.Table("circle_ships").
				Select("USER_ID,ROLE,TIME").
				Where("CIRCLE_ID = ? AND (STATUS =0 OR STATUS = 1)", circle_ship.Circle_ID).
				Find(&member_ships).Error
			if err_ship != nil && !errors.Is(err_ship, gorm.ErrRecordNotFound) {
				wrong_tip(&fdbk, ctx, "<Detail Find Error>")
				ctx.Abort()
				return
			}
			if len(circle_ships) == 0 {
				wrong_tip(&fdbk, ctx, "<Detail Not Found>")
				ctx.Abort()
				return
			}
			for n, member_ship := range member_ships {
				fdbk += fmt.Sprintf("   |-User[%d]  ID[%d] Role<%s>  --Join Time{%s}\n",
					n+1, member_ship.User_ID,
					role_get(member_ship.Role),
					member_ship.Time.Format("2006-01-02 15:04:05"))
			}
		}
	}
	fdbk += "--------------------------\n"

	normal_tip(&fdbk, ctx, "Checked Success")
}

// manage-退出圈子
// Action=exit & Circle_id=... 退出圈子
func chatroom_manage_exit(ctx *gin.Context,
	sql *gorm.DB,
	cncts_manage *config.Connect_manage,
	Secret []byte,
	news_spaces []chan config.Circle_message) {
	var (
		circle_ship     config.Circle_ship
		circle_ship_dlt config.Circle_ship
		circle_id_str   = ctx.Query("Circle_id")
		circle_id       int
		msg_exit        config.Circle_message
		user_id         = controllers.Get_UID(ctx, Secret)
		user            config.User
		chan_id         = controllers.Chan_ID_get(int64(len(news_spaces)))
		fdbk            = "======================\n" +
			"<SITE>url/chat/chatroom/manage\n" +
			"<ACTION>exit\n"
	)
	//判断circle_id是否正确
	circle_id, err := strconv.Atoi(circle_id_str)
	if err != nil {
		wrong_tip(&fdbk, ctx, "<Circle_id Error>")
		ctx.Abort()
		return
	}
	//查找对应ship
	err_s := sql.Table("circle_ships").
		Select("ID,STATUS").
		Where("CIRCLE_ID = ? AND USER_ID = ?", circle_id, user_id).
		First(&circle_ship).Error
	if err_s != nil && !errors.Is(err_s, gorm.ErrRecordNotFound) {
		wrong_tip(&fdbk, ctx, "<Circle Search Error>")
		ctx.Abort()
		return
	}
	if circle_ship.ID == 0 {
		wrong_tip(&fdbk, ctx, "<Circle Not Found>")
		ctx.Abort()
		return
	}
	if circle_ship.Status == 0 || circle_ship.Status == 1 {
		//删除前先备份
		err_find := sql.Table("circle_ships").
			Select("*").
			Where("USER_ID = ? AND CIRCLE_ID = ?", user_id, circle_id).
			First(&circle_ship_dlt).Error
		if err_find != nil {
			wrong_tip(&fdbk, ctx, "<Circle_ship Search Error>")
			ctx.Abort()
			return
		}
		err_dlt := sql.Table("circle_ships").
			Where("USER_ID = ? AND CIRCLE_ID = ?", user_id, circle_id).
			Delete(&config.Circle_ship{}).Error
		if err_dlt != nil {
			wrong_tip(&fdbk, ctx, "<Circle_ship Delete Error>")
			ctx.Abort()
			return
		}
		cncts_manage.Lock()
		if cncts_manage.Connects[circle_id][user_id] != nil {
			cncts_manage.Connects[circle_id][user_id].Conn.Close()
		}
		cncts_manage.Connects[circle_id][user_id] = nil
		cncts_manage.Unlock()
	} else {
		wrong_tip(&fdbk, ctx, "<You Are Not in the Circle>")
		ctx.Abort()
		return
	}
	//更新群聊人数
	err = sql.Table("circles").
		Where("CIRCLE_ID = ?", circle_id).
		Update("NUM", gorm.Expr("NUM-1")).Error
	if err != nil {
		//回滚
		er := sql.Table("circle_ships").
			Create(circle_ship_dlt).Error
		if er != nil {
			fdbk += "Roll Back Fail\n"
		} else {
			fdbk += "Roll Back Success\n"
		}
		wrong_tip(&fdbk, ctx, "Circle Update Error")
		ctx.Abort()
		return
	}
	sql.Table("users").
		Select("*").
		Where("ID=?", user_id).
		First(&user)
	//将消息推送给群聊
	msg_exit.User_ID = user_id
	msg_exit.Time = time.Now()
	msg_exit.Content = "|<Leave the Circle>|"
	msg_exit.Circle_ID = circle_id
	msg_exit.User_name = user.Name
	msg_exit.Message_ID = msg_exit.Message_ID
	news_spaces[chan_id] <- msg_exit
	//结束
	normal_tip(&fdbk, ctx, "Exit Success")
}

// manage-踢人
// Action=kick & Target_id=... & Circle_id=... 踢出圈子
func chatroom_manage_kick(ctx *gin.Context,
	sql *gorm.DB,
	cncts_manage *config.Connect_manage,
	Secret []byte,
	news_spaces []chan config.Circle_message) {
	var (
		msg_kick        config.Circle_message
		circle_id       int
		target_id       int
		user            config.User
		circle_ship_tar config.Circle_ship
		circle_ship_usr config.Circle_ship
		user_id         = controllers.Get_UID(ctx, Secret)
		chan_id         = controllers.Chan_ID_get(int64(len(news_spaces)))
		fdbk            = "======================\n" +
			"<SITE>url/chat/chatroom/manage\n" +
			"<ACTION>Kick\n"
	)
	sql.Table("users").
		Select("NAME").
		Where("ID=?", user_id).
		First(&user)
	//判断参数
	circle_id, err_cid := strconv.Atoi(ctx.Query("Circle_id"))
	target_id, err_tid := strconv.Atoi(ctx.Query("Target_id"))
	if err_cid != nil || err_tid != nil {
		wrong_tip(&fdbk, ctx, "<Wrong ID>")
		ctx.Abort()
		return
	}
	if target_id == user_id {
		wrong_tip(&fdbk, ctx, "<Can Not Kick YourSelf>")
		ctx.Abort()
		return
	}
	//找目标
	err1 := sql.Table("circle_ships").
		Select("*").
		Where("CIRCLE_ID = ? and USER_ID=?", circle_id, target_id).
		First(&circle_ship_tar).Error
	if err1 != nil {
		wrong_tip(&fdbk, ctx, "<Target_ship Search Error>")
		ctx.Abort()
		return
	}
	if circle_ship_tar.ID == 0 {
		wrong_tip(&fdbk, ctx, "<Target_ship Not Found>")
		ctx.Abort()
		return
	}
	if circle_ship_tar.Role == 1 {
		wrong_tip(&fdbk, ctx, "<You Have No Privileges>")
		ctx.Abort()
		return
	}

	//看用户有无权限
	err2 := sql.Table("circle_ships").
		Select("*").
		Where("CIRCLE_ID = ? and USER_ID=? and STATUS!=3", circle_id, user_id).
		First(&circle_ship_usr).Error
	if err2 != nil && !errors.Is(err2, gorm.ErrRecordNotFound) {
		wrong_tip(&fdbk, ctx, "<User_ship Search Error>")
		ctx.Abort()
		return
	}
	if circle_ship_usr.ID == 0 {
		wrong_tip(&fdbk, ctx, "<User_ship Not Found>")
		ctx.Abort()
		return
	}
	if circle_ship_usr.Status == 4 {
		wrong_tip(&fdbk, ctx, "<Circle Dissolved>")
		ctx.Abort()
		return
	}
	if circle_ship_usr.Status == 2 {
		wrong_tip(&fdbk, ctx, "<You Are Blocked>")
		ctx.Abort()
		return
	}
	if circle_ship_usr.Role == 2 {
		if circle_ship_tar.Role == 1 || circle_ship_tar.Role == 2 {
			wrong_tip(&fdbk, ctx, "<You Have No Privilege>")
			ctx.Abort()
			return
		} else if circle_ship_tar.Role == 3 {
			wrong_tip(&fdbk, ctx, "<You Have No Privilege>")
		}
		ctx.Abort()
		return
	}
	//踢
	err3 := sql.Table("circle_ships").
		Where("CIRCLE_ID = ? and USER_ID=?", circle_id, target_id).
		Delete(&circle_ship_tar).Error
	if err3 != nil {
		wrong_tip(&fdbk, ctx, "<Kick Error>")
		ctx.Abort()
		return
	}
	//更新circle
	err4 := sql.Table("circles").
		Where("CIRCLE_ID = ?", circle_id).
		Update("NUM", gorm.Expr("NUM-1")).Error
	if err4 != nil {
		err5 := sql.Table("circle_ships").
			Where("CIRCLE_ID = ? and USER_ID = ?", circle_id, user_id).
			Create(&circle_ship_tar).Error
		if err5 != nil {
			fdbk += "Roll Back Fail\n"
		} else {
			fdbk += "Roll Back Success\n"
		}
		wrong_tip(&fdbk, ctx, "<Circle Update Error>")
		ctx.Abort()
		return
	}
	//断掉连接
	cncts_manage.Lock()
	if cncts_manage.Connects[circle_id][target_id].Conn != nil {
		cncts_manage.Connects[circle_id][target_id].Conn.Close()
		fmt.Println("Kick Close Success")
	}
	cncts_manage.Connects[circle_id][target_id] = nil
	cncts_manage.Unlock()

	msg_kick.Time = time.Now()
	msg_kick.Circle_ID = circle_id
	msg_kick.User_name = user.Name
	msg_kick.Content = "|<Had Been Kicked Out by Master>|"
	msg_kick.User_ID = user_id
	news_spaces[chan_id] <- msg_kick

	fd := fmt.Sprintf("Kicked out User(%d) in circle[%d]", target_id, circle_id)
	normal_tip(&fdbk, ctx, fd)
}

// manage-更改身份
// Action=role & Target_id=user_id & Circle_id=c_id & Role=...
func chatroom_manage_role(
	ctx *gin.Context,
	sql *gorm.DB,
	Secret []byte,
	news_spaces []chan config.Circle_message) {
	var (
		user_id          = controllers.Get_UID(ctx, Secret)
		target_str       = ctx.Query("Target_id")
		circle_str       = ctx.Query("Circle_id")
		role_str         = ctx.Query("Role")
		target_id, err_1 = strconv.Atoi(target_str)
		circle_id, err_2 = strconv.Atoi(circle_str)
		role, err_3      = strconv.Atoi(role_str)
		user_ship        config.Circle_ship
		target_ship      config.Circle_ship
		chan_id          = controllers.Chan_ID_get(int64(len(news_spaces)))
		circle_message   config.Circle_message
		user             config.User

		fdbk = "======================\n" +
			"<SITE>url/chat/chatroom/manage\n" +
			"<ACTION>Role\n"
	)
	//先看参数有无问题
	if err_1 != nil || err_2 != nil || err_3 != nil {
		wrong_tip(&fdbk, ctx, "<Arguments Error>")
		ctx.Abort()
		return
	}
	if user_id == target_id {
		wrong_tip(&fdbk, ctx, "<Target Can't Be Yourself>")
		ctx.Abort()
		return
	}
	//看用户有无权限
	err_user := sql.Table("circle_ships").
		Select("USER_ID,CIRCLE_ID,ROLE,STATUS").
		Where("CIRCLE_ID = ? and USER_ID=? and STATUS!=3", circle_id, user_id).
		First(&user_ship).Error
	if err_user != nil {
		wrong_tip(&fdbk, ctx, "<User_ship Search Error>")
		ctx.Abort()
		return
	}
	if user_ship.User_ID == 0 {
		wrong_tip(&fdbk, ctx, "<User_ship Not Found>")
		ctx.Abort()
		return
	}
	if user_ship.Status == 4 {
		wrong_tip(&fdbk, ctx, "<Circle Dissolved>")
		ctx.Abort()
		return
	}
	if user_ship.Role != 1 {
		wrong_tip(&fdbk, ctx, "<You Have No Privilege>")
		ctx.Abort()
		return
	}
	circle_id = user_ship.Circle_ID
	//找用户名
	err := sql.Table("users").
		Where("ID=?", user_id).
		Select("NAME").
		First(&user).Error
	if err != nil {
		wrong_tip(&fdbk, ctx, "<User Not Found>")
		ctx.Abort()
		return
	}
	circle_message.User_name = user.Name
	//看对象是否正常存在于群聊
	err_target := sql.Table("circle_ships").
		Select("USER_ID,ROLE,STATUS").
		Where("CIRCLE_ID = ? and USER_ID=? AND STATUS !=3 and STATUS!=2", circle_id, target_id).
		First(&target_ship).Error
	if err_target != nil {
		wrong_tip(&fdbk, ctx, "<Target_ship Search Error>")
		ctx.Abort()
		return
	}
	if target_ship.User_ID == 0 {
		fmt.Println("<Target_ship Not Found>")
		ctx.Abort()
		return
	}
	//执行修改
	err_udt := sql.Table("circle_ships").
		Where("CIRCLE_ID = ? and USER_ID=?", circle_id, target_id).
		Update("ROLE", gorm.Expr(strconv.Itoa(role))).Error
	if err_udt != nil {
		wrong_tip(&fdbk, ctx, "<Role Change Error>")
		ctx.Abort()
		return
	}

	circle_message.Time = time.Now()
	circle_message.Circle_ID = circle_id
	circle_message.User_ID = user_id
	circle_message.User_name = user.Name
	circle_message.Content = fmt.Sprintf(
		"\n--==[ROLE CHANGE]==--\n<|User_ID:%d|>\n[%s]=>[%s]\n---------------------\n",
		target_id, role_get(target_ship.Role), role_get(role))
	err_cr := sql.Table("circle_messages").
		Omit("MESSAGE_ID").
		Create(&circle_message).Error
	if err_cr != nil {
		wrong_tip(&fdbk, ctx, "<Message Create Error>")
		ctx.Abort()
		return
	}
	news_spaces[chan_id] <- circle_message
	normal_tip(&fdbk, ctx, "<Role Change Success>")
}

// ==============invitations=============
// invitations-创建邀请
// Action=create & Circle_id & Target_id
func chatroom_invitations_create(ctx *gin.Context,
	sql *gorm.DB,
	Secret []byte) {
	var (
		circle_id       int
		target_id       int
		circle_ship_tar config.Circle_ship
		circle_ship_usr config.Circle_ship
		ship_crt        config.Circle_ship
		user_id         = controllers.Get_UID(ctx, Secret)
		fdbk            = "======================\n" +
			"<SITE>url/chat/chatroom/manage\n" +
			"<ACTION>Create\n"
	)
	//判断参数
	circle_id, err_cid := strconv.Atoi(ctx.Query("Circle_id"))
	target_id, err_tid := strconv.Atoi(ctx.Query("Target_id"))
	if err_cid != nil || err_tid != nil {
		wrong_tip(&fdbk, ctx, "<Wrong ID>")
		ctx.Abort()
		return
	}
	if target_id == user_id {
		wrong_tip(&fdbk, ctx, "<Can Not Invite YourSelf>")
		ctx.Abort()
		return
	}
	//分别检查两个人的群关系
	ur_rst := sql.Table("circle_ships"). //邀请者
		Select("STATUS").
		Where("CIRCLE_ID = ? and USER_ID=?", circle_id, user_id).
		First(&circle_ship_usr)
	if ur_rst.Error != nil {
		wrong_tip(&fdbk, ctx, "<User_ship Searched Error>")
		ctx.Abort()
		return
	}
	if ur_rst.RowsAffected == 0 || (circle_ship_usr.Status != 0 && circle_ship_usr.Status != 1) {
		wrong_tip(&fdbk, ctx, "<You Are Not in Circle>")
		ctx.Abort()
		return
	}
	tr_rst := sql.Table("circle_ships"). //被邀请者
		Select("STATUS").
		Where("CIRCLE_ID = ? and USER_ID=?", circle_id, target_id).
		First(&circle_ship_tar)
	if tr_rst.Error != nil && !errors.Is(tr_rst.Error, gorm.ErrRecordNotFound) {
		wrong_tip(&fdbk, ctx, "<Target_ship Searched Error>")
		ctx.Abort()
		return
	}
	//如果有ship，则分类讨论为何不能创建邀请
	if tr_rst.RowsAffected != 0 {
		if circle_ship_tar.Status == 0 || circle_ship_tar.Status == 1 {
			wrong_tip(&fdbk, ctx, "<Target is already in Circle>")
			ctx.Abort()
			return
		} else if circle_ship_tar.Status == 2 {
			wrong_tip(&fdbk, ctx, "<Target is Blocked in Circle>")
			ctx.Abort()
			return
		} else if circle_ship_tar.Status == 3 {
			wrong_tip(&fdbk, ctx, "<Invitation Already Exist>")
			ctx.Abort()
			return
		} else if circle_ship_tar.Status == 4 {
			wrong_tip(&fdbk, ctx, "<Circle was Dissolved>")
			ctx.Abort()
			return
		}
	}
	//创建邀请
	ship_crt.Status = 3
	ship_crt.User_ID = target_id
	ship_crt.Circle_ID = circle_id
	ship_crt.Time = time.Now()
	ship_crt.Role = 3
	err_crt := sql.Table("circle_ships").
		Omit("ID").
		Create(&ship_crt).Error
	if err_crt != nil {
		wrong_tip(&fdbk, ctx, "<Invitation Create Error>")
		ctx.Abort()
		return
	}
	//成功
	normal_tip(&fdbk, ctx, "<Invitation Success>")
}

// invitations-查看请求
// Action=check & Range=all/circle_id
func chatroom_invitations_check(ctx *gin.Context,
	sql *gorm.DB,
	Secret []byte) {
	var (
		user_ship        config.Circle_ship
		own_circle_ships []config.Circle_ship
		own_circle_ids   []int
		wait_ships       []config.Circle_ship
		rag_str          = ctx.Query("Range")
		rag_int, err_rg  = strconv.Atoi(rag_str)
		is_num           = true
		user_id          = controllers.Get_UID(ctx, Secret)
		fdbk             = "======================\n" +
			"<SITE>url/chat/chatroom/manage\n" +
			"<ACTION>Check\n"
	)
	if rag_str == "all" {
		is_num = false
	} else {
		if err_rg != nil {
			wrong_tip(&fdbk, ctx, "<Wrong Range>")
			ctx.Abort()
			return
		}
	}
	//定点搜索的情况
	if is_num {
		//判断用户是不是有权限
		rst := sql.Table("circle_ships").
			Select("ROLE,STATUS").
			Where("CIRCLE_ID = ? AND USER_ID = ?", rag_int, user_id).
			First(&user_ship)
		if rst.Error != nil && !errors.Is(rst.Error, gorm.ErrRecordNotFound) {
			wrong_tip(&fdbk, ctx, "<User_ship Searched Error>")
			ctx.Abort()
			return
		}
		if rst.RowsAffected == 0 {
			wrong_tip(&fdbk, ctx, "<User_ship Not Found>")
			ctx.Abort()
			return
		}
		if user_ship.Role != 1 {
			wrong_tip(&fdbk, ctx, "<You Have No Privilege>")
			ctx.Abort()
			return
		}
		if user_ship.Status == 4 {
			wrong_tip(&fdbk, ctx, "<Circle Already Dissolved>")
			ctx.Abort()
			return
		}
		rst_ws := sql.Table("circle_ships").
			Select("ID,USER_ID,CIRCLE_ID,TIME").
			Where("CIRCLE_ID = ? AND STATUS = ?", rag_int, 3).
			Find(&wait_ships)
		for i, wait_ship := range wait_ships {
			fd := fmt.Sprintf("---APL_ID[%d]\n<%d>-User(%d)=>Circle{%d}---<%s>\n--------------------\n",
				wait_ship.ID, i+1, wait_ship.User_ID, wait_ship.Circle_ID, wait_ship.Time.Format("2006-01-02 15:04:05"))
			fdbk += fd
		}
		if rst_ws.Error != nil && !errors.Is(rst_ws.Error, gorm.ErrRecordNotFound) {
			wrong_tip(&fdbk, ctx, "<Applies Searched Error>")
			ctx.Abort()
			return
		}
	} else {
		//找用户有权限的所有Circle
		rst_cc := sql.Table("circle_ships").
			Select("CIRCLE_ID").
			Where("USER_ID = ? AND ROLE = 1 AND (STATUS = 0 OR STATUS = 1)", user_id).
			Find(&own_circle_ships)
		if rst_cc.Error != nil && !errors.Is(rst_cc.Error, gorm.ErrRecordNotFound) {
			wrong_tip(&fdbk, ctx, "<Circles Searched Error>")
			ctx.Abort()
			return
		}
		if rst_cc.RowsAffected == 0 {
			wrong_tip(&fdbk, ctx, "<You Have No Own Circle>")
			ctx.Abort()
			return
		}
		for _, circle_ship := range own_circle_ships {
			own_circle_ids = append(own_circle_ids, circle_ship.Circle_ID)
		}
		//找出所有相关ship，判断是不是用户的
		w_er := sql.Table("circle_ships").
			Select("ID,USER_ID,CIRCLE_ID,TIME").
			Where("STATUS = 3 AND ROLE!=1 AND USER_ID != ?", user_id).
			Find(&wait_ships)
		if w_er.Error != nil && !errors.Is(w_er.Error, gorm.ErrRecordNotFound) {
			wrong_tip(&fdbk, ctx, "<Applies Searched Error>")
			ctx.Abort()
			return
		}
		n := 0
		if len(wait_ships) != 0 {
			for _, wait_ship := range wait_ships {
				if slices.Contains(own_circle_ids, wait_ship.Circle_ID) {
					n += 1
					fd := fmt.Sprintf("---APL_ID[%d]---\n<%d>-User(%d)=>Circle{%d}---<%s>\n--------------------\n",
						wait_ship.ID, n, wait_ship.User_ID, wait_ship.Circle_ID, wait_ship.Time.Format("2006-01-02 15:04:05"))
					fdbk += fd
				}
			}
		}
	}
	if len(wait_ships) == 0 {
		fdbk += "(Applies Not Found)\n"
	}
	normal_tip(&fdbk, ctx, "<Check Success>")
}

// invitations-处理加群请求
// Action=solve & ID= & Answer=Y/N
func chatroom_invitations_solve(ctx *gin.Context,
	sql *gorm.DB,
	cncts_manage *config.Connect_manage,
	Secret []byte) {
	var ( //拿ID->用circle搜Circle的master->操作ship
		user_ship       config.Circle_ship //请求的
		app_ship        config.Circle_ship //用户的
		answer          = ctx.Query("Answer")
		id_str          = ctx.Query("ID")
		id_ship, err_id = strconv.Atoi(id_str)
		circle_id       int
		user_id         = controllers.Get_UID(ctx, Secret)
		fdbk            = "======================\n" +
			"<SITE>url/chat/chatroom/manage\n" +
			"<ACTION>Check\n"
	)
	//判断参数
	if err_id != nil {
		wrong_tip(&fdbk, ctx, "<ID Error>")
	}
	if answer != "Y" && answer != "N" {
		wrong_tip(&fdbk, ctx, "<Answer Error>")
	}
	//找出ID对应的circle
	rst_ship := sql.Table("circle_ships").
		Select("CIRCLE_ID,STATUS,USER_ID").
		Where("ID = ?", id_ship).
		First(&app_ship)
	if rst_ship.Error != nil && !errors.Is(rst_ship.Error, gorm.ErrRecordNotFound) {
		wrong_tip(&fdbk, ctx, "<Apply Searched Error>")
		ctx.Abort()
		return
	}
	if rst_ship.RowsAffected == 0 {
		wrong_tip(&fdbk, ctx, "<No Such Apply")
		ctx.Abort()
		return
	}
	if app_ship.Status != 3 {
		wrong_tip(&fdbk, ctx, "<Status is Not Wait to Solve>")
		ctx.Abort()
		return
	}
	//判断user是否有权利
	circle_id = app_ship.Circle_ID
	rst_usr := sql.Table("circle_ships").
		Select("ROLE,STATUS").
		Where("USER_ID = ? AND CIRCLE_ID = ?", user_id, circle_id).
		First(&user_ship)
	if rst_usr.Error != nil && !errors.Is(rst_usr.Error, gorm.ErrRecordNotFound) {
		wrong_tip(&fdbk, ctx, "<User_ship Searched Error>")
		ctx.Abort()
		return
	}
	if rst_usr.RowsAffected == 0 {
		wrong_tip(&fdbk, ctx, "<You Are Not in the Circle>")
		ctx.Abort()
		return
	}
	if user_ship.Status == 4 {
		wrong_tip(&fdbk, ctx, "<Circle Dissolved>")
		ctx.Abort()
		return
	}
	if user_ship.Role != 1 {
		wrong_tip(&fdbk, ctx, "<You Have No Privilege>")
		ctx.Abort()
		return
	}
	//根据answer来更改ship
	if answer == "Y" {
		udt_err := sql.Table("circle_ships").
			Where("ID = ?", id_ship).
			Update("STATUS", user_ship.Status).Error
		if udt_err != nil {
			wrong_tip(&fdbk, ctx, "<Solve Error>")
			ctx.Abort()
			return
		}
		udt_cc_num := sql.Table("circles").
			Where("CIRCLE_ID = ?", circle_id).
			Update("NUM", gorm.Expr("NUM+1")).Error
		upt_cc_num := sql.Table("circles").
			Where("CIRCLE_ID = ?", circle_id).
			Update("TIME", time.Now()).Error
		if udt_cc_num != nil || upt_cc_num != nil {
			er := sql.Table("circle_ships").
				Where("ID = ?", id_ship).
				Update("STATUS", 3).Error
			if er != nil {
				fdbk += "(Ship Roll Back Failed)\n"
			} else {
				fdbk += "(Ship Roll Back Success)\n"
			}
			wrong_tip(&fdbk, ctx, "<Circle Num Update Error>")
			ctx.Abort()
			return
		}

	} else {
		udt_err := sql.Table("circle_ships").
			Where("ID = ?", id_ship).
			Delete(&app_ship).Error
		if udt_err != nil {
			wrong_tip(&fdbk, ctx, "<Solve Error>")
			ctx.Abort()
			return
		}
	}
	if answer == "Y" {
		fdbk += "Apply Has been Approved\n"
	} else {
		fdbk += "Apply Has Been Refused\n"
	}
	cncts_manage.Lock()
	cncts_manage.Connects[circle_id][app_ship.User_ID] = new(config.Connect)
	cncts_manage.Connects[circle_id][app_ship.User_ID].Status = user_ship.Status
	cncts_manage.Connects[circle_id][app_ship.User_ID].User_id = app_ship.User_ID
	cncts_manage.Connects[circle_id][app_ship.User_ID].Conn = nil
	cncts_manage.Unlock()
	normal_tip(&fdbk, ctx, "<Solve Success>")
}
