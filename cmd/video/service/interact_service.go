package service

import (
	"context"
	"my-tiktok/cmd/api/rpc"
	"my-tiktok/cmd/video/dal/db"
	"my-tiktok/pkg/constants"
	"my-tiktok/pkg/errno"
	"my-tiktok/pkg/kitex_gen/base"
	"my-tiktok/pkg/kitex_gen/user"
	"my-tiktok/pkg/kitex_gen/video"
)

type FavoriteService struct {
	ctx context.Context
}

// NewFavoriteService create favorite service
func NewFavoriteService(ctx context.Context) *FavoriteService {
	return &FavoriteService{ctx: ctx}
}

// FavoriteAction like a video and return result
func (r *FavoriteService) FavoriteAction(req *video.FavoriteActionRequest) (flag bool, err error) {
	_, err = db.CheckVideoExistById(req.VideoId)
	if err != nil {
		return false, err
	}
	if req.ActionType != constants.FavoriteActionType && req.ActionType != constants.UnFavoriteActionType {
		return false, errno.ParamErr
	}
	current_user_id := req.CurrentUserId

	new_favorite_relation := &db.Favorites{
		UserId:  current_user_id,
		VideoId: req.VideoId,
	}
	favorite_exist, _ := db.QueryFavoriteExist(new_favorite_relation.UserId, new_favorite_relation.VideoId)
	if req.ActionType == constants.FavoriteActionType {
		if favorite_exist {
			return false, errno.FavoriteRelationAlreadyExistErr
		}
		flag, err = db.AddNewFavorite(new_favorite_relation)
	} else {
		if !favorite_exist {
			return false, errno.FavoriteRelationNotExistErr
		}
		flag, err = db.DeleteFavorite(new_favorite_relation)
	}
	return flag, err
}

// GetFavoriteList query favorite list by the user id in the request
func (r *FavoriteService) GetFavoriteList(req *video.FavoriteListRequest) (favoritelist []*base.Video, err error) {
	query_user_id := req.UserId
	resp, err := rpc.UserClient.CheckUserExitsById(context.Background(), &user.CheckUserExitsByIdRequset{UserId: query_user_id})
	if resp.BaseResp.StatusCode != 0 {
		err = errno.NewErrNo(resp.BaseResp.StatusCode, resp.BaseResp.StatusMsg)
	}
	if err != nil {
		return nil, err
	}
	current_user_id := req.CurrentUserId

	video_id_list, err := db.GetFavoriteIdList(query_user_id)
	if err != nil {
		return nil, err
	}
	dbVideos, err := db.GetVideoListByVideoIDList(video_id_list)
	if err != nil {
		return nil, err
	}
	var videos []*base.Video
	f := NewFeedService(r.ctx)
	err = f.CopyVideos(&videos, &dbVideos, current_user_id)
	for _, item := range videos {
		video := &base.Video{
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
		favoritelist = append(favoritelist, video)
	}
	return favoritelist, err
}
func (s *FavoriteService) QueryTotalFavoritedByAuthorIDs(userId int64) (int64, error) {
	return db.GetFavoriteCountByUserID(userId)
}

type CommentService struct {
	ctx context.Context
}

// NewCommentService create comment service
func NewCommentService(ctx context.Context) *CommentService {
	return &CommentService{
		ctx: ctx,
	}
}

// CommentAction add a comment to a video and return result
func (r *CommentService) CommentAction(req *video.CommentActionRequest) (*base.Comment, error) {
	current_user_id := req.CurrentUserId
	video_id := req.VideoId
	action_type := req.ActionType
	comment_text := req.CommentText
	comment_id := req.CommentId
	comment := &base.Comment{}
	// 发表评论
	if action_type == 1 {
		db_comment := &db.Comment{
			UserId:      current_user_id,
			VideoId:     video_id,
			CommentText: comment_text,
		}
		err := db.AddNewComment(db_comment)
		if err != nil {
			return comment, err
		}
		comment.Id = db_comment.ID
		comment.CreateDate = db_comment.CreatedAt.Format("01-02")
		comment.Content = db_comment.CommentText
		resp, err := rpc.UserClient.GetUserInfo(context.Background(), &user.GetUserInfoRequest{UserId: current_user_id, CurrentUserId: current_user_id})
		if err != nil {
			return comment, err
		}
		if resp.BaseResp.StatusCode != 0 {
			return comment, errno.NewErrNo(resp.BaseResp.StatusCode, resp.BaseResp.StatusMsg)
		}
		comment.User = resp.User
		if err != nil {
			return comment, err
		}
		return comment, nil
	} else {
		err := db.DeleteCommentById(comment_id)
		if err != nil {
			return comment, err
		}
		return comment, nil
	}
}

// GetCommentList query comment list by the video id in the request
func (r *CommentService) GetCommentList(req *video.CommentListRequest) (commentlist []*base.Comment, err error) {
	video_id := req.VideoId
	current_user_id := req.CurrentUserId
	dbcomments, err := db.GetCommentListByVideoID(video_id)
	if err != nil {
		return nil, err
	}
	err = r.copyComment(&commentlist, &dbcomments, current_user_id)
	if err != nil {
		return nil, err
	}
	return commentlist, nil
}

func (c *CommentService) copyComment(result *[]*base.Comment, data *[]*db.Comment, current_user_id int64) error {
	for _, item := range *data {
		comment, err := c.createComment(item, current_user_id)
		if err != nil {
			return err
		}
		*result = append(*result, comment)
	}
	return nil
}

// createComment convert single comment from db to model
func (c *CommentService) createComment(data *db.Comment, userId int64) (*base.Comment, error) {
	comment := &base.Comment{
		Id:         data.ID,
		Content:    data.CommentText,
		CreateDate: data.CreatedAt.Format("01-02"),
	}

	resp, err := rpc.UserClient.GetUserInfo(context.Background(), &user.GetUserInfoRequest{UserId: userId, CurrentUserId: userId})
	if err != nil {
		return nil, err
	}
	if resp.BaseResp.StatusCode != 0 {
		return nil, errno.NewErrNo(resp.BaseResp.StatusCode, resp.BaseResp.StatusMsg)
	}
	comment.User = resp.User
	return comment, nil
}
