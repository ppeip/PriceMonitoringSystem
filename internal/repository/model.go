package repository

import (
	"fmt"

	"gorm.io/gorm"
)

type MetalExchangeRecords struct {
	Id            int64   `json:"id" gorm:"primaryKey;autoIncrement " ` /*日交易记录id:"1"*/
	Variety       string  `json:"variety"`                              /*品种:"Au"*/
	Openpri       float32 `json:"openpri"`                              /*开盘价:"6712.00"*/
	Maxpri        float32 `json:"maxpri"`                               /*最高价:"6721.00"*/
	Minpri        float32 `json:"minpri"`                               /*最低价:"6581.00"*/
	Limit         string  `json:"limit"`                                /*涨跌幅:"-1.98%"*/
	Yespri        float32 `json:"yespri"`                               /*昨收价:"6718.00"*/
	Totalvol      float64 `json:"totalvol"`                             /*日成交量:"1564524.0000"*/
	CreatedAtTime string  `json:"time" gorm:"type:datetime"`            /*创建时间:"2012-12-19"*/
}

type MetalType struct {
	VarietyId     int64   `json:"variety_id"`                /*金属种类记录id:"1"*/
	Variety       string  `json:"variety" gorm:"primaryKey"` /*品种:"Au"*/
	Totalvol      float64 `json:"totalvol"`                  /*总成交量:"1564524.0000"*/
	Varietycnname string  `json:"cnname"`                    /*品种中文名:"金"*/
}

type MetalPrice struct {
	RecordsId   int64   `json:"records_id" gorm:"primaryKey"`
	Variety     string  `json:"variety"`                   /*品种:"Au"*/
	Exchangepri float32 `json:"exchangepri"`               /*成交价:"6712.00"*/
	Vol         float32 `json:"vol"`                       /*成交量:"5000.0000"*/
	Time        string  `json:"time" gorm:"type:datetime"` /*成交时间:"2012-12-19"*/
}

type ExchangRecords struct {
	RecordsId   int64 `json:"records_id" gorm:"primaryKey;"` /*交易记录id:"1"*/
	UserId      int64 `json:"user_id" `                      /*用户所id:"1"*/
	ExchangesId int64 `json:"exchanges_id" `                 /*交易所id:"1"*/
}
type User struct {
	UserId      int64  `json:"user_id" gorm:"primaryKey;"` /*用户所id:"1"*/
	ExchangesId int64  `json:"exchanges_id" `              /*交易所id:"1"*/
	Name        string `json:"name" `                      /*用户姓名:"1"*/
}
type ExchangsInfo struct {
	ExchangesId   int64  `json:"exchanges_id" gorm:"primaryKey;"` /*交易所id:"1"*/
	ExchangesName string `json:"exchanges_name"`                  /*交易名字id:"上海证券交易所"*/
}

func (metalPrice *MetalPrice) AfterCreate(db *gorm.DB) (err error) {
	var metalType MetalType
	fmt.Println(metalPrice.Vol)
	db.Model(&metalType).Where("variety = ?", metalPrice.Variety).Update("totalvol", gorm.Expr("totalvol + ?", metalType.Totalvol))
	fmt.Println("修改成功")
	return nil
}

func (metalPrice *MetalPrice) BeforeDelete(db *gorm.DB) (err error) {
	var metalType MetalType
	fmt.Println(metalPrice.Vol)
	db.Model(&metalType).Where("variety = ?", metalPrice.Variety).Update("totalvol", gorm.Expr("totalvol - ?", metalType.Totalvol))
	fmt.Println("修改成功")
	return nil
}

func (metalExchangeRecords *MetalExchangeRecords) AfterCreate(db *gorm.DB) (err error) {
	var metalType MetalType
	fmt.Println(metalExchangeRecords.Totalvol)
	db.Model(&metalType).Where("variety = ?", metalExchangeRecords.Variety).Update("totalvol", gorm.Expr("totalvol + ?", metalType.Totalvol))
	fmt.Println("修改成功")
	return nil
}

func (metalExchangeRecords *MetalExchangeRecords) BeforeDelete(db *gorm.DB) (err error) {
	var metalType MetalType
	fmt.Println(metalExchangeRecords.Totalvol)
	db.Model(&metalType).Where("variety = ?", metalExchangeRecords.Variety).Update("totalvol", gorm.Expr("totalvol - ?", metalType.Totalvol))
	fmt.Println("修改成功")
	return nil
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&MetalExchangeRecords{}, &MetalType{}, &MetalPrice{}, &ExchangRecords{}, &User{}, &ExchangsInfo{})
}
