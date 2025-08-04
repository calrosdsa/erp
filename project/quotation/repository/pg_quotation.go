package quotation_repo

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

type QuotationRepository interface {
	GetQuotation(req *common.RequestContext, d *dto.RequestEntityWithParty) (
		res dto.ResultEntity[dto.QuotationDetailDto], err error)
	CreateQuotation(tx *query.QueryTx, req *common.RequestContext, d dto.QuotationBody) (
		res model.Quotation, err error)
	GetQuotations(req *common.RequestContext, d *dto.RequestQuotations) (
		res dto.PaginationResult[[]dto.QuotationDto], err error)
	UpdateStatus(tx *query.QueryTx, req *common.RequestContext,
		id, prevState, nextState string) (res model.Quotation, err error)
	EditQuotation(tx *query.QueryTx, req *common.RequestContext, d dto.QuotationBody) (err error)
	GetFilterOptions() []dto.FilterOptionDto
}

const QUOTATION_SUPPLIER_CODE_TEMPLATE = "CP-#######"
const QUOTATION_CODE_TEMPLATE = "COT-#######"

type quotationRepository struct {
	Q         *query.Query
	convertor helpers.ConvertorHelper
	query     helpers.QueryHelper
	generator helpers.Generator
	dbHelper  db.DbHelper
}

func NewQuotationRepository(
	db db.Connection,
	helpers *helpers.Helpers,
) QuotationRepository {
	return &quotationRepository{
		Q:         db.GetQ(),
		convertor: helpers.Convertor,
		dbHelper:  db.GetDbHelper(),
		query:     helpers.Query,
		generator: helpers.Generator,
	}
}

func (r *quotationRepository) UpdateStatus(tx *query.QueryTx, req *common.RequestContext,
	id, prevState, nextState string) (res model.Quotation, err error) {
	quotationQ := tx.Quotation
	_, err = tx.Quotation.WithContext(req.Ctx).Where(
		quotationQ.CompanyID.Eq(req.ActiveCompany.ID),
		quotationQ.Status.Eq(prevState),
		quotationQ.Code.Eq(id),
	).UpdateSimple(quotationQ.Status.Value(nextState))
	if err != nil {
		return
	}
	stockEntry, err := tx.Quotation.WithContext(req.Ctx).Where(
		quotationQ.CompanyID.Eq(req.ActiveCompany.ID),
	).First()
	if err != nil {
		return
	}

	return *stockEntry, err
}

func (r *quotationRepository) GetQuotation(req *common.RequestContext, d *dto.RequestEntityWithParty) (
	res dto.ResultEntity[dto.QuotationDetailDto], err error) {
	quotationQ := r.Q.Quotation
	partyQ := r.Q.Party
	supplierQ := r.Q.Supplier
	customerQ := r.Q.Customer
	projectQ := r.Q.Project
	costCenterQ := r.Q.CostCenter
	priceListQ := r.Q.PriceList
	err = quotationQ.WithContext(req.Ctx).Select(
		quotationQ.ID, quotationQ.Code, quotationQ.Status, quotationQ.PostingTime,
		quotationQ.PostingDate, quotationQ.Currency, quotationQ.Tz, quotationQ.ValidTill,
		quotationQ.PartyID, partyQ.PartyTypeCode.As("party_type"),
		supplierQ.Name.As("party_name"), supplierQ.UUID.As("party_uuid"),
		customerQ.Name.As("party_name"), customerQ.UUID.As("party_uuid"),
		projectQ.Name.As("project"), projectQ.ID.As("project_id"), projectQ.UUID.As("project_uuid"),
		costCenterQ.Name.As("cost_center"), costCenterQ.ID.As("cost_center_id"), costCenterQ.UUID.As("cost_center_uuid"),
		priceListQ.Name.As("price_list"), priceListQ.ID.As("price_list_id"), priceListQ.UUID.As("price_list_uuid"),
	).
		Join(partyQ, partyQ.ID.EqCol(quotationQ.PartyID)).
		LeftJoin(supplierQ, partyQ.PartyTypeCode.Eq(proto.PartyType_supplier.String()), supplierQ.ID.EqCol(partyQ.ID)).
		LeftJoin(customerQ, partyQ.PartyTypeCode.Eq(proto.PartyType_customer.String()), customerQ.ID.EqCol(partyQ.ID)).
		LeftJoin(projectQ, projectQ.ID.EqCol(quotationQ.ProjectID)).
		LeftJoin(costCenterQ, costCenterQ.ID.EqCol(quotationQ.CostCenterID)).
		LeftJoin(priceListQ,priceListQ.ID.EqCol(quotationQ.PriceListID)).
		Where(quotationQ.CompanyID.Eq(req.ActiveCompany.ID), quotationQ.Code.Eq(d.ID)).
		Scan(&res.Entity.Quotation)
	return
}

func (r *quotationRepository) CreateQuotation(tx *query.QueryTx, req *common.RequestContext, d dto.QuotationBody) (
	res model.Quotation, err error) {
	res.CompanyID = req.ActiveCompany.ID
	quotationPartyID, err := tx.Quotation.InsertParty(d.Quotation.QuotationPartyType)
	if err != nil {
		return
	}
	res.ID = quotationPartyID
	fields := d.Quotation.Fields
	res.Status = proto.State_DRAFT.String()
	
	if err  = r.convertor.CopyStructData(fields,&res);err != nil {
		return 
	}
	fmt.Println("QUOTATION PARTY",d.Quotation.QuotationPartyType)
	count, err := tx.Quotation.WithContext(req.Ctx).Where(
		tx.Quotation.CompanyID.Eq(req.ActiveCompany.ID),
		tx.Party.PartyTypeCode.Eq(d.Quotation.QuotationPartyType),
	).Join(
		tx.Party, tx.Party.ID.EqCol(tx.Quotation.ID),
	).Count()
	if err != nil {
		return
	}
	var quotationTemplate string 
	if proto.PartyType_supplierQuotation.String() == d.Quotation.QuotationPartyType {
		quotationTemplate = QUOTATION_SUPPLIER_CODE_TEMPLATE
	}
	if proto.PartyType_salesQuotation.String() == d.Quotation.QuotationPartyType {
		quotationTemplate = QUOTATION_CODE_TEMPLATE
	}
	res.Code, err = r.generator.GenerateCodeAutoIncrement(quotationTemplate, count)
	if err != nil {
		return
	}
	
	err = tx.Quotation.WithContext(req.Ctx).Save(&res)
	if err != nil {
		return
	}
	err = tx.Quotation.InsertActivity(res.ID, req.Profile.ID, proto.ActivityType_CREATE.String(), nil)
	if err != nil {
		return
	}
	//Create Reference to party
	references := d.Quotation.References
	references = append(references, &fields.PartyID,
		fields.CostCenterID, fields.ProjectID)

	err = r.dbHelper.InsertReferences(req.Ctx, tx, res.ID, references)
	if err != nil {
		return
	}
	return
}

func (r *quotationRepository) EditQuotation(tx *query.QueryTx, req *common.RequestContext, d dto.QuotationBody) (err error) {
	// if err = r.dbHelper.DeleteReferences(req.Ctx, tx, d.Quotation.ID); err != nil {
	// 	return
	// }
	data, err := r.convertor.DataMap(d.Quotation.Fields)
	if err != nil {
		return
	}
	err = tx.Quotation.UnderlyingDB().WithContext(req.Ctx).Model(
		&model.Pricing{ID: d.Quotation.ID},
	).Updates(data).Error
	if err != nil {
		return
	}

	err = tx.Quotation.InsertActivity(d.Quotation.ID, req.Profile.ID, proto.ActivityType_EDIT.String(), nil)
	if err != nil {
		return
	}
	var references []*int64
	references = append(references,
		&d.Quotation.Fields.PartyID,
		d.Quotation.Fields.CostCenterID,
		d.Quotation.Fields.ProjectID,
	)
	references = append(references, d.Quotation.References...)
	if err = r.dbHelper.InsertReferences(req.Ctx, tx, d.Quotation.ID, references,true); err != nil {
		return
	}

	return
}

func (r *quotationRepository) GetQuotations(req *common.RequestContext, d *dto.RequestQuotations) (
	res dto.PaginationResult[[]dto.QuotationDto], err error) {
	var (
		generateSQL strings.Builder
	)
	builder := r.Q.WithContext(req.Ctx).Invoice
	queryData := r.convertor.GenerateQueryMap(d)
	params := r.quotationsQuery(req, queryData, &generateSQL, d.PartyType)
	err = builder.UnderlyingDB().Raw(generateSQL.String(), params...).Scan(&res.Results).Error
	return
}

func (r *quotationRepository) quotationsQuery(req *common.RequestContext, d map[string]string, generateSQL *strings.Builder,
	docPartyType string) (
	params []interface{}) {
	var (
		whereSQL strings.Builder
	)
	generateSQL.WriteString(`SELECT 
		e.id,e.code,e.created_at,e.status,e.posting_date,e.posting_time,e.tz,e.valid_till,
		e.currency,
		p.party_type_code as party_type,party.name as party_name,party.uuid as party_uuid
		from quotations as e 
		join parties as doc_party on doc_party.id = e.id
		join parties as p on p.id = e.party_id
	`)
	if docPartyType == proto.PartyType_supplierQuotation.String() {
		generateSQL.WriteString(`join suppliers as party on party.id = e.party_id `)
	}
	if docPartyType == proto.PartyType_salesQuotation.String() {
		generateSQL.WriteString(`join customers as party on party.id = e.party_id `)
	}
	whereSQL.WriteString(` e.deleted_at is null and e.company_id = ? and doc_party.party_type_code = ?`)
	params = append(params, req.ActiveCompany.ID, docPartyType)
	// if id, ok := d["id"]; ok {
	// 	whereSQL.WriteString(r.convertor.GetConditionFromQuery(id, "e.id", &params))
	// }
	// if status, ok := d["status"]; ok {
	// 	whereSQL.WriteString(r.convertor.GetConditionFromQuery(status, "e.status", &params))
	// }
	// if code, ok := d["code"]; ok {
	// 	whereSQL.WriteString(r.convertor.GetConditionFromQuery(code, "e.code", &params))
	// }
	// if party, ok := d["party_id"]; ok {
	// 	whereSQL.WriteString(r.convertor.GetConditionFromQuery(party, "e.party_id", &params))
	// }
	// if party, ok := d["posting_date"]; ok {
	// 	whereSQL.WriteString(r.convertor.GetConditionFromQuery(party, "e.posting_date", &params))
	// }
	// if party, ok := d["valid_till"]; ok {
	// 	whereSQL.WriteString(r.convertor.GetConditionFromQuery(party, "e.valid_till", &params))
	// }
	// if pricing, ok := d["pricing"]; ok {
	// 	whereSQL.WriteString(r.convertor.GetConditionFromQuery(pricing, "pr.party_id", &params))
	// 	generateSQL.WriteString(`join party_references as pr on pr.reference_id = e.id `)
	// }

	columnFilters := []string{
		"id",
		"project_id",
		"cost_center_id",
		"status",
		"code",
		"party_id",
		"posting_date",
		"valid_till",
	}
	r.query.FilterBuilder(&whereSQL, &params, d, columnFilters...)

	r.query.ReferenceFilterBuilder(generateSQL, &whereSQL, &params, d, "pricing_id")

	// if invoiceID, ok := d["invoice_id"]; ok {
	// 	whereSQL.WriteString(r.convertor.GetConditionFromQuery(invoiceID, "invoice_id", &params))
	// }
	helper.JoinWhereBuilder(generateSQL, whereSQL)

	r.query.OrderAndLimitBuilder(generateSQL, d)
	return
}

func (r *quotationRepository) GetFilterOptions() []dto.FilterOptionDto {
	filterOptions := []dto.FilterOptionDto{}
	status := dto.FilterOptionDto{
		Param:     "status",
		Name:      "Estado",
		Type:      dto.FILTER_TYPE_STRING,
		Operators: dto.StringOperators,
		Options: []string{proto.State_DRAFT.String(), proto.State_SUBMITTED.String(),
			proto.State_CANCELLED.String(),
		},
	}
	dueDate := dto.FilterOptionDto{
		Name:      "Fecha de Validez",
		Param:     "valid_till",
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

	pricing := dto.FilterOptionDto{
		Name:      "Pricing",
		Param:     "pricing_id",
		Type:      dto.FILTER_TYPE_STRING,
		Operators: dto.StringOperators,
		PartyType: proto.PartyType_pricing.String(),
	}

	filterOptions = append(filterOptions, status, dueDate, postingDate, code, pricing)
	return filterOptions
}
