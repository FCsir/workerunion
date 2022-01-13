package pkg

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

type RedisHelper struct {
	*redis.Client
}

var RedisClient *redis.Client

func SetUpRedis() {
	redisConfig := viper.GetStringMap("redis")
	fmt.Println("redis: ", redisConfig)
	port := redisConfig["port"]
	host := redisConfig["host"]
	addr := host.(string) + ":" + port.(string)
	RedisClient = redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   0,
	})
}
