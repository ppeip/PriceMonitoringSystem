package router

import (
	"database/internal/repository"
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"net/http"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()
	router.Use(CORSMiddleware())
	//根据id获取一条数据
	router.GET("api/records/:id", func(context *gin.Context) {
		id := context.Param("id")
		var agPrice repository.AgPrice
		err := db.First(&agPrice, "id = ?", id).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.String(http.StatusNotFound, "pricerecord not found")
			return
		}
		context.JSON(http.StatusOK, agPrice)
	})
	//获取所有的数据
	router.GET("api/allrecords", func(context *gin.Context) {
		var agPrice []repository.AgPrice
		err := db.Find(&agPrice).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.String(http.StatusNotFound, "No pricerecord")
			return
		}
		context.JSON(http.StatusOK, agPrice)
	})
	//获取多少种类金属以及是什么
	router.GET("api/records/varieties", func(context *gin.Context) {
		var metalType []repository.MetalType
		type Res struct {
			Variety       string `json:"variety" `
			Varietycnname string `json:"cnname"`
		}
		var res []Res
		result := db.Model(&metalType).Find(&res)
		err := result.Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.String(http.StatusNotFound, "pricerecord not found")
			return
		}
		context.JSON(http.StatusOK, gin.H{
			"variety_num": result.RowsAffected,
			"variety":     res,
		})
	})
	//根据金属种类获取所有数据
	router.GET("api/records/variety/:variety", func(context *gin.Context) {
		variety := context.Param("variety")
		var agPrice []repository.AgPrice
		err := db.Find(&agPrice, "variety = ?", variety).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.String(http.StatusNotFound, "pricerecord not found")
			return
		}
		context.JSON(http.StatusOK, agPrice)
	})
	//返回七天内平均交易量
	router.GET("/api/aver/:variety", func(context *gin.Context) {
		variety := context.Param("variety")
		var agPrice []repository.AgPrice
		type Avg_Totalvols struct {
			Avg_totalvol float64
		}
		var avg_totalvols Avg_Totalvols
		//明天用gorm重构一下
		//err := db.Raw("SELECT avg(totalvol) AS avg_totalvol FROM (SELECT variety,created_at_time,totalvol  FROM ag_prices   WHERE `variety` = ? ORDER BY created_at_time DESC LIMIT 7) TEMP ", variety).Scan(&avg_totalvols).Error
		err := db.Table("(?) as temp ", db.Select("variety", "created_at_time", "totalvol").Find(&agPrice).Where("variety = ?", variety).
			Order("created_at_time DESC").Limit(7)).
			Select("avg(totalvol) AS avg_totalvol").
			Scan(&avg_totalvols).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.String(http.StatusNotFound, "pricerecord not found")
			return
		}
		context.JSON(http.StatusOK, gin.H{
			"vol": avg_totalvols.Avg_totalvol,
		})
	})
	//返回总交易量
	router.GET("api/total/:variety", func(context *gin.Context) {
		variety := context.Param("variety")
		var metalType repository.MetalType
		err := db.First(&metalType, "variety = ?", variety).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.String(http.StatusNotFound, "pricerecord not found")
			return
		}
		context.JSON(http.StatusOK, gin.H{
			"vol": metalType.Totalvol,
		})
	})
	//向数据库中添加新的数据
	router.POST("api/records", func(context *gin.Context) {
		var agPrice repository.AgPrice
		err := context.BindJSON(&agPrice)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid signin json", "json": agPrice})
			context.Abort()
			return
		}
		db.Create(&agPrice)
		context.Status(http.StatusCreated)
	})
	//根据id修改现有数据
	router.PUT("api/records/:id", func(context *gin.Context) {
		var agPrice repository.AgPrice
		id := context.Param("id")
		err := db.First(&agPrice, "id =?", id).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.String(http.StatusNotFound, "pricerecord not found")
			return
		}
		err = context.BindJSON(&agPrice)
		if err != nil {
			return
		}
		db.Save(&agPrice)
		context.Status(http.StatusOK)
	})
	//根据id删除现有的某条数据
	router.DELETE("api/records/:id", func(context *gin.Context) {
		id := context.Param("id")
		var agPrice repository.AgPrice
		err := db.First(&agPrice, "id =?", id).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.String(http.StatusNotFound, "pricerecord not found")
			return
		}
		err = db.Where("id = ? ", id).Delete(&agPrice).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.String(http.StatusNotFound, "pricerecord not found")
			return
		}
	})

	return router
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT,DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
