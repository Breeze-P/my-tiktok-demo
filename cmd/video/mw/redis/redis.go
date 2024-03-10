package redis

import (
	"my-tiktok/pkg/constants"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	expireTime  = time.Hour * 1
	rdbFollows  *redis.Client
	rdbFavorite *redis.Client
)

func Init() {
	rdbFollows = redis.NewClient(&redis.Options{
		Addr:     constants.RedisAddr,
		Password: constants.RedisPassword,
		DB:       0,
	})
	rdbFavorite = redis.NewClient(&redis.Options{
		Addr:     constants.RedisAddr,
		Password: constants.RedisPassword,
		DB:       1,
	})
}
