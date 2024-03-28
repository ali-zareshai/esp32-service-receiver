package controller

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"sensor_iot/Util"
	"sensor_iot/domain"
	"strconv"
	"time"
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
		context.JSON(http.StatusNotAcceptable, gin.H{"token": "", "error": err.Error()})
		return
	}
	user, err := checkLogin(req.Username, req.Password)
	if err != nil {
		context.JSON(http.StatusNotAcceptable, gin.H{"token": "", "error": err.Error()})
		return
	}
	token, err := generateToken(user)
	if err != nil {
		context.JSON(http.StatusNotAcceptable, gin.H{"token": "", "error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"token": token, "error": ""})
}

func verifyPassword(password, hashPass string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(password))
}

func checkLogin(username string, password string) (user *domain.UserModel, err error) {
	err = Util.MyDataBase.Model(domain.UserModel{}).Where("username=?", username).First(&user).Error
	if err != nil {
		log.Println(err.Error())
		return
	}

	err = verifyPassword(password, user.Password)
	if err != nil && errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		log.Println(err.Error())
		return
	}

	return
}

func generateToken(user *domain.UserModel) (string, error) {
	tokenLife, err := strconv.Atoi(os.Getenv("TOKEN_HOUR_LIFESPAN"))
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["name"] = user.Name
	claims["user_id"] = user.ID
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(tokenLife)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtToke, err := token.SignedString([]byte(os.Getenv("API_SECRET")))
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	return jwtToke, nil
}
