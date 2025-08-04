package cash_outflow

import (
	"context"
	"erp/pkg/system"
	cash_outflow_rest "erp/project/accounting/cash_outflow/handler/rest"
	cash_outflow_fsm "erp/project/accounting/cash_outflow/pkg/fsm"
	cash_outflow_pdf "erp/project/accounting/cash_outflow/pkg/pdf"
	cash_outflow_repo "erp/project/accounting/cash_outflow/repository"
	cash_outflow_ucase "erp/project/accounting/cash_outflow/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct{}

func (m *Module) Startup(ctx context.Context, svc system.Service) error {
	return Root(ctx, svc)
}

func Root(ctx context.Context, svc system.Service) error {
	cashOutflowPdf := cash_outflow_pdf.NewCashOutflowPdf(
		svc.Helpers(),svc.DBConn().GetQ(),svc.DocumentService(),
	)
	cashOutflowFsm := cash_outflow_fsm.NewFsm()
	cashOutflowRepo := cash_outflow_repo.NewRepository(svc.DBConn(), svc.Helpers())
	cashOutflowUcase := cash_outflow_ucase.NewCashOutflowUseCase(
		svc.Logger(), svc.PermissionService(), cashOutflowRepo, cashOutflowFsm,
		svc.EventBus(), svc.Container(), svc.CoreService(), svc.DocumentService(),
		cashOutflowPdf,
	)
	cash_outflow_rest.NewHandler(svc.HumaApi(), svc.Helpers(),
		huma.Middlewares{svc.Middlewares().Authenticate}, svc.PermissionService(), cashOutflowUcase)
	return nil
}
