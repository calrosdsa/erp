package connection

import (
	"context"
	"erp/pkg/system"
	connection_rest "erp/project/core/connection/handler/rest"
	connection_repo "erp/project/core/connection/repository"
	connection_ucase "erp/project/core/connection/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct{}

func (m *Module) Startup(ctx context.Context, svc system.Service) error {
	return Root(ctx, svc)
}

func Root(ctx context.Context, svc system.Service) error {
	connectionRepo := connection_repo.NewConnectionRepo(svc.DBConn(), svc.Helpers())
	connectionUcase := connection_ucase.NewConnectionUcase(svc.Logger(), connectionRepo)
	connection_rest.NewConnectionHandler(svc.HumaApi(), svc.Helpers(), connectionUcase,
		huma.Middlewares{svc.Middlewares().Authenticate}, svc.PermissionService())
	return nil
}
