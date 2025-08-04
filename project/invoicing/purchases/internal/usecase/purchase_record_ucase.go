package purchase_record_ucase

import (
	"bytes"
	"erp/api/common"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	"erp/pkg/exporter"
	"erp/pkg/fsm"
	"erp/pkg/logger"
	purchase_record_pdf "erp/project/invoicing/purchases/internal/pkg/pdf"
	purchase_record_repo "erp/project/invoicing/purchases/internal/repository"
	"fmt"
)

type PurchaseRecordUcase interface {
	GetPurchaseRecord(req *common.RequestContext, d *dto.RequestEntity) (
		res dto.ResultEntity[dto.PurchaseRecordDto], err error)
	CreatePurchaseRecord(req *common.RequestContext, d *dto.CreatePurchaseRecordRequest) (
		res dto.PurchaseRecordDto, err error)
	GetPurchaseRecords(req *common.RequestContext, d *dto.PurchaseRecordsRequest) (
		res dto.PaginationResult[[]dto.PurchaseRecordDto], err error)
	EditPurchaseRecord(req *common.RequestContext, d *dto.EditPurchaseRecordRequest) (err error)
	UpdateStatus(req *common.RequestContext, d *dto.UpdateStatusWithEvent) error
	ExportData(req *common.RequestContext, d *dto.ExportDataRequest) (res *bytes.Buffer, err error)
	ExportPurchaseRecord(req *common.RequestContext, i dto.ExportDocumentData) (res []byte, err error) 
}

type purchaseRecordUcase struct {
	emitLog            logger.EmitLog
	purchaseRecordRepo purchase_record_repo.PurchaseRecordRepo
	permission         repository.PermissionService
	core               repository.CoreService
	fsm                fsm.FsmState
	excelExporter      exporter.ExcelExporter
	currency           helpers.CurrencyHelper
	pdfGenerator purchase_record_pdf.PurchaseRecordPDF
}

func NewPurchaseRecordUcase(
	logger logger.Logger,
	purchaseRecordRepo purchase_record_repo.PurchaseRecordRepo,
	permission repository.PermissionService,
	core repository.CoreService,
	fsm fsm.FsmState,
	helpers *helpers.Helpers,
	pdfGenerator purchase_record_pdf.PurchaseRecordPDF,
) PurchaseRecordUcase {
	return &purchaseRecordUcase{
		emitLog:            logger.EmitLog("purchase-record-usecase"),
		purchaseRecordRepo: purchaseRecordRepo,
		permission:         permission,
		core:               core,
		fsm:                fsm,
		excelExporter:      helpers.ExcelExporter,
		currency:           helpers.Currency,
		pdfGenerator: pdfGenerator,
	}
}

func(u *purchaseRecordUcase) ExportPurchaseRecord(req *common.RequestContext, i dto.ExportDocumentData) (res []byte, err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("ExportPurchaseRecord"))
		}
	}()
	purchaseRecord,err := u.GetPurchaseRecord(req,&dto.RequestEntity{
		ID: i.ID,
	})
	if err != nil {
		return
	}
	res,err =u.pdfGenerator.GeneratePurchaseRecordDocument(req,purchaseRecord.Entity)
	return
}

func (u *purchaseRecordUcase) UpdateStatus(req *common.RequestContext, d *dto.UpdateStatusWithEvent) (err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("UpdateStatus"))
		}
	}()
	if err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.PURCHASE_RECORD, domain.EDIT); err != nil {
		return
	}
	nextState, err := u.fsm.NextState(d.Body.CurrentState, d.Body.Events)
	if err != nil {
		return err
	}
	err = u.purchaseRecordRepo.UpdateStatus(req, d, nextState)
	return
}

func (u *purchaseRecordUcase) EditPurchaseRecord(req *common.RequestContext, d *dto.EditPurchaseRecordRequest) (err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("EditPurchaseRecord"))
		}
	}()
	err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.PURCHASE_RECORD, domain.EDIT)
	if err != nil {
		return
	}
	err = u.purchaseRecordRepo.EditPurchaseRecord(req, d)
	return
}

func (u *purchaseRecordUcase) GetPurchaseRecord(req *common.RequestContext, d *dto.RequestEntity) (
	res dto.ResultEntity[dto.PurchaseRecordDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetPurchaseRecord"))
		}
	}()
	err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.PURCHASE_RECORD, domain.VIEW)
	if err != nil {
		return
	}
	res, err = u.purchaseRecordRepo.GetPurchaseRecord(req, d)
	if err != nil {
		return
	}
	res.Activities = u.core.GerActivitiesByPartyID(req, res.Entity.ID)
	return
}
func (u *purchaseRecordUcase) CreatePurchaseRecord(req *common.RequestContext, d *dto.CreatePurchaseRecordRequest) (
	res dto.PurchaseRecordDto, err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("CreatePurchaseRecord"))
		}
	}()
	err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.PURCHASE_RECORD, domain.CREATE)
	if err != nil {
		return
	}
	purchaseRecord, err := u.purchaseRecordRepo.CreatePurchaseRecord(req, d)
	if err != nil {
		return
	}
	res = dto.PurchaseRecordDtoFromModel(&purchaseRecord)
	return
}

func (u *purchaseRecordUcase) ExportData(req *common.RequestContext, d *dto.ExportDataRequest) (
	res *bytes.Buffer, err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("ExportData"))
		}
	}()
	if err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.PURCHASE_RECORD, domain.EDIT); err != nil {
		return
	}
	purchaseRecords, err := u.purchaseRecordRepo.ExportData(req, d)
	if err != nil {
		return
	}
	fmt.Println("SALES RECORDS", d.Body.Data)
	data := make([][]interface{}, len(purchaseRecords))
	for idx, d := range purchaseRecords {
		//COmplete
		data[idx] = []interface{}{
			d.SupplierNit,
			d.SupplierBusinessName,
			d.AuthorizationCode,
			d.InvoiceNo,
			d.DuiDimNo,
			d.InvoiceDuiDimDate,
			u.currency.IntToFloat(int64(d.TotalPurchaseAmount)),
			u.currency.IntToFloat(int64(d.IceAmount)),
			u.currency.IntToFloat(int64(d.IehdAmount)),
			u.currency.IntToFloat(int64(d.IpjAmount)),
			u.currency.IntToFloat(int64(d.TaxRates)),
			u.currency.IntToFloat(int64(d.OtherNotSubjectToTaxCredit)),
			u.currency.IntToFloat(int64(d.ExemptAmounts)),
			u.currency.IntToFloat(int64(d.ZeroRateTaxablePurchasesAmount)),
			u.currency.IntToFloat(int64(d.Subtotal)),
			u.currency.IntToFloat(int64(d.DiscountsBonusRebatesSubjectToVat)),
			u.currency.IntToFloat(int64(d.GiftCardAmount)),
			u.currency.IntToFloat(int64(d.CfBaseAmount)),
			u.currency.IntToFloat(int64(d.TaxCredit)),
			d.PurchaseType,
			d.ControlCode,
			d.WithTaxCreditRight,
			d.ConsolidationStatus,
		}
	}
	headers := []interface{}{"NIT PROVEEDOR", "RAZON SOCIAL PROVEEDOR", "CODIGO DE AUTORIZACIÓN",
		"Nº DE LA FACTURA", "NUMERO DUI/DIM", "FECHA DE FACTURA/DUI/DIM", "IMPORTE TOTAL DE LA COMPRA",
		"IMPORTE ICE", "IMPORTE IEHD", "IMPORTE IPJ",
		"TASAS", "OTROS NO SUJETOS AL IVA", "IMPORTES EXENTOS", "VENTAS GRAVADAS A TASA CERO",
		"SUBTOTAL", "DESCUENTOS BONIFICACIONES Y REBAJAS SUJETAS AL IVA", "IMPORTE GIFT CARD",
		"IMPORTE BASE CF",
		"CREDITO FISCAL", "TIPO COMPRA", "CODIGO DE CONTROL", "CON DERECHO A CREDITO FISCAL", "ESTADO CONSOLIDACION"}
	res, err = u.excelExporter.Export("Sheet", headers, data)
	return
}

func (u *purchaseRecordUcase) GetPurchaseRecords(req *common.RequestContext, d *dto.PurchaseRecordsRequest) (
	res dto.PaginationResult[[]dto.PurchaseRecordDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetPurchaseRecords"))
		}
	}()
	err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.PURCHASE_RECORD, domain.VIEW)
	if err != nil {
		return res, domain.ACTION_NOT_ALLOWED
	}

	res, err = u.purchaseRecordRepo.GetPurchaseRecords(req, d)
	if err != nil {
		return
	}
	res.FilterOptions = u.purchaseRecordRepo.GetFilterOptions()
	return
}
