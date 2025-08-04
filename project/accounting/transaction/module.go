package transaction

import (
	"context"
	"erp/pkg/system"
	transaction_event "erp/project/accounting/transaction/internal/handler/event"
	transaction_repo "erp/project/accounting/transaction/internal/repository"
)

type Module struct{}

func (m Module) Startup(ctx context.Context, svc system.Service) error {
	return Root(ctx, svc)
}

func Root(ctx context.Context, svc system.Service) error {
	transactionRepo := transaction_repo.NewTransactionRepository(svc.Helpers())
	// transactionEventRepo := transaction_repo.NewTransactionEventRepo(svc.DBConn())
	atxBuyingEventRepo := transaction_repo.NewAtxBuyingRepo(
		transactionRepo, svc.CoreService(), svc.AccountingService(), svc.Helpers(),
	)
	atxSellingEventRepo := transaction_repo.NewAtxSellingEventRepo(
		transactionRepo, svc.CoreService(), svc.AccountingService(), svc.Helpers(),
	)

	atxStockEntryEventRepo := transaction_repo.NewAtxStockEntryEventRepo(svc.SettingService(),
		transactionRepo, svc.AccountingService(), svc.Helpers())
	atxAccountingRepo := transaction_repo.NewAtxAccountingEventRepo(transactionRepo, svc.AccountingService(),
		svc.Helpers())
	// transaction_event.NewTransactionEventHandler(
	// 	svc.EventBus(), svc.Logger(), transactionEventRepo,
	// )
	transaction_event.NewAtxEBuyingventHandler(
		svc.EventBus(),
		svc.Logger(),
		atxBuyingEventRepo,
		atxSellingEventRepo,
		atxStockEntryEventRepo,
		svc.AccountingService(),
	)
	transaction_event.NewAtxAccountingEventHandler(
		svc.EventBus(),
		svc.Logger(),
		atxAccountingRepo,
	)
	return nil
}
