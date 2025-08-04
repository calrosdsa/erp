package role

import (
	"context"
	"erp/api/middlewares"
	"erp/pkg/system"
	role_event "erp/project/manage/role/internal/handler/event"
	role_rest "erp/project/manage/role/internal/handler/rest"
	role_repo "erp/project/manage/role/internal/repository"
	role_ucase "erp/project/manage/role/internal/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct{}

func (m Module) Startup(ctx context.Context, svc system.Service) error {
	return Root(ctx, svc)
}

func Root(ctx context.Context, svc system.Service) error {
	roleEventRepo := role_repo.NewRoleEventRepository()
	roleRepo := role_repo.NewRoleRepository(
		svc.DBConn(), svc.Logger(), svc.Helpers(), svc.PermissionService())
	roleUseCase := role_ucase.NewRoleUseCase(
		svc.Logger(), svc.PermissionService(), roleRepo,svc.Container(),
	)
	m := middlewares.NewMiddlewares(svc.SessionService(), svc.HumaApi(), svc.Helpers().Jwt)
	role_rest.NewRoleHandler(svc.HumaApi(), svc.Helpers(), svc.PermissionService(), huma.Middlewares{m.Authenticate},
		roleUseCase)
	role_event.NewRoleEventHandler(svc.Logger(), svc.EventBus(), roleEventRepo)
	return nil
}
