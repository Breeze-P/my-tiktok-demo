package service

import (
	"context"
	"my-tiktok/biz/dal/db"
	"my-tiktok/biz/model/basic/publish"
	"my-tiktok/biz/model/common"
	"my-tiktok/biz/mw/minio"
	feed_service "my-tiktok/biz/service/feed"
	"my-tiktok/pkg/constants"
	"my-tiktok/pkg/utils"
	"path"
	"strconv"
	"time"

	"github.com/cloudwego/hertz-examples/bizdemo/tiktok_demo/biz/mw/ffmpeg"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type PublishService struct {
	ctx context.Context
	c   *app.RequestContext
}

// NewPublishService create publish service
func NewPublishService(ctx context.Context, c *app.RequestContext) *PublishService {
	return &PublishService{ctx: ctx, c: c}
}

func (s *PublishService) PublishAction(req *publish.DouyinPublishActionRequest) (err error) {
	v, _ := s.c.Get("current_user_id")
	title := s.c.PostForm("title")
	user_id := v.(int64)
	nowTime := time.Now()
	filename := utils.NewFileName(user_id, nowTime.Unix())
	req.Data.Filename = filename + path.Ext(req.Data.Filename)
	uploadInfo, err := minio.PutToBucket(s.ctx, constants.MinioVideoBucketName, req.Data)
	hlog.CtxInfof(s.ctx, "video upload size:"+strconv.FormatInt(uploadInfo.Size, 10))
	PlayURL := constants.MinioImgBucketName + "/" + req.Data.Filename
	buf, err := ffmpeg.GetSnapshot(utils.URLconvert(s.ctx, s.c, PlayURL))
	uploadInfo, err = minio.PutToBucketByBuf(s.ctx, constants.MinioImgBucketName, filename+".png", buf)
	hlog.CtxInfof(s.ctx, "image upload size:"+strconv.FormatInt(uploadInfo.Size, 10))
	if err != nil {
		hlog.CtxInfof(s.ctx, "err:"+err.Error())
	}
	_, err = db.CreateVideo(&db.Video{
		AuthorID:    user_id,
		PlayURL:     PlayURL,
		CoverURL:    constants.MinioImgBucketName + "/" + filename + ".png",
		PublishTime: nowTime,
		Title:       title,
	})
	return err
}

func (s *PublishService) PublishList(req *publish.DouyinPublishListRequest) (resp *publish.DouyinPublishListResponse, err error) {
	resp = &publish.DouyinPublishListResponse{}
	query_user_id := req.UserId
	current_user_id, exist := s.c.Get("current_user_id")
	if !exist {
		current_user_id = int64(0)
	}
	dbVideos, err := db.GetVideoByUserID(query_user_id)
	if err != nil {
		return resp, err
	}
	var videos []*common.Video

	f := feed_service.NewFeedService(s.ctx, s.c)
	err = f.CopyVideos(&videos, &dbVideos, current_user_id.(int64))
	if err != nil {
		return resp, err
	}
	for _, item := range videos {
		video := common.Video{
			Id: item.Id,
			Author: &common.User{
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
		resp.VideoList = append(resp.VideoList, video)
	}
	return resp, nil
}
