package responses

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"supxming.com/my_project/WeCircle_Project/Circle_root/config"
	"supxming.com/my_project/WeCircle_Project/Circle_root/controllers"
)

func rps_goods_vote(ctx *gin.Context, fdbk string, sql *gorm.DB, Secret []byte) {
	fdbk += "ACTION:Vote\n"
	var goods config.Goods_data //已有商品信息
	var chances int64           //当前已投票数
	var vote_num int64
	code := ctx.Query("Code")                                    //指定的商品编号
	todayst := time.Now().Truncate(24 * time.Hour)               //今日开始
	todayed := todayst.Add(24 * time.Hour).Add(-1 * time.Second) //今日结束
	var voter config.Voter
	var u_id = controllers.Get_UID(ctx, Secret) //用户id
	if u_id == -1 {
		fd_updt(&fdbk, u_id, code, "Wrong Token!")
		ctx.String(http.StatusForbidden, fdbk)
		ctx.Abort()
		return
	}

	//判断已投票数
	sql.Table("voter").
		Where("USER_ID = ? and TIME between ? and ?", u_id, todayst, todayed).
		Count(&chances)
	sql.Table("voter").
		Count(&vote_num)
	if chances >= 10 {
		fd_updt(&fdbk, u_id, code, "Vote chances exhausted!")
		ctx.String(http.StatusForbidden, fdbk)
		ctx.Abort()
		return
	}

	//判断商品是否存在
	err := sql.Table("goods_datas").
		Select("*").
		Where("GOODS_CODE = ?", code).
		First(&goods).Error
	if err != nil || goods.Goods_Code == "" {
		fd_updt(&fdbk, u_id, code, "Invalid goods code!")
		ctx.String(http.StatusForbidden, fdbk)
		ctx.Abort()
		return
	}
	//操作合法时
	//检查更新点赞数有没有问题
	err = sql.Table("goods_datas").
		Where("GOODS_CODE = ?", code).
		Update("LIKES", "LIKES-1").Error
	if err != nil {
		fd_updt(&fdbk, u_id, code, "Update Error!")
		ctx.String(http.StatusForbidden, fdbk)
		ctx.Abort()
		return
	}

	voter.Time = time.Now()
	voter.User_ID = u_id
	voter.Goods_code = code
	voter.Vote_ID = vote_num + 1
	err = sql.Table("voter").Create(&voter).Error
	if err != nil {
		fd_updt(&fdbk, u_id, code, "Create Error!")
		err = sql.Table("goods_datas").
			Where("GOODS_CODE = ?", code).
			Update("LIKES", "LIKES-1").Error
		if err != nil {
			ctx.String(http.StatusForbidden, fdbk+"(roll back fail)\n")
			ctx.Abort()
		} else {
			ctx.String(http.StatusForbidden, fdbk+"(roll back success)\n")
		}
		return
	}
	fdbk += "User:" + strconv.Itoa(u_id) + "\n" +
		"Code:" + code + "\n" +
		"STATUS:Applied\n" +
		"Support remain:" + strconv.Itoa(9-int(chances)) + "\n" +
		"=====================\n"
	fmt.Println("<goods vote success>")
	ctx.String(http.StatusOK, fdbk)
}

// tally的响应  rank - voters - goods(直接用goods的)
func rps_tally_rank(ctx *gin.Context, fdbk string, sql *gorm.DB) {
	var goods_10 []config.Goods_data

	goods_num, err_i := strconv.Atoi(ctx.Query("Goods_num"))
	if err_i != nil {
		fmt.Println("<Goods_num Error>")
		fdbk += "<Error>:Goods_num Error!\n" +
			"=====================\n"
		ctx.String(http.StatusBadRequest, fdbk)
		ctx.Abort()
		return
	}

	err := sql.Table("goods_datas").
		Order("LIKES desc").
		Limit(goods_num).
		Find(&goods_10).Error
	if err != nil {
		fdbk += "<STATUS>:Failed\n" +
			"<ERROR>" + err.Error() + "\n" +
			"=====================\n"
		ctx.String(http.StatusBadRequest, fdbk)
		ctx.Abort()
		fmt.Println("<goods rank fail>")
		return
	}
	fdbk += "STATUS:Success\n" +
		"<RANK DATA BELOW>\n"
	var n = 0
	for i, good := range goods_10 {
		fdbk += fmt.Sprintf("<%d> - %s(%s) Support[%d]\n",
			i+1, good.Name, good.Goods_Code, good.Likes)
		if good.Goods_Code != "" {
			n += 1
		}
	}
	if n < goods_num {
		fdbk += "-\n<<TIPS>>:There is only " + strconv.Itoa(n) + " records in table.\n"
	}
	fdbk += "=====================\n"
	fmt.Println("<goods rank success>")
	ctx.String(http.StatusOK, fdbk)
}

func rps_tally_voters(ctx *gin.Context, fdbk string, sql *gorm.DB) {
	var v_check = ctx.Query("Voters")
	var voters []config.Voter
	if v_check == "all" {
		fdbk += "<STATUS>:Applied\n" +
			"<NUM LIMIT>:100\n"
		sql.Table("voter").
			Select("*").
			Limit(100).
			Find(&voters)
		for i, voter := range voters {
			fdbk += fmt.Sprintf("<%d> - User:[%d]--Goods(%s) Time-%s\n",
				i+1, voter.User_ID, voter.Goods_code, voter.Time.Format("2006-01-02 15:04:05"))
		}
		fmt.Println("<voters check success>")
		ctx.String(http.StatusOK, fdbk)
		return
	} else if v_check == "" {
		fdbk += "<STATUS>:Refused\n" +
			"<REASON>:Empty range\n" +
			"=====================\n"
		fmt.Println("<voters check fail>:empty range")
		ctx.String(http.StatusBadRequest, fdbk)
		ctx.Abort()
		return
	} else {
		err := sql.Table("voter").
			Where("GOODS_CODE = ?", v_check).
			Limit(100).
			Find(&voters).Error
		if err != nil {
			fdbk += "<STATUS>:Failed\n" +
				"<ERROR>:Wrong code\n" +
				"=====================\n"
			fmt.Println("<voters check fail>:Wrong code")
			ctx.String(http.StatusBadRequest, fdbk)
			ctx.Abort()
			return
		}
		fdbk += "<STATUS>:Success\n" +
			"<CODE>:" + v_check + "\n" +
			"<NUM LIMIT>:100\n"
		for i, voter := range voters {
			fdbk += fmt.Sprintf("<%d> - User:[%d]--Goods(%s) Time-%v\n", i+1, voter.User_ID, voter.Goods_code, voter.Time)
		}
		if len(voters) == 0 {
			fdbk += "<WARNING>:Empty Records\n"
		}
		fdbk += "=====================\n"
		ctx.String(http.StatusOK, fdbk)
		fmt.Println("<voters check success>")
	}
}

func wrong_request(fdbk string, ctx *gin.Context, path string) {
	fdbk += "ACTION:Wrong\n" +
		"STATUS:Refused\n" +
		"======================\n"
	fmt.Println("SITE:url/rank/" +
		path + "\n" +
		"<ACTION WRONG>\n")
	ctx.String(http.StatusBadRequest, fdbk)
	ctx.Abort()
}

func Response_rank(typ string, sql *gorm.DB, Secret []byte) func(ctx *gin.Context) {

	if typ == "goods" {
		return func(ctx *gin.Context) {
			var fdbk = "======================\n" +
				"SITE:url/rank/goods\n"
			action := ctx.Query("Action")
			if action == "check" {
				rps_goods_check(ctx, fdbk, sql)
			} else if action == "vote" {
				rps_goods_vote(ctx, fdbk, sql, Secret)
			} else {
				wrong_request(fdbk, ctx, "goods")
				return
			}
		}
	} else if typ == "tally" {
		return func(ctx *gin.Context) {
			var fdbk = "======================\n" +
				"SITE:url/rank/tally\n"
			action := ctx.Query("Action")
			if action == "rank" {
				rps_tally_rank(ctx, fdbk, sql)
			} else if action == "goods" {
				rps_goods_check(ctx, fdbk, sql)
			} else if action == "voters" {
				rps_tally_voters(ctx, fdbk, sql)
			} else {
				wrong_request(fdbk, ctx, "tally")
			}
		}
	} else if typ == "rule" {
		return func(ctx *gin.Context) {
			var fdbk = "======================\n" +
				"SITE:url/rank/rule\n" +
				"<RULE BELOW>\n" +
				"Hey legends, ready to make your rig famous?\n" +
				"Welcome to [Gear Wars]—our living, breathing, ever-changing hall of fame for the sickest gaming setups on the planet.\n" +
				"Here’s the deal in bite-size, meme-friendly bullet points:\n" +
				"<1> 10 shots a day, every day!\n" +
				"Wake up, grab your coffee, spam those 10 votes on the mice, keyboards, headsets or GPUs that make you go “TAKE MY MONEY!”\n" +
				"<2> Spread the love or go all-in!\n" +
				"Split your 10 among ten different beasts or dump the whole stack on one dream machine—your call, your power.\n" +
				"<3> Midnight = refresh!\n" +
				"When the clock hits 00:00 server time, the ballot box resets. Miss a day? The hype train leaves without you—choo choo.\n" +
				"<4> No sock-puppets, no bots, no cap!\n" +
				"<5> Climb, fall, repeat!\n" +
				"Rankings update in real time. Yesterday’s king can be today’s meme—keep voting and watch the chaos unfold.\n" +
				"<6> Flex, discuss, meme it up!\n" +
				"Drop comments, post unboxings, roast or toast gear in the threads. Your vote + your voice = ultimate community combo.\n" +
				"So jump in, punch that vote button and let the world know which hardware truly slays.\n" +
				"May your frames be high and your latency low—happy voting!\n" +
				"======================\n"
			fmt.Println(fdbk)
			normal_tip(&fdbk, ctx, "<Rule Check Success>")
		}
	} else if typ == "public" {
		return func(ctx *gin.Context) {
			rps_goods_public(ctx, sql)
		}
	} else {
		return func(ctx *gin.Context) {
			var fdbk = "======================\n" +
				"SITE:url/rank\n"
			defer func() {
				if err := recover(); err != nil {
					fmt.Println(err)
					fdbk += "There is something wrong in your request!"
					wrong_tip(&fdbk, ctx, "<Wrong Action>")
					ctx.Abort()
					return
				}
			}()
			panic("Unsupported type of response: " + typ)
		}
	}
}

// public的响应
func rps_goods_public(ctx *gin.Context, sql *gorm.DB) {
	var (
		goods_code               = ctx.Query("Goods_code")
		name, profile, logo, url string
		price                    float32
		tp                       = ctx.Query("Type")
		score, err_s             = strconv.Atoi(ctx.Query("Score"))
		goods_data               config.Goods_data
		fdbk                     = "\n=====================\n" +
			"<SITE>url/chat/chatroom/ctcheck\n"
	)
	_ = ctx.ShouldBindJSON(&goods_data)
	if err_s != nil || score <= 0 {
		wrong_tip(&fdbk, ctx, "<Wrong Parameters>")
		ctx.Abort()
		return
	}

	name = goods_data.Name
	profile = goods_data.Profile
	logo = goods_data.Logo
	url = goods_data.Url
	price = goods_data.Price
	if name == "" || profile == "" ||
		logo == "" || url == "" ||
		tp == "" || goods_code == "" || price <= 0 {
		wrong_tip(&fdbk, ctx, "<Empty Parameters>")
		ctx.Abort()
		return
	}
	slct := sql.Table("goods_datas").
		Where("GOODS_CODE = ?", goods_code).
		Select("GOODS_CODE").
		Find(&goods_data)
	if slct.Error != nil && !errors.Is(slct.Error, gorm.ErrRecordNotFound) {
		wrong_tip(&fdbk, ctx, "<Sql Error>")
		ctx.Abort()
		return
	}
	if slct.RowsAffected != 0 {
		wrong_tip(&fdbk, ctx, "<GOODS_CODE Already Exist>")
		ctx.Abort()
		return
	}
	goods_data.Goods_Code = goods_code
	goods_data.Name = name
	goods_data.Profile = profile
	goods_data.Logo = logo
	goods_data.Score = score
	goods_data.Price = price
	goods_data.Url = url
	goods_data.Type = tp
	goods_data.Likes = 0

	add := sql.Table("goods_datas").
		Create(&goods_data)
	if add.Error != nil {
		wrong_tip(&fdbk, ctx, "<Upload Error>")
		ctx.Abort()
		return
	}
	normal_tip(&fdbk, ctx, "<Goods Public Success>")
}

// goods的响应  vote - check
func rps_goods_check(ctx *gin.Context, fdbk string, sql *gorm.DB) {
	fdbk += "ACTION:Check\n"
	var goods config.Goods_data
	code := ctx.Query("Code")
	err := sql.Table("goods_datas").Select("*").Where("GOODS_CODE = ?", code).First(&goods).Error
	if err != nil || goods.Goods_Code == "" {
		fmt.Println("<goods check error>:Something wrong in goods_code.")
		fdbk += "STATUS:Refused\n" +
			"======================\n"
		ctx.String(http.StatusBadRequest, fdbk)
		ctx.Abort()
		return
	}

	fdbk += "STATUS:OK\n" +
		"NAME:" + goods.Name + "\n" +
		"<CODE> " + goods.Goods_Code + "\n" +
		"<SUPPORT>" + strconv.Itoa(goods.Likes) + "\n" +
		"INTRODUCTION:" + goods.Profile + "\n" +
		"======================\n"
	ctx.String(http.StatusOK, fdbk)
	fmt.Println("<goods check success>")
}

func fd_updt(fdbk *string, u_id int, code string, reason string) {
	*fdbk += "User:" + strconv.Itoa(u_id) + "\n" +
		"Code:" + code + "\n" +
		"STATUS:Refused\n" +
		"<Reason>" + reason + "\n" +
		"=====================\n"
	fmt.Println("<goods vote fail>:" + reason + "\n")
}
