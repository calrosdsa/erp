package payment

import (
	"context"
	"erp/internal/domain"
	"erp/pkg/di"
	"erp/pkg/system"
	payment_event "erp/project/accounting/payment/internal/handler/event"
	payment_rest "erp/project/accounting/payment/internal/handler/rest"
	payment_fsm "erp/project/accounting/payment/internal/pkg/fsm"
	payment_pdf "erp/project/accounting/payment/internal/pkg/pdf"
	payment_repo "erp/project/accounting/payment/internal/repository"
	payment_ucase "erp/project/accounting/payment/internal/usecase"
	"fmt"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct{}

func (m Module) Startup(ctx context.Context, svc system.Service) error {
	return Root(ctx, svc)
}

func Root(ctx context.Context, svc system.Service) error {
	fmt.Println("INIT PAYMENT MODULE")
	container := di.New()
	container.AddScoped(domain.DatabaseTransactionKey, func(c di.Container) (any, error) {
		return svc.DBConn().GetQ().Begin(), nil
	})
	paymentFsm := payment_fsm.NewPaymentFsm()
	paymentPdf := payment_pdf.NewPaymentPdf(
		svc.Helpers(),
		svc.DBConn().GetQ(),
	)
	paymentEventRepo := payment_repo.NewPaymentEventRepository()
	paymentRepo := payment_repo.NewPaymentRepository(svc.DBConn(), svc.Helpers(), svc.SettingService(),
		svc.AccountingService())
	paymentUcase := payment_ucase.NewPaymentUseCase(
		paymentRepo, svc.Logger(), svc.PermissionService(), svc.EventBus(),
		container, paymentFsm,paymentPdf,svc.CoreService(),
	)
	payment_rest.NewPaymentHandler(
		svc.HumaApi(), huma.Middlewares{svc.Middlewares().Authenticate},
		svc.Helpers(),
		paymentUcase,
		svc.PermissionService(),
	)

	payment_event.NewPaymentEventHandler(
		svc.EventBus(), svc.Logger(), paymentEventRepo,
	)
	return nil
}
