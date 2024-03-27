package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sensor_iot/domain"
	"sensor_iot/service"
)

func DataController(engine *gin.Engine) {
	r := engine.Group("/data")
	{
		r.POST("", func(context *gin.Context) {
			var req domain.DataJsonRequest

			if err := context.ShouldBindJSON(&req); err != nil {
				context.JSON(http.StatusNotAcceptable, gin.H{"status": "error", "error": err})
			}
			if service.AddData(req) {
				context.JSON(http.StatusCreated, gin.H{"status": "success", "error": ""})
			}
		})
	}
}
