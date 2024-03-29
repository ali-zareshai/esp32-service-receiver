package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sensor_iot/Util"
	"sensor_iot/domain"
)

func UserController(engine *gin.Engine) {
	userRoute := engine.Group("/user")
	{
		userRoute.POST("", addUser)
	}
}

func addUser(context *gin.Context) {
	res := Util.Gin{C: context}
	var user domain.UserModel
	if err := context.ShouldBindJSON(&user); err != nil {
		res.Response(http.StatusNotAcceptable, "error", err.Error())
		return
	}
	user.Save()
	res.Response(http.StatusOK, "success", user)
}
