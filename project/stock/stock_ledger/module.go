package stocke_ledger

import (
	"context"
	"erp/pkg/system"
	stock_ledger_event "erp/project/stock/stock_ledger/internal/handler/event"
	stock_ledger_rest "erp/project/stock/stock_ledger/internal/handler/rest"
	stock_ledger_repo "erp/project/stock/stock_ledger/internal/repository"
	stock_ledger_ucase "erp/project/stock/stock_ledger/internal/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct{}

func (m *Module) Startup(ctx context.Context, svc system.Service) error {
	return Root(ctx, svc)
}

func Root(ctx context.Context, svc system.Service) error {

	stockTxRepo := stock_ledger_repo.NewStockLedgeTX(svc.DBConn(), svc.Helpers())
	stockTxSellingRepo := stock_ledger_repo.NewStockTxSellingRepo(stockTxRepo,svc.Helpers(),svc.AccountingService())
	stockTxBuyingRepo := stock_ledger_repo.NewStockTxBuyingRepo(
		stockTxRepo,svc.AccountingService(),svc.Helpers(),
	)
	stockTxEntryRepo := stock_ledger_repo.NewStockTxStockEntryRepo(stockTxRepo,svc.Helpers(),svc.AccountingService())
	stockLederUcase := stock_ledger_ucase.NewStockLedgerUseCase(
		svc.Logger(), stockTxRepo,
	)

	stock_ledger_event.NewStockTxBuyingHandler(svc.EventBus(), svc.Logger(),
	stockTxBuyingRepo, stockTxSellingRepo,stockTxEntryRepo,stockTxRepo)
	stock_ledger_rest.NewStockLedgerHandler(svc.HumaApi(), huma.Middlewares{svc.Middlewares().Authenticate},
		svc.Helpers(), stockLederUcase, svc.PermissionService())
	return nil
}
