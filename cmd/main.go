package main

import (
	"database/internal/repository"
	"database/internal/router"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	//日志
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // Log level
			Colorful:      false,       // 禁用彩色打印
		},
	)

	// 连接数据库
	dsn := "wawaro:123qwe@tcp(127.0.0.1:3306)/monitoringsystem?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		panic(err)
	}
	err = repository.AutoMigrate(db)
	if err != nil {
		panic(err)
	}
	router := router.SetupRouter(db)
	router.Run()

}
