package sales_record

import (
	"context"
	"erp/pkg/system"
	sales_record_rest "erp/project/invoicing/sales/internal/handler/rest"
	sales_record_fsm "erp/project/invoicing/sales/internal/pkg/fsm"
	sales_record_repo "erp/project/invoicing/sales/internal/repository"
	sales_record_ucase "erp/project/invoicing/sales/internal/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct {
}

func (m Module) Startup(ctx context.Context, svc system.Service) error {
	return Root(ctx, svc)
}

func Root(ctx context.Context, svc system.Service) (err error) {
	fsm := sales_record_fsm.NewSalesRecordFsm()
	salesRecordRepo := sales_record_repo.NewSaleRecordRepo(svc.DBConn(), svc.Helpers())
	salesRecordUcase := sales_record_ucase.NewSalesRecordUcase(svc.Logger(), salesRecordRepo,
		svc.PermissionService(), svc.CoreService(), fsm,svc.Helpers())
	sales_record_rest.NewHandler(svc.HumaApi(), svc.Helpers(), salesRecordUcase,
		huma.Middlewares{svc.Middlewares().Authenticate}, svc.PermissionService())
	return nil
}
