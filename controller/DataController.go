package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sensor_iot/Util"
	"sensor_iot/domain"
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
	var dataModels []domain.DataModel
	device := context.DefaultQuery("device", "")
	query := Util.MyDataBase.Model(&domain.DataModel{}).Scopes(Util.Paginate(context))
	if device != "" {
		query.Where(domain.DataModel{Device: device})
	}
	query.Find(&dataModels)
	res.Response(http.StatusOK, "success", dataModels)

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
