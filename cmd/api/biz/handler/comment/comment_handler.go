// Code generated by hertz generator.

package comment

import (
	"context"

	comment "my-tiktok/cmd/api/biz/model/interact/comment"
	"my-tiktok/cmd/api/rpc"
	"my-tiktok/pkg/errno"
	"my-tiktok/pkg/kitex_gen/video"
	"my-tiktok/pkg/utils"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// CommentAction .
// @router /douyin/comment/action/ [POST]
func CommentAction(ctx context.Context, c *app.RequestContext) {
	var err error
	var req comment.DouyinCommentActionRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		resp := utils.BuildBaseResp(err)
		c.JSON(consts.StatusOK, comment.DouyinCommentActionResponse{
			StatusCode: resp.StatusCode,
			StatusMsg:  resp.StatusMsg,
		})
		return
	}

	current_user_id, exists := c.Get("current_user_id")
	if !exists {
		current_user_id = int64(0)
	}
	// comment_, err := comment_service.NewCommentService(ctx, c).AddNewComment(&req)
	resp, err := rpc.VideoClient.CommentAction(context.Background(), &video.CommentActionRequest{CurrentUserId: current_user_id.(int64), VideoId: req.VideoId, ActionType: req.ActionType, CommentText: req.CommentText, CommentId: req.CommentId})
	if err != nil {
		resp := utils.BuildBaseResp(err)
		c.JSON(consts.StatusOK, comment.DouyinCommentActionResponse{
			StatusCode: resp.StatusCode,
			StatusMsg:  resp.StatusMsg,
		})
		return
	}
	if resp.BaseResp.StatusCode != 0 {
		resp := utils.BuildBaseResp(errno.NewErrNo(resp.BaseResp.StatusCode, resp.BaseResp.StatusMsg))
		c.JSON(consts.StatusOK, comment.DouyinCommentActionResponse{
			StatusCode: resp.StatusCode,
			StatusMsg:  resp.StatusMsg,
		})
		return
	}

	c.JSON(consts.StatusOK, comment.DouyinCommentActionResponse{
		StatusCode: errno.SuccessCode,
		StatusMsg:  errno.SuccessMsg,
		Comment:    utils.ReverseComment(resp.Comment),
	})
}

// CommentList .
// @router /douyin/comment/list/ [GET]
func CommentList(ctx context.Context, c *app.RequestContext) {
	var err error
	var req comment.DouyinCommentListRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp, err := rpc.VideoClient.GetCommentList(context.Background(), &video.CommentListRequest{VideoId: req.VideoId})
	if err != nil {
		resp := utils.BuildBaseResp(err)
		c.JSON(consts.StatusOK, comment.DouyinCommentActionResponse{
			StatusCode: resp.StatusCode,
			StatusMsg:  resp.StatusMsg,
		})
		return
	}
	if resp.BaseResp.StatusCode != 0 {
		c.JSON(consts.StatusOK, comment.DouyinCommentActionResponse{
			StatusCode: resp.BaseResp.StatusCode,
			StatusMsg:  resp.BaseResp.StatusMsg,
		})
		return
	}
	c.JSON(consts.StatusOK, resp)
}
