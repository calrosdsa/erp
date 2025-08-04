package address

import (
	"context"
	"erp/pkg/system"
	address_rest "erp/project/core/address/handler/rest"
	address_fsm "erp/project/core/address/pkg/fsm"
	address_repo "erp/project/core/address/repository"
	address_ucase "erp/project/core/address/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct{}

func (m *Module) Startup(ctx context.Context, svc system.Service) error {
	return Root(ctx, svc)
}

func Root(ctx context.Context, svc system.Service) error {
	addressFsm := address_fsm.NewFsm()

	addressRepo := address_repo.NewAddressRepository(svc.DBConn(), svc.Helpers(),)
	addressUcase := address_ucase.NewAddressUseCase(svc.Logger(), addressRepo,
		svc.PermissionService(), svc.CoreService(),addressFsm)
	address_rest.NewAddressHandler(svc.HumaApi(), svc.Helpers(),
		huma.Middlewares{svc.Middlewares().Authenticate}, addressUcase, svc.PermissionService())
	return nil
}
