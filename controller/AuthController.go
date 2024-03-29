package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"sensor_iot/Util"
	"sensor_iot/domain"
	"strconv"
	"strings"
	"time"
)

var logger = Util.Logger

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
	res := Util.Gin{C: context}
	var req LoginRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		res.Response(http.StatusNotAcceptable, err.Error(), nil)
		return
	}
	user, err := checkLogin(req.Username, req.Password)
	if err != nil {
		res.Response(http.StatusNotAcceptable, err.Error(), nil)
		return
	}
	user.Password = ""
	token, err := generateToken(user)
	if err != nil {
		res.Response(http.StatusNotAcceptable, err.Error(), nil)
		return
	}
	res.Response(http.StatusOK, "success", gin.H{"token": token})
}

func FindCurrentUser(c *gin.Context) {
	user, err := ExtractUser(c)
	res := Util.Gin{C: c}
	if err != nil {
		res.Response(http.StatusUnauthorized, err.Error(), nil)
		return
	}
	res.Response(http.StatusOK, "success", gin.H{"user": user})
}

func verifyPassword(password, hashPass string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(password))
}

func checkLogin(username string, password string) (user *domain.UserModel, err error) {
	err = Util.MyDataBase.Model(domain.UserModel{}).Where("username=?", username).First(&user).Error
	if err != nil {
		logger.Error(err)
		return
	}

	err = verifyPassword(password, user.Password)
	if err != nil && errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		logger.Error(err)
		return
	}

	return
}

func generateToken(user *domain.UserModel) (string, error) {
	tokenLife, err := strconv.Atoi(os.Getenv("TOKEN_HOUR_LIFESPAN"))
	if err != nil {
		logger.Error(err)
		return "", err
	}
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	jsonUser, err := json.Marshal(user)
	if err != nil {
		logger.Error(err)
		return "", err
	}
	claims["user"] = string(jsonUser)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(tokenLife)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtToke, err := token.SignedString([]byte(os.Getenv("API_SECRET")))
	if err != nil {
		logger.Error(err)
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
		err = json.Unmarshal([]byte(userJson.(string)), &user)
	}
	return
}

func JwtMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		user, err := ExtractUser(context)
		if err != nil {
			res := Util.Gin{C: context}
			res.Response(http.StatusUnauthorized, err.Error(), nil)
			context.Abort()
			return
		}
		context.Set("user_name", user.Name)
		context.Set("user_id", user.ID)
		context.Set("user_role", user.Role)
		context.Next()
	}
}
