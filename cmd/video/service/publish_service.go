package service

import (
	"context"
	"my-tiktok/cmd/video/dal/db"
	"my-tiktok/pkg/kitex_gen/base"
	"my-tiktok/pkg/kitex_gen/video"
	"time"
)

type PublishService struct {
	ctx context.Context
}

// NewPublishService create feed service
func NewPublishService(ctx context.Context) *PublishService {
	return &PublishService{ctx: ctx}
}

func (s *PublishService) PublishAction(req *video.PublishActionRequest) (err error) {
	user_id := req.CurrentUserId
	title := req.Title
	nowTime := req.PublishTime
	_, err = db.CreateVideo(&db.Video{
		AuthorID:    user_id,
		PlayURL:     req.PlayUrl,
		CoverURL:    req.CoverUrl,
		PublishTime: time.Unix(nowTime, 0),
		Title:       title,
	})
	return err
}

func (s *PublishService) PublishList(req *video.PublishListRequest) (resp *video.PublishListResponse, err error) {
	resp = &video.PublishListResponse{}
	query_user_id := req.UserId
	current_user_id := req.CurrentUserId
	dbVideos, err := db.GetVideoByUserID(query_user_id)
	if err != nil {
		return resp, err
	}
	var videos []*base.Video

	f := NewFeedService(s.ctx)
	err = f.CopyVideos(&videos, &dbVideos, current_user_id)
	if err != nil {
		return resp, err
	}
	for _, item := range videos {
		video := base.Video{
			Id: item.Id,
			Author: &base.User{
				Id:              item.Author.Id,
				Name:            item.Author.Name,
				FollowCount:     item.Author.FollowCount,
				FollowerCount:   item.Author.FollowerCount,
				Avatar:          item.Author.Avatar,
				BackgroundImage: item.Author.BackgroundImage,
				Signature:       item.Author.Signature,
				TotalFavorited:  item.Author.TotalFavorited,
				WorkCount:       item.Author.WorkCount,
			},
			PlayUrl:       item.PlayUrl,
			CoverUrl:      item.CoverUrl,
			FavoriteCount: item.FavoriteCount,
			CommentCount:  item.CommentCount,
			IsFavorite:    item.IsFavorite,
			Title:         item.Title,
		}
		resp.VideoList = append(resp.VideoList, &video)
	}
	return resp, nil
}
