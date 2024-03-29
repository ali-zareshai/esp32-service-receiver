package Util

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	MyDataBase *gorm.DB
	Logger     *logrus.Logger
)

type Response struct {
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type Gin struct {
	C *gin.Context
}

func (g *Gin) Response(codeStatus int, msg string, data interface{}) {
	g.C.JSON(codeStatus, Response{
		Msg:  msg,
		Data: data,
	})
}
