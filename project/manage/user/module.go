package user

import (
	"context"
	"erp/api/middlewares"
	"erp/pkg/system"
	user_rest "erp/project/manage/user/internal/handler/rest"
	user_repo "erp/project/manage/user/internal/repository"
	user_ucase "erp/project/manage/user/internal/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct{}

func (m Module)Startup(ctx context.Context, svc system.Service) error {
	return Root(ctx, svc)
}

func Root(ctx context.Context, svc system.Service) error {
	userRepo := user_repo.NewUserRepository(svc.DBConn(), svc.Helpers(), svc.Config().PG)
	userUseCase := user_ucase.NewUserUseCase(
		svc.EventBus(),svc.Logger(), svc.PermissionService(), userRepo,svc.Container(),
	)
	m := middlewares.NewMiddlewares(svc.SessionService(), svc.HumaApi(), svc.Helpers().Jwt)
	user_rest.NewUserHandler(svc.HumaApi(), svc.Helpers(),
		svc.PermissionService(), huma.Middlewares{m.Authenticate}, userUseCase)
	return nil
}
