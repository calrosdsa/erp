package invoice

import (
	"context"
	"erp/api/middlewares"
	"erp/internal/domain"
	"erp/pkg/di"
	"erp/pkg/system"
	invoice_event "erp/project/invoice/internal/handler/event"
	invoice_rest "erp/project/invoice/internal/handler/rest"
	invoice_fsm "erp/project/invoice/internal/pkg/fsm"
	invoice_repo "erp/project/invoice/internal/repository"
	invoice_ucase "erp/project/invoice/internal/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct{}

func (m Module) Startup(ctx context.Context, svc system.Service) error {
	return Root(ctx, svc)
}

func Root(ctx context.Context, svc system.Service) error {
	container := di.New()

	container.AddScoped(domain.DatabaseTransactionKey, func(c di.Container) (any, error) {
		return svc.DBConn().GetQ().Begin(),nil
	})

	invoiceFsm := invoice_fsm.NewInvoiceFsm()
	invoiceEventRepo := invoice_repo.NewInvoiceEventRepo(
		svc.Helpers(),
	)
	invoiceRepo := invoice_repo.NewInvoiceRepository(svc.DBConn(), svc.Helpers())
	invoiceUCase := invoice_ucase.NewInvoiceUseCase(
		svc.Logger(),svc.PermissionService(),invoiceRepo,invoiceFsm,svc.EventBus(),container,
		svc.CoreService(),svc.DocumentService(),
	)
	m := middlewares.NewMiddlewares(svc.SessionService(), svc.HumaApi(), svc.Helpers().Jwt)

	invoice_rest.NewInvoiceHandler(svc.HumaApi(),svc.Helpers(),
	svc.PermissionService(),huma.Middlewares{m.Authenticate},invoiceUCase)

	invoice_event.NewInvoiceEventHandler(
		svc.EventBus(),svc.Logger(),invoiceEventRepo)

	return nil
}
