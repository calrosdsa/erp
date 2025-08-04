package purchase_record

import (
	"context"
	"erp/pkg/system"
	purchase_record_rest "erp/project/invoicing/purchases/internal/handler/rest"
	purchase_record_fsm "erp/project/invoicing/purchases/internal/pkg/fsm"
	purchase_record_pdf "erp/project/invoicing/purchases/internal/pkg/pdf"
	purchase_record_repo "erp/project/invoicing/purchases/internal/repository"
	purchase_record_ucase "erp/project/invoicing/purchases/internal/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct{}

func (m Module) Startup(ctx context.Context, svc system.Service) error {
	return Root(ctx, svc)
}

func Root(ctx context.Context, svc system.Service) error {
	pdfGenerator := purchase_record_pdf.NewPurchaseRecordPdf(svc.Helpers(),svc.DBConn().GetQ())
	fsm := purchase_record_fsm.NewPurchaseRecordFsm()
	purchaseRecordRepo := purchase_record_repo.NewPurchaseRecordRepo(svc.DBConn(), svc.Helpers())
	purchaseRecordUcase := purchase_record_ucase.NewPurchaseRecordUcase(svc.Logger(), purchaseRecordRepo,
		svc.PermissionService(), svc.CoreService(), fsm, svc.Helpers(),pdfGenerator)
	purchase_record_rest.NewHandler(svc.HumaApi(), svc.Helpers(), purchaseRecordUcase,
		huma.Middlewares{svc.Middlewares().Authenticate}, svc.PermissionService())
	return nil
}
