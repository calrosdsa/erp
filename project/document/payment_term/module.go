package payment_terms

import (
	"context"
	"erp/pkg/system"
	payment_terms_event "erp/project/document/payment_term/handler/event"
	payment_terms_rest "erp/project/document/payment_term/handler/rest"
	payment_terms_fsm "erp/project/document/payment_term/internal/pkg/fsm"
	payment_terms_repo "erp/project/document/payment_term/repository"
	payment_terms_ucase "erp/project/document/payment_term/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct{}

func (m *Module) Startup(ctx context.Context, svc system.Service) error {
	return Root(ctx, svc)
}

func Root(ctx context.Context, svc system.Service) error {
	fsm := payment_terms_fsm.NewFsm()
	paymentTermsRepo := payment_terms_repo.NewRepository(
		svc.DBConn(), svc.Helpers(),
	)
	paymentTermsLineRepo := payment_terms_repo.NewPaymentTermsLineRepo(
		svc.DBConn(),svc.Helpers(),
	)
	paymentTermsUcase := payment_terms_ucase.NewUseCase(
		svc.Logger(), svc.CoreService(), svc.PermissionService(), paymentTermsRepo, fsm,
		paymentTermsLineRepo,
	)
	payment_terms_rest.NewHandler(
		svc.HumaApi(), svc.Helpers(), paymentTermsUcase, huma.Middlewares{svc.Middlewares().Authenticate},
		svc.PermissionService(),
	)
	payment_terms_event.NewPtLineEventHandler(
		svc.Logger(),
		svc.EventBus(),
		paymentTermsLineRepo,
	)
	return nil
}
