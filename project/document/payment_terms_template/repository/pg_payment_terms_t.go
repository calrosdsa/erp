package payment_terms_t_repo

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

type PaymentTermsTemplateRepo interface {
	GetPaymentTermsTemplates(req *common.RequestContext, d dto.PaymentTermsTemplateRequest) (
		res []dto.PaymentTermsTemplateDto, err error)
	GetPaymentTermsTemplateDetail(req *common.RequestContext, d dto.RequestEntity) (
		res dto.PaymentTermsTemplateDto, err error)
	CreatePaymentTermsTemplate(tx *query.QueryTx,req *common.RequestContext, d dto.PaymentTermsTemplateData) (res model.PaymentTermsTemplate, err error)
	EditPaymentTermsTemplate(tx *query.QueryTx,req *common.RequestContext, d dto.PaymentTermsTemplateData) (err error)
	UpdateStatus(req *common.RequestContext, d dto.UpdateStatusWithEvent, nextState string) (err error)
	GetFilterOptions(lng string) []dto.FilterOptionDto
	Greet(name string) (string, error)
}
type paymentTermsRepo struct {
	convertor helpers.ConvertorHelper
	Q         *query.Query
	query     helpers.QueryHelper
}

func NewRepository(
	conn db.Connection,
	helpers *helpers.Helpers,
) PaymentTermsTemplateRepo {
	return &paymentTermsRepo{
		convertor: helpers.Convertor,
		Q:         conn.GetQ(),
		query:     helpers.Query,
	}
}

func (r *paymentTermsRepo) Greet(name string) (string, error) {
	return fmt.Sprintf("Hello, %s",name),nil
}


func (r *paymentTermsRepo) GetPaymentTermsTemplateDetail(req *common.RequestContext, d dto.RequestEntity) (
	res dto.PaymentTermsTemplateDto, err error) {
	id := r.convertor.StrtoInt(d.ID)
	qI := r.Q.PaymentTermsTemplate
	err = qI.WithContext(req.Ctx).Select(
		qI.ID, qI.UUID, qI.Name, qI.Status, qI.CreatedAt,qI.CreatedAt,
	).
		Where(
			qI.CompanyID.Eq(req.ActiveCompany.ID),
			qI.ID.Eq(id),
		).Scan(&res)
	return
}

func (r *paymentTermsRepo) CreatePaymentTermsTemplate(tx *query.QueryTx,req *common.RequestContext, d dto.PaymentTermsTemplateData) (
	res model.PaymentTermsTemplate, err error) {
	id, err := tx.PaymentTermsTemplate.InsertParty(proto.PartyType_paymentTermsTemplate.String())
	if err != nil {
		return
	}
	fields := d.Fields
	res.ID = id
	res.CompanyID = req.ActiveCompany.ID
	res.Status = proto.State_ENABLED.String()
	if err = r.convertor.CopyStructData(fields, &res); err != nil {
		return
	}
	err = tx.WithContext(req.Ctx).PaymentTermsTemplate.Save(&res)
	if err != nil {
		return
	}
	if err = tx.Invoice.InsertActivity(id, req.Profile.ID, proto.ActivityType_CREATE.String(), nil); err != nil {
		return
	}
	return
}
func (r *paymentTermsRepo) EditPaymentTermsTemplate(tx *query.QueryTx,req *common.RequestContext, d dto.PaymentTermsTemplateData) (
	err error) {
	data, err := r.convertor.DataMap(d.Fields)
	if err != nil {
		return
	}
	err = tx.PaymentTermsTemplate.UnderlyingDB().WithContext(req.Ctx).Model(
		&model.PaymentTermsTemplate{ID: d.ID},
	).Updates(data).Error
	if err != nil {
		return
	}
	err = tx.PaymentTermsTemplate.InsertActivity(d.ID, req.Profile.ID, proto.ActivityType_EDIT.String(), nil)
	if err != nil {
		return
	}
	return
}

func (r *paymentTermsRepo) UpdateStatus(req *common.RequestContext, d dto.UpdateStatusWithEvent,
	nextState string) (err error) {
	qI := r.Q.PaymentTermsTemplate
	_, err = r.Q.PaymentTermsTemplate.WithContext(req.Ctx).Where(
		qI.CompanyID.Eq(req.ActiveCompany.ID),
		qI.Status.Eq(d.Body.CurrentState),
		qI.UUID.Eq(d.Body.PartyID),
	).UpdateSimple(
		qI.Status.Value(nextState),
	)
	return
}

func (r *paymentTermsRepo) GetPaymentTermsTemplates(req *common.RequestContext, d dto.PaymentTermsTemplateRequest) (
	res []dto.PaymentTermsTemplateDto, err error) {
	var (
		generateSQL strings.Builder
	)
	builder := r.Q.WithContext(req.Ctx).PaymentTermsTemplate
	queryData := r.convertor.GenerateQueryMap(d)
	params := r.paymentTermsTemplateQuery(req, queryData, &generateSQL)
	err = builder.UnderlyingDB().Raw(generateSQL.String(), params...).Scan(&res).Error
	return
}

func (r *paymentTermsRepo) paymentTermsTemplateQuery(req *common.RequestContext, d map[string]string,
	generateSQL *strings.Builder) (
	params []interface{}) {
	var (
		whereSQL strings.Builder
	)
	generateSQL.WriteString(`SELECT 
	e.id,e.uuid,e.created_at,e.name,e.status
	from payment_terms_templates as e `)
	whereSQL.WriteString(` e.deleted_at is null and e.company_id = ?`)
	params = append(params, req.ActiveCompany.ID)
	columnFilters := []string{
		"name",
		"created_at",
		"status",
	}
	r.query.FilterBuilder(&whereSQL, &params, d, columnFilters...)
	helper.JoinWhereBuilder(generateSQL, whereSQL)
	r.query.OrderAndLimitBuilder(generateSQL, d)
	return
}

func (r *paymentTermsRepo) GetFilterOptions(lng string) []dto.FilterOptionDto {
	filterOptions := []dto.FilterOptionDto{}
	status := dto.FilterOptionDto{
		Param:     "status",
		Name:      "Estado",
		Type:      dto.FILTER_TYPE_STRING,
		Operators: dto.StringOperators,
		Options:   []string{proto.State_ENABLED.String(), proto.State_DISABLED.String()},
	}
	createdAt := dto.FilterOptionDto{
		Name:      "Fecha de Creacion",
		Param:     "created_at",
		Type:      dto.FILTER_TYPE_DATE,
		Operators: dto.DateOperators,
	}

	name := dto.FilterOptionDto{
		Name:      "Nombre",
		Param:     "name",
		Type:      dto.FILTER_TYPE_STRING,
		Operators: dto.StringOperators,
	}

	filterOptions = append(filterOptions, status, createdAt, name)
	return filterOptions
}
