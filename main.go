package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"sensor_iot/Util"
	"sensor_iot/controller"
	"sensor_iot/domain"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err.Error())
	}

	var err error
	Util.MyDataBase, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	Util.MyDataBase.AutoMigrate(&domain.DataModel{}, &domain.UserModel{})

	r := gin.Default()
	controller.DataController(r)
	controller.UserController(r)
	controller.AuthController(r)

	r.Run(":" + os.Getenv("SERVER_PORT"))
}
