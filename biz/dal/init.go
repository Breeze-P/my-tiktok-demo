package dal

import (
	"my-tiktok/biz/dal/db"
	"my-tiktok/biz/mw/redis"
)

// Init init dal
func Init() {
	db.Init()
	redis.Init()
}
