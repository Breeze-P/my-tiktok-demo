// Code generated by hertz generator. DO NOT EDIT.

package register

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	router "my-tiktok/biz/router"
)

// register registers all routers.
func Register(r *server.Hertz) {

	router.GeneratedRegister(r)

	customizedRegister(r)
}
