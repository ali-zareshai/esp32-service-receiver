package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"sensor_iot/Util"
	"sensor_iot/domain"
	"sensor_iot/service"
	"strconv"
)

func DataController(engine *gin.Engine) {
	r := engine.Group("/data")
	{
		r.Use(JwtMiddleware()).GET("", getData)
		r.POST("", addData)

	}
}

func getData(context *gin.Context) {
	res := Util.Gin{C: context}
	count := context.DefaultQuery("count", "100")
	device := context.DefaultQuery("device", "")
	if c, err := strconv.Atoi(count); err == nil {
		data := service.FindData(c, device)
		res.Response(http.StatusOK, "success", data)
		return
	}
	res.Response(http.StatusNotAcceptable, "error", errors.New("wrong count number"))
}

func addData(context *gin.Context) {
	res := Util.Gin{C: context}
	var model domain.DataModel

	if err := context.ShouldBindJSON(&model); err != nil {
		res.Response(http.StatusNotAcceptable, "error", err.Error())
		return
	}

	model.Save()
	res.Response(http.StatusCreated, "success", model)

}
