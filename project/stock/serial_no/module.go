package serial_no

import (
	"context"
	"erp/pkg/system"
	serial_no_event "erp/project/stock/serial_no/internal/handler/event"
	rest_serial_no "erp/project/stock/serial_no/internal/handler/rest"
	serial_no_repo "erp/project/stock/serial_no/internal/repository"
	serial_no_ucase "erp/project/stock/serial_no/internal/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct {
}

func (m *Module) Startup(ctx context.Context, svc system.Service) error {
	return Root(ctx, svc)
}

func Root(ctx context.Context, svc system.Service) error {
	serialNoEventRepo := serial_no_repo.NewSerialEventRepo(svc.Helpers(),svc.AccountingService())
	serialNoRepo := serial_no_repo.NewSerialRepository(svc.DBConn(), svc.Helpers())
	serialNoUcase := serial_no_ucase.NewSerialUcase(svc.Logger(), serialNoRepo, svc.PermissionService(),
		svc.CoreService())
	rest_serial_no.NewSerialHandler(svc.HumaApi(), svc.Helpers(), serialNoUcase,
		huma.Middlewares{svc.Middlewares().Authenticate}, svc.PermissionService())
	serial_no_event.NewSerialNoEventHandler(svc.Logger(),svc.EventBus(),serialNoEventRepo)
	return nil
}
