package Util

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
	"strconv"
)

func ConnectToRedis() {
	dbNumber, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		dbNumber = 0
	}
	MyRedis = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_ADDRESS"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       dbNumber,
	})
}
