package receipt_repo

import (
	"context"
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/gen/proto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/pkg/db"
	"erp/pkg/di"
	"strings"

	"gorm.io/gen/helper"
)

type ReceiptRepository interface {
	CreateReceipt(req *common.RequestContext, tx *query.QueryTx, i dto.ReceiptBody) (model.Receipt, error)
	GetReceipts(req *common.RequestContext, i *dto.RequestReceipts) (
		dto.PaginationResult[[]dto.ReceiptDto], error)
	GetReceiptDetail(req *common.RequestContext, i *dto.RequestEntityWithParty) (
		dto.ResultEntity[dto.ReceiptDetailDto], error)
	EditReceipt(tx *query.QueryTx, req *common.RequestContext, d dto.ReceiptBody) (err error)
	UpdateReceiptState(ctx context.Context, req *common.RequestContext, id string, prevState, nexState string) (
		*model.Receipt, error)
	GetFilterOptions() []dto.FilterOptionDto
}

const RECEIPT_CODE_TEMPLATE = "RC-#######"
const DN_CODE_TEMPLATE = "NE-#######"


type receiptRepository struct {
	Q         *query.Query
	convertor helpers.ConvertorHelper
	dbHelper  db.DbHelper
	query     helpers.QueryHelper
	generator helpers.Generator
}

func NewReceiptRepository(
	conn db.Connection,
	helpers *helpers.Helpers,
) ReceiptRepository {
	return &receiptRepository{
		Q:         conn.GetQ(),
		convertor: helpers.Convertor,
		query:     helpers.Query,
		dbHelper:  conn.GetDbHelper(),
		// currency:  helpers.Currency,
		generator: helpers.Generator,
	}
}

func (r *receiptRepository) UpdateReceiptState(ctx context.Context, req *common.RequestContext, id string, prevState, nexState string) (
	res *model.Receipt, err error) {
	// tx :=
	tx, ok := di.Get(ctx, domain.DatabaseTransactionKey).(*query.QueryTx)
	if !ok {
		return res, domain.FAIL_TYPE_ASSERTION
	}

	receiptQ := r.Q.Receipt
	_, err = tx.Receipt.WithContext(req.Ctx).Where(
		receiptQ.CompanyID.Eq(req.ActiveCompany.ID),
		receiptQ.Status.Eq(prevState),
		receiptQ.Code.Eq(id),
	).UpdateSimple(receiptQ.Status.Value(nexState))
	if err != nil {
		return
	}
	res, err = tx.Receipt.WithContext(ctx).Where(
		receiptQ.CompanyID.Eq(req.ActiveCompany.ID),
		receiptQ.Code.Eq(id),
	).First()
	return
}

func (r *receiptRepository) GetReceipts(req *common.RequestContext, d *dto.RequestReceipts) (
	res dto.PaginationResult[[]dto.ReceiptDto], err error) {
	var (
		generateSQL strings.Builder
	)
	builder := r.Q.WithContext(req.Ctx).Receipt
	queryData := r.convertor.GenerateQueryMap(d)
	params := r.receiptQuery(req, queryData, &generateSQL, d.PartyType)
	err = builder.UnderlyingDB().Raw(generateSQL.String(), params...).Scan(&res.Results).Error
	return
}

func (r *receiptRepository) receiptQuery(req *common.RequestContext, d map[string]string, generateSQL *strings.Builder,
	docPartyType string) (
	params []interface{}) {
	var (
		whereSQL strings.Builder
	)
	generateSQL.WriteString(`SELECT 
		e.id,e.code,e.created_at,e.status,e.posting_date,e.posting_time,e.tz,e.currency,
		p.party_type_code as party_type,party.name as party_name,party.uuid as party_uuid
		from receipts as e 
		join parties as doc_party on doc_party.id = e.id
		join parties as p on p.id = e.party_id
	`)
	if docPartyType == proto.PartyType_purchaseReceipt.String() {
		generateSQL.WriteString(`join suppliers as party on party.id = e.party_id `)
	}
	if docPartyType == proto.PartyType_deliveryNote.String() {
		generateSQL.WriteString(`join customers as party on party.id = e.party_id `)
	}
	whereSQL.WriteString(` e.deleted_at is null and e.company_id = ? and doc_party.party_type_code = ?`)
	params = append(params, req.ActiveCompany.ID, docPartyType)
	if status, ok := d["status"]; ok {
		whereSQL.WriteString(r.convertor.GetConditionFromQuery(status, "e.status", &params))
	}
	if query, ok := d["query"]; ok {
		whereSQL.WriteString(r.convertor.GetConditionFromQuery(query, "e.code", &params))
	}
	if party, ok := d["party_id"]; ok {
		whereSQL.WriteString(r.convertor.GetConditionFromQuery(party, "e.party_id", &params))
	}
	if party, ok := d["posting_date"]; ok {
		whereSQL.WriteString(r.convertor.GetConditionFromQuery(party, "e.posting_date", &params))
	}

	// if invoiceID, ok := d["invoice_id"]; ok {
	// 	whereSQL.WriteString(r.convertor.GetConditionFromQuery(invoiceID, "invoice_id", &params))
	// }
	helper.JoinWhereBuilder(generateSQL, whereSQL)

	r.query.OrderAndLimitBuilder(generateSQL, d)
	return
}

func (r *receiptRepository) GetReceiptDetail(req *common.RequestContext, d *dto.RequestEntityWithParty) (
	res dto.ResultEntity[dto.ReceiptDetailDto], err error) {
	receiptQ := r.Q.Receipt
	partyQ := r.Q.Party
	supplierQ := r.Q.Supplier
	customerQ := r.Q.Customer
	projectQ := r.Q.Project
	costCenterQ := r.Q.CostCenter
	priceListQ := r.Q.PriceList
	warehouseQ := r.Q.WareHouse
	receiptQ.WithContext(req.Ctx).Select(
		receiptQ.ID, receiptQ.Code, receiptQ.CreatedAt, receiptQ.PostingDate,
		receiptQ.PostingTime, receiptQ.Tz, receiptQ.Status, receiptQ.Currency,
		receiptQ.DocReferenceID,
		receiptQ.PartyID, partyQ.PartyTypeCode.As("party_type"),
		supplierQ.Name.As("party_name"), supplierQ.UUID.As("party_uuid"),
		customerQ.Name.As("party_name"), customerQ.UUID.As("party_uuid"),
		projectQ.Name.As("project"), projectQ.ID.As("project_id"), projectQ.UUID.As("project_uuid"),
		costCenterQ.Name.As("cost_center"), costCenterQ.ID.As("cost_center_id"), costCenterQ.UUID.As("cost_center_uuid"),
		priceListQ.Name.As("price_list"), priceListQ.ID.As("price_list_id"), priceListQ.UUID.As("price_list_uuid"),
		warehouseQ.Name.As("warehouse"), warehouseQ.ID.As("warehouse_id"), warehouseQ.UUID.As("warehouse_uuid"),
	).
		Join(partyQ, partyQ.ID.EqCol(receiptQ.PartyID)).
		LeftJoin(supplierQ, partyQ.PartyTypeCode.Eq(proto.PartyType_supplier.String()), supplierQ.ID.EqCol(partyQ.ID)).
		LeftJoin(customerQ, partyQ.PartyTypeCode.Eq(proto.PartyType_customer.String()), customerQ.ID.EqCol(partyQ.ID)).
		LeftJoin(projectQ, projectQ.ID.EqCol(receiptQ.ProjectID)).
		LeftJoin(costCenterQ, costCenterQ.ID.EqCol(receiptQ.CostCenterID)).
		LeftJoin(warehouseQ, warehouseQ.ID.EqCol(receiptQ.WarehouseID)).
		LeftJoin(priceListQ,priceListQ.ID.EqCol(receiptQ.PriceListID)).
		Where(
			receiptQ.CompanyID.Eq(req.ActiveCompany.ID),
			receiptQ.Code.Eq(d.ID),
		).Scan(&res.Entity.Receipt)
	// fmt.Println(itemLines)
	return res, err
}

func (r *receiptRepository) CreateReceipt(req *common.RequestContext, tx *query.QueryTx, d dto.ReceiptBody) (
	receipt model.Receipt, err error) {
	count, err := tx.Receipt.WithContext(req.Ctx).Where(
		tx.Receipt.CompanyID.Eq(req.ActiveCompany.ID),
		tx.Party.PartyTypeCode.Eq(d.Receipt.ReceiptPartyType),
	).Join(
		tx.Party, tx.Party.ID.EqCol(tx.Receipt.ID),
	).Count()
	if err != nil {
		return
	}
	var receiptTemplateCode string 
	if proto.PartyType_purchaseReceipt.String() == d.Receipt.ReceiptPartyType {
		receiptTemplateCode = RECEIPT_CODE_TEMPLATE
	}
	if proto.PartyType_deliveryNote.String() == d.Receipt.ReceiptPartyType {
		receiptTemplateCode = DN_CODE_TEMPLATE
	}
	receipt.Code, err = r.generator.GenerateCodeAutoIncrement(receiptTemplateCode, count)
	if err != nil {
		return
	}
	receipt.ID, err = tx.Receipt.InsertParty(d.Receipt.ReceiptPartyType)
	if err != nil {
		return
	}
	
	receipt.Status = proto.State_DRAFT.String()
	receipt.CompanyID = req.ActiveCompany.ID
	fields := d.Receipt.Fields
	if err = r.convertor.CopyStructData(fields, &receipt); err != nil {
		return
	}
	err = tx.Receipt.WithContext(req.Ctx).Save(&receipt)
	if err != nil {
		return
	}
	err = tx.Receipt.InsertActivity(receipt.ID, req.Profile.ID, proto.ActivityType_CREATE.String(), nil)
	if err != nil {
		return
	}
	var references []*int64
	references = append(references,
		&d.Receipt.Fields.PartyID,
		d.Receipt.Fields.DocReferenceID,
		d.Receipt.Fields.CostCenterID,
		d.Receipt.Fields.ProjectID,
	)
	if err = r.dbHelper.InsertReferences(req.Ctx, tx,receipt.ID, references); err != nil {
		return
	}
	return
}
func (r *receiptRepository) EditReceipt(tx *query.QueryTx, req *common.RequestContext, d dto.ReceiptBody) (err error) {
	if err = r.dbHelper.DeleteReferences(req.Ctx, tx, d.Receipt.ID); err != nil {
		return
	}
	data, err := r.convertor.DataMap(d.Receipt.Fields)
	if err != nil {
		return
	}
	err = tx.Receipt.UnderlyingDB().WithContext(req.Ctx).Model(
		&model.Pricing{ID: d.Receipt.ID},
	).Updates(data).Error
	if err != nil {
		return
	}

	err = tx.Receipt.InsertActivity(d.Receipt.ID, req.Profile.ID, proto.ActivityType_EDIT.String(), nil)
	if err != nil {
		return
	}
	var references []*int64
	references = append(references,
		&d.Receipt.Fields.PartyID,
		d.Receipt.Fields.DocReferenceID,
		d.Receipt.Fields.CostCenterID,
		d.Receipt.Fields.ProjectID,
	)
	if err = r.dbHelper.InsertReferences(req.Ctx, tx, d.Receipt.ID, references); err != nil {
		return
	}
	return
}

func (r *receiptRepository) GetFilterOptions() []dto.FilterOptionDto {

	filterOptions := []dto.FilterOptionDto{}

	status := dto.FilterOptionDto{
		Param:     "status",
		Name:      "Estado",
		Type:      dto.FILTER_TYPE_STRING,
		Operators: dto.StringOperators,
		Options: []string{proto.State_DRAFT.String(), proto.State_COMPLETED.String(),
			proto.State_CANCELLED.String(), proto.State_TO_BILL.String(), proto.State_PAID.String(),
		},
	}

	postingDate := dto.FilterOptionDto{
		Name:      "Fecha de Publicacion",
		Param:     "posting_date",
		Type:      dto.FILTER_TYPE_DATE,
		Operators: dto.DateOperators,
	}

	code := dto.FilterOptionDto{
		Name:      "ID",
		Param:     "code",
		Type:      dto.FILTER_TYPE_STRING,
		Operators: dto.StringOperators,
	}

	filterOptions = append(filterOptions, status, postingDate, code)
	return filterOptions
}
