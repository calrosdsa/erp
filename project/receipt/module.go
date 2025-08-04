package receipt

import (
	"context"
	"erp/api/middlewares"
	"erp/internal/domain"
	"erp/pkg/di"
	"erp/pkg/system"
	receipt_event "erp/project/receipt/internal/handler/event"
	receipt_rest "erp/project/receipt/internal/handler/rest"
	receipt_repo "erp/project/receipt/internal/repository"
	receipt_ucase "erp/project/receipt/internal/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct{}

func (m Module) Startup(ctx context.Context, mono system.Service) (err error) {
	return Root(ctx, mono)
}

func Root(ctx context.Context, svc system.Service) (err error) {
	container := di.New()

	container.AddScoped(domain.DatabaseTransactionKey, func(c di.Container) (any, error) {
		return svc.DBConn().GetQ().Begin(),nil
	})
	receiptEventRepo := receipt_repo.NewRecieptEventRepo()
	receiptRepo := receipt_repo.NewReceiptRepository(
		svc.DBConn(),
		svc.Helpers(),
	)
	receiptUseCase := receipt_ucase.NewReceiptUseCase(
		svc.PermissionService(),
		svc.Logger(),
		receiptRepo,
		svc.EventBus(),
		container,
		svc.CoreService(),
		svc.StockService(),
		svc.DocumentService(),
	)
	m := middlewares.NewMiddlewares(svc.SessionService(), svc.HumaApi(), svc.Helpers().Jwt)
	receipt_rest.NewReceiptHandler(svc.HumaApi(),huma.Middlewares{m.Authenticate},
    svc.Helpers(),receiptUseCase,svc.PermissionService())

	receipt_event.NewReceiptEventHandler(
		svc.EventBus(),svc.Logger(),receiptEventRepo,
	)

	return nil
}