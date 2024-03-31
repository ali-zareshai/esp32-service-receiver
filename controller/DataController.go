package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sensor_iot/domain"
	"sensor_iot/utils"
)

func DataController(engine *gin.Engine) {
	r := engine.Group("/data")
	{
		r.Use(JwtMiddleware()).GET("", getData)
		r.POST("", addData)

	}
}

func getData(context *gin.Context) {
	res := utils.Gin{C: context}
	cached, _ := utils.GetRedis(context.Request.URL.String())
	if cached != nil {
		res.Response(http.StatusOK, "success", cached)
		return
	}
	var dataModels []domain.DataModel
	device := context.DefaultQuery("device", "")
	query := utils.MyDataBase.Model(&domain.DataModel{}).Scopes(utils.Paginate(context))
	if device != "" {
		query.Where(domain.DataModel{Device: device})
	}
	query.Find(&dataModels)
	utils.SetRedis(context.Request.URL.String(), dataModels, 300)
	res.Response(http.StatusOK, "success", dataModels)

}

func addData(context *gin.Context) {
	res := utils.Gin{C: context}
	var model domain.DataModel

	if err := context.ShouldBindJSON(&model); err != nil {
		res.Response(http.StatusNotAcceptable, "error", err.Error())
		return
	}

	model.Save()
	res.Response(http.StatusCreated, "success", model)

}
