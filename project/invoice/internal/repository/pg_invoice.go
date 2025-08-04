package invoice_repo

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/gen/proto"
	"strings"

	// "erp/internal/app/domain"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/pkg/db"

	"gorm.io/gen/helper"
	"gorm.io/gorm"
)

type InvoiceRepository interface {
	CreateInvoice(req *common.RequestContext, tx *query.QueryTx, i dto.InvoiceBody) (model.Invoice, error)
	GetInvoices(req *common.RequestContext, i *dto.RequestInvoices) (dto.PaginationResult[[]dto.InvoiceDto], error)
	GetInvoiceDetail(req *common.RequestContext, d *dto.RequestEntityWithParty) (
		dto.ResultEntity[dto.InvoiceDetailDto], error,
	)
	EditInvoice(tx *query.QueryTx, req *common.RequestContext, d dto.InvoiceBody) (err error)
	UpdateInvoiceState(tx *query.QueryTx, req *common.RequestContext, id string, prevState, nexState string) (
		*model.Invoice, error)
	GetFilterOptions(lng string, invoicePartyType string) []dto.FilterOptionDto
}

const PI_CODE_TEMPLATE = "FC-#######"
const SI_CODE_TEMPLATE = "FV-#######"



type invoiceRepository struct {
	Q           *query.Query
	dbHelper    db.DbHelper
	DB          *gorm.DB
	convertor   helpers.ConvertorHelper
	currency    helpers.CurrencyHelper
	errorHelper helpers.ErrorHelper
	query       helpers.QueryHelper
	generator   helpers.Generator
	locale      helpers.Locale
}

func NewInvoiceRepository(
	conn db.Connection,
	helpers *helpers.Helpers,
) InvoiceRepository {
	return &invoiceRepository{
		Q:           conn.GetQ(),
		DB:          conn.GetDB(),
		dbHelper:    conn.GetDbHelper(),
		convertor:   helpers.Convertor,
		currency:    helpers.Currency,
		errorHelper: helpers.Error,
		query:       helpers.Query,
		generator:   helpers.Generator,
		locale:      helpers.Locale,
	}
}

func (r *invoiceRepository) UpdateInvoiceState(tx *query.QueryTx, req *common.RequestContext, invoiceID string, prevState,
	nexState string) (res *model.Invoice, err error) {
	invoiceQ := tx.Invoice
	_, err = tx.Invoice.WithContext(req.Ctx).Where(
		invoiceQ.CompanyID.Eq(req.ActiveCompany.ID),
		invoiceQ.Status.Eq(prevState),
		invoiceQ.Code.Eq(invoiceID),
	).UpdateSimple(invoiceQ.Status.Value(nexState))
	if err != nil {
		return
	}
	res, err = tx.Invoice.WithContext(req.Ctx).Where(
		invoiceQ.Code.Eq(invoiceID),
		invoiceQ.CompanyID.Eq(req.ActiveCompany.ID),
	).First()
	if err != nil {
		return
	}

	return
}

func (r *invoiceRepository) GetInvoiceDetail(req *common.RequestContext, d *dto.RequestEntityWithParty) (
	dto.ResultEntity[dto.InvoiceDetailDto], error,
) {
	var (
		res dto.ResultEntity[dto.InvoiceDetailDto]
		err error
	)
	invoiceQ := r.Q.Invoice
	progressInvoiceQ := r.Q.ProgressInvoice
	partyQ := r.Q.Party
	supplierQ := r.Q.Supplier
	customerQ := r.Q.Customer
	projectQ := r.Q.Project
	costCenterQ := r.Q.CostCenter
	priceListQ := r.Q.PriceList
	warehouseQ := r.Q.WareHouse
	err = invoiceQ.WithContext(req.Ctx).Select(
		invoiceQ.ID, invoiceQ.PostingDate, invoiceQ.PostingTime, invoiceQ.Tz,
		invoiceQ.Code, invoiceQ.CreatedAt, invoiceQ.PartyID,
		invoiceQ.DueDate, invoiceQ.Status, invoiceQ.Currency,
		invoiceQ.UpdateStock,invoiceQ.DocReferenceID, partyQ.PartyTypeCode.As("party_type"),
		progressInvoiceQ.PaidAmount.As("paid"), progressInvoiceQ.TotalAmount.As("total"),
		supplierQ.Name.As("party_name"), supplierQ.UUID.As("party_uuid"),
		customerQ.Name.As("party_name"), customerQ.UUID.As("party_uuid"),
		projectQ.Name.As("project"), projectQ.ID.As("project_id"), projectQ.UUID.As("project_uuid"),
		costCenterQ.Name.As("cost_center"), costCenterQ.ID.As("cost_center_id"), costCenterQ.UUID.As("cost_center_uuid"),
		priceListQ.Name.As("price_list"), priceListQ.ID.As("price_list_id"), priceListQ.UUID.As("price_list_uuid"),
		warehouseQ.Name.As("warehouse"), warehouseQ.ID.As("warehouse_id"), warehouseQ.UUID.As("warehouse_uuid"),
	).
		Join(partyQ, partyQ.ID.EqCol(invoiceQ.PartyID)).
		LeftJoin(supplierQ, partyQ.PartyTypeCode.Eq(proto.PartyType_supplier.String()), supplierQ.ID.EqCol(partyQ.ID)).
		LeftJoin(customerQ, partyQ.PartyTypeCode.Eq(proto.PartyType_customer.String()), customerQ.ID.EqCol(partyQ.ID)).
		LeftJoin(projectQ, projectQ.ID.EqCol(invoiceQ.ProjectID)).
		LeftJoin(costCenterQ, costCenterQ.ID.EqCol(invoiceQ.CostCenterID)).
		LeftJoin(priceListQ, priceListQ.ID.EqCol(invoiceQ.PriceListID)).
		LeftJoin(warehouseQ, warehouseQ.ID.EqCol(invoiceQ.WarehouseID)).
		LeftJoin(progressInvoiceQ, progressInvoiceQ.InvoiceID.EqCol(invoiceQ.ID)).
		Where(
			invoiceQ.CompanyID.Eq(req.ActiveCompany.ID),
			invoiceQ.Code.Eq(d.ID),
		).Scan(&res.Entity.Invoice)
	if err != nil {
		return res, err
	}
	// totals, err := r.getTotals(req.Ctx, res.Entity.Invoice.ID)

	// res.Entity.Totals = totals
	return res, err
}
	
func (r *invoiceRepository) GetInvoices(req *common.RequestContext, d *dto.RequestInvoices) (
	res dto.PaginationResult[[]dto.InvoiceDto], err error) {
	var (
		generateSQL strings.Builder
	)
	builder := r.Q.WithContext(req.Ctx).Invoice
	queryData := r.convertor.GenerateQueryMap(d)
	params := r.invoicesQuery(req, queryData, &generateSQL, d.PartyType)
	err = builder.UnderlyingDB().Raw(generateSQL.String(), params...).Scan(&res.Results).Error
	return
}

func (r *invoiceRepository) invoicesQuery(req *common.RequestContext, d map[string]string, generateSQL *strings.Builder,
	docPartyType string) (
	params []interface{}) {
	var (
		whereSQL strings.Builder
	)
	generateSQL.WriteString(`SELECT 
		e.id,e.code,e.due_date,e.created_at,e.status,e.posting_date,e.posting_time,e.tz,
		e.currency,
		p.party_type_code as party_type,party.name as party_name,party.uuid as party_uuid,
		pi.total_amount,pi.paid_amount
		from invoices as e 
		join parties as invoice_party on invoice_party.id = e.id
		join parties as p on p.id = e.party_id
		join progress_invoices as pi on pi.invoice_id = e.id
	`)
	if docPartyType == proto.PartyType_purchaseInvoice.String() {
		generateSQL.WriteString(`join suppliers as party on party.id = e.party_id `)
	}
	if docPartyType == proto.PartyType_saleInvoice.String() {
		generateSQL.WriteString(`join customers as party on party.id = e.party_id `)
	}
	whereSQL.WriteString(` e.deleted_at is null and e.company_id = ? and invoice_party.party_type_code = ?`)
	params = append(params, req.ActiveCompany.ID, docPartyType)
	columnFilters := []string{
		"project_id",
		"cost_center_id",
		"status",
		"code",
		"party_id",
		"posting_date",
		"due_date",
	}
	r.query.FilterBuilder(&whereSQL, &params, d, columnFilters...)
	// if orderID,ok := d["order_id"];ok {
	// 	whereSQL.WriteString(r.convertor.GetConditionFromQuery(orderID,"i.doc_reference_id",&params))
	// }
	r.query.ReferenceFilterBuilder(generateSQL, &whereSQL, &params, d, "order_id")

	helper.JoinWhereBuilder(generateSQL, whereSQL)

	r.query.OrderAndLimitBuilder(generateSQL, d)
	return
}

func (r *invoiceRepository) CreateInvoice(req *common.RequestContext, tx *query.QueryTx, d dto.InvoiceBody) (
	invoice model.Invoice, err error) {
	var references []*int64
	invoicePartyID, err := tx.Invoice.InsertParty(d.Invoice.InvoicePartyType)
	if err != nil {
		return
	}
	invoice.ID = invoicePartyID
	invoice.CompanyID = req.ActiveCompany.ID
	invoice.Status = proto.State_DRAFT.String()
	count, err := tx.Invoice.WithContext(req.Ctx).Where(
		tx.Invoice.CompanyID.Eq(req.ActiveCompany.ID),
		tx.Party.PartyTypeCode.Eq(d.Invoice.InvoicePartyType),
	).Join(
		tx.Party, tx.Party.ID.EqCol(tx.Invoice.ID),
	).Count()
	if err != nil {
		return
	}
	var invoiceTemplateCode string
	if proto.PartyType_purchaseInvoice.String() == d.Invoice.InvoicePartyType {
		invoiceTemplateCode = PI_CODE_TEMPLATE
	}
	if proto.PartyType_saleInvoice.String() == d.Invoice.InvoicePartyType {
		invoiceTemplateCode = SI_CODE_TEMPLATE
	}
	invoice.Code, err = r.generator.GenerateCodeAutoIncrement(invoiceTemplateCode, count)
	if err != nil {
		return
	}
	fields := d.Invoice.Fields
	if err = r.convertor.CopyStructData(fields, &invoice); err != nil {
		return
	}
	err = tx.WithContext(req.Ctx).Invoice.Save(&invoice)
	if err != nil {
		return
	}

	totalAmount := r.calculateTotal(d.CreateItemLines,d.CreateTaxAndCharges)
	progressInvoice := model.ProgressInvoice{}
	progressInvoice.TotalAmount = r.currency.FloatToInt(totalAmount)
	progressInvoice.InvoiceID = invoice.ID
	err = tx.ProgressInvoice.WithContext(req.Ctx).Save(&progressInvoice)
	if err != nil {
		return
	}
	if err = tx.Invoice.InsertActivity(invoice.ID, req.Profile.ID, proto.ActivityType_CREATE.String(), nil); err != nil {
		return
	}

	references = append(references, &fields.PartyID, fields.DocReferenceID,
		fields.CostCenterID, fields.ProjectID)
	if err = r.dbHelper.InsertReferences(req.Ctx, tx, invoice.ID, references); err != nil {
		return
	}

	return
}

func (s *invoiceRepository) calculateTotal(
	d dto.CreateItemLines,c dto.CreateTaxAndChanges) (float64) {
	var totalAmount float64
	var totalAmountCharges float64
	for _, line := range d.Lines {
		totalAmount += line.Rate * float64(line.Quantity)
	}
	for _, line := range c.TaxAndCharges {
		if line.IsDeducted {
			totalAmount -= line.Amount
		}else {
			totalAmount += line.Amount

		}
	}
	return totalAmount+totalAmountCharges
}

func (r *invoiceRepository) EditInvoice(tx *query.QueryTx, req *common.RequestContext, d dto.InvoiceBody) (err error) {
	data, err := r.convertor.DataMap(d.Invoice.Fields)
	if err != nil {
		return
	}
	err = tx.Invoice.UnderlyingDB().WithContext(req.Ctx).Model(
		&model.Invoice{ID: d.Invoice.ID},
	).Updates(data).Error
	if err != nil {
		return
	}
	totalAmount := r.calculateTotal(d.CreateItemLines,d.CreateTaxAndCharges)
	_,err = tx.ProgressInvoice.WithContext(req.Ctx).Where(
		tx.ProgressInvoice.InvoiceID.Eq(d.Invoice.ID),
	).UpdateSimple(
		tx.ProgressInvoice.TotalAmount.Value(r.currency.FloatToInt(totalAmount)),
	)
	if err != nil {
		return
	}
	err = tx.Invoice.InsertActivity(d.Invoice.ID, req.Profile.ID, proto.ActivityType_EDIT.String(), nil)
	if err != nil {
		return
	}
	var references []*int64
	references = append(references,
		&d.Invoice.Fields.PartyID,
		d.Invoice.Fields.DocReferenceID,
		d.Invoice.Fields.CostCenterID,
		d.Invoice.Fields.ProjectID,
	)
	if err = r.dbHelper.InsertReferences(req.Ctx, tx, d.Invoice.ID, references,true); err != nil {
		return
	}
	return
}
func (r *invoiceRepository) GetFilterOptions(lng string, invoicePartyType string) []dto.FilterOptionDto {
	filterOptions := []dto.FilterOptionDto{}
	status := dto.FilterOptionDto{
		Param:     "status",
		Name:      "Estado",
		Type:      dto.FILTER_TYPE_STRING,
		Operators: dto.StringOperators,
		Options: []string{proto.State_DRAFT.String(), proto.State_PAID.String(),
			proto.State_UNPAID.String(), proto.State_CANCELLED.String()},
	}
	dueDate := dto.FilterOptionDto{
		Name:      "Fecha de Vencimiento",
		Param:     "due_date",
		Type:      dto.FILTER_TYPE_DATE,
		Operators: dto.DateOperators,
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

	order := dto.FilterOptionDto{
		Param:     "order_id",
		Type:      dto.FILTER_TYPE_STRING,
		Operators: dto.StringOperators,
	}

	t := r.locale.Translate(lng)
	if invoicePartyType == proto.PartyType_purchaseInvoice.String() {
		order.PartyType = proto.PartyType_purchaseOrder.String()
		order.Name = t("Entity." + domain.PURCHASE_ORDER.Name)
	}
	if invoicePartyType == proto.PartyType_saleInvoice.String() {
		order.PartyType = proto.PartyType_saleOrder.String()
		order.Name = t("Entity." + domain.SALE_ORDER.Name)
	}
	filterOptions = append(filterOptions, status, dueDate, postingDate, code, order)
	return filterOptions
}
