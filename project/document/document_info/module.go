package documentinfo

import (
	"context"
	"erp/pkg/system"
	documentinfo_event "erp/project/document/document_info/handler/event"
	documentinfo_rest "erp/project/document/document_info/handler/rest"
	documentinfo_repo "erp/project/document/document_info/repository"
	documentinfo_ucase "erp/project/document/document_info/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct{}

func (m *Module) Startup(ctx context.Context, svc system.Service) error {
	return Root(ctx, svc)
}

func Root(ctx context.Context, svc system.Service) error {
	documentInfoEventRepo := documentinfo_repo.NewDocumentEventRepository()
	documentInfoRepo := documentinfo_repo.NewDocumentInfoRepository(
		svc.DBConn(), svc.Helpers(),
	)
	documentUcase := documentinfo_ucase.NewDocumentInfoUseCase(
		svc.Logger(), documentInfoRepo, svc.PermissionService(),
	)
	documentinfo_rest.NewDocumentInfoHandler(svc.HumaApi(),
		svc.Helpers(), huma.Middlewares{svc.Middlewares().Authenticate}, documentUcase,
		svc.PermissionService())
	documentinfo_event.NewDocumentEventHandler(svc.Logger(), documentInfoEventRepo, svc.EventBus())

	return nil
}
