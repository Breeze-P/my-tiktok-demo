package utils

import (
	"errors"
	"my-tiktok/pkg/errno"
	"my-tiktok/pkg/kitex_gen/base"
)

// type BaseResp struct {
// 	StatusCode int32
// 	StatusMsg  string
// }

// BuildBaseResp convert error and build BaseResp
func BuildBaseResp(err error) *base.BaseResponse {
	if err == nil {
		return baseResp(errno.Success)
	}

	e := errno.ErrNo{}
	if errors.As(err, &e) {
		return baseResp(e)
	}

	s := errno.ServiceErr.WithMessage(err.Error())
	return baseResp(s)
}

// baseResp build BaseResp from error
func baseResp(err errno.ErrNo) *base.BaseResponse {
	return &base.BaseResponse{
		StatusCode: err.ErrCode,
		StatusMsg:  err.ErrMsg,
	}
}
