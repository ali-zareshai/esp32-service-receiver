package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
	"strconv"
	"time"
)

var ctx = context.Background()

const (
	AlertRedisChannels = "alert"
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

func SetRedis(key string, data interface{}, expireTime time.Duration) error {
	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = MyRedis.Set(ctx, key, value, expireTime*time.Second).Err()

	if err != nil {
		return err
	}

	return nil
}

func GetRedis(key string) (interface{}, error) {
	var result interface{}
	value, err := MyRedis.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("%v does not exists", key)
	} else if err != nil {
		return nil, err
	} else {
		json.Unmarshal([]byte(value), &result)
		return result, nil
	}
}

func PublishRedis(channel string, msg interface{}) {
	MyRedis.Publish(ctx, channel, msg)
}
