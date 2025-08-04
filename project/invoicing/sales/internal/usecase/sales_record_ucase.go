package sales_record_ucase

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
	sales_record_repo "erp/project/invoicing/sales/internal/repository"
	"fmt"
)

type SalesRecordUcase interface {
	GetSalesRecord(req *common.RequestContext, d *dto.RequestEntity) (
		res dto.ResultEntity[dto.SalesRecordDto], err error)
	CreateSalesRecord(req *common.RequestContext, d *dto.CreateSalesRecordRequest) (
		res dto.SalesRecordDto, err error)
	GetSalesRecords(req *common.RequestContext, d *dto.SalesRecordsRequest) (
		res dto.PaginationResult[[]dto.SalesRecordDto], err error)
	EditSalesRecord(req *common.RequestContext, d *dto.EditSalesRecordRequest) (err error)
	UpdateStatus(req *common.RequestContext, d *dto.UpdateStatusWithEvent) error
	ExportData(req *common.RequestContext, d *dto.ExportDataRequest) (res *bytes.Buffer, err error)
}

type salesRecordUcase struct {
	emitLog         logger.EmitLog
	salesRecordRepo sales_record_repo.SalesRecordRepo
	permission      repository.PermissionService
	core            repository.CoreService
	fsm             fsm.FsmState
	currency        helpers.CurrencyHelper
	excelExporter   exporter.ExcelExporter
}

func NewSalesRecordUcase(
	logger logger.Logger,
	salesRecordRepo sales_record_repo.SalesRecordRepo,
	permission repository.PermissionService,
	core repository.CoreService,
	fsm fsm.FsmState,
	helpers *helpers.Helpers,
) SalesRecordUcase {
	return &salesRecordUcase{
		emitLog:         logger.EmitLog("sales-record-usecase"),
		salesRecordRepo: salesRecordRepo,
		permission:      permission,
		core:            core,
		fsm:             fsm,
		currency:        helpers.Currency,
		excelExporter:   helpers.ExcelExporter,
	}
}

func (u *salesRecordUcase) ExportData(req *common.RequestContext, d *dto.ExportDataRequest) (
	res *bytes.Buffer, err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("ExportData"))
		}
	}()
	if err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.SALES_RECORD, domain.EDIT); err != nil {
		return
	}
	salesRecords, err := u.salesRecordRepo.ExportData(req, d)
	if err != nil {
		return
	}
	fmt.Println("SALES RECORDS", d.Body.Data)
	data := make([][]interface{}, len(salesRecords))
	for idx, d := range salesRecords {
		//COmplete
		data[idx] = []interface{}{
			d.InvoiceDate,
			d.InvoiceNo,
			d.AuthorizationCode,
			d.CustomerNitCi,
			d.Supplement,
			d.NameOrBusinessName,
			u.currency.IntToFloat(int64(d.TotalSaleAmount)),
			u.currency.IntToFloat(int64(d.IceAmount)),
			u.currency.IntToFloat(int64(d.IehdAmount)),
			u.currency.IntToFloat(int64(d.IpjAmount)),
			u.currency.IntToFloat(int64(d.TaxRates)),
			u.currency.IntToFloat(int64(d.OtherNotSubjectToVat)),
			u.currency.IntToFloat(int64(d.ExportsAndExemptOperations)),
			u.currency.IntToFloat(int64(d.ZeroRateTaxableSales)),
			u.currency.IntToFloat(int64(d.Subtotal)),
			u.currency.IntToFloat(int64(d.DiscountsBonusAndRebatesSubjectToVat)),
			u.currency.IntToFloat(int64(d.GiftCardAmount)),
			u.currency.IntToFloat(int64(d.BaseAmountForTaxDebit)),
			u.currency.IntToFloat(int64(d.TaxDebit)),
			d.State,
			d.ControlCode,
			d.SaleType,
			d.WithTaxCreditRight,
			d.ConsolidationStatus,
		}
	}
	headers := []interface{}{"FECHA DE LA FACTURA", "Nº DE LA FACTURA", "CODIGO DE AUTORIZACIÓN", "NIT / CI CLIENTE",
		"COMPLEMENTO", "NOMBRE O RAZON SOCIAL", "IMPORTE TOTAL DE LA VENTA", "IMPORTE ICE", "IMPORTE IEHD", "IMPORTE IPJ",
		"TASAS", "OTROS NO SUJETOS AL IVA", "EXPORTACIONES Y OPERACIONES EXENTAS", "VENTAS GRAVADAS A TASA CERO",
		"SUBTOTAL", "DESCUENTOS BONIFICACIONES Y REBAJAS SUJETAS AL IVA", "IMPORTE GIFT CARD", "IMPORTE BASE PARA DEBITO FISCAL",
		"DEBITO FISCAL", "ESTADO", "CODIGO DE CONTROL", "TIPO DE VENTA", "CON DERECHO A CREDITO FISCAL", "ESTADO CONSOLIDACION"}
	res, err = u.excelExporter.Export("Sheet", headers, data)
	return
}

func (u *salesRecordUcase) UpdateStatus(req *common.RequestContext, d *dto.UpdateStatusWithEvent) (err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("UpdateStatus"))
		}
	}()
	if err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.SALES_RECORD, domain.EDIT); err != nil {
		return
	}
	nextState, err := u.fsm.NextState(d.Body.CurrentState, d.Body.Events)
	if err != nil {
		return err
	}
	err = u.salesRecordRepo.UpdateStatus(req, d, nextState)
	return
}

func (u *salesRecordUcase) EditSalesRecord(req *common.RequestContext, d *dto.EditSalesRecordRequest) (err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("EditSalesRecord"))
		}
	}()
	err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.SALES_RECORD, domain.EDIT)
	if err != nil {
		return
	}
	err = u.salesRecordRepo.EditSalesRecord(req, d)
	return
}

func (u *salesRecordUcase) GetSalesRecord(req *common.RequestContext, d *dto.RequestEntity) (
	res dto.ResultEntity[dto.SalesRecordDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetSalesRecord"))
		}
	}()
	err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.SALES_RECORD, domain.VIEW)
	if err != nil {
		return
	}
	res, err = u.salesRecordRepo.GetSalesRecord(req, d)
	if err != nil {
		return
	}
	res.Activities = u.core.GerActivitiesByPartyID(req, res.Entity.ID)
	return
}
func (u *salesRecordUcase) CreateSalesRecord(req *common.RequestContext, d *dto.CreateSalesRecordRequest) (
	res dto.SalesRecordDto, err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("CreateSalesRecord"))
		}
	}()
	err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.SALES_RECORD, domain.CREATE)
	if err != nil {
		return
	}
	salesRecord, err := u.salesRecordRepo.CreateSalesRecord(req, d)
	if err != nil {
		return
	}
	res = dto.SalesRecordDtoFromModel(&salesRecord)
	return
}
func (u *salesRecordUcase) GetSalesRecords(req *common.RequestContext, d *dto.SalesRecordsRequest) (
	res dto.PaginationResult[[]dto.SalesRecordDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetSalesRecords"))
		}
	}()
	err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.SALES_RECORD, domain.VIEW)
	if err != nil {
		return res, domain.ACTION_NOT_ALLOWED
	}

	res, err = u.salesRecordRepo.GetSalesRecords(req, d)
	if err != nil {
		return
	}
	res.FilterOptions = u.salesRecordRepo.GetFilterOptions(req)
	return
}
