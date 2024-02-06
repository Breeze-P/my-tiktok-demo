package service

import (
	"context"
	"my-tiktok/biz/dal/db"
	"my-tiktok/biz/model/basic/user"
	"my-tiktok/biz/model/common"
	"my-tiktok/pkg/constants"
	"my-tiktok/pkg/errno"
	"my-tiktok/pkg/utils"
	"sync"

	"github.com/cloudwego/hertz/pkg/app"
)

type UserService struct {
	ctx context.Context
	c   *app.RequestContext
}

// 使用面向对象，我觉得是因为不同重复传cc，但是每次调用需要初始化一个service，也挺烦的...
// NewUserService create user service
func NewUserService(ctx context.Context, c *app.RequestContext) *UserService {
	return &UserService{ctx: ctx, c: c}
}

func (s *UserService) UserRegister(req *user.DouyinUserRegisterRequest) (user_id int64, err error) {
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

func (s *UserService) UserInfo(req *user.DouyinUserRequest) (*common.User, error) {
	query_user_id := req.UserId
	v, exists := s.c.Get("current_user_id")
	var current_user_id int64
	if !exists {
		current_user_id = 0 // means 游客访问
	}
	current_user_id = v.(int64)
	return s.GetUserInfo(query_user_id, current_user_id)
}

func (s *UserService) GetUserInfo(query_user_id, user_id int64) (*common.User, error) {
	u := &common.User{}
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
		WorkCount, err := db.GetWorkCount(query_user_id)
		if err != nil {
			errChan <- err
		} else {
			u.WorkCount = WorkCount
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
		FavoriteCount, err := db.GetFavoriteCountByUserID(query_user_id)
		if err != nil {
			errChan <- err
		} else {
			u.FavoriteCount = FavoriteCount
		}
		wg.Done()
	}()

	go func() {
		TotalFavorited, err := db.QueryTotalFavoritedByAuthorID(query_user_id)
		if err != nil {
			errChan <- err
		} else {
			u.TotalFavorited = TotalFavorited
		}
		wg.Done()
	}()

	wg.Wait()
	select {
	case result := <-errChan:
		return &common.User{}, result
	default:
	}
	u.Id = query_user_id
	return u, nil
}
