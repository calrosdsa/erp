package sales_record_repo

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/gen/proto"
	"erp/internal/app/service/helpers"
	"erp/pkg/db"
	"fmt"
	"strings"

	"gorm.io/gen/field"
	"gorm.io/gen/helper"
)

type SalesRecordRepo interface {
	GetSalesRecord(req *common.RequestContext, d *dto.RequestEntity) (
		res dto.ResultEntity[dto.SalesRecordDto], err error)
	CreateSalesRecord(req *common.RequestContext, d *dto.CreateSalesRecordRequest) (
		res model.SalesRecord, err error)
	GetSalesRecords(req *common.RequestContext, d *dto.SalesRecordsRequest) (
		res dto.PaginationResult[[]dto.SalesRecordDto], err error)
	EditSalesRecord(req *common.RequestContext, d *dto.EditSalesRecordRequest) (err error)
	UpdateStatus(req *common.RequestContext, d *dto.UpdateStatusWithEvent, nextState string) (err error)
	ExportData(req *common.RequestContext, d *dto.ExportDataRequest) (res []dto.SalesRecordDto, err error)
	GetFilterOptions(req *common.RequestContext) []dto.FilterOptionDto
}

type salesRecordRepo struct {
	Q         *query.Query
	convertor helpers.ConvertorHelper
	dbHelper  db.DbHelper
	currency  helpers.CurrencyHelper
	query     helpers.QueryHelper
}

func NewSaleRecordRepo(
	db db.Connection,
	helpers *helpers.Helpers,
) SalesRecordRepo {
	return &salesRecordRepo{
		Q:         db.GetQ(),
		convertor: helpers.Convertor,
		dbHelper:  db.GetDbHelper(),
		currency:  helpers.Currency,
		query:     helpers.Query,
	}
}

func (r *salesRecordRepo) UpdateStatus(req *common.RequestContext, d *dto.UpdateStatusWithEvent,
	nextState string) (err error) {
	salesRecordQ := r.Q.SalesRecord
	_, err = r.Q.SalesRecord.WithContext(req.Ctx).Where(
		salesRecordQ.CompanyID.Eq(req.ActiveCompany.ID),
		salesRecordQ.Status.Eq(d.Body.CurrentState),
		salesRecordQ.UUID.Eq(d.Body.PartyID),
	).UpdateSimple(
		salesRecordQ.Status.Value(nextState),
	)
	return
}

func (r *salesRecordRepo) GetSalesRecord(req *common.RequestContext, d *dto.RequestEntity) (
	res dto.ResultEntity[dto.SalesRecordDto], err error) {
	id := r.convertor.StrtoInt(d.ID)
	salesRecordQ := r.Q.SalesRecord
	customerQ := r.Q.Customer
	invoiceQ := r.Q.Invoice
	err = salesRecordQ.WithContext(req.Ctx).
		Select(salesRecordQ.ALL, salesRecordQ.CustomerID, customerQ.Name.As("customer"), customerQ.UUID.As("customer_uuid"),
			invoiceQ.ID.As("invoice_id"), invoiceQ.Code.As("invoice_code"),
		).
		Where(salesRecordQ.CompanyID.Eq(req.ActiveCompany.ID), salesRecordQ.ID.Eq(id)).
		Join(customerQ, customerQ.ID.EqCol(salesRecordQ.CustomerID)).
		LeftJoin(invoiceQ, invoiceQ.ID.EqCol(salesRecordQ.InvoiceID)).
		Scan(&res.Entity)
	return
}

func (r *salesRecordRepo) CreateSalesRecord(req *common.RequestContext, d *dto.CreateSalesRecordRequest) (
	res model.SalesRecord, err error) {
	tx := r.Q.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	var references []*int64
	// Populate the fields of the SalesRecord struct
	res.InvoiceDate = d.Body.InvoiceDate
	res.InvoiceNo = d.Body.InvoiceNo
	res.AuthorizationCode = d.Body.AuthorizationCode
	res.CustomerNitCi = d.Body.CustomerNitCi
	res.Supplement = d.Body.Supplement
	res.NameOrBusinessName = d.Body.NameOrBusinessName

	// Convert floating-point fields to int32
	res.TotalSaleAmount = r.currency.FloatToInt(d.Body.TotalSaleAmount)
	res.IceAmount = int32(r.currency.FloatToInt(d.Body.IceAmount))
	res.IehdAmount = int32(r.currency.FloatToInt(d.Body.IehdAmount))
	res.IpjAmount = int32(r.currency.FloatToInt(d.Body.IpjAmount))
	res.TaxRates = int32(r.currency.FloatToInt(d.Body.TaxRates))
	res.OtherNotSubjectToVat = int32(r.currency.FloatToInt(d.Body.OtherNotSubjectToVat))
	res.ExportsAndExemptOperations = int32(r.currency.FloatToInt(d.Body.ExportsAndExemptOperations))
	res.ZeroRateTaxableSales = int32(r.currency.FloatToInt(d.Body.ZeroRateTaxableSales))
	res.Subtotal = r.currency.FloatToInt(d.Body.Subtotal)
	res.DiscountsBonusAndRebatesSubjectToVat = int32(r.currency.FloatToInt(d.Body.DiscountsBonusAndRebatesSubjectToVat))
	res.GiftCardAmount = int32(r.currency.FloatToInt(d.Body.GiftCardAmount))
	res.BaseAmountForTaxDebit = int32(r.currency.FloatToInt(d.Body.BaseAmountForTaxDebit))
	res.TaxDebit = int32(r.currency.FloatToInt(d.Body.TaxDebit))

	res.State = d.Body.State
	res.ControlCode = d.Body.ControlCode
	res.SaleType = d.Body.SaleType
	res.WithTaxCreditRight = d.Body.WithTaxCreditRight
	res.ConsolidationStatus = d.Body.ConsolidationStatus
	res.CustomerID = d.Body.CustomerID
	res.CompanyID = req.ActiveCompany.ID
	if d.Body.InvoiceID != 0 {
		references = append(references, &d.Body.InvoiceID)
		res.InvoiceID = &d.Body.InvoiceID
	}
	// Insert into the database and get the newly generated ID
	currencyExchangeID, err := tx.SalesRecord.InsertParty(proto.PartyType_salesRecord.String())
	if err != nil {
		return
	}
	res.ID = currencyExchangeID
	res.Status = proto.State_DRAFT.String()

	// Save the record in the database
	err = tx.SalesRecord.WithContext(req.Ctx).Save(&res)
	if err != nil {
		return
	}
	err = tx.SalesRecord.InsertActivity(res.ID, req.Profile.ID, proto.ActivityType_CREATE.String(), nil)
	if err != nil {
		return
	}
	references = append(references, &d.Body.CustomerID)
	err = r.dbHelper.InsertReferences(req.Ctx, tx, res.ID, references)
	if err != nil {
		return
	}
	// Commit the transaction
	err = tx.Commit()
	return
}
func (r *salesRecordRepo) EditSalesRecord(req *common.RequestContext, d *dto.EditSalesRecordRequest) (
	err error) {
	tx := r.Q.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	salesRecordQ := tx.SalesRecord
	var (
		columns    []field.AssignExpr
		references []*int64
	)
	err = r.dbHelper.DeleteReferences(req.Ctx, tx, d.Body.ID)
	if err != nil {
		return
	}
	// Append fields with type conversion for floating-point values
	columns = append(columns,
		salesRecordQ.InvoiceDate.Value(d.Body.InvoiceDate),
		salesRecordQ.InvoiceNo.Value(d.Body.InvoiceNo),
		salesRecordQ.AuthorizationCode.Value(d.Body.AuthorizationCode),
		salesRecordQ.CustomerNitCi.Value(d.Body.CustomerNitCi),
		salesRecordQ.Supplement.Value(d.Body.Supplement),
		salesRecordQ.NameOrBusinessName.Value(d.Body.NameOrBusinessName),
		salesRecordQ.TotalSaleAmount.Value(r.currency.FloatToInt(d.Body.TotalSaleAmount)),
		salesRecordQ.IceAmount.Value(int32(r.currency.FloatToInt(d.Body.IceAmount))),
		salesRecordQ.IehdAmount.Value(int32(r.currency.FloatToInt(d.Body.IehdAmount))),
		salesRecordQ.IpjAmount.Value(int32(r.currency.FloatToInt(d.Body.IpjAmount))),
		salesRecordQ.TaxRates.Value(int32(r.currency.FloatToInt(d.Body.TaxRates))),
		salesRecordQ.OtherNotSubjectToVat.Value(int32(r.currency.FloatToInt(d.Body.OtherNotSubjectToVat))),
		salesRecordQ.ExportsAndExemptOperations.Value(int32(r.currency.FloatToInt(d.Body.ExportsAndExemptOperations))),
		salesRecordQ.ZeroRateTaxableSales.Value(int32(r.currency.FloatToInt(d.Body.ZeroRateTaxableSales))),
		salesRecordQ.Subtotal.Value(r.currency.FloatToInt(d.Body.Subtotal)),
		salesRecordQ.DiscountsBonusAndRebatesSubjectToVat.Value(int32(r.currency.FloatToInt(d.Body.DiscountsBonusAndRebatesSubjectToVat))),
		salesRecordQ.GiftCardAmount.Value(int32(r.currency.FloatToInt(d.Body.GiftCardAmount))),
		salesRecordQ.BaseAmountForTaxDebit.Value(int32(r.currency.FloatToInt(d.Body.BaseAmountForTaxDebit))),
		salesRecordQ.TaxDebit.Value(int32(r.currency.FloatToInt(d.Body.TaxDebit))),
		salesRecordQ.State.Value(d.Body.State),
		salesRecordQ.ControlCode.Value(d.Body.ControlCode),
		salesRecordQ.SaleType.Value(d.Body.SaleType),
		salesRecordQ.WithTaxCreditRight.Value(d.Body.WithTaxCreditRight),
		salesRecordQ.ConsolidationStatus.Value(d.Body.ConsolidationStatus),
		salesRecordQ.CustomerID.Value(d.Body.CustomerID))

	if d.Body.InvoiceID != 0 {
		columns = append(columns, salesRecordQ.InvoiceID.Value(d.Body.InvoiceID))
		references = append(references, &d.Body.InvoiceID)
	} else {
		columns = append(columns, salesRecordQ.InvoiceID.Null())
	}
	// Execute the update query
	_, err = tx.SalesRecord.WithContext(req.Ctx).Where(
		salesRecordQ.ID.Eq(d.Body.ID), salesRecordQ.CompanyID.Eq(req.ActiveCompany.ID),
	).UpdateSimple(columns...)
	if err != nil {
		return
	}

	// Insert activity log
	err = tx.SalesRecord.InsertActivity(d.Body.ID, req.Profile.ID, proto.ActivityType_EDIT.String(), nil)
	if err != nil {
		return
	}
	references = append(references, &d.Body.CustomerID)
	err = r.dbHelper.InsertReferences(req.Ctx, tx, d.Body.ID, references)
	if err != nil {
		return
	}
	// Commit the transaction
	err = tx.Commit()
	return
}

func (r *salesRecordRepo) ExportData(req *common.RequestContext, d *dto.ExportDataRequest) (
	res []dto.SalesRecordDto, err error) {
	var (
		generateSQL strings.Builder
	)
	builder := r.Q.WithContext(req.Ctx).SalesRecord
	if fromDate, ok := d.Body.Data["from_date"]; ok {
		fmt.Println("FromDate", fromDate)
	}
	params := r.builder(req, d.Body.Data, &generateSQL)
	fmt.Println("QUERY", generateSQL.String(), params)
	err = builder.UnderlyingDB().Raw(generateSQL.String(), params...).Scan(&res).Error
	return
}

func (r *salesRecordRepo) GetSalesRecords(req *common.RequestContext, d *dto.SalesRecordsRequest) (
	res dto.PaginationResult[[]dto.SalesRecordDto], err error) {
	var (
		generateSQL strings.Builder
	)
	builder := r.Q.WithContext(req.Ctx).SalesRecord
	// queryData := r.convertor.GenerateQueryMap(d.PaginationParams)
	queryData := r.convertor.GenerateQueryMap(d)
	params := r.builder(req, queryData, &generateSQL)

	err = builder.UnderlyingDB().Raw(generateSQL.String(), params...).Scan(&res.Results).Error
	return
}
func (r *salesRecordRepo) builder(req *common.RequestContext, d map[string]string, generateSQL *strings.Builder) (
	params []interface{}) {
	var (
		whereSQL strings.Builder
	)
	generateSQL.WriteString(`
		SELECT e.*, c.name AS customer 
			FROM sales_records AS e
			JOIN customers AS c 
			ON c.id = e.customer_id 
	`)
	whereSQL.WriteString(` e.deleted_at is null and e.company_id = ?`)
	params = append(params, req.ActiveCompany.ID)
	if invoiceDate, ok := d["invoice_date"]; ok {
		whereSQL.WriteString(r.convertor.GetConditionFromQuery(invoiceDate, "e.invoice_date", &params))
	}
	if invoiceID, ok := d["invoice_id"]; ok {
		whereSQL.WriteString(r.convertor.GetConditionFromQuery(invoiceID, "e.invoice_id", &params))
	}
	if status, ok := d["status"]; ok {
		whereSQL.WriteString(r.convertor.GetConditionFromQuery(status, "e.status", &params))
	}
	if customerID, ok := d["customer_id"]; ok {
		whereSQL.WriteString(r.convertor.GetConditionFromQuery(customerID, "e.customer_id", &params))
	}
	helper.JoinWhereBuilder(generateSQL, whereSQL)
	r.query.OrderAndLimitBuilder(generateSQL, d)

	return
}

func (r *salesRecordRepo) GetFilterOptions(req *common.RequestContext) []dto.FilterOptionDto {
	filterOptions := []dto.FilterOptionDto{}
	status := dto.FilterOptionDto{
		Param:     "status",
		Name:      "Estado",
		Type:      "string",
		Operators: []string{"=", "!=", "in"},
		Options:   []string{proto.State_DRAFT.String(), proto.State_SUBMITTED.String()},
	}
	invoiceDate := dto.FilterOptionDto{
		Param:     "invoice_date",
		Type:      "date",
		Name:      "Fecha de Factura",
		Operators: []string{"=", "!=", ">", "<", ">=", "<=", "between", "in"},
	}
	filterOptions = append(filterOptions, status, invoiceDate)
	return filterOptions
}
