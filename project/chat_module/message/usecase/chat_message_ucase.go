package chat_message_ucase

import (
	"context"
	"encoding/json"
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/proto"
	"erp/pkg/logger"
	"erp/pkg/ws"
	chat_message_repo "erp/project/chat_module/message/repository"
)

type ChatMessageUcase interface {
	GetMessages(req *common.RequestContext, d dto.ChatMessagesRequest) (res []dto.ChatMessageDto, err error)
	CreateMessage(req *common.RequestContext, d dto.ChatMessageData) (res dto.ChatMessageDto, err error)
}

type chatMessageUcase struct {
	emitLog logger.EmitLog
	repo    chat_message_repo.ChatMessageRepo
	wsConn  ws.WsConn
}

func NewChatMessageUcase(
	logger logger.Logger,
	repo chat_message_repo.ChatMessageRepo,
	wsConn ws.WsConn,
) ChatMessageUcase {
	return &chatMessageUcase{
		repo:    repo,
		emitLog: logger.EmitLog("chat-message-usecase"),
		wsConn:  wsConn,
	}
}

func (u *chatMessageUcase) GetMessages(req *common.RequestContext, d dto.ChatMessagesRequest) (res []dto.ChatMessageDto, err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetMessages"))
		}
	}()
	res, err = u.repo.GetMessages(req, d)
	return
}
func (u *chatMessageUcase) CreateMessage(req *common.RequestContext, d dto.ChatMessageData) (res dto.ChatMessageDto, err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("CreateMessages"))
		}
	}()
	res, err = u.repo.CreateMessage(req, d)
	if err != nil {
		return
	}

	go func(message dto.ChatMessageDto) {
		memberProfileIds, err := u.repo.GetMembersIds(d.Fields.ChatID)
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetMembersIds"))
			return
		}

		msg,_ := json.Marshal(res)
		u.wsConn.PublishToSubscribers(context.Background(), memberProfileIds, string(msg), proto.MessageType_ChatMessage)
	}(res)

	return
}
