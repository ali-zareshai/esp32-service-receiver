package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"sensor_iot/Util"
	"sensor_iot/domain"
	"strconv"
	"strings"
	"time"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func AuthController(engine *gin.Engine) {
	authRoute := engine.Group("/auth")
	authRoute.POST("", LoginUser)
	authRoute.GET("/user", FindCurrentUser)
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

func FindCurrentUser(c *gin.Context) {
	user, err := ExtractUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
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
	jsonUser, err := json.Marshal(user)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	claims["user"] = string(jsonUser)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(tokenLife)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtToke, err := token.SignedString([]byte(os.Getenv("API_SECRET")))
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	return jwtToke, nil
}

func extractToken(c *gin.Context) string {
	token := c.Request.Header.Get("Authorization")
	if len(strings.Split(token, " ")) == 2 {
		return strings.Split(token, " ")[1]
	}
	return ""
}

func ExtractUser(c *gin.Context) (user domain.UserModel, err error) {
	tokenString := extractToken(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		userJson := claims["user"]
		log.Println(userJson)
		err = json.Unmarshal([]byte(userJson.(string)), &user)
	}
	return
}
