package quotation

import (
	"context"
	"erp/internal/domain"
	"erp/pkg/di"
	"erp/pkg/system"
	rest_quotation "erp/project/quotation/handler/rest"
	quotation_fsm "erp/project/quotation/internal/pkg/fsm"
	quotation_repo "erp/project/quotation/repository"
	quotation_ucase "erp/project/quotation/usecase"

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
	quotationFsm := quotation_fsm.NewQuotationFsm()
	quotationRepo := quotation_repo.NewQuotationRepository(svc.DBConn(), svc.Helpers())
	quotationUcase := quotation_ucase.NewQuotationUseCase(svc.Logger(),
		svc.PermissionService(), quotationRepo, quotationFsm, svc.EventBus(),
		container, svc.CoreService(), svc.DocumentService())
	rest_quotation.NewQuotationHandler(svc.HumaApi(), svc.Helpers(), quotationUcase,
		huma.Middlewares{svc.Middlewares().Authenticate}, svc.PermissionService())

	svc.Container().AddSingleton(domain.QuotationUseCase, func(c di.Container) (any, error) {
		return quotationUcase, nil
	})
	return nil
}
