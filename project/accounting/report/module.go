package acct_report

import (
	"context"
	"erp/pkg/system"
	acct_report_rest "erp/project/accounting/report/internal/handler/rest"
	acct_report_repo "erp/project/accounting/report/internal/repository"
	acct_report_ucase "erp/project/accounting/report/internal/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct{}

func (m *Module) Startup(ctx context.Context, svc system.Service) error {
	return Root(ctx, svc)
}

func Root(ctx context.Context, svc system.Service) error {
	financialStmtRepo := acct_report_repo.NewAcctFinancialStatementRepo(svc.DBConn(),svc.Helpers())
	acctReportRepo := acct_report_repo.NewAcctReportRepository(svc.DBConn(), svc.Helpers())
	acctReportUcase := acct_report_ucase.NewAcctReportUseCase(
		svc.Logger(), svc.PermissionService(), acctReportRepo,
	)
	financialStmtUcase := acct_report_ucase.NewFinancialStmtUcase(svc.Logger(), financialStmtRepo)
	acct_report_rest.NewAcctReportHandler(
		svc.HumaApi(), huma.Middlewares{svc.Middlewares().Authenticate},
		svc.Helpers(), acctReportUcase, svc.PermissionService(),
	)
	acct_report_rest.NewFinancialStatementHandler(
		svc.HumaApi(), huma.Middlewares{svc.Middlewares().Authenticate},
		svc.Helpers(), financialStmtUcase, svc.PermissionService(),
	)
	return nil
}
