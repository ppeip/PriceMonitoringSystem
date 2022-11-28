package repository

import (
	"gorm.io/gorm"
)

type AgPrice struct {
	No            int64   `json:"no" gorm:"primaryKey;autoIncrement " `
	Variety       string  `json:"variety"`   /*品种:"Ag(T+D)"*/
	Latestpri     float32 `json:"latestpri"` /*最新价:"6585.00"*/
	Openpri       float32 `json:"openpri"`   /*开盘价:"6712.00"*/
	Maxpri        float32 `json:"maxpri"`    /*最高价:"6721.00"*/
	Minpri        float32 `json:"minpri"`    /*最低价:"6581.00"*/
	Limit         string  `json:"limit"`     /*涨跌幅:"-1.98%"*/
	Yespri        float32 `json:"yespri"`    /*昨收价:"6718.00"*/
	Totalvol      float64 `json:"totalvol"`  /*总成交量:"1564524.0000"*/
	CreatedAtTime string  `json:"time"`      /*更新时间:"2012-12-19 15:29:59"*/
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&AgPrice{})
}
