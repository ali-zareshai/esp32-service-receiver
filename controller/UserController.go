package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
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
	var user domain.UserModel
	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}
	user.Save()
	context.JSON(http.StatusOK, gin.H{"error": ""})
}

func VerifyPassword(password, hashPass string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(password))
}

func CheckLogin(username string, password string) (user *domain.UserModel, err error) {
	err = Util.MyDataBase.Model(domain.UserModel{}).Where("username=?", username).First(&user).Error
	if err != nil {
		log.Println(err.Error())
		return
	}

	err = VerifyPassword(password, user.Password)
	if err != nil && errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		log.Println(err.Error())
		return
	}

	return
}
