package service

import (
	"context"
	"fmt"
	"log"
	"my-tiktok/cmd/video/dal/db"
	"my-tiktok/cmd/video/rpc"
	"my-tiktok/pkg/constants"
	"my-tiktok/pkg/errno"
	"my-tiktok/pkg/kitex_gen/base"
	"my-tiktok/pkg/kitex_gen/user"
	"my-tiktok/pkg/kitex_gen/video"
	"sync"
	"time"
)

type FeedService struct {
	ctx context.Context
}

// NewFeedService create feed service
func NewFeedService(ctx context.Context) *FeedService {
	return &FeedService{ctx: ctx}
}

// Feed get the last ten videos until the deadline
func (s *FeedService) Feed(req *video.FeedRequest) (*video.FeedResponse, error) {
	resp := &video.FeedResponse{}
	var lastTime time.Time
	if req.LatestTime == 0 {
		lastTime = time.Now()
	} else {
		lastTime = time.Unix(req.LatestTime/1000, 0)
	}
	fmt.Printf("LastTime: %v\n", lastTime)
	current_user_id := req.CurrentUserId
	if current_user_id != 0 {
		checkUserResp, err := rpc.UserClient.CheckUserExitsById(context.Background(), &user.CheckUserExitsByIdRequset{UserId: current_user_id})
		if err != nil {
			return resp, err
		}
		if checkUserResp.BaseResp.StatusCode != 0 {
			return resp, errno.NewErrNo(checkUserResp.BaseResp.StatusCode, checkUserResp.BaseResp.StatusMsg)
		}
		if !checkUserResp.Exits {
			return resp, errno.UserIsNotExistErr
		}
	}
	dbVideos, err := db.GetVideosByLastTime(lastTime)
	if err != nil {
		return resp, err
	}

	videos := make([]*base.Video, 0, constants.VideoFeedCount)
	err = s.CopyVideos(&videos, &dbVideos, current_user_id)
	if err != nil {
		return resp, nil
	}
	resp.VideoList = videos
	if len(dbVideos) != 0 {
		resp.NextTime = dbVideos[len(dbVideos)-1].PublishTime.Unix()
	}
	return resp, nil
}

// CopyVideos use db.Video make feed.Video
func (s *FeedService) CopyVideos(result *[]*base.Video, data *[]*db.Video, userId int64) error {
	for _, item := range *data {
		video := s.createVideo(item, userId)
		*result = append(*result, video)
	}
	return nil
}

// createVideo get video info by concurrent query
func (s *FeedService) createVideo(data *db.Video, userId int64) *base.Video {
	video := &base.Video{
		Id: data.ID,
		// convert path in the db into a complete url accessible by the front end
		// todo 在http端转换一下
		// PlayUrl:  utils.URLconvert(s.ctx, s.c, data.PlayURL),
		// CoverUrl: utils.URLconvert(s.ctx, s.c, data.CoverURL),
		PlayUrl:  data.PlayURL,
		CoverUrl: data.CoverURL,
		Title:    data.Title,
	}

	var wg sync.WaitGroup
	wg.Add(4)

	// Get author information
	go func() {
		// author, err := user_service.NewUserService(s.ctx, s.c).GetUserInfo(data.AuthorID, userId)
		resp, err := rpc.UserClient.GetUserInfo(context.Background(), &user.GetUserInfoRequest{UserId: data.AuthorID, CurrentUserId: userId})
		if err != nil {
			log.Printf("GetUserInfo rpc client error:" + err.Error())
		}
		if resp.BaseResp.StatusCode != 0 {
			log.Printf("GetUserInfo rpc server error:" + err.Error())
		}
		author := resp.User
		video.Author = &base.User{
			Id:              author.Id,
			Name:            author.Name,
			FollowCount:     author.FollowCount,
			FollowerCount:   author.FollowerCount,
			IsFollow:        author.IsFollow,
			Avatar:          author.Avatar,
			BackgroundImage: author.BackgroundImage,
			Signature:       author.Signature,
			TotalFavorited:  author.TotalFavorited,
			WorkCount:       author.WorkCount,
			FavoriteCount:   author.FavoriteCount,
		}

		wg.Done()
	}()

	// Get the number of video likes
	go func() {
		err := *new(error)
		video.FavoriteCount, err = db.GetFavoriteCount(data.ID)
		if err != nil {
			log.Printf("GetFavoriteCount func error:" + err.Error())
		}
		wg.Done()
	}()

	// Get comment count
	go func() {
		err := *new(error)
		video.CommentCount, err = db.GetCommentCountByVideoID(data.ID)
		if err != nil {
			log.Printf("GetCommentCountByVideoID func error:" + err.Error())
		}
		wg.Done()
	}()

	// Get favorite exist
	go func() {
		err := *new(error)
		video.IsFavorite, err = db.QueryFavoriteExist(userId, data.ID)
		if err != nil {
			log.Printf("QueryFavoriteExist func error:" + err.Error())
		}
		wg.Done()
	}()

	wg.Wait()
	return video
}

func (s *FeedService) GetWorkCount(userId int64) (int64, error) {
	return db.GetWorkCount(userId)
}

func (s *FeedService) GetFavoriteCountByUserID(userId int64) (int64, error) {
	return db.GetFavoriteCountByUserID(userId)
}
