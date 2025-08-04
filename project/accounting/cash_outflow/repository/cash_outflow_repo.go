package cash_outflow_repo

import (
	common "erp/api/common"
	"erp/internal/app/service/helpers"
	"erp/pkg/db"
	"strings"

	dto "erp/api/dto"

	model "erp/gen/db/model"

	"erp/gen/db/query"
	"erp/gen/proto"

	"gorm.io/gen/helper"
)

type CashOutflowRepository interface {
	Create(tx *query.QueryTx, req *common.RequestContext, d dto.CashOutflowData) (res model.CashOutflow, err error)
	UpdateStatus(tx *query.QueryTx, req *common.RequestContext, d dto.UpdateStatusWithEvent, nextState string) (
		res *model.CashOutflow, err error)
	Edit(tx *query.QueryTx, req *common.RequestContext, d dto.CashOutflowData) (err error)
	GetCashOutflow(req *common.RequestContext, d dto.RequestEntity) (
		res dto.CashOutflowDto, err error)
	GetCashOutflows(req *common.RequestContext, d dto.CashOutflowsRequest) (
		res []dto.CashOutflowDto, err error)
	GetFilterOptions() []dto.FilterOptionDto
}

type repository struct {
	convertor helpers.ConvertorHelper
	generator helpers.Generator
	Q         *query.Query
	query     helpers.QueryHelper
	dbHelper  db.DbHelper
}

const EC_CODE_TEMPLATE = "EC-#######"

func NewRepository(
	conn db.Connection,
	helpers *helpers.Helpers,
) CashOutflowRepository {
	return &repository{
		convertor: helpers.Convertor,
		Q:         conn.GetQ(),
		query:     helpers.Query,
		generator: helpers.Generator,
		dbHelper:  conn.GetDbHelper(),
	}
}

func (r *repository) GetCashOutflow(req *common.RequestContext, d dto.RequestEntity) (
	res dto.CashOutflowDto, err error) {
	e := r.Q.CashOutflow
	partyQ := r.Q.Party
	supplierQ := r.Q.Supplier
	projectQ := r.Q.Project
	costCenterQ := r.Q.CostCenter
	err = e.WithContext(req.Ctx).Select(
		e.ID, e.PostingDate, e.PostingTime, e.Tz,e.Status,
		e.Code, e.CreatedAt, e.PartyID,
		e.CashOutflowType, e.Concept, e.InvoiceNo, e.Nit, e.AuthCode,
		e.CtrlCode, e.EmisionDate, e.Amount,
		partyQ.PartyTypeCode.As("party_type"), e.PartyID, supplierQ.Name.As("party"), supplierQ.UUID.As("party_uuid"),
		projectQ.Name.As("project"), projectQ.ID.As("project_id"), projectQ.UUID.As("project_uuid"),
		costCenterQ.Name.As("cost_center"), costCenterQ.ID.As("cost_center_id"), costCenterQ.UUID.As("cost_center_uuid"),
	).
		Join(partyQ, partyQ.ID.EqCol(e.PartyID)).
		LeftJoin(supplierQ, partyQ.PartyTypeCode.Eq(proto.PartyType_supplier.String()), supplierQ.ID.EqCol(partyQ.ID)).
		LeftJoin(projectQ, projectQ.ID.EqCol(e.ProjectID)).
		LeftJoin(costCenterQ, costCenterQ.ID.EqCol(e.CostCenterID)).
		Where(
			e.CompanyID.Eq(req.ActiveCompany.ID),
			e.Code.Eq(d.ID),
		).Scan(&res)
	if err != nil {
		return res, err
	}

	return res, err
}

func (r *repository) GetCashOutflows(req *common.RequestContext, d dto.CashOutflowsRequest) (
	res []dto.CashOutflowDto, err error) {
	var (
		generateSQL strings.Builder
	)
	builder := r.Q.WithContext(req.Ctx).CashOutflow
	queryData := r.convertor.GenerateQueryMap(d)
	params := r.cashOutflowsQuery(req, queryData, &generateSQL)
	err = builder.UnderlyingDB().Raw(generateSQL.String(), params...).Scan(&res).Error
	return
}

func (r *repository) cashOutflowsQuery(req *common.RequestContext, d map[string]string, generateSQL *strings.Builder,
) (params []interface{}) {
	var (
		whereSQL strings.Builder
	)
	generateSQL.WriteString(`SELECT 
		e.id,e.code,e.created_at,e.status,e.posting_date,e.posting_time,e.tz,
		e.concept,e.cash_outflow_type,e.invoice_no,e.nit,e.auth_code,
		e.ctrl_code,e.emision_date,e.amount,
		e.party_type,e.party_id,
		coalesce(sp.name,'') as party,sp.uuid as party_uuid
		from cash_outflows as e 
		join suppliers as sp on sp.id = e.party_id 
	`)
	whereSQL.WriteString(` e.deleted_at is null and e.company_id = ? `)
	params = append(params, req.ActiveCompany.ID)
	columnFilters := []string{
		"project_id",
		"cost_center_id",
		"status",
		"code",
		"party_id",
		"posting_date",
	}
	r.query.FilterBuilder(&whereSQL, &params, d, columnFilters...)

	helper.JoinWhereBuilder(generateSQL, whereSQL)

	r.query.OrderAndLimitBuilder(generateSQL, d)
	return
}

func (m *repository) Create(tx *query.QueryTx, req *common.RequestContext, d dto.CashOutflowData) (res model.CashOutflow, err error) {
	id, err := tx.CashOutflow.InsertParty(proto.PartyType_cashOutflow.String())
	if err != nil {
		return
	}
	fields := d.Fields
	res.ID = id
	res.CompanyID = req.ActiveCompany.ID
	res.Status = proto.State_DRAFT.String()
	if err = m.convertor.CopyStructData(fields, &res); err != nil {
		return
	}
	count, err := tx.CashOutflow.WithContext(req.Ctx).Where(
		tx.CashOutflow.CompanyID.Eq(req.ActiveCompany.ID),
	).Count()
	if err != nil {
		return
	}
	res.Code, err = m.generator.GenerateCodeAutoIncrement(EC_CODE_TEMPLATE, count)
	if err != nil {
		return
	}
	err = tx.WithContext(req.Ctx).CashOutflow.Save(&res)
	if err != nil {
		return
	}
	if err = tx.CashOutflow.InsertActivity(id, req.Profile.ID, proto.ActivityType_CREATE.String(), nil); err != nil {
		return
	}
	var references []*int64
	references = append(references, &fields.PartyID,
		fields.CostCenterID, fields.ProjectID)
	if err = m.dbHelper.InsertReferences(req.Ctx, tx, res.ID, references); err != nil {
		return
	}
	return
}

func (m *repository) Edit(tx *query.QueryTx, req *common.RequestContext, d dto.CashOutflowData) (err error) {

	data, err := m.convertor.DataMap(d.Fields)
	if err != nil {
		return
	}
	err = tx.CashOutflow.UnderlyingDB().WithContext(req.Ctx).Model(
		&model.CashOutflow{ID: d.ID},
	).Updates(data).Error
	if err != nil {
		return
	}
	err = tx.CashOutflow.InsertActivity(d.ID, req.Profile.ID, proto.ActivityType_EDIT.String(), nil)
	if err != nil {
		return
	}
	var references []*int64
	references = append(references,
		&d.Fields.PartyID,
		d.Fields.CostCenterID,
		d.Fields.ProjectID,
	)
	if err = m.dbHelper.InsertReferences(req.Ctx, tx, d.ID, references, true); err != nil {
		return
	}

	return

}

func (m *repository) UpdateStatus(tx *query.QueryTx, req *common.RequestContext, d dto.UpdateStatusWithEvent, nextState string) (
	res *model.CashOutflow, err error) {

	qI := tx.CashOutflow
	res, err = tx.CashOutflow.WithContext(req.Ctx).Where(
		qI.Code.Eq(d.Body.PartyID),
		qI.CompanyID.Eq(req.ActiveCompany.ID),
	).First()
	if err != nil {
		return
	}
	_, err = tx.CashOutflow.WithContext(req.Ctx).Where(
		qI.CompanyID.Eq(req.ActiveCompany.ID),
		qI.Status.Eq(d.Body.CurrentState),
		qI.ID.Eq(res.ID),
	).UpdateSimple(
		qI.Status.Value(nextState),
	)
	if err != nil {
		return
	}
	res.Status = nextState
	return

}

func (r *repository) GetFilterOptions() []dto.FilterOptionDto {
	filterOptions := []dto.FilterOptionDto{}
	status := dto.FilterOptionDto{
		Param:     "status",
		Name:      "Estado",
		Type:      dto.FILTER_TYPE_STRING,
		Operators: dto.StringOperators,
		Options: []string{proto.State_DRAFT.String(), proto.State_PAID.String(),
			proto.State_UNPAID.String(), proto.State_CANCELLED.String()},
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
