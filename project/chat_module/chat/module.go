package chat

import (
	"context"
	"erp/pkg/system"
	chat_event "erp/project/chat_module/chat/handler/event"
	chat_rest "erp/project/chat_module/chat/handler/rest"
	chat_repo "erp/project/chat_module/chat/repository"
	chat_ucase "erp/project/chat_module/chat/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct{}

func (m *Module) Startup(ctx context.Context, svc system.Service) error {
	return Root(ctx, svc)
}

func Root(ctx context.Context, svc system.Service) error {
	chatEventRepo := chat_repo.NewChatEventRepo()
	chatRepo := chat_repo.NewChatRepo(svc.DBConn(), svc.Helpers())
	chatUcase := chat_ucase.NewChatUseCase(chatRepo, svc.Logger())
	chat_rest.NewHandler(svc.HumaApi(), svc.Helpers(),
		huma.Middlewares{svc.Middlewares().Authenticate}, chatUcase)
	chat_event.NewChatEventRepo(svc.Logger(),svc.EventBus(),chatEventRepo)
	return nil
}
