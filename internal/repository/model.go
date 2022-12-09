package repository

import (
	"fmt"

	"gorm.io/gorm"
)

type AgPrice struct {
	Id            int64   `json:"id" gorm:"primaryKey;autoIncrement " `
	Variety       string  `json:"variety"` /*品种:"Ag(T+D)"*/
	Varietycnname string  `json:"cnname"`
	Latestpri     float32 `json:"latestpri"`                 /*最新价:"6585.00"*/
	Openpri       float32 `json:"openpri"`                   /*开盘价:"6712.00"*/
	Maxpri        float32 `json:"maxpri"`                    /*最高价:"6721.00"*/
	Minpri        float32 `json:"minpri"`                    /*最低价:"6581.00"*/
	Limit         string  `json:"limit"`                     /*涨跌幅:"-1.98%"*/
	Yespri        float32 `json:"yespri"`                    /*昨收价:"6718.00"*/
	Totalvol      float64 `json:"totalvol"`                  /*总成交量:"1564524.0000"*/
	CreatedAtTime string  `json:"time" gorm:"type:datetime"` /*更新时间:"2012-12-19 15:29:59"*/
}

type MetalType struct {
	Id            int64   `json:"id"`
	Variety       string  `json:"variety" gorm:"primaryKey"`
	Totalvol      float64 `json:"totalvol"`
	Varietycnname string  `json:"cnname"`
}

func (agPrice *AgPrice) AfterCreate(db *gorm.DB) (err error) {
	var metalType MetalType
	fmt.Println(agPrice.Totalvol)
	db.Model(&metalType).Where("variety = ?", agPrice.Variety).Update("totalvol", gorm.Expr("totalvol + ?", agPrice.Totalvol))
	fmt.Println("修改成功")
	return nil
}

func (agPrice *AgPrice) BeforeDelete(db *gorm.DB) (err error) {
	var metalType MetalType
	fmt.Println(agPrice.Totalvol)
	db.Model(&metalType).Where("variety = ?", agPrice.Variety).Update("totalvol", gorm.Expr("totalvol - ?", agPrice.Totalvol))
	fmt.Println("修改成功")
	return nil
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&AgPrice{}, &MetalType{})
}
