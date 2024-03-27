package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sensor_iot/domain"
	"sensor_iot/service"
	"strconv"
)

func DataController(engine *gin.Engine) {
	r := engine.Group("/data")
	{
		r.GET("", func(context *gin.Context) {
			count := context.DefaultQuery("count", "100")
			if c, err := strconv.Atoi(count); err == nil {
				data := service.FindData(c)
				context.JSON(http.StatusOK, data)
				return
			}

			context.JSON(http.StatusNotAcceptable, gin.H{"error": "error"})
		})

		r.POST("", func(context *gin.Context) {
			var req domain.DataJsonRequest

			if err := context.ShouldBindJSON(&req); err != nil {
				context.JSON(http.StatusNotAcceptable, gin.H{"status": "error", "error": err})
				return
			}
			if service.AddData(req) {
				context.JSON(http.StatusCreated, gin.H{"status": "success", "error": ""})
			}
		})

	}
}
