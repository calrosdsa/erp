package payment_terms_template

import (
	"context"
	"erp/pkg/system"
	payment_terms_t_rest "erp/project/document/payment_terms_template/handler/rest"
	payment_terms_t_fsm "erp/project/document/payment_terms_template/internal/pkg/fsm"
	payment_terms_t_repo "erp/project/document/payment_terms_template/repository"
	payment_terms_t_ucase "erp/project/document/payment_terms_template/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct{}

func (m *Module) Startup(ctx context.Context, svc system.Service) error {
	return Root(ctx, svc)
}

func Root(ctx context.Context, svc system.Service) error {
	fsm := payment_terms_t_fsm.NewFsm()
	paymentTermsTemplateRepo := payment_terms_t_repo.NewRepository(
		svc.DBConn(), svc.Helpers(),
	)
	paymentTermsTemplateUcase := payment_terms_t_ucase.NewUseCase(
		svc.Logger(), svc.CoreService(), svc.PermissionService(), paymentTermsTemplateRepo, fsm,
		svc.EventBus(), svc.Container(),
	)
	payment_terms_t_rest.NewHandler(
		svc.HumaApi(), svc.Helpers(), paymentTermsTemplateUcase, huma.Middlewares{svc.Middlewares().Authenticate},
		svc.PermissionService(),
	)
	return nil
}
