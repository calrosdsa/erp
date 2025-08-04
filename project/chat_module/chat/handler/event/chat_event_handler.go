package chat_event

import (
	"context"
	"erp/internal/domain"
	"erp/internal/domain/event"
	"erp/pkg/bus"
	"erp/pkg/logger"
	chat_repo "erp/project/chat_module/chat/repository"
)

type chatEvent struct {
	emitLog logger.EmitLog
	repo chat_repo.ChatEventRepo
}

func NewChatEventRepo(
	logger logger.Logger,
	bus bus.Bus,
	repo chat_repo.ChatEventRepo,
){
	h := chatEvent{
		emitLog: logger.EmitLog("chat-event"),
		repo: repo,
	}
	bus.RegisterHandler(domain.UserCreatedEvent,h.OnUserCreated())
}

func (h *chatEvent) OnUserCreated() bus.Handler{
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func(){
				if err != nil {
					h.emitLog.Err(err,logger.OptionsLog.WithMethod("OnUserCreated"))
				}
			}()
			payload,ok := e.Data.(event.UserCreatedEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.repo.OnUserAdded(ctx,payload)
			return
		},
		AbortOnError: true,
		Matcher: domain.UserCreatedEvent,
	}
}