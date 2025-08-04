package chat_message

import (
	"context"
	"erp/pkg/system"
	chat_message_rest "erp/project/chat_module/message/handler/rest"
	chat_message_repo "erp/project/chat_module/message/repository"
	chat_message_ucase "erp/project/chat_module/message/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct{}

func (m *Module) Startup(ctx context.Context, svc system.Service) error {
	return Root(ctx, svc)
}

func Root(ctx context.Context, svc system.Service) error {
	chatMessageRepo := chat_message_repo.NewChatMessageRepo(svc.DBConn(), svc.Helpers())
	chatMessageUcase := chat_message_ucase.NewChatMessageUcase(svc.Logger(), chatMessageRepo,
		svc.WsConn())
	chat_message_rest.NewHandler(svc.HumaApi(), svc.Helpers(),
		huma.Middlewares{svc.Middlewares().Authenticate}, chatMessageUcase)
	return nil
}
