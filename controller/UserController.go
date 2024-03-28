package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sensor_iot/domain"
)

func UserController(engine *gin.Engine) {
	userRoute := engine.Group("/user")
	{
		userRoute.POST("", addUser)
	}
}

func addUser(context *gin.Context) {
	var user domain.UserModel
	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}
	user.Save()
	context.JSON(http.StatusOK, gin.H{"error": ""})
}
