package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io"
	"os"
	"path/filepath"
	"sensor_iot/Util"
	"sensor_iot/controller"
	"sensor_iot/domain"
	"time"
)

func main() {
	InitLogger()
	logger := Util.Logger
	//logger.Error(errors.New("test shod22f"))
	if err := godotenv.Load(); err != nil {
		logger.Error(err)
	}

	var err error
	Util.MyDataBase, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		logger.Error(err)
	}
	Util.MyDataBase.AutoMigrate(&domain.DataModel{}, &domain.UserModel{})

	r := gin.Default()
	r.Use(cors.Default(), Util.RateLimitMiddleware())

	controller.DataController(r)
	controller.UserController(r)
	controller.AuthController(r)

	r.Run(":" + os.Getenv("SERVER_PORT"))
}

func InitLogger() {
	lumberjackLogger := &lumberjack.Logger{
		Filename:   filepath.ToSlash("./logs/" + time.Now().Format("2006-01-02") + "_app.log"),
		MaxSize:    1, // MB
		MaxBackups: 10,
		MaxAge:     30,   // days
		Compress:   true, // disabled by default
	}

	multiWriter := io.MultiWriter(os.Stderr, lumberjackLogger)

	logFormatter := new(logrus.TextFormatter)
	logFormatter.TimestampFormat = time.RFC3339
	logFormatter.FullTimestamp = true

	Util.Logger = logrus.New()
	Util.Logger.SetFormatter(logFormatter)
	Util.Logger.SetLevel(logrus.InfoLevel)
	Util.Logger.SetOutput(multiWriter)
}
