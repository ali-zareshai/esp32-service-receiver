package main

import (
	"github.com/gin-gonic/gin"
	"sensor_iot/controller"
)

func main() {
	r := gin.Default()
	controller.DataController(r)

	r.Run()
}
