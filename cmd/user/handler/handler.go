package handler

import (
	"context"
	"my-tiktok/cmd/user/service"
	"my-tiktok/pkg/errno"
	user "my-tiktok/pkg/kitex_gen/user"
	"my-tiktok/pkg/utils"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// Login implements the UserServiceImpl interface.
func (s *UserServiceImpl) Login(ctx context.Context, req *user.LoginRequest) (resp *user.LoginResponse, err error) {
	// TODO: Your code here...
	resp = new(user.LoginResponse)
	if len(req.Username) == 0 || len(req.Password) == 0 {
		resp.BaseResp = utils.BuildBaseResp(errno.ParamErr)
		return resp, nil
	}
	uid, err := service.NewAuthService(ctx).Login(req)
	if err != nil {
		resp.BaseResp = utils.BuildBaseResp(err)
		return resp, nil
	}
	resp.UserId = uid
	resp.BaseResp = utils.BuildBaseResp(errno.Success)
	return resp, nil
}

// AdminLogin implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, req *user.RegisterRequest) (resp *user.RegisterResponse, err error) {
	// TODO: Your code here...
	resp = new(user.RegisterResponse)
	if len(req.Username) == 0 || len(req.Password) == 0 {
		resp.BaseResp = utils.BuildBaseResp(errno.ParamErr)
		return resp, nil
	}
	uid, err := service.NewAuthService(ctx).Register(req)
	if err != nil {
		resp.BaseResp = utils.BuildBaseResp(err)
		return resp, nil
	}
	resp.UserId = uid
	resp.BaseResp = utils.BuildBaseResp(errno.Success)
	return resp, nil
}

// ChangeAdminPassword implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetUserInfo(ctx context.Context, req *user.GetUserInfoRequest) (resp *user.GetUserInfoResponse, err error) {
	// TODO: Your code here...
	resp = new(user.GetUserInfoResponse)
	if req.UserId == 0 || req.CurrentUserId == 0 {
		resp.BaseResp = utils.BuildBaseResp(errno.ParamErr)
		return resp, nil
	}
	user, err := service.NewAuthService(ctx).UserInfo(req)
	if err != nil {
		resp.BaseResp = utils.BuildBaseResp(err)
		return resp, nil
	}
	resp.User = user
	resp.BaseResp = utils.BuildBaseResp(errno.Success)
	return resp, nil
}

// RelationAction implements the UserServiceImpl interface.
func (s *UserServiceImpl) RelationAction(ctx context.Context, req *user.RelationActionRequest) (resp *user.RelationActionResponse, err error) {
	// TODO: Your code here...
	resp = new(user.RelationActionResponse)
	if req.CurrentUserId == 0 || req.ToUserId == 0 || req.ActionType != 1 && req.ActionType != 2 {
		resp.BaseResp = utils.BuildBaseResp(errno.ParamErr)
		return resp, nil
	}
	_, err = service.NewRelationService(ctx).RelationAction(req)
	if err != nil {
		resp.BaseResp = utils.BuildBaseResp(err)
		return resp, nil
	}
	resp.BaseResp = utils.BuildBaseResp(errno.Success)
	return resp, nil
}

// GetFollowList implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetFollowList(ctx context.Context, req *user.RelationFollowListRequest) (resp *user.RelationFollowListResponse, err error) {
	// TODO: Your code here...
	resp = new(user.RelationFollowListResponse)
	if req.CurrentUserId == 0 || req.UserId == 0 {
		resp.BaseResp = utils.BuildBaseResp(errno.ParamErr)
		return resp, nil
	}
	followList, err := service.NewRelationService(ctx).GetFollowList(req)
	if err != nil {
		resp.BaseResp = utils.BuildBaseResp(err)
		return resp, nil
	}
	resp.BaseResp = utils.BuildBaseResp(errno.Success)
	resp.UserList = followList
	return resp, nil
}

// GetFollowerList implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetFollowerList(ctx context.Context, req *user.RelationFollowerListRequest) (resp *user.RelationFollowerListResponse, err error) {
	// TODO: Your code here...
	resp = new(user.RelationFollowerListResponse)
	if req.CurrentUserId == 0 || req.UserId == 0 {
		resp.BaseResp = utils.BuildBaseResp(errno.ParamErr)
		return resp, nil
	}
	followList, err := service.NewRelationService(ctx).GetFollowerList(req)
	if err != nil {
		resp.BaseResp = utils.BuildBaseResp(err)
		return resp, nil
	}
	resp.BaseResp = utils.BuildBaseResp(errno.Success)
	resp.UserList = followList
	return resp, nil
}

// GetFriendList implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetFriendList(ctx context.Context, req *user.RelationFriendListRequest) (resp *user.RelationFriendListResponse, err error) {
	// TODO: Your code here...
	resp = new(user.RelationFriendListResponse)
	if req.CurrentUserId == 0 || req.UserId == 0 {
		resp.BaseResp = utils.BuildBaseResp(errno.ParamErr)
		return resp, nil
	}
	followList, err := service.NewRelationService(ctx).GetFriendList(req)
	if err != nil {
		resp.BaseResp = utils.BuildBaseResp(err)
		return resp, nil
	}
	resp.BaseResp = utils.BuildBaseResp(errno.Success)
	resp.UserList = followList
	return resp, nil
}

// Chat implements the UserServiceImpl interface.
func (s *UserServiceImpl) Chat(ctx context.Context, req *user.MessageChatRequest) (resp *user.MessageChatResponse, err error) {
	// TODO: Your code here...
	resp = new(user.MessageChatResponse)
	if req.CurrentUserId == 0 || req.ToUserId == 0 {
		resp.BaseResp = utils.BuildBaseResp(errno.ParamErr)
		return resp, nil
	}
	messageList, err := service.NewMessageService(ctx).GetMessageChat(req)
	if err != nil {
		resp.BaseResp = utils.BuildBaseResp(err)
		return resp, nil
	}
	resp.BaseResp = utils.BuildBaseResp(errno.Success)
	resp.MessageList = messageList
	return resp, nil
}

// MessageAction implements the UserServiceImpl interface.
func (s *UserServiceImpl) MessageAction(ctx context.Context, req *user.MessageActionRequest) (resp *user.MessageActionResponse, err error) {
	// TODO: Your code here...
	resp = new(user.MessageActionResponse)
	if req.CurrentUserId == 0 || req.ToUserId == 0 {
		resp.BaseResp = utils.BuildBaseResp(errno.ParamErr)
		return resp, nil
	}
	err = service.NewMessageService(ctx).MessageAction(req)
	if err != nil {
		resp.BaseResp = utils.BuildBaseResp(err)
		return resp, nil
	}
	resp.BaseResp = utils.BuildBaseResp(errno.Success)
	return resp, nil
}

func (s *UserServiceImpl) CheckUserExitsById(ctx context.Context, req *user.CheckUserExitsByIdRequset) (resp *user.CheckUserExitsByIdResponse, err error) {
	resp = new(user.CheckUserExitsByIdResponse)
	if req.UserId == 0 {
		resp.BaseResp = utils.BuildBaseResp(errno.ParamErr)
		return resp, nil
	}
	exits, err := service.NewAuthService(ctx).CheckUserExistById(req.UserId)
	if err != nil {
		resp.BaseResp = utils.BuildBaseResp(err)
		return resp, nil
	}
	resp.BaseResp = utils.BuildBaseResp(errno.Success)
	resp.Exits = exits
	return resp, nil
}
