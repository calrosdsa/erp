package chat_ucase

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/pkg/logger"
	chat_repo "erp/project/chat_module/chat/repository"
)

type ChatUseCase interface {
	GetChat(req *common.RequestContext, d dto.RequestEntity) (res dto.ChatDetailDto, err error)
	GetChats(req *common.RequestContext,d dto.ChatsRequest) (res []dto.ChatDto, err error)
	UpdateMemberLastRead(req *common.RequestContext,d dto.RequestEntity)(err error)
}

type chatUcase struct {
	repo chat_repo.ChatRepository
	emitLog logger.EmitLog
}

func NewChatUseCase(
	repo chat_repo.ChatRepository,
	logger logger.Logger,
) ChatUseCase {
	return &chatUcase{
		repo: repo,
		emitLog: logger.EmitLog("chat-usecase"),
	}
}

func(u *chatUcase) UpdateMemberLastRead(req *common.RequestContext,d dto.RequestEntity)(err error){
	defer func(){
		if err != nil {
			u.emitLog.Err(err,logger.OptionsLog.WithMethod("UpdateMemberLastRead"))
		}
	}()
	err = u.repo.UpdateMemberLastRead(req,d)
	return
}

func (u *chatUcase) GetChat(req *common.RequestContext, d dto.RequestEntity) (res dto.ChatDetailDto, err error) {
	defer func(){
		if err != nil {
			u.emitLog.Err(err,logger.OptionsLog.WithMethod("GetChat"))
		}
	}()
	res,err = u.repo.GetChat(req,d)
	
	return
}
func (u *chatUcase) GetChats(req *common.RequestContext,d dto.ChatsRequest) (res []dto.ChatDto, err error) {
	defer func(){
		if err != nil {
			u.emitLog.Err(err,logger.OptionsLog.WithMethod("GetChats"))
		}
	}()
	res,err = u.repo.GetChats(req,d)
	return
}