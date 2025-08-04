package charges_template

import (
	"context"
	"erp/internal/domain"
	"erp/pkg/di"
	"erp/pkg/system"
	charges_template_rest "erp/project/accounting/charges_template/internal/handler/rest"
	charges_template_repo "erp/project/accounting/charges_template/internal/repository"
	charge_template_ucase "erp/project/accounting/charges_template/internal/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct{}

func (m *Module) Startup(ctx context.Context,svc system.Service) error {
	return Root(ctx,svc)
}

func Root(ctx context.Context,svc system.Service) error {
	container := di.New()
	container.AddScoped(domain.DatabaseTransactionKey, func(c di.Container) (any, error) {
		return svc.DBConn().GetQ().Begin(), nil
	})
	chargeTemplateRepo := charges_template_repo.NewChargesTemplateRepo(
		svc.DBConn(),svc.Helpers(),)
	chargeTemplateUcase := charge_template_ucase.NewChargesTemplateUcase(
		svc.Logger(),chargeTemplateRepo,svc.PermissionService(),
		svc.CoreService(),svc.EventBus(),container,
	)
	charges_template_rest.NewHandler(svc.HumaApi(),svc.Helpers(),chargeTemplateUcase,
	huma.Middlewares{svc.Middlewares().Authenticate},svc.PermissionService())
	return nil
}