package stock_entry

import (
	"context"
	"erp/internal/domain"
	"erp/pkg/di"
	"erp/pkg/system"
	rest_stock_entry "erp/project/stock/stock_entry/internal/handler/rest"
	stock_entry_fsm "erp/project/stock/stock_entry/internal/pkg/fsm"
	stock_entry_repo "erp/project/stock/stock_entry/internal/repository"
	stock_entry_ucase "erp/project/stock/stock_entry/internal/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct{}

func (m *Module) Startup(ctx context.Context, svc system.Service) error {
	return Root(ctx, svc)
}

func Root(ctx context.Context, svc system.Service) error {
	container := di.New()
	container.AddScoped(domain.DatabaseTransactionKey, func(c di.Container) (any, error) {
		return svc.DBConn().GetQ().Begin(), nil
	})
	stockEntryFsm := stock_entry_fsm.NewStockEntryFsm()
	stockEntry := stock_entry_repo.NewStockEntryRepository(svc.DBConn(), svc.Helpers())
	stockUcase := stock_entry_ucase.NewStockEntrytUcase(svc.Logger(), stockEntry,
		svc.PermissionService(), svc.CoreService(),
		container, svc.EventBus(), stockEntryFsm,svc.DocumentService())
	rest_stock_entry.NewStockEntryHandler(svc.HumaApi(), svc.Helpers(), stockUcase,
		huma.Middlewares{svc.Middlewares().Authenticate}, svc.PermissionService())
	return nil
}
