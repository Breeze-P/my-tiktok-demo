package dal

import (
	"my-tiktok/cmd/user/dal/db"
	"my-tiktok/cmd/user/mw/redis"
)

// Init init dal
func Init() {
	db.Init()
	redis.Init()
}
