package terms_and_conditions_repo

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

type TermsAndConditionsRepository interface {
	GetTermsAndConditions(req *common.RequestContext, d *dto.TermsAndConditionsRequest) (
		res []dto.TermsAndConditionsDto, err error)
	GetTermsAndConditionsDetial(req *common.RequestContext, d *dto.RequestEntity) (
		res dto.TermsAndConditionsDto, err error)
	CreateTermsAndConditions(req *common.RequestContext, d dto.TermsAndConditionsData) (res model.TermsAndCondition, err error)
	EditTermsAndConditions(req *common.RequestContext, d dto.TermsAndConditionsData) (err error)
	UpdateStatus(req *common.RequestContext, d *dto.UpdateStatusWithEvent, nextState string) (err error)
	GetFilterOptions(lng string) []dto.FilterOptionDto
}
type termsAndConditionsRepository struct {
	convertor helpers.ConvertorHelper
	Q         *query.Query
	query     helpers.QueryHelper
}

func NewRepository(
	conn db.Connection,
	helpers *helpers.Helpers,
) TermsAndConditionsRepository {
	return &termsAndConditionsRepository{
		convertor: helpers.Convertor,
		Q:         conn.GetQ(),
		query:     helpers.Query,
	}
}

func (r *termsAndConditionsRepository) GetTermsAndConditionsDetial(req *common.RequestContext, d *dto.RequestEntity) (
	res dto.TermsAndConditionsDto, err error) {
	id := r.convertor.StrtoInt(d.ID)
	qI := r.Q.TermsAndCondition
	err = qI.WithContext(req.Ctx).Select(
		qI.ID, qI.UUID, qI.Name, qI.Status, qI.TermsAndConditions, qI.CreatedAt,
	).
		Where(
			qI.CompanyID.Eq(req.ActiveCompany.ID),
			qI.ID.Eq(id),
		).Scan(&res)
	return
}

func (r *termsAndConditionsRepository) CreateTermsAndConditions(req *common.RequestContext, d dto.TermsAndConditionsData) (
	res model.TermsAndCondition, err error) {
	tx := r.Q.Begin()
	defer func() {
		if err != nil {
			err = tx.Rollback()
		}
	}()
	id, err := tx.TermsAndCondition.InsertParty(proto.PartyType_termsAndConditions.String())
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
	err = tx.WithContext(req.Ctx).TermsAndCondition.Save(&res)
	if err != nil {
		return
	}
	if err = tx.Invoice.InsertActivity(id, req.Profile.ID, proto.ActivityType_CREATE.String(), nil); err != nil {
		return
	}
	err = tx.Commit()
	return
}
func (r *termsAndConditionsRepository) EditTermsAndConditions(req *common.RequestContext, d dto.TermsAndConditionsData) (
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
	err = tx.TermsAndCondition.UnderlyingDB().WithContext(req.Ctx).Model(
		&model.TermsAndCondition{ID: d.ID},
	).Updates(data).Error
	if err != nil {
		return
	}
	err = tx.TermsAndCondition.InsertActivity(d.ID, req.Profile.ID, proto.ActivityType_EDIT.String(), nil)
	if err != nil {
		return
	}
	err = tx.Commit()
	return
}

func (r *termsAndConditionsRepository) UpdateStatus(req *common.RequestContext, d *dto.UpdateStatusWithEvent,
	nextState string) (err error) {
	termsAndConditionsQ := r.Q.TermsAndCondition
	_, err = r.Q.TermsAndCondition.WithContext(req.Ctx).Where(
		termsAndConditionsQ.CompanyID.Eq(req.ActiveCompany.ID),
		termsAndConditionsQ.Status.Eq(d.Body.CurrentState),
		termsAndConditionsQ.UUID.Eq(d.Body.PartyID),
	).UpdateSimple(
		termsAndConditionsQ.Status.Value(nextState),
	)
	return
}

func (r *termsAndConditionsRepository) GetTermsAndConditions(req *common.RequestContext, d *dto.TermsAndConditionsRequest) (
	res []dto.TermsAndConditionsDto, err error) {
	var (
		generateSQL strings.Builder
	)
	builder := r.Q.WithContext(req.Ctx).TermsAndCondition
	queryData := r.convertor.GenerateQueryMap(d)
	params := r.termsAndConditionsQuery(req, queryData, &generateSQL)
	err = builder.UnderlyingDB().Raw(generateSQL.String(), params...).Scan(&res).Error
	return
}

func (r *termsAndConditionsRepository) termsAndConditionsQuery(req *common.RequestContext, d map[string]string,
	generateSQL *strings.Builder) (
	params []interface{}) {
	var (
		whereSQL strings.Builder
	)
	generateSQL.WriteString(`SELECT 
	e.id,e.uuid,e.created_at,e.name,e.terms_and_conditions,e.status
	from terms_and_conditions as e `)
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

func (r *termsAndConditionsRepository) GetFilterOptions(lng string) []dto.FilterOptionDto {
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
