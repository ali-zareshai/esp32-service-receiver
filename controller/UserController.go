package controller

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/validator.v2"
	"net/http"
	"sensor_iot/domain"
	"sensor_iot/utils"
)

func UserController(engine *gin.Engine) {
	userRoute := engine.Group("/user")
	{
		userRoute.POST("", addUser)
	}
}

func addUser(context *gin.Context) {
	res := utils.Gin{C: context}
	var user domain.UserModel
	if err := context.ShouldBindJSON(&user); err != nil {
		res.Response(http.StatusNotAcceptable, "error", err.Error())
		return
	}

	if err := validator.Validate(user); err != nil {
		res.Response(http.StatusBadRequest, "error", err.Error())
		return
	}
	user.Save()
	res.Response(http.StatusOK, "success", user)
}
