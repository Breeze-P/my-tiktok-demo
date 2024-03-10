package jwt

import (
	"context"
	"my-tiktok/cmd/api/biz/model/basic/user"
	"my-tiktok/cmd/api/rpc"
	"my-tiktok/pkg/errno"
	"my-tiktok/pkg/utils"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/jwt"

	userservice "my-tiktok/pkg/kitex_gen/user"
)

var (
	JwtMiddleware *jwt.HertzJWTMiddleware // 单例
	identity      = "user_id"
)

func Init() {
	JwtMiddleware, _ = jwt.New(&jwt.HertzJWTMiddleware{
		Key:         []byte("tiktok demo secret"),
		TokenLookup: "query:token,form:token",
		Timeout:     24 * time.Hour,
		IdentityKey: identity,
		Authenticator: func(ctx context.Context, c *app.RequestContext) (interface{}, error) {
			var loginRequest user.DouyinUserLoginRequest
			if err := c.BindAndValidate(&loginRequest); err != nil {
				return nil, err
			}
			resp, err := rpc.UserClient.Login(context.Background(), &userservice.LoginRequest{Username: loginRequest.Username, Password: loginRequest.Password})
			if err != nil {
				return nil, err
			}
			if resp.BaseResp.StatusCode != 0 {
				return nil, errno.NewErrNo(resp.BaseResp.StatusCode, resp.BaseResp.StatusMsg)
			}
			// user, err := db.QueryUser(loginRequest.Username)
			// if err != nil {
			// 	return nil, err
			// }
			// if !utils.VerifyPassword(loginRequest.Password, user.Password) {
			// 	return nil, errno.PasswordIsNotVerified
			// }
			c.Set("user_id", resp.UserId)
			return resp.UserId, nil
		},
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(int64); ok {
				return jwt.MapClaims{
					identity: v,
				}
			}
			return jwt.MapClaims{}
		},
		LoginResponse: func(ctx context.Context, c *app.RequestContext, code int, token string, expire time.Time) {
			hlog.CtxInfof(ctx, "Login success ，token is issued clientIP: "+c.ClientIP())
			c.Set("token", token)
		},
		// Verify token and get the id of logged-in user
		Authorizator: func(data interface{}, ctx context.Context, c *app.RequestContext) bool {
			if v, ok := data.(float64); ok {
				current_user_id := int64(v)
				c.Set("current_user_id", current_user_id)
				hlog.CtxInfof(ctx, "Token is verified clientIP: "+c.ClientIP())
				return true
			}
			return false
		},
		// Validation failed, build the message
		Unauthorized: func(ctx context.Context, c *app.RequestContext, code int, message string) {
			c.JSON(consts.StatusOK, user.DouyinUserLoginResponse{
				StatusCode: errno.AuthorizationFailedErrCode,
				StatusMsg:  message,
			})
		},
		HTTPStatusMessageFunc: func(e error, ctx context.Context, c *app.RequestContext) string {
			resp := utils.BuildBaseResp(e)
			return resp.StatusMsg
		},
	})
}
