package purchase_record_repo

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

	"gorm.io/gen/helper"
)

type PurchaseRecordRepo interface {
	GetPurchaseRecord(req *common.RequestContext, d *dto.RequestEntity) (
		res dto.ResultEntity[dto.PurchaseRecordDto], err error)
	CreatePurchaseRecord(req *common.RequestContext, d *dto.CreatePurchaseRecordRequest) (
		res model.PurchaseRecord, err error)
	GetPurchaseRecords(req *common.RequestContext, d *dto.PurchaseRecordsRequest) (
		res dto.PaginationResult[[]dto.PurchaseRecordDto], err error)
	ExportData(req *common.RequestContext, d *dto.ExportDataRequest) (res []dto.PurchaseRecordDto, err error)
	EditPurchaseRecord(req *common.RequestContext, d *dto.EditPurchaseRecordRequest) (err error)
	UpdateStatus(req *common.RequestContext, d *dto.UpdateStatusWithEvent, nextState string) (err error)
	GetFilterOptions() []dto.FilterOptionDto
}

type purchaseRecordRepo struct {
	Q         *query.Query
	convertor helpers.ConvertorHelper
	dbHelper  db.DbHelper
	currency  helpers.CurrencyHelper
	query     helpers.QueryHelper
}

func NewPurchaseRecordRepo(
	db db.Connection,
	helpers *helpers.Helpers,
) PurchaseRecordRepo {
	return &purchaseRecordRepo{
		Q:         db.GetQ(),
		convertor: helpers.Convertor,
		query:     helpers.Query,
		dbHelper:  db.GetDbHelper(),
		currency:  helpers.Currency,
	}
}
func (r *purchaseRecordRepo) UpdateStatus(req *common.RequestContext, d *dto.UpdateStatusWithEvent,
	nextState string) (err error) {
	purchaseRecordQ := r.Q.PurchaseRecord
	_, err = r.Q.PurchaseRecord.WithContext(req.Ctx).Where(
		purchaseRecordQ.CompanyID.Eq(req.ActiveCompany.ID),
		purchaseRecordQ.Status.Eq(d.Body.CurrentState),
		purchaseRecordQ.UUID.Eq(d.Body.PartyID),
	).UpdateSimple(
		purchaseRecordQ.Status.Value(nextState),
	)
	return
}

func (r *purchaseRecordRepo) GetPurchaseRecord(req *common.RequestContext, d *dto.RequestEntity) (
	res dto.ResultEntity[dto.PurchaseRecordDto], err error) {
	id := r.convertor.StrtoInt(d.ID)
	purchaseRecordQ := r.Q.PurchaseRecord
	supplierQ := r.Q.Supplier
	invoiceQ := r.Q.Invoice
	err = purchaseRecordQ.WithContext(req.Ctx).
		Select(purchaseRecordQ.ALL, supplierQ.Name.As("supplier"), supplierQ.UUID.As("supplier_uuid"),
			invoiceQ.ID.As("invoice_id"), invoiceQ.Code.As("invoice_code")).
		Join(supplierQ, supplierQ.ID.EqCol(purchaseRecordQ.SupplierID)).
		LeftJoin(invoiceQ, invoiceQ.ID.EqCol(purchaseRecordQ.InvoiceID)).
		Where(purchaseRecordQ.CompanyID.Eq(req.ActiveCompany.ID), purchaseRecordQ.ID.Eq(id)).
		Scan(&res.Entity)
	return
}

func (r *purchaseRecordRepo) CreatePurchaseRecord(req *common.RequestContext, d *dto.CreatePurchaseRecordRequest) (
	res model.PurchaseRecord, err error) {
	tx := r.Q.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	var references []*int64
	// Insert into the database and get the newly generated ID
	purchaseRecordID, err := tx.PurchaseRecord.InsertParty(proto.PartyType_purchaseRecord.String())
	if err != nil {
		return
	}
	res.ID = purchaseRecordID
	res.Status = proto.State_DRAFT.String()
	res.CompanyID = req.ActiveCompany.ID
	fields := d.Body.Fields
	if err = r.convertor.CopyStructData(fields, &res); err != nil {
		return
	}
	// Save the record in the database
	err = tx.PurchaseRecord.WithContext(req.Ctx).Save(&res)
	if err != nil {
		return
	}
	err = tx.SalesRecord.InsertActivity(res.ID, req.Profile.ID, proto.ActivityType_CREATE.String(), nil)
	if err != nil {
		return
	}
	references = append(references, &fields.SupplierID)
	err = r.dbHelper.InsertReferences(req.Ctx, tx, res.ID, references)
	if err != nil {
		return
	}
	// Commit the transaction
	err = tx.Commit()
	return
}
func (r *purchaseRecordRepo) EditPurchaseRecord(req *common.RequestContext, d *dto.EditPurchaseRecordRequest) (
	err error) {
	tx := r.Q.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	var references []*int64
	fields := d.Body.Fields
	data, err := r.convertor.DataMap(fields)
	if err != nil {
		return
	}
	err = tx.PurchaseRecord.UnderlyingDB().WithContext(req.Ctx).Model(
		&model.Payment{ID: d.Body.ID},
	).Updates(data).Error
	if err != nil {
		return
	}

	err = tx.PurchaseRecord.InsertActivity(d.Body.ID, req.Profile.ID, proto.ActivityType_EDIT.String(), nil)
	if err != nil {
		return
	}
	references = append(references, &fields.SupplierID)
	err = r.dbHelper.InsertReferences(req.Ctx, tx, d.Body.ID, references, true)
	if err != nil {
		return
	}
	// Commit the transaction
	err = tx.Commit()
	return
}

func (r *purchaseRecordRepo) ExportData(req *common.RequestContext, d *dto.ExportDataRequest) (
	res []dto.PurchaseRecordDto, err error) {
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

func (r *purchaseRecordRepo) GetPurchaseRecords(req *common.RequestContext, d *dto.PurchaseRecordsRequest) (
	res dto.PaginationResult[[]dto.PurchaseRecordDto], err error) {
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
func (r *purchaseRecordRepo) builder(req *common.RequestContext, d map[string]string, generateSQL *strings.Builder) (
	params []interface{}) {
	var (
		whereSQL strings.Builder
	)
	generateSQL.WriteString(`SELECT * FROM purchase_records `)
	whereSQL.WriteString(` deleted_at is null and company_id = ?`)
	params = append(params, req.ActiveCompany.ID)
	if invoiceDate, ok := d["invoice_date"]; ok {
		whereSQL.WriteString(r.convertor.GetConditionFromQuery(invoiceDate, "invoice_date", &params))
	}
	if invoiceID, ok := d["invoice_id"]; ok {
		whereSQL.WriteString(r.convertor.GetConditionFromQuery(invoiceID, "invoice_id", &params))
	}
	if status, ok := d["status"]; ok {
		whereSQL.WriteString(r.convertor.GetConditionFromQuery(status, "status", &params))
	}
	if supplierID, ok := d["supplier_id"]; ok {
		whereSQL.WriteString(r.convertor.GetConditionFromQuery(supplierID, "supplier_id", &params))
	}
	helper.JoinWhereBuilder(generateSQL, whereSQL)
	r.query.OrderAndLimitBuilder(generateSQL, d)

	return
}
func (r *purchaseRecordRepo) GetFilterOptions() []dto.FilterOptionDto {
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
