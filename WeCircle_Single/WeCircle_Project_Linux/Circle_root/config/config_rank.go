package config

import "time"

// 排行榜相关========================================
type Voter struct {
	Vote_ID    int64     `gorm:"column:VOTE_ID;type:bigint;primary key" json:"vote_id"`
	User_ID    int       `gorm:"column:USER_ID;not null" json:"user_id"`
	Goods_code string    `gorm:"column:GOODS_CODE;not null" json:"goods_code"`
	Time       time.Time `gorm:"column:TIME;not null" json:"time"`
}
type Goods_data struct {
	Goods_Code string  `gorm:"column:GOODS_CODE;type:varchar(50);primary_key" json:"goods_code"`
	Name       string  `gorm:"column:NAME;type:varchar(50);not null;unique_index" json:"name"`
	Logo       string  `gorm:"column:LOGO;type:varchar(20);not null" json:"logo"`
	Price      float32 `gorm:"column:PRICE;not null" json:"price"`
	Type       string  `gorm:"column:TYPE;type:varchar(30);not null" json:"type"`
	Likes      int     `gorm:"column:LIKES;type:int;not null" json:"likes"`
	Score      int     `gorm:"column:SCORE;type:int;not null" json:"score"` //between 1 and 5
	Profile    string  `gorm:"column:PROFILE;type:varchar(150);not null" json:"profile"`
	Url        string  `gorm:"column:URL;type:varchar(150);not null" json:"url"`
}
