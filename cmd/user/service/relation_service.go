package service

import (
	"context"
	"log"
	"my-tiktok/cmd/user/dal/db"
	"my-tiktok/pkg/errno"
	"my-tiktok/pkg/kitex_gen/base"
	"my-tiktok/pkg/kitex_gen/user"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

const (
	FOLLOW   = 1
	UNFOLLOW = 2
)

type RelationService struct {
	ctx context.Context
}

func NewRelationService(ctx context.Context) *RelationService {
	return &RelationService{
		ctx: ctx,
	}
}

func (r *RelationService) RelationAction(req *user.RelationActionRequest) (flag bool, err error) {
	_, err = db.CheckUserExistById(req.ToUserId)
	if err != nil {
		return false, err
	}
	if req.ActionType != FOLLOW && req.ActionType != UNFOLLOW {
		return false, errno.ParamErr
	}
	// current_user_id, _ := r.c.Get("current_user_id")
	// Not allowed to pay attention to oneself
	if req.ToUserId == req.CurrentUserId {
		return false, errno.ParamErr
	}
	new_follow_relation := &db.Follows{
		UserId:     req.ToUserId,
		FollowerId: req.CurrentUserId,
	}
	follow_exist, _ := db.QueryFollowExist(new_follow_relation.UserId, new_follow_relation.FollowerId)
	if req.ActionType == FOLLOW {
		if follow_exist {
			return false, errno.FollowRelationAlreadyExistErr
		}
		flag, err = db.AddNewFollow(new_follow_relation)
	} else {
		if !follow_exist {
			return false, errno.FollowRelationNotExistErr
		}
		flag, err = db.DeleteFollow(new_follow_relation)
	}
	return flag, err
}

func (r *RelationService) GetFollowList(req *user.RelationFollowListRequest) (followList []*base.User, err error) {
	_, err = db.CheckUserExistById(req.UserId)
	if err != nil {
		return nil, err
	}

	// var followList []*base.User
	followIdList, err := db.GetFollowIdList(req.UserId)
	if err != nil {
		return followList, err
	}

	for _, follow := range followIdList {
		user_info, err := NewAuthService(r.ctx).GetUserInfo(follow, req.CurrentUserId)
		if err != nil {
			continue
		}
		user := base.User{
			Id:              user_info.Id,
			Name:            user_info.Name,
			FollowCount:     user_info.FollowCount,
			FollowerCount:   user_info.FollowerCount,
			IsFollow:        user_info.IsFollow,
			Avatar:          user_info.Avatar,
			BackgroundImage: user_info.BackgroundImage,
			Signature:       user_info.Signature,
			TotalFavorited:  user_info.TotalFavorited,
			WorkCount:       user_info.WorkCount,
			FavoriteCount:   user_info.FavoriteCount,
		}
		followList = append(followList, &user)
	}
	return followList, nil
}

// GetFollowerList get follower list by the user id in req
func (r *RelationService) GetFollowerList(req *user.RelationFollowerListRequest) ([]*base.User, error) {
	user_id := req.UserId
	var followerList []*base.User

	followerIdList, err := db.GetFollowerIdList(user_id)
	if err != nil {
		return followerList, err
	}

	for _, follower := range followerIdList {
		user_info, err := NewAuthService(r.ctx).GetUserInfo(follower, req.CurrentUserId)
		if err != nil {
			hlog.Error("func error: GetFollowerList -> GetUserInfo")
		}

		user := &base.User{
			Id:              user_info.Id,
			Name:            user_info.Name,
			FollowCount:     user_info.FollowCount,
			FollowerCount:   user_info.FollowerCount,
			IsFollow:        user_info.IsFollow,
			Avatar:          user_info.Avatar,
			BackgroundImage: user_info.BackgroundImage,
			Signature:       user_info.Signature,
			TotalFavorited:  user_info.TotalFavorited,
			WorkCount:       user_info.WorkCount,
			FavoriteCount:   user_info.FavoriteCount,
		}
		followerList = append(followerList, user)
	}
	return followerList, nil
}

// GetFriendList get friend list by the user id in req
func (r *RelationService) GetFriendList(req *user.RelationFriendListRequest) ([]*base.FriendUser, error) {
	user_id := req.UserId
	current_user_id := req.CurrentUserId

	if current_user_id != user_id {
		return nil, errno.FriendListNoPermissionErr
	}

	var friendList []*base.FriendUser

	friendIdList, err := db.GetFriendIdList(user_id)
	if err != nil {
		return friendList, err
	}

	for _, id := range friendIdList {
		user_info, err := NewAuthService(r.ctx).GetUserInfo(id, user_id)
		if err != nil {
			log.Printf("func error: GetFriendList -> GetUserInfo")
		}
		message, err := db.GetLatestMessageByIdPair(user_id, id)
		if err != nil {
			log.Printf("func error: GetFriendList -> GetLatestMessageByIdPair")
		}
		var msgType int64
		if message == nil { // No chat history
			msgType = 2
			message = &db.Messages{}
		} else if user_id == message.FromUserId {
			msgType = 1
		} else {
			msgType = 0
		}
		friendList = append(friendList, &base.FriendUser{
			User: &base.User{
				Id:              user_info.Id,
				Name:            user_info.Name,
				FollowCount:     user_info.FollowCount,
				FollowerCount:   user_info.FollowerCount,
				IsFollow:        user_info.IsFollow,
				Avatar:          user_info.Avatar,
				BackgroundImage: user_info.BackgroundImage,
				Signature:       user_info.Signature,
				TotalFavorited:  user_info.TotalFavorited,
				WorkCount:       user_info.WorkCount,
				FavoriteCount:   user_info.FavoriteCount,
			},
			Message: message.Content,
			MsgType: msgType,
		})
	}

	return friendList, nil
}
