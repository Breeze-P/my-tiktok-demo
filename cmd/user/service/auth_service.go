package service

import (
	"context"
	"my-tiktok/cmd/user/dal/db"
	"my-tiktok/cmd/user/rpc"
	"my-tiktok/pkg/constants"
	"my-tiktok/pkg/errno"
	"my-tiktok/pkg/kitex_gen/base"
	"my-tiktok/pkg/kitex_gen/user"
	"my-tiktok/pkg/kitex_gen/video"
	"my-tiktok/pkg/utils"
	"sync"
)

type AuthService struct {
	ctx context.Context
}

func NewAuthService(ctx context.Context) *AuthService {
	return &AuthService{
		ctx: ctx,
	}
}

// Login implements the UserServiceImpl interface.
func (s *AuthService) Login(req *user.LoginRequest) (uid int64, err error) {
	user, err := db.QueryUser(req.Username)
	if err != nil {
		return 0, err
	}
	if !utils.VerifyPassword(req.Password, user.Password) {
		return 0, errno.PasswordIsNotVerified
	}
	return user.ID, nil
}

func (s *AuthService) Register(req *user.RegisterRequest) (user_id int64, err error) {
	// 已存在判断
	user, err := db.QueryUser(req.Username)
	if err != nil {
		return 0, err
	}
	if *user != (db.User{}) {
		return 0, errno.UserAlreadyExistErr
	}

	passWord, err := utils.Crypt(req.Password)

	if err != nil {
		return 0, err
	}
	user_id, err = db.CreateUser(&db.User{
		UserName:        req.Username,
		Password:        passWord,
		Avatar:          constants.TestAva,
		BackgroundImage: constants.TestBackground,
		Signature:       constants.TestSign,
	})
	if err != nil {
		return 0, err
	}
	return user_id, nil
}

func (s *AuthService) UserInfo(req *user.GetUserInfoRequest) (user *base.User, err error) {
	query_user_id := req.UserId
	// v, exists := s.ctx.Get("current_user_id")
	// var current_user_id int64
	// if !exists {
	// 	current_user_id = 0 // means 游客访问
	// }
	current_user_id := req.CurrentUserId
	return s.GetUserInfo(query_user_id, current_user_id)
}

func (s *AuthService) GetUserInfo(query_user_id, user_id int64) (*base.User, error) {
	u := &base.User{}
	// 这里的并发调用一定背诵默写
	errChan := make(chan error, 7)
	var wg sync.WaitGroup
	wg.Add(7)
	go func() {
		dbUser, err := db.QueryUserById(query_user_id)
		if err != nil {
			errChan <- err
		} else {
			u.Name = dbUser.UserName
			u.Avatar = dbUser.Avatar
			u.BackgroundImage = dbUser.BackgroundImage
			u.Signature = dbUser.Signature
		}
		wg.Done()
	}()

	go func() {
		// 两种方案：1、直接访问数据库; 2、请求rpc服务;
		resp, err := rpc.VideoClient.GetWorkCount(context.Background(), &video.GetWorkCountRequest{CurrentUserId: query_user_id})
		if err != nil {
			errChan <- err
		} else if resp.BaseResp.StatusCode != 0 {
			errChan <- errno.NewErrNo(resp.BaseResp.StatusCode, resp.BaseResp.StatusMsg)
		} else {
			u.WorkCount = resp.Count
		}
		wg.Done()
	}()

	go func() {
		FollowCount, err := db.GetFollowCount(query_user_id)
		if err != nil {
			errChan <- err
			return
		} else {
			u.FollowCount = FollowCount
		}
		wg.Done()
	}()

	go func() {
		FollowerCount, err := db.GetFollowerCount(query_user_id)
		if err != nil {
			errChan <- err
		} else {
			u.FollowerCount = FollowerCount
		}
		wg.Done()
	}()

	go func() {
		if user_id != 0 {
			IsFollow, err := db.QueryFollowExist(user_id, query_user_id)
			if err != nil {
				errChan <- err
			} else {
				u.IsFollow = IsFollow
			}
		} else {
			u.IsFollow = false
		}
		wg.Done()
	}()

	go func() {
		resp, err := rpc.VideoClient.GetFavoriteCountByUserID(context.Background(), &video.GetFavoriteCountByUserIDRequest{CurrentUserId: query_user_id})
		if err != nil {
			errChan <- err
		} else if resp.BaseResp.StatusCode != 0 {
			errChan <- errno.NewErrNo(resp.BaseResp.StatusCode, resp.BaseResp.StatusMsg)
		} else {
			u.FavoriteCount = resp.Count
		}
		wg.Done()
	}()

	go func() {
		resp, err := rpc.VideoClient.QueryTotalFavoritedByAuthorID(context.Background(), &video.QueryTotalFavoritedByAuthorIDRequest{CurrentUserId: query_user_id})
		if err != nil {
			errChan <- err
		} else if resp.BaseResp.StatusCode != 0 {
			errChan <- errno.NewErrNo(resp.BaseResp.StatusCode, resp.BaseResp.StatusMsg)
		} else {
			u.TotalFavorited = resp.Count
		}
		wg.Done()
	}()

	wg.Wait()
	select {
	case result := <-errChan:
		return &base.User{}, result
	default:
	}
	u.Id = query_user_id
	return u, nil
}

func (s *AuthService) CheckUserExistById(userId int64) (exits bool, err error) {
	return db.CheckUserExistById(userId)
}
