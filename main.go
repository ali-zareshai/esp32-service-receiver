package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"sensor_iot/Util"
	"sensor_iot/controller"
	"sensor_iot/domain"
)

func main() {
	var err error
	Util.MyDataBase, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	Util.MyDataBase.AutoMigrate(&domain.DataModel{})

	r := gin.Default()
	controller.DataController(r)

	r.Run()
}
