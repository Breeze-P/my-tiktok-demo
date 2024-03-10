package utils

import (
	"my-tiktok/cmd/api/biz/model/common"
	"my-tiktok/cmd/api/biz/model/interact/comment"
	"my-tiktok/cmd/api/biz/model/social/message"
	"my-tiktok/cmd/api/biz/model/social/relation"
	"my-tiktok/pkg/kitex_gen/base"
)

// 避免重新生成hertz代码使得新写的代码被复写，reverse新的base里的data struct to common
// todo 实际上rpc的返回可以直接由api返回的，尝试重新用hertz生成

func ReverseUser(user *base.User) *common.User {
	return &common.User{
		Id:              user.Id,
		Name:            user.Name,
		FollowCount:     user.FollowCount,
		FollowerCount:   user.FollowerCount,
		IsFollow:        user.IsFollow,
		Avatar:          user.Avatar,
		BackgroundImage: user.BackgroundImage,
		Signature:       user.Signature,
		TotalFavorited:  user.TotalFavorited,
		WorkCount:       user.WorkCount,
		FavoriteCount:   user.FavoriteCount,
	}
}

func ReverseUsers(users []*base.User) []*common.User {
	res := []*common.User{}
	for _, u := range users {
		res = append(res, ReverseUser(u))
	}
	return res
}

func ReverseFriend(uu *base.FriendUser) *relation.FriendUser {
	user := uu.User
	u := common.User{
		Id:              user.Id,
		Name:            user.Name,
		FollowCount:     user.FollowCount,
		FollowerCount:   user.FollowerCount,
		IsFollow:        user.IsFollow,
		Avatar:          user.Avatar,
		BackgroundImage: user.BackgroundImage,
		Signature:       user.Signature,
		TotalFavorited:  user.TotalFavorited,
		WorkCount:       user.WorkCount,
		FavoriteCount:   user.FavoriteCount,
	}
	return &relation.FriendUser{User: u}
}

func ReverseFriends(users []*base.FriendUser) []*relation.FriendUser {
	res := []*relation.FriendUser{}
	for _, u := range users {
		res = append(res, ReverseFriend(u))
	}
	return res
}

func ReverseComment(c *base.Comment) *comment.Comment {
	return &comment.Comment{
		Id:      c.Id,
		Content: c.Content,
		User:    ReverseUser(c.User),
	}
}

func ReverseVideo(v *base.Video) *common.Video {
	return &common.Video{
		Id:            v.Id,
		Title:         v.Title,
		Author:        ReverseUser(v.Author),
		PlayUrl:       v.PlayUrl,
		CoverUrl:      v.CoverUrl,
		FavoriteCount: v.FavoriteCount,
		CommentCount:  v.CommentCount,
		IsFavorite:    v.IsFavorite,
	}
}

func ReverseVideos(videos []*base.Video) []common.Video {
	res := []common.Video{}
	for _, v := range videos {
		res = append(res, *ReverseVideo(v))
	}
	return res
}

func ReverseVideosPtr(videos []*base.Video) []*common.Video {
	res := []*common.Video{}
	for _, v := range videos {
		res = append(res, ReverseVideo(v))
	}
	return res
}

func ReverseMessage(m *base.Message) *message.Message {
	return &message.Message{
		Id:      m.Id,
		Content: m.Content,
	}
}

func ReverseMessages(messages []*base.Message) []*message.Message {
	res := []*message.Message{}
	for _, m := range messages {
		res = append(res, ReverseMessage(m))
	}
	return res
}
