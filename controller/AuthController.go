package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func AuthController(engine *gin.Engine) {
	authRoute := engine.Group("/auth")
	authRoute.POST("", LoginUser)
}

func LoginUser(context *gin.Context) {
	var req LoginRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}
	user, err := CheckLogin(req.Username, req.Password)
	if err != nil {
		context.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"user": user})
}
