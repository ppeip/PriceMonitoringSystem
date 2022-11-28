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
	router.GET("/pricerecord/:no", func(context *gin.Context) {
		no := context.Param("no")
		var agPrice repository.AgPrice
		err := db.First(&agPrice, "no = ?", no).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.String(http.StatusNotFound, "pricerecord not found")
			return
		}
		context.JSON(http.StatusOK, agPrice)
	})

	router.POST("/pricerecord", func(context *gin.Context) {
		var agPrice repository.AgPrice
		err := context.BindJSON(&agPrice)
		if err != nil {
			println("json格式不正确")
			context.JSON(406, gin.H{"message": "Invalid signin json", "json": agPrice})
			context.Abort()
			return
		}
		db.Create(&agPrice)
		context.Status(http.StatusCreated)
	})

	router.PUT("/pricerecord/:no", func(context *gin.Context) {
		var agPrice repository.AgPrice
		no := context.Param("no")
		err := db.First(&agPrice, "no =?", no).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.String(http.StatusNotFound, "pricerecord not found")
			return
		}
		err = context.BindJSON(&agPrice)
		if err != nil {
			return
		}
		db.Save(&agPrice)
		context.Status(http.StatusNoContent)
	})

	router.DELETE("/pricerecord/:no", func(context *gin.Context) {
		no := context.Param("no")
		var agPrice repository.AgPrice
		err := db.Where("no = ? ", no).Delete(&agPrice).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.String(http.StatusNotFound, "pricerecord not found")
			return
		}
	})

	return router
}
