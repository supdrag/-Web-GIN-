package responses

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"supxming.com/my_project/WeCircle_Distributed/Circle_root/config"
	"supxming.com/my_project/WeCircle_Distributed/Circle_root/controllers"
	"supxming.com/my_project/WeCircle_Distributed/Circle_root/webskt"
)

func pair_sort(user_id int, target_id int) (int, int) {
	if user_id > target_id {
		return target_id, user_id
	} else {
		return user_id, target_id
	}
}

func gender_get(gender bool) string {
	if gender {
		return "Boy"
	} else {
		return "Girl"
	}
}

func status_frd_get(
	status_small int,
	status_big int,
	small_user int,
	big_user int) string {
	if status_small == 1 && status_big == 1 {
		return "Common"
	} else if status_small == 1 && status_big == 0 {
		return fmt.Sprintf("User(%d) is Muted by User(%d)",
			small_user, big_user)
	} else if status_small == 0 && status_big == 1 {
		return fmt.Sprintf("User(%d) is Muted by User(%d)",
			big_user, small_user)
	} else if status_small == 0 && status_big == 0 {
		return "Muted Each Other"
	} else if status_small == 2 && status_big == 2 {
		return "Blocked Each Other"
	} else if status_small == 2 && status_big == 1 {
		return fmt.Sprintf("User(%d) Blocked User(%d)",
			small_user, big_user)
	} else if status_small == 1 && status_big == 2 {
		return fmt.Sprintf("User(%d) Blocked User(%d)",
			big_user, small_user)
	} else if status_small == 0 && status_big == 2 {
		return fmt.Sprintf("User(%d) Blocked User(%d),User(%d) Muted User(%d)",
			big_user, small_user, small_user, big_user)
	} else if status_small == 2 && status_big == 0 {
		return fmt.Sprintf("User(%d) Blocked User(%d),User(%d) Muted User(%d)",
			small_user, big_user, big_user, small_user)
	} else {
		return "False Ship"
	}
}

func frd_id(min_id int, max_id int, user_id int) int {
	if user_id == min_id {
		return max_id
	} else {
		return min_id
	}
}

func Response_user(typ string,
	sql *gorm.DB,
	Secret []byte,
	cncts_manage *config.Solo_connect_manage) func(ctx *gin.Context) {
	if typ == "profile" {
		return func(ctx *gin.Context) {
			user_profile(ctx, sql, Secret)
		}
	} else if typ == "friends" {
		return func(ctx *gin.Context) {
			user_friends(ctx, sql, Secret)
		}
	} else if typ == "cancel" {
		return func(ctx *gin.Context) {
			user_cancel(ctx, sql, Secret, cncts_manage)
		}
	} else if typ == "history" {
		return func(ctx *gin.Context) {
			user_history(ctx, sql, Secret)
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

func Response_user_friends(typ string,
	sql *gorm.DB,
	Secret []byte,
	cncts_manage *config.Solo_connect_manage,
	upgrader *websocket.Upgrader) func(ctx *gin.Context) {
	if typ == "manage" {
		return func(ctx *gin.Context) {
			user_friends_manage(ctx, sql, Secret)
		}
	} else if typ == "contact" {
		return func(ctx *gin.Context) {
			user_friends_contact(ctx, sql, Secret, cncts_manage, upgrader)
		}
	} else if typ == "applies" {
		return func(ctx *gin.Context) {
			user_friends_applies(ctx, sql, Secret)
		}
	} else if typ == "interaction" {
		return func(ctx *gin.Context) {
			user_friends_interaction(ctx, sql, Secret)
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

func user_profile(
	ctx *gin.Context,
	sql *gorm.DB,
	Secret []byte) {
	var (
		// 基本信息
		action  = ctx.Query("Action")
		user_id = controllers.Get_UID(ctx, Secret)
		fdbk    = "\n=====================\n" +
			"<SITE>url/user/profile\n"
	)
	if action == "check" {
		user_profile_check(ctx, sql, user_id, &fdbk)
	} else if action == "update" {
		user_profile_update(ctx, sql, user_id, &fdbk)
	} else {
		wrong_tip(&fdbk, ctx, "Wrong Action")
		ctx.Abort()
		return
	}
}

// 查询用户主页
// Action = check/update 行为
// Target_id = user_id
// Moments = true/false 动态
// Comments = true/false 动态的评论
// Goods = true/false 推荐商品
// Game = true/false 推荐游戏
// Collect = true/false 收藏
// Playing = true/false 最近在玩
func user_profile_check(
	ctx *gin.Context,
	sql *gorm.DB,
	user_id int,
	fdbk *string) {
	var (
		// 基本信息
		profile          config.User_profile
		user             config.User
		target_str       = ctx.Query("Target_id")
		target_id, err_t = strconv.Atoi(target_str)

		//参数信息
		moments    []config.Moment
		moment_str = ctx.Query("Moments")

		comments     []config.Moment_Comment
		comments_str = ctx.Query("Comments")

		goods     []config.Goods_recommend
		goods_str = ctx.Query("Goods")

		game     []config.Game_recommend
		game_str = ctx.Query("Game")

		collect     []config.Collect
		collect_str = ctx.Query("Collect")

		playing     config.Playing
		playing_str = ctx.Query("Playing")
	)
	// 判断行为合法
	if err_t != nil {
		wrong_tip(fdbk, ctx, "Wrong Target_id")
		ctx.Abort()
		return
	}
	//用户
	err_user := controllers.User_select(target_id, &user, sql)
	if err_user != nil {
		wrong_tip(fdbk, ctx, err_user.Error())
		ctx.Abort()
		return
	}
	// 搜索用户的主页基本信息
	pro_rst := sql.Table("user_profiles").
		Select("*").
		Where("USER_ID = ?", target_id).
		Find(&profile)
	if pro_rst.Error != nil && !errors.Is(pro_rst.Error, gorm.ErrRecordNotFound) {
		wrong_tip(fdbk, ctx, "User Profile Find Error")
		ctx.Abort()
		return
	}
	if pro_rst.RowsAffected == 0 {
		wrong_tip(fdbk, ctx, "User Profile Not Found")
		ctx.Abort()
		return
	}
	*fdbk += "__--=======================--__\n" +
		fmt.Sprintf("|||   %s's PROFILE   |||\n", user.Name) +
		fmt.Sprintf("[USWE_ID]:%d  [NAME]:%s  [GENDER]:%s  [AGE]:%d\n",
			target_id, (user).Name, gender_get(profile.Gender), profile.Age) +
		fmt.Sprintf("[POPULARITY]:%d\n", profile.Popularity) +
		fmt.Sprintf("[SIGNATURE]:%s\n", profile.Signature) +
		fmt.Sprintf("[JOB]:%s  [LOCATION]:%s  [TIME]:%s\n",
			profile.Job, profile.Location, profile.Time.Format("2006-01-02 15:04:05")) +
		"-----------------------------\n"
	//===各个参数判断===
	//推荐商品
	if goods_str == "true" {
		*fdbk += "-----------------------------\n"
		*fdbk += fmt.Sprintf("---GOODS RECOMMENDS OF %s---\n", user.Name)
		rst_goods := sql.Table("goods_recommends").
			Select("GOODS_URL,REASON,SCORE,TIME").
			Where("USER_ID = ?", target_id).
			Find(&goods)
		if rst_goods.Error != nil && !errors.Is(rst_goods.Error, gorm.ErrRecordNotFound) {
			wrong_tip(fdbk, ctx, "Goods Recommend Find Error")
			ctx.Abort()
			return
		}
		if rst_goods.RowsAffected == 0 {
			*fdbk += "(NO Goods Recommend)\n"
		} else {
			for i, good := range goods {
				*fdbk += fmt.Sprintf("   -\n    RCMD (%d) [SCORE:%d --%s]\n",
					i+1, good.Score, good.Time.Format("2006-01-02 15:04:05")) +
					fmt.Sprintf("    [URL]%s\n    [REASON]%s\n", good.Goods_url, good.Reason)
			}
		}
		*fdbk += "-----------------------------\n"
	}

	//推荐游戏
	if game_str == "true" {
		*fdbk += "------------------------------\n"
		*fdbk += fmt.Sprintf("---GAME RECOMMENDS OF %s---\n", user.Name)
		rst_game := sql.Table("game_recommends").
			Select("*").
			Where("USER_ID = ?", target_id).
			Find(&game)
		if rst_game.Error != nil && !errors.Is(rst_game.Error, gorm.ErrRecordNotFound) {
			wrong_tip(fdbk, ctx, "Game Recommends Find Error")
			ctx.Abort()
			return
		}
		if rst_game.RowsAffected == 0 {
			*fdbk += "(NO Game Recommends)\n"
		} else {
			for i, gm := range game {
				*fdbk += fmt.Sprintf("   -\n    RCMD (%d) [SCORE:%d --%s]\n",
					i+1, gm.Score, gm.Time.Format("2006-01-02 15:04:05")) +
					fmt.Sprintf("    [URL]%s\n    [REASON]%s\n", gm.Game_url, gm.Reason)
			}
		}
		*fdbk += "-----------------------------\n"
	}

	//收藏
	if collect_str == "true" {
		*fdbk += "-----------------------------\n"
		*fdbk += fmt.Sprintf("---COLLECTS OF %s---\n", user.Name)
		rst_collect := sql.Table("collects").
			Select("DATA,TIME").
			Where("USER_ID = ?", target_id).
			Find(&collect)
		if rst_collect.Error != nil && !errors.Is(rst_collect.Error, gorm.ErrRecordNotFound) {
			wrong_tip(fdbk, ctx, "Collect Recommend Find Error")
			ctx.Abort()
			return
		}
		if rst_collect.RowsAffected == 0 {
			*fdbk += "(NO Collect)\n"
		} else {
			for i, clt := range collect {
				*fdbk += fmt.Sprintf("   -\n    COLLECT (%d) [%s]\n    CONTENT%s\n",
					i+1, clt.Data, clt.Time.Format("2006-01-02 15:04:05"))
			}
		}
		*fdbk += "-----------------------------\n"
	}

	//最近在玩
	if playing_str == "true" {
		*fdbk += "-----------------------------\n"
		*fdbk += fmt.Sprintf("---%s PLAYING---\n", user.Name)
		rst_play := sql.Table("playings").
			Select("GAME_NAME").
			Where("USER_ID = ?", target_id).
			Find(&playing)
		if rst_play.Error != nil && !errors.Is(rst_play.Error, gorm.ErrRecordNotFound) {
			wrong_tip(fdbk, ctx, "Playing Find Error")
			ctx.Abort()
			return
		}
		if rst_play.RowsAffected == 0 {
			*fdbk += "(NO Playing)\n"
		} else {
			*fdbk += fmt.Sprintf("-\n    Playing[%s]\n", playing.Game_name)

		}
		*fdbk += "-----------------------------\n"
	}
	//用户动态 & 动态评论
	if moment_str == "true" {
		rst_mmts := sql.Table("moments").
			Select("*").
			Where("USER_ID = ?", target_id).
			Find(&moments)
		if rst_mmts.Error != nil && !errors.Is(rst_mmts.Error, gorm.ErrRecordNotFound) {
			wrong_tip(fdbk, ctx, "User Moments Find Error")
			ctx.Abort()
			return
		}
		*fdbk += "-----------------------------\n"
		*fdbk += fmt.Sprintf("---%s MOMENTS---\n", user.Name)

		if len(moments) == 0 {
			*fdbk += "(NO MOMENTS)\n"
			*fdbk += "-----------------------------\n"
		} else {
			for i, moment := range moments {
				if moment.Status == 0 && target_id != user_id {
					continue
				}
				*fdbk += fmt.Sprintf("  ___\n  M[%d]========\n", i+1)
				*fdbk += fmt.Sprintf("  ->NAME[%s]  STATUS[%d]  -TIME[%s]\n",
					user.Name, moment.Status, moment.Time.Format("2006-01-02 15:04:05")) +
					fmt.Sprintf("    [CONTENT]:%s\n", moment.Content)
				*fdbk += fmt.Sprintf("    LIKES (%d)  COMMENT(%d)\n",
					moment.Likes, moment.Comment_num)
				if comments_str == "true" {
					rst_comments := sql.Table("moment_comments").
						Select("USER_ID,TIME,CONTENT").
						Where("MOMENT_ID = ?", moment.Moment_ID).
						Find(&comments)
					if rst_comments.Error != nil && !errors.Is(rst_comments.Error, gorm.ErrRecordNotFound) {
						wrong_tip(fdbk, ctx, "User Moment Comment Find Error")
						ctx.Abort()
						return
					}
					if rst_comments.RowsAffected != 0 {
						for n, comment := range comments {
							*fdbk += fmt.Sprintf("    COMMENT (%d) [ID:%d Time:%s]\n    %s\n",
								n+1, comment.User_ID,
								comment.Time.Format("2006-01-02 15:04:05"),
								comment.Content)
						}
						*fdbk += "============\n"
					}
				}
				*fdbk += "-----------------------------\n"
			}
		}
	}
	normal_tip(fdbk, ctx, "<Profile Check Success>")
}

// 更新主页信息
// Update = profile & -d {signature,age,gender,location,job}
// Update = moment & -d {content,status}
// Update = goods_recommend & -d {goods_url,reason,score}
// Update = game_recommend & -d {game_url,reason,score}
// Update = collect & -d {data}
// Update = playing & -d {game_name}
func user_profile_update(
	ctx *gin.Context,
	sql *gorm.DB,
	user_id int,
	fdbk *string) {
	var (
		//更新对象
		update = ctx.Query("Update")

		//主页
		profile      config.User_profile
		age_str      = ctx.Query("Age")
		age, err_age = strconv.Atoi(age_str)
		gender       = ctx.Query("Gender")
		//动态
		moment config.Moment
		//商品推荐
		goods_recommend config.Goods_recommend
		//游戏推荐
		game_recommend config.Game_recommend
		//收藏品
		collect config.Collect
		//最近在玩
		playing config.Playing

		//获取请求值
		updts = map[string]interface{}{}
	)
	if update == "profile" {
		_ = ctx.ShouldBindJSON(&profile)
		updts["TIME"] = time.Now()
		if gender == "true" {
			updts["GENDER"] = true
		} else if gender == "false" {
			updts["GENDER"] = false
		}
		if err_age == nil {
			updts["AGE"] = age
		}
		if profile.Location != "" {
			updts["LOCATION"] = profile.Location
		}
		if profile.Signature != "" {
			updts["SIGNATURE"] = profile.Signature
		}
		if profile.Job != "" {
			updts["JOB"] = profile.Job
		}

		rst := sql.Table("user_profiles").
			Where("USER_ID = ?", user_id).
			Updates(updts)
		if rst.Error != nil {
			wrong_tip(fdbk, ctx, "User Profile Update Error")
			ctx.Abort()
			return
		}
	} else if update == "moment" {
		_ = ctx.ShouldBindJSON(&moment)
		moment.User_ID = user_id
		moment.Time = time.Now()
		moment.Likes = 0
		moment.Comment_num = 0
		moment.Status = 1
		rst := sql.Table("moments").
			Omit("MOMENT_ID").
			Create(&moment)
		if rst.Error != nil {
			wrong_tip(fdbk, ctx, "User Moment Create Error")
			ctx.Abort()
			return
		}
	} else if update == "goods_recommend" {
		_ = ctx.ShouldBindJSON(&goods_recommend)
		goods_recommend.Time = time.Now()
		goods_recommend.User_id = user_id
		rst := sql.Table("goods_recommends").
			Omit("ID").
			Create(&goods_recommend)
		if rst.Error != nil {
			wrong_tip(fdbk, ctx, "User Goods Recommend Create Error")
			ctx.Abort()
			return
		}
	} else if update == "game_recommend" {
		_ = ctx.ShouldBindJSON(&game_recommend)
		game_recommend.Time = time.Now()
		game_recommend.User_id = user_id
		rst := sql.Table("game_recommends").
			Omit("ID").
			Create(&game_recommend)
		if rst.Error != nil {
			wrong_tip(fdbk, ctx, "User Game Recommend Create Error")
			ctx.Abort()
			return
		}
	} else if update == "collect" {
		_ = ctx.ShouldBindJSON(&collect)
		collect.User_id = user_id
		collect.Time = time.Now()
		rst := sql.Table("collects").
			Omit("ID").
			Create(&collect)
		if rst.Error != nil {
			wrong_tip(fdbk, ctx, "User Collect Create Error")
			ctx.Abort()
			return
		}
	} else if update == "playing" {
		_ = ctx.ShouldBindJSON(&playing)
		playing.User_id = user_id
		rst := sql.Table("playings").
			Create(&playing)
		if rst.Error != nil {
			wrong_tip(fdbk, ctx, "User Playing Create Error")
			ctx.Abort()
			return
		}
	} else {
		wrong_tip(fdbk, ctx, "Wrong Action")
		ctx.Abort()
		return
	}
	normal_tip(fdbk, ctx, "<Profile Update Success>")
}

// 好友列表
// Range = all/Friend_ID
func user_friends(
	ctx *gin.Context,
	sql *gorm.DB,
	Secret []byte) {
	var (
		rag_str          = ctx.Query("Range")
		rag_int, err_0   = strconv.Atoi(rag_str)
		user_id          = controllers.Get_UID(ctx, Secret)
		friend_ships     []config.Friend_ship
		min_num, max_num int
		fdbk             = "\n=====================\n" +
			"<SITE>url/user/friends\n"
	)
	//判断范围格式
	if err_0 != nil && rag_str != "all" {
		wrong_tip(&fdbk, ctx, "<Wrong Range>")
		ctx.Abort()
		return
	}
	//查找ship
	if rag_str == "all" {
		rst_all := sql.Table("friend_ships").
			Select("*").
			Where("(SMALL_ID = ? OR BIG_ID = ?)", user_id, user_id).
			Find(&friend_ships)
		if rst_all.Error != nil && !errors.Is(rst_all.Error, gorm.ErrRecordNotFound) {
			wrong_tip(&fdbk, ctx, "<Friends Find Error>")
			ctx.Abort()
			return
		}
		if len(friend_ships) == 0 {
			wrong_tip(&fdbk, ctx, "<Friends Not Found>")
			ctx.Abort()
			return
		}
	} else {
		if rag_int == user_id {
			wrong_tip(&fdbk, ctx, "<Friend Can't Be Yourself>")
			ctx.Abort()
			return
		}
		min_num, max_num = pair_sort(user_id, rag_int)
		rst_rag := sql.Table("friend_ships").
			Select("*").
			Where("SMALL_ID=? AND BIG_ID =?", min_num, max_num).
			Find(&friend_ships)
		if rst_rag.Error != nil && !errors.Is(rst_rag.Error, gorm.ErrRecordNotFound) {
			wrong_tip(&fdbk, ctx, "<Friend Find Error>")
			ctx.Abort()
			return
		}
		if len(friend_ships) == 0 {
			wrong_tip(&fdbk, ctx, "<Friend Not Found>")
			ctx.Abort()
			return
		}
	}
	fdbk += "==-----__\n"
	for i, friend_ship := range friend_ships {
		friend_id := frd_id(friend_ship.Small_ID, friend_ship.Big_ID, user_id)
		fdbk += fmt.Sprintf("<%d>=-  User_id(%d) Status{%s} -Time[%s]\n",
			i+1, friend_id,
			status_frd_get(friend_ship.Status_small, friend_ship.Status_big, friend_ship.Small_ID, friend_ship.Big_ID),
			friend_ship.Time.Format("2006-01-02 15:04:05"))
	}
	fdbk += "_______--\n"
	normal_tip(&fdbk, ctx, "<Friends Check Success>")
}

// 用户注销(软删除，仍然保留用户数据)
// Passwd = passwd
// 把solo关闭，
func user_cancel(
	ctx *gin.Context,
	sql *gorm.DB,
	Secret []byte,
	cncts_manage *config.Solo_connect_manage) {
	var (
		user_id = controllers.Get_UID(ctx, Secret)
		passwd  = ctx.Query("Passwd")
		user    config.User
		fdbk    = "\n=====================\n" +
			"<SITE>url/user/cancel\n"
	)
	//核对密码
	err_user := controllers.User_select(user_id, &user, sql)
	if err_user != nil {
		wrong_tip(&fdbk, ctx, "<User Find Error>")
		ctx.Abort()
		return
	}
	if user.Passwd != passwd {
		wrong_tip(&fdbk, ctx, "<User Password Error>")
		ctx.Abort()
		return
	}
	//修改用户状态-注销
	rst_updt := sql.Table("users").
		Where("ID = ?", user_id).
		Update("STATUS", gorm.Expr("0"))
	if rst_updt.Error != nil {
		wrong_tip(&fdbk, ctx, "<User Status Update Error>")
		ctx.Abort()
		return
	}
	//关掉solo通道
	cncts_manage.Lock()
	if cncts_manage.Connects[user_id] != nil {
		cncts_manage.Connects[user_id].Conn = nil
		cncts_manage.Connects[user_id].Status = 0
	}
	cncts_manage.Unlock()
	normal_tip(&fdbk, ctx, "<User Cancel Success>")
}

// 查看浏览历史
// Time_min = time/* & Time_max = time/*
// 2025-10-10 无需时分秒
func user_history(
	ctx *gin.Context,
	sql *gorm.DB,
	Secret []byte) {
	var (
		user_id            = controllers.Get_UID(ctx, Secret)
		history            []config.Visit_log
		time_min, time_max = ctx.Query("Time_min"), ctx.Query("Time_max")
		rst_hist           *gorm.DB
		fdbk               = "\n=====================\n" +
			"<SITE>url/user/history\n"
	)
	//检查时间格式
	if (!controllers.Time_check(time_min) && time_min != "*") || (!controllers.Time_check(time_max) && time_max != "*") {
		wrong_tip(&fdbk, ctx, "<Time Format Error>")
		ctx.Abort()
		return
	}
	//检查时间范围
	if !controllers.Time_less(time_min, time_max) {
		wrong_tip(&fdbk, ctx, "<Time Range Error>")
		ctx.Abort()
		return
	}
	//搜历史记录
	fmt.Println(time_min, time_max)
	if time_min != "*" && time_max != "*" {
		time_min += " 00:00:00"
		time_max += " 00:00:00"
		rst_hist = sql.Table("visit_logs").
			Select("PATH,TYPE,TIME").
			Where("USER_ID = ? AND TIME >= ? AND TIME <= ?", user_id, time_min, time_max).
			Find(&history)
	} else if time_min != "*" {
		time_min += " 00:00:00"
		rst_hist = sql.Table("visit_logs").
			Select("PATH,TYPE,TIME").
			Where("USER_ID = ? AND TIME >= ?").
			Find(&history)
	} else if time_max != "*" {
		time_max += " 00:00:00"
		rst_hist = sql.Table("visit_logs").
			Select("PATH,TYPE,TIME").
			Where("USER_ID = ? AND TIME <= ?", user_id, time_max).
			Find(&history)
	} else {
		rst_hist = sql.Table("visit_logs").
			Select("PATH,TYPE,TIME").
			Where("USER_ID = ?", user_id).
			Find(&history)
	}
	//判断
	if rst_hist.Error != nil && !errors.Is(rst_hist.Error, gorm.ErrRecordNotFound) {
		wrong_tip(&fdbk, ctx, "<Visit Log Find Error>")
		ctx.Abort()
		return
	}
	//显示
	fdbk += "===---__\n"
	if rst_hist.RowsAffected == 0 {
		fdbk += "(No history)\n"
	} else {
		for i, log := range history {
			fdbk += fmt.Sprintf("---Log (%d) Path-|%s|- Type<%s>  -Time[%s]\n",
				i+1, log.Path, log.Type, log.Time.Format("2006-01-02 15:04:05"))
		}
	}
	fdbk += "===___--\n"
	normal_tip(&fdbk, ctx, "<History Find Success>")
}

// 管理好友
// Friend_id = id & Action = mute/block/delete
func user_friends_manage(
	ctx *gin.Context,
	sql *gorm.DB,
	Secret []byte) {
	var (
		user_id        = controllers.Get_UID(ctx, Secret)
		friend_str     = ctx.Query("Friend_id")
		friend_id, err = strconv.Atoi(friend_str)
		friend_ship    config.Friend_ship
		status_user    int
		target_status  int
		isbig          bool
		updt           = true
		sm_id, bg_id   = pair_sort(user_id, friend_id)
		action         = ctx.Query("Action")
		rst_updt       *gorm.DB
		fdbk           = "\n=====================\n" +
			"<SITE>url/user/friends/manage\n"
	)
	//参数判断
	if err != nil {
		wrong_tip(&fdbk, ctx, "<Friend ID Error>")
		ctx.Abort()
		return
	}
	if action != "mute" && action != "block" && action != "delete" {
		wrong_tip(&fdbk, ctx, "<Action Error>")
		ctx.Abort()
		return
	}
	//先查关系
	rst_ship := sql.Table("friend_ships").
		Select("*").
		Where("SMALL_ID=? and BIG_ID=?", sm_id, bg_id).
		Find(&friend_ship)
	if rst_ship.Error != nil && !errors.Is(rst_ship.Error, gorm.ErrRecordNotFound) {
		wrong_tip(&fdbk, ctx, "<Friend Ship Find Error>")
		ctx.Abort()
		return
	}
	if rst_ship.RowsAffected == 0 {
		wrong_tip(&fdbk, ctx, "<Friend Ship Not Found>")
		ctx.Abort()
		return
	}
	if friend_ship.Status_small == 3 && friend_ship.Status_big == 3 {
		wrong_tip(&fdbk, ctx, "<Friend Ship Have Not Confirmed>")
		ctx.Abort()
		return
	}
	//给状态赋值
	if user_id == sm_id {
		status_user = friend_ship.Status_small
		isbig = false
	} else {
		status_user = friend_ship.Status_big
		isbig = true
	}
	//对应不同操作
	if action == "mute" {
		if status_user == 0 {
			fdbk += "(Had Muted Before Now)\n"
			updt = false
		} else {
			target_status = 0
		}
	} else if action == "block" {
		if status_user == 2 {
			fdbk += "(Had Blocked Before Now)\n"
			updt = false
		} else {
			target_status = 2
		}
	} else {
		rst_dlt := sql.Table("friend_ships").
			Where("SMALL_ID=? and BIG_ID = ?", sm_id, bg_id).
			Delete(&config.Friend_ship{})
		if rst_dlt.Error != nil {
			wrong_tip(&fdbk, ctx, "<Friend Ship Delete Error>")
			ctx.Abort()
			return
		}
	}
	if (action == "block" || action == "mute") && updt {
		if !isbig {
			rst_updt = sql.Table("friend_ships").
				Where("SMALL_ID = ? and BIG_ID = ?", sm_id, bg_id).
				Update("STATUS_SMALL", target_status)
		} else {
			rst_updt = sql.Table("friend_ships").
				Where("SMALL_ID = ? and BIG_ID = ?", sm_id, bg_id).
				Update("STATUS_BIG", target_status)
		}
		if rst_updt.Error != nil {
			wrong_tip(&fdbk, ctx, "<Friend Ship Update Error>")
			ctx.Abort()
			return
		}
	}
	normal_tip(&fdbk, ctx, "<Friend Ship Manage Success>")
}

// 联系好友
// Friend_id=id
func user_friends_contact(
	ctx *gin.Context,
	sql *gorm.DB,
	Secret []byte,
	cncts_manage *config.Solo_connect_manage,
	upgrader *websocket.Upgrader) {
	var (
		user_id              = controllers.Get_UID(ctx, Secret)
		friend_id, err0      = strconv.Atoi(ctx.Query("Friend_id"))
		sm_id, bg_id         int
		status_sm, status_bg int
		err                  error
		messages             []config.Solo_message
		user                 config.User
		friend               config.User
		friend_ship          config.Friend_ship
		fdbk                 = "\n=====================\n" +
			"<SITE>url/user/friends/contact\n"
		data_extra = ""
	)
	if err0 != nil {
		fmt.Println("Friend_id Error" + err0.Error())
		ctx.Abort()
		return
	}
	//找出好友
	rst := sql.Table("users").
		Select("ID,STATUS").
		Where("ID = ?", friend_id).
		First(&friend)
	if rst.Error != nil {
		fmt.Println("Friend Search Error" + rst.Error.Error())
		ctx.Abort()
		return
	}
	if friend.ID == 0 {
		fmt.Println("<Friend Not Found>")
		ctx.Abort()
		return
	}
	//更新用户ID对
	sm_id, bg_id = pair_sort(user_id, friend_id)
	//找出好友关系
	rst = sql.Table("friend_ships").
		Select("STATUS_SMALL,STATUS_BIG").
		Where("SMALL_ID = ? and BIG_ID = ?", sm_id, bg_id).
		First(&friend_ship)
	if rst.Error != nil {
		fmt.Println("Friend ship Search Error" + rst.Error.Error())
		ctx.Abort()
		return
	}
	if rst.RowsAffected == 0 {
		fmt.Println("Friend Ship Not Found")
		ctx.Abort()
		return
	}
	if cncts_manage.Connects[user_id] == nil {
		fmt.Println("Server Data Sync Error")
		ctx.Abort()
		return
	}
	//获取关系
	status_sm, status_bg = friend_ship.Status_small, friend_ship.Status_big
	//查找用户信息
	err_user := sql.Table("users").
		Select("ID,NAME,STATUS").
		Where("ID = ?", user_id).
		First(&user).Error
	if err_user != nil && !errors.Is(err_user, gorm.ErrRecordNotFound) {
		fmt.Println("User Find Error" + err_user.Error())
		ctx.Abort()
		return
	}
	if user.Name == "" {
		fmt.Println("<User Not Found>")
		ctx.Abort()
		return
	}
	//建立连接
	cncts_manage.Lock()
	if user_id < friend_id {
		cncts_manage.Connects[user_id].Status = status_sm
		cncts_manage.Connects[user_id].Status_friend = status_bg
	} else {
		cncts_manage.Connects[user_id].Status = status_bg
		cncts_manage.Connects[user_id].Status_friend = status_sm
	}
	cncts_manage.Connects[user_id].Conn = new(websocket.Conn)
	cncts_manage.Connects[user_id].Conn, err = upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	cncts_manage.Unlock()
	if err != nil {
		fmt.Println("Connect Error" + err.Error())
		ctx.Abort()
		return
	}

	//找出历史记录，嵌套查询，先找降序五十条，再升序排列
	err = sql.Table("(?) as sub", (sql.Table("solo_messages").
		Select("*").
		Order("TIME desc").
		Where("(USER_ID=? AND FRIEND_ID=?) OR (USER_ID=? AND FRIEND_ID=?)", sm_id, bg_id, bg_id, sm_id).
		Limit(50))).
		Order("TIME ASC").
		Find(&messages).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("Solo Message Search Error" + err.Error())
		ctx.Abort()
		return
	}
	if len(messages) == 0 {
		fmt.Println("<Solo Message Not Found>")
	}
	fdbk += "<STATUS>Success\n" +
		"<CONNECT>" +
		fmt.Sprintf("user[%d]->friend[%d]\n", user_id, friend_id) +
		"---------------------\n"

	if cncts_manage.Connects[user_id].Status_friend == 1 {
		if len(messages) == 0 {
			fdbk += "(No News Between You)\n"
		} else {
			for _, msg := range messages {
				fdbk += (string)(webskt.Solo_message_send(msg))
			}
		}
	} else {
		fdbk += "(Your Account Status is Abnormal)\n"
		data_extra = "(abnormal)"
	}
	cncts_manage.Connects[user_id].Lock()
	err = cncts_manage.Connects[user_id].Conn.WriteMessage(websocket.TextMessage, ([]byte)(fdbk))
	cncts_manage.Connects[user_id].Unlock()
	if err != nil {
		fmt.Println("<Websocket Welcome-Write Error>")
	}
	//fmt.Println(fdbk)  //释放聊天记录到终端

	fmt.Println("<Connection Success>" + fmt.Sprintf("user[%d]->friend[%d]%s\n", user_id, friend_id, data_extra))

	go webskt.Solo_read_worker(cncts_manage, user_id, friend_id, user.Name, sql)
}

// 好友申请
// Action = create/solve/check/recommend
func user_friends_applies(
	ctx *gin.Context,
	sql *gorm.DB,
	Secret []byte) {
	var (
		user_id = controllers.Get_UID(ctx, Secret)
		action  = ctx.Query("Action")
		fdbk    = "\n=====================\n" +
			"<SITE>url/user/friends/applies\n"
	)
	if action == "create" {
		applies_create(ctx, sql, user_id, &fdbk)
	} else if action == "solve" {
		applies_solve(ctx, sql, &fdbk)
	} else if action == "check" {
		applies_check(ctx, sql, user_id, &fdbk)
	} else if action == "recommend" {
		applies_recommend(ctx, sql, &fdbk)
	} else {
		wrong_tip(&fdbk, ctx, "Wrong Action")
		ctx.Abort()
		return
	}
}

// 创建好友申请
// Friend_id = id
func applies_create(
	ctx *gin.Context,
	sql *gorm.DB,
	user_id int,
	fdbk *string) {
	var (
		friend_str     = ctx.Query("Friend_id")
		friend_id, err = strconv.Atoi(friend_str)
		sm_id, bg_id   = pair_sort(user_id, friend_id)
		friend_ship    config.Friend_ship
	)
	//判断参数
	if err != nil {
		wrong_tip(fdbk, ctx, "Wrong Friend ID")
		ctx.Abort()
		return
	}
	//查ship
	sql.Table("friend_ships").
		Select("*").
		Where("SMALL_ID = ? AND BIG_ID = ?", sm_id, bg_id).
		Find(&friend_ship)
	if friend_ship.ID != 0 {
		wrong_tip(fdbk, ctx, "Friend Ship Exists")
		ctx.Abort()
		return
	}
	//执行创建
	friend_ship.Time = time.Now()
	friend_ship.Status_small = 3
	friend_ship.Status_big = 3
	friend_ship.Big_ID = bg_id
	friend_ship.Small_ID = sm_id
	rst_cr := sql.Table("friend_ships").
		Omit("ID").
		Create(&friend_ship)
	if rst_cr.Error != nil {
		wrong_tip(fdbk, ctx, "Apply Create Error")
		ctx.Abort()
		return
	}
	normal_tip(fdbk, ctx, "Apply Create Success")
}

// 处理好友申请
// ID=id & Solve = agree/refuse
func applies_solve(
	ctx *gin.Context,
	sql *gorm.DB,
	fdbk *string) {
	var (
		ID_str  = ctx.Query("ID")
		ID, err = strconv.Atoi(ID_str)
		solve   = ctx.Query("Solve")
		ship    config.Friend_ship
	)
	//检查参数
	if err != nil {
		wrong_tip(fdbk, ctx, "Wrong Apply ID")
		ctx.Abort()
		return
	}
	if solve != "agree" && solve != "refuse" {
		wrong_tip(fdbk, ctx, "Wrong Solve Method")
		ctx.Abort()
		return
	}
	//检查ship
	rst_ship := sql.Table("friend_ships").
		Select("*").
		Where("ID = ? AND STATUS_SMALL = 3 AND STATUS_BIG = 3", ID).
		Find(&ship)
	if rst_ship.Error != nil && !errors.Is(rst_ship.Error, gorm.ErrRecordNotFound) {
		wrong_tip(fdbk, ctx, "Apply Find Error")
		ctx.Abort()
		return
	}
	if rst_ship.RowsAffected == 0 {
		wrong_tip(fdbk, ctx, "Apply Find Error")
		ctx.Abort()
		return
	}
	//执行修改
	if solve == "agree" {
		rst_updt := sql.Table("friend_ships").
			Where("ID = ?", ID).
			Updates(map[string]int{"STATUS_SMALL": 1, "STATUS_BIG": 1})
		if rst_updt.Error != nil {
			wrong_tip(fdbk, ctx, "Apply Agree Error")
			ctx.Abort()
			return
		}
	} else {
		rst_delete := sql.Table("friend_ships").
			Where("ID = ?", ID).
			Delete(&config.Friend_ship{})
		if rst_delete.Error != nil {
			wrong_tip(fdbk, ctx, "Apply Refused Error")
			ctx.Abort()
			return
		}
	}
	normal_tip(fdbk, ctx, "Apply Solve Success")
}

// 查看好友申请
func applies_check(
	ctx *gin.Context,
	sql *gorm.DB,
	user_id int,
	fdbk *string) {
	var (
		friend_ships []config.Friend_ship
		sm_id, bg_id int
		friend_id    int
	)
	//直接查
	rst := sql.Table("friend_ships").
		Select("*").
		Where("STATUS_SMALL=3 AND STATUS_BIG=3").
		Find(&friend_ships)
	if rst.Error != nil && !errors.Is(rst.Error, gorm.ErrRecordNotFound) {
		wrong_tip(fdbk, ctx, "Applies Check Error")
		ctx.Abort()
		return
	}
	if len(friend_ships) == 0 {
		*fdbk += "(No Applies\n)"
	} else {
		for _, ship := range friend_ships {
			sm_id, bg_id = ship.Small_ID, ship.Big_ID
			if user_id == sm_id {
				friend_id = bg_id
			} else {
				friend_id = sm_id
			}
			*fdbk += fmt.Sprintf("---Apply_ID (%d) User(%d)->Friend(%d) --Time[%s]\n",
				ship.ID, user_id, friend_id, ship.Time.Format("2006-01-02 15:04:05"))
		}
	}
	normal_tip(fdbk, ctx, "Apply Check Success")
}

// 推荐好友
// Friend_id = id & Target_id = id
func applies_recommend(
	ctx *gin.Context,
	sql *gorm.DB,
	fdbk *string) {
	var (
		friend_str       = ctx.Query("Friend_id")
		target_str       = ctx.Query("Target_id")
		friend_id, err_f = strconv.Atoi(friend_str)
		target_id, err_t = strconv.Atoi(target_str)
		sm_id, bg_id     int
		friend_ship      config.Friend_ship
	)
	//判断参数准确
	if err_f != nil || err_t != nil {
		wrong_tip(fdbk, ctx, "ID Error")
		ctx.Abort()
		return
	}
	//判断关系状态
	sm_id, bg_id = pair_sort(friend_id, target_id)
	rst_ship := sql.Table("friend_ships").
		Select("*").
		Where("SMALL_ID = ? AND BIG_ID = ?", sm_id, bg_id).
		Find(&friend_ship)
	if rst_ship.Error != nil {
		wrong_tip(fdbk, ctx, "Ship Check Error")
		ctx.Abort()
		return
	}
	if friend_ship.ID != 0 {
		wrong_tip(fdbk, ctx, "Ship Exist Already")
		ctx.Abort()
		return
	}
	//新增ship
	friend_ship.Time = time.Now()
	friend_ship.Small_ID = sm_id
	friend_ship.Big_ID = bg_id
	friend_ship.Status_big = 3
	friend_ship.Status_small = 3
	rst_add := sql.Table("friend_ships").
		Omit("ID").
		Create(&friend_ship)
	if rst_add.Error != nil {
		wrong_tip(fdbk, ctx, "Recommend Add Error")
		ctx.Abort()
		return
	}
	normal_tip(fdbk, ctx, "Recommend Add Success")
}

// 好友非实时互动
// Action=like/comment 点赞或者评论
func user_friends_interaction(
	ctx *gin.Context,
	sql *gorm.DB,
	Secret []byte) {
	var (
		action  = ctx.Query("Action")
		user_id = controllers.Get_UID(ctx, Secret)
		fdbk    = "\n=====================\n" +
			"<SITE>url/user/friends/interaction\n"
	)
	if action == "like" {
		interaction_like(ctx, sql, &fdbk)
	} else if action == "comment" {
		interaction_comment(ctx, sql, user_id, &fdbk)
	} else {
		wrong_tip(&fdbk, ctx, "Wrong Action")
		ctx.Abort()
		return
	}
}

// 点赞 Friend_id = id
func interaction_like(
	ctx *gin.Context,
	sql *gorm.DB,
	fdbk *string) {
	var (
		friend_str       = ctx.Query("Friend_id")
		moment_str       = ctx.Query("Moment_id")
		friend_id, err_f = strconv.Atoi(friend_str)
		moment_id, err_m = strconv.Atoi(moment_str)
	)
	//检查参数
	if err_f != nil && err_m != nil {
		wrong_tip(fdbk, ctx, "ID Error")
		ctx.Abort()
		return
	}

	if err_f == nil {
		//执行主页点赞更新
		rst_like_f := sql.Table("user_profiles").
			Where("USER_ID = ?", friend_id).
			Update("POPULARITY", gorm.Expr("POPULARITY+1"))
		if rst_like_f.Error != nil {
			wrong_tip(fdbk, ctx, "Like Post Error")
			ctx.Abort()
			return
		}
		*fdbk += fmt.Sprintf("---You Have Given A Like to The Profile of User(%d)\n",
			friend_id)
	}

	if err_m == nil {
		//执行动态点赞更新
		rst_like_m := sql.Table("moments").
			Where("MOMENT_ID = ?", moment_id).
			Update("LIKES", gorm.Expr("LIKES+1"))
		if rst_like_m.Error != nil {
			wrong_tip(fdbk, ctx, "Like Post Error")
			ctx.Abort()
			return
		}
		*fdbk += fmt.Sprintf("---You Have Given A Like to The Profile of Moment(%d)\n",
			moment_id)
	}
	normal_tip(fdbk, ctx, "Post Like Success")
}

// 评论Moment_id = id  -d{content:string}
func interaction_comment(
	ctx *gin.Context,
	sql *gorm.DB,
	user_id int,
	fdbk *string) {
	var (
		moment_str       = ctx.Query("Moment_id")
		moment_id, err_m = strconv.ParseInt(moment_str, 10, 64)
		comment          config.Moment_Comment
	)
	//判断参数
	if err_m != nil {
		wrong_tip(fdbk, ctx, "ID Error")
		ctx.Abort()
		return
	}
	//直接更新动态的评论
	rst_cmt := sql.Table("moments").
		Where("MOMENT_ID = ? and STATUS = 1", moment_id).
		Update("COMMENT_NUM", gorm.Expr("COMMENT_NUM+1"))
	if rst_cmt.Error != nil {
		wrong_tip(fdbk, ctx, "Comment Post Error")
		ctx.Abort()
		return
	}
	//保存评论
	_ = ctx.ShouldBindJSON(&comment)
	comment.Time = time.Now()
	comment.Moment_ID = moment_id
	comment.User_ID = user_id
	rst_sav := sql.Table("moment_comments").
		Omit("COMMENT_ID").
		Create(&comment)
	if rst_sav.Error != nil {
		rst_rb := sql.Table("moments").
			Where("MOMENT_ID = ? and STATUS = 1", moment_id).
			Update("COMMENT_NUM", gorm.Expr("COMMENT_NUM-1"))
		if rst_rb.Error != nil {
			*fdbk += "(Comment_num Roll back Error)\n"
		} else {
			*fdbk += "(Comment_num Roll Back Success)\n"
		}
		wrong_tip(fdbk, ctx, "Comment Post Error")
		ctx.Abort()
		return
	}

	normal_tip(fdbk, ctx, "Post Comment Success")
}
