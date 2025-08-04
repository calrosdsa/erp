package contact

import (
	"context"
	"erp/internal/domain"
	"erp/pkg/di"
	"erp/pkg/system"
	contact_event "erp/project/core/contact/handler/event"
	contact_rest "erp/project/core/contact/handler/rest"
	contact_repo "erp/project/core/contact/repository"
	contact_ucase "erp/project/core/contact/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct{}

func (m Module) Startup(ctx context.Context, svc system.Service) error {
	return Root(ctx, svc)
}

func Root(ctx context.Context, svc system.Service) error {
	contactRepo := contact_repo.NewContactRepository(
		svc.DBConn(), svc.Helpers(),
	)
	contactUseCase := contact_ucase.NewContactUseCase(
		svc.Helpers(), contactRepo,svc.PermissionService(),
		svc.CoreService(),
	)
	contact_rest.NewContactHandler(svc.HumaApi(), svc.Helpers(),
		huma.Middlewares{svc.Middlewares().Authenticate}, contactUseCase, svc.PermissionService())

	contact_event.NewContactEventHandler(svc.Logger(),
		svc.EventBus(), contactRepo)

	
	svc.Container().AddSingleton(domain.ContactUseCase, func(c di.Container) (any, error) {
		return contactUseCase,nil
	})
	return nil
}
