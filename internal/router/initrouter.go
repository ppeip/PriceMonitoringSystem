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
		var metalExchangeRecords repository.MetalExchangeRecords
		err := db.First(&metalExchangeRecords, "id = ?", id).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.String(http.StatusNotFound, "pricerecord not found")
			return
		}
		context.JSON(http.StatusOK, metalExchangeRecords)
	})
	//获取所有的数据
	router.GET("api/allrecords", func(context *gin.Context) {
		var metalExchangeRecords []repository.MetalExchangeRecords
		err := db.Find(&metalExchangeRecords).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.String(http.StatusNotFound, "No pricerecord")
			return
		}
		context.JSON(http.StatusOK, metalExchangeRecords)
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
		var metalExchangeRecords []repository.MetalExchangeRecords
		err := db.Find(&metalExchangeRecords, "variety = ?", variety).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.String(http.StatusNotFound, "pricerecord not found")
			return
		}
		context.JSON(http.StatusOK, metalExchangeRecords)
	})
	//返回七天内平均交易量
	router.GET("/api/aver/:variety", func(context *gin.Context) {
		variety := context.Param("variety")
		var metalExchangeRecords []repository.MetalExchangeRecords
		type Avg_Totalvols struct {
			Avg_totalvol float64
		}
		var avg_totalvols Avg_Totalvols
		//明天用gorm重构一下
		//err := db.Raw("SELECT avg(totalvol) AS avg_totalvol FROM (SELECT variety,created_at_time,totalvol  FROM ag_prices   WHERE `variety` = ? ORDER BY created_at_time DESC LIMIT 7) TEMP ", variety).Scan(&avg_totalvols).Error
		err := db.Table("(?) as temp ", db.Select("variety", "created_at_time", "totalvol").Find(&metalExchangeRecords).Where("variety = ?", variety).
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
	//向数据库MetalExchangeRecords表中添加新的数据
	router.POST("api/records", func(context *gin.Context) {
		var metalExchangeRecords repository.MetalExchangeRecords
		err := context.BindJSON(&metalExchangeRecords)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid signin json", "json": metalExchangeRecords})
			context.Abort()
			return
		}
		db.Create(&metalExchangeRecords)
		context.Status(http.StatusCreated)
	})
	//向数据库ExchangRecords表中添加新的数据
	router.POST("api/exchangrecords", func(context *gin.Context) {
		var exchangrecords repository.ExchangRecords
		err := context.BindJSON(&exchangrecords)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid signin json", "json": exchangrecords})
			context.Abort()
			return
		}
		db.Create(&exchangrecords)
		context.Status(http.StatusCreated)
	})
	//添加新账户
	router.POST("api/account", func(context *gin.Context) {
		var acconut repository.Account
		err := context.BindJSON(&acconut)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid signin json", "json": acconut})
			context.Abort()
			return
		}
		db.Create(&acconut)
		context.Status(http.StatusCreated)
	})
	//添加新用户
	router.POST("api/user", func(context *gin.Context) {
		var user repository.User
		err := context.BindJSON(&user)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid signin json", "json": user})
			context.Abort()
			return
		}
		db.Create(&user)
		context.Status(http.StatusCreated)
	})
	//根据id修改MetalExchangeRecords现有数据
	router.PUT("api/records/:id", func(context *gin.Context) {
		var metalExchangeRecords repository.MetalExchangeRecords
		id := context.Param("id")
		err := db.First(&metalExchangeRecords, "id =?", id).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.String(http.StatusNotFound, "pricerecord not found")
			return
		}
		err = context.BindJSON(&metalExchangeRecords)
		if err != nil {
			return
		}
		db.Save(&metalExchangeRecords)
		context.Status(http.StatusOK)
	})
	//根据records_id修改MetalPrice现有数据
	router.PUT("api/MetalPrice/:records_id", func(context *gin.Context) {
		var exchangrecords repository.ExchangRecords
		id := context.Param("records_id")
		err := db.First(&exchangrecords, "records_id =?", id).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.String(http.StatusNotFound, "pricerecord not found")
			return
		}
		err = context.BindJSON(&exchangrecords)
		if err != nil {
			return
		}
		db.Save(&exchangrecords)
		context.Status(http.StatusOK)
	})
	//根据id删除MetalExchangeRecords现有的某条数据
	router.DELETE("api/records/:id", func(context *gin.Context) {
		id := context.Param("id")
		var metalExchangeRecords repository.MetalExchangeRecords
		err := db.First(&metalExchangeRecords, "id =?", id).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.String(http.StatusNotFound, "pricerecord not found")
			return
		}
		err = db.Where("id = ? ", id).Delete(&metalExchangeRecords).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.String(http.StatusNotFound, "pricerecord not found")
			return
		}
	})
	//根据id删除ExchangRecords现有的某条数据
	router.DELETE("api/exchangrecords/:records_id", func(context *gin.Context) {
		id := context.Param("records_id")
		var exchangrecords repository.ExchangRecords
		err := db.First(&exchangrecords, "records_id =?", id).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.String(http.StatusNotFound, "pricerecord not found")
			return
		}
		err = db.Where("records_id = ? ", id).Delete(&exchangrecords).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.String(http.StatusNotFound, "pricerecord not found")
			return
		}
	})
	//根据id删除某个账户
	router.DELETE("api/account/:account_id", func(context *gin.Context) {
		id := context.Param("account_id")
		var user repository.User
		err := db.Where("account_id = ? ", id).Delete(&user).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.String(http.StatusNotFound, "account not found")
			return
		}
	})
	//根据id删除某个用户
	router.DELETE("api/user/:user_id", func(context *gin.Context) {
		id := context.Param("user_id")
		var user repository.User
		err := db.Where("user_id = ? ", id).Delete(&user).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.String(http.StatusNotFound, "user not found")
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
