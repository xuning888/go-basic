package dal

import (
	"github.com/redis/go-redis/v9"
	"go-basic/webbook/config"
)

func InitRedis() redis.Cmdable {
	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Add,
	})
	return redisClient
}
