package payment_terms_repo

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/gen/proto"
	"erp/internal/app/service/helpers"
	"erp/pkg/db"
	"strings"

	"gorm.io/gen/helper"
)

type PaymentTermsRepo interface {
	GetPaymentTerms(req *common.RequestContext, d *dto.PaymentTermsRequest) (
		res []dto.PaymentTermsDto, err error)
	GetPaymentTermsDetail(req *common.RequestContext, d *dto.RequestEntity) (
		res dto.PaymentTermsDto, err error)
	CreatePaymentTerms(req *common.RequestContext, d dto.PaymentTermsData) (res model.PaymentTerm, err error)
	EditPaymentTerms(req *common.RequestContext, d dto.PaymentTermsData) (err error)
	UpdateStatus(req *common.RequestContext, d *dto.UpdateStatusWithEvent, nextState string) (err error)
	GetFilterOptions(lng string) []dto.FilterOptionDto
}
type paymentTermsRepo struct {
	convertor helpers.ConvertorHelper
	Q         *query.Query
	query     helpers.QueryHelper
}

func NewRepository(
	conn db.Connection,
	helpers *helpers.Helpers,
) PaymentTermsRepo {
	return &paymentTermsRepo{
		convertor: helpers.Convertor,
		Q:         conn.GetQ(),
		query:     helpers.Query,
	}
}

func (r *paymentTermsRepo) GetPaymentTermsDetail(req *common.RequestContext, d *dto.RequestEntity) (
	res dto.PaymentTermsDto, err error) {
	qI := r.Q.PaymentTerm
	id := r.convertor.StrtoInt(d.ID)
	err = qI.WithContext(req.Ctx).Select(
		qI.ID, qI.UUID, qI.Name, qI.Status, qI.CreatedAt, qI.CreditDays, qI.Description,
		qI.Discount, qI.DiscountType, qI.DiscountValidity, qI.DiscountValidityBaseOn,
		qI.DueDateBaseOn, qI.InvoicePortion,
	).
		Where(
			qI.CompanyID.Eq(req.ActiveCompany.ID),
			qI.ID.Eq(id),
		).Scan(&res)
	return
}

func (r *paymentTermsRepo) CreatePaymentTerms(req *common.RequestContext, d dto.PaymentTermsData) (
	res model.PaymentTerm, err error) {
	tx := r.Q.Begin()
	defer func() {
		if err != nil {
			err = tx.Rollback()
		}
	}()
	id, err := tx.PaymentTerm.InsertParty(proto.PartyType_paymentTerms.String())
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
	err = tx.WithContext(req.Ctx).PaymentTerm.Save(&res)
	if err != nil {
		return
	}
	if err = tx.Invoice.InsertActivity(id, req.Profile.ID, proto.ActivityType_CREATE.String(), nil); err != nil {
		return
	}
	err = tx.Commit()
	return
}
func (r *paymentTermsRepo) EditPaymentTerms(req *common.RequestContext, d dto.PaymentTermsData) (
	err error) {
	tx := r.Q.Begin()
	defer func() {
		if err != nil {
			err = tx.Rollback()
		}
	}()
	data, err := r.convertor.DataMap(d.Fields)
	if err != nil {
		return
	}
	err = tx.PaymentTerm.UnderlyingDB().WithContext(req.Ctx).Model(
		&model.PaymentTerm{ID: d.ID},
	).Updates(data).Error
	if err != nil {
		return
	}
	err = tx.PaymentTerm.InsertActivity(d.ID, req.Profile.ID, proto.ActivityType_EDIT.String(), nil)
	if err != nil {
		return
	}
	err = tx.Commit()
	return
}

func (r *paymentTermsRepo) UpdateStatus(req *common.RequestContext, d *dto.UpdateStatusWithEvent,
	nextState string) (err error) {
	qI := r.Q.PaymentTerm
	_, err = r.Q.PaymentTerm.WithContext(req.Ctx).Where(
		qI.CompanyID.Eq(req.ActiveCompany.ID),
		qI.Status.Eq(d.Body.CurrentState),
		qI.UUID.Eq(d.Body.PartyID),
	).UpdateSimple(
		qI.Status.Value(nextState),
	)
	return
}

func (r *paymentTermsRepo) GetPaymentTerms(req *common.RequestContext, d *dto.PaymentTermsRequest) (
	res []dto.PaymentTermsDto, err error) {
	var (
		generateSQL strings.Builder
	)
	builder := r.Q.WithContext(req.Ctx).PaymentTerm
	queryData := r.convertor.GenerateQueryMap(d)
	params := r.paymentTermsQuery(req, queryData, &generateSQL)
	err = builder.UnderlyingDB().Raw(generateSQL.String(), params...).Scan(&res).Error
	return
}

func (r *paymentTermsRepo) paymentTermsQuery(req *common.RequestContext, d map[string]string,
	generateSQL *strings.Builder) (
	params []interface{}) {
	var (
		whereSQL strings.Builder
	)
	generateSQL.WriteString(`SELECT 
	e.id,e.uuid,e.created_at,e.name,e.status,
	e.invoice_portion,e.credit_days,e.due_date_base_on,e.description,
	e.discount_type,e.discount,e.discount_validity_base_on,e.discount_validity
	from payment_terms as e `)
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
