package service

import (
	"context"
	"my-tiktok/cmd/user/dal/db"
	"my-tiktok/pkg/errno"
	"my-tiktok/pkg/kitex_gen/base"
	"my-tiktok/pkg/kitex_gen/user"
	"my-tiktok/pkg/utils"
)

type MessageService struct {
	ctx context.Context
}

// NewMessageService create message service
func NewMessageService(ctx context.Context) *MessageService {
	return &MessageService{ctx: ctx}
}

// GetMessageChat get chat records
func (m *MessageService) GetMessageChat(req *user.MessageChatRequest) ([]*base.Message, error) {
	messages := make([]*base.Message, 0)
	from_user_id := req.CurrentUserId
	to_user_id := req.ToUserId
	pre_msg_time := req.PreMsgTime
	db_messages, err := db.GetMessageByIdPair(from_user_id, to_user_id, utils.MillTimeStampToTime(pre_msg_time))
	if err != nil {
		return messages, err
	}
	for _, db_message := range db_messages {
		messages = append(messages, &base.Message{
			Id:         db_message.ID,
			ToUserId:   db_message.ToUserId,
			FromUserId: db_message.FromUserId,
			Content:    db_message.Content,
			CreateTime: db_message.CreatedAt.UnixNano() / 1000000,
		})
	}
	return messages, nil
}

// MessageAction add a message
func (m *MessageService) MessageAction(req *user.MessageActionRequest) error {
	from_user_id := req.CurrentUserId
	to_user_id := req.ToUserId
	content := req.Content

	ok, err := db.AddNewMessage(&db.Messages{
		FromUserId: from_user_id,
		ToUserId:   to_user_id,
		Content:    content,
	})
	if err != nil {
		return err
	}
	if !ok {
		return errno.MessageAddFailedErr
	}
	return nil
}
