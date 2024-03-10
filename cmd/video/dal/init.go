package dal

import (
	"my-tiktok/cmd/video/dal/db"
	"my-tiktok/cmd/video/mw/redis"
)

func Init() {
	db.Init()
	redis.Init()
}
