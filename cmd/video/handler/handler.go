package handler

import (
	"context"
	"my-tiktok/cmd/video/service"
	"my-tiktok/pkg/errno"
	video "my-tiktok/pkg/kitex_gen/video"
	"my-tiktok/pkg/utils"
)

// VideoServiceImpl implements the last service interface defined in the IDL.
type VideoServiceImpl struct{}

// GetFeed implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) GetFeed(ctx context.Context, req *video.FeedRequest) (resp *video.FeedResponse, err error) {
	// TODO: Your code here...
	// req.CurrentUserId == 0 非登陆用户
	// req.LatestTime 可选
	resp, err = service.NewFeedService(ctx).Feed(req)
	return resp, err
}

// PublishAction implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) PublishAction(ctx context.Context, req *video.PublishActionRequest) (resp *video.PublishActionResponse, err error) {
	// TODO: Your code here...
	resp = new(video.PublishActionResponse)
	if req.CurrentUserId == 0 || len(req.PlayUrl) == 0 || len(req.CoverUrl) == 0 || len(req.Title) == 0 {
		// req.CurrentUserId == 0 未登陆用户
		resp.BaseResp = utils.BuildBaseResp(errno.ParamErr)
		return resp, nil
	}
	err = service.NewPublishService(ctx).PublishAction(req)
	if err != nil {
		resp.BaseResp = utils.BuildBaseResp(err)
		return resp, nil
	}
	resp.BaseResp = utils.BuildBaseResp(errno.Success)
	return resp, nil
}

// GetPublishList implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) GetPublishList(ctx context.Context, req *video.PublishListRequest) (resp *video.PublishListResponse, err error) {
	// TODO: Your code here...
	// req.CurrentUserId == 0 非登陆用户
	// req.UserId == 0 查不到数据
	resp, err = service.NewPublishService(ctx).PublishList(req)
	return
}

// GetWorkCount implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) GetWorkCount(ctx context.Context, req *video.GetWorkCountRequest) (resp *video.GetWorkCountResponse, err error) {
	// TODO: Your code here...
	resp = new(video.GetWorkCountResponse)
	count, err := service.NewFeedService(ctx).GetWorkCount(req.CurrentUserId)
	if err != nil {
		resp.BaseResp = utils.BuildBaseResp(err)
		return resp, nil
	}
	resp.BaseResp = utils.BuildBaseResp(errno.Success)
	resp.Count = count
	return resp, nil
}

// FavoriteAction implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) FavoriteAction(ctx context.Context, req *video.FavoriteActionRequest) (resp *video.FavoriteActionResponse, err error) {
	// TODO: Your code here...
	resp = new(video.FavoriteActionResponse)
	if req.ActionType != 2 && req.ActionType != 1 || req.CurrentUserId == 0 || req.VideoId == 0 {
		resp.BaseResp = utils.BuildBaseResp(errno.ParamErr)
		return resp, nil
	}
	flag, err := service.NewFavoriteService(ctx).FavoriteAction(req)
	if err != nil {
		resp.BaseResp = utils.BuildBaseResp(err)
		return resp, nil
	}
	resp.BaseResp = utils.BuildBaseResp(errno.Success)
	resp.Success = flag
	return
}

// GetFavoriteList implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) GetFavoriteList(ctx context.Context, req *video.FavoriteListRequest) (resp *video.FavoriteListResponse, err error) {
	// TODO: Your code here...
	resp = new(video.FavoriteListResponse)
	if req.UserId == 0 {
		resp.BaseResp = utils.BuildBaseResp(errno.ParamErr)
		return resp, nil
	}
	favoritelist, err := service.NewFavoriteService(ctx).GetFavoriteList(req)
	if err != nil {
		resp.BaseResp = utils.BuildBaseResp(err)
		return resp, nil
	}
	resp.BaseResp = utils.BuildBaseResp(errno.Success)
	resp.VideoList = favoritelist
	return
}

// GetFavoriteCountByUserID implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) GetFavoriteCountByUserID(ctx context.Context, req *video.GetFavoriteCountByUserIDRequest) (resp *video.GetFavoriteCountByUserIDResponse, err error) {
	// TODO: Your code here...
	resp = new(video.GetFavoriteCountByUserIDResponse)
	count, err := service.NewFeedService(ctx).GetFavoriteCountByUserID(req.CurrentUserId)
	if err != nil {
		resp.BaseResp = utils.BuildBaseResp(err)
		return resp, nil
	}
	resp.BaseResp = utils.BuildBaseResp(errno.Success)
	resp.Count = count
	return
}

// QueryTotalFavoritedByAuthorID implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) QueryTotalFavoritedByAuthorID(ctx context.Context, req *video.QueryTotalFavoritedByAuthorIDRequest) (resp *video.QueryTotalFavoritedByAuthorIDResponse, err error) {
	// TODO: Your code here...
	resp = new(video.QueryTotalFavoritedByAuthorIDResponse)
	count, err := service.NewFavoriteService(ctx).QueryTotalFavoritedByAuthorIDs(req.CurrentUserId)
	if err != nil {
		resp.BaseResp = utils.BuildBaseResp(err)
		return resp, nil
	}
	resp.BaseResp = utils.BuildBaseResp(errno.Success)
	resp.Count = count
	return
}

// CommentAction implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) CommentAction(ctx context.Context, req *video.CommentActionRequest) (resp *video.CommentActionResponse, err error) {
	// TODO: Your code here...
	if req.ActionType != 1 && req.ActionType != 2 || req.CurrentUserId == 0 || req.VideoId == 0 || len(req.CommentText) == 0 {
		resp.BaseResp = utils.BuildBaseResp(errno.ParamErr)
		return resp, nil
	}
	resp = new(video.CommentActionResponse)
	comment, err := service.NewCommentService(ctx).CommentAction(req)
	if err != nil {
		resp.BaseResp = utils.BuildBaseResp(err)
		return resp, nil
	}
	resp.BaseResp = utils.BuildBaseResp(errno.Success)
	resp.Comment = comment
	return
}

// GetCommentList implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) GetCommentList(ctx context.Context, req *video.CommentListRequest) (resp *video.CommentListResponse, err error) {
	// TODO: Your code here...
	resp = new(video.CommentListResponse)
	if req.VideoId == 0 {
		resp.BaseResp = utils.BuildBaseResp(errno.ParamErr)
		return resp, nil
	}
	commentlist, err := service.NewCommentService(ctx).GetCommentList(req)
	if err != nil {
		resp.BaseResp = utils.BuildBaseResp(err)
		return resp, nil
	}
	resp.BaseResp = utils.BuildBaseResp(errno.Success)
	resp.CommentList = commentlist
	return
}
