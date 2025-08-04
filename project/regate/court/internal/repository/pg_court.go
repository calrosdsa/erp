package court_repo

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

type CourtRepository interface {
	CreateCourt(req *common.RequestContext, i dto.CreateCourtBody) (model.Court, error)
	GetCourt(req *common.RequestContext, i *dto.RequestEntity) (dto.ResultEntity[dto.CourtDto], error)
	GetCourts(req *common.RequestContext, i dto.CourtsRequest) ([]dto.CourtDto, error)
	EditCourt(req *common.RequestContext,d dto.EditCourtBody)(err error)
	GetFilterOptions(lng string) []dto.FilterOptionDto
}
type courtRepository struct {
	conn      db.Connection
	Q         *query.Query
	convertor helpers.ConvertorHelper
	query helpers.QueryHelper
}

func NewCourtRepositository(
	conn db.Connection,
	helpers *helpers.Helpers,
) CourtRepository {
	return &courtRepository{
		conn:      conn,
		convertor: helpers.Convertor,
		Q:         conn.GetQ(),
		query:     helpers.Query,
	}
}
func (r *courtRepository) EditCourt(req *common.RequestContext, d dto.EditCourtBody) (err error) {
	tx := r.Q.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	courtQ := tx.Court
	var columns []field.AssignExpr
	columns = append(columns, courtQ.Name.Value(d.Name))
	_, err = tx.Court.WithContext(req.Ctx).Where(
		courtQ.ID.Eq(d.CourtID), courtQ.CompanyID.Eq(req.ActiveCompany.ID),
	).UpdateSimple(columns...)
	if err != nil {
		return
	}
	err = tx.Court.InsertActivity(d.CourtID, req.Profile.ID, proto.ActivityType_EDIT.String(), nil)
	if err != nil {
		return
	}
	err = tx.Commit()
	return
}

func (r *courtRepository) CreateCourt(req *common.RequestContext, i dto.CreateCourtBody) (court model.Court, err error) {
	fmt.Println("CRAETE COURT..")
	tx := r.Q.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	partyId, err := tx.Court.InsertParty(proto.RegatePartyType_court.String())
	if err != nil {
		return
	}
	court.ID = partyId
	court.Name = i.Name
	// court.Status = i.Status
	court.CompanyID  = req.ActiveCompany.ID
	err = tx.Court.Save(&court)
	if err != nil {
		return
	}
	err = tx.Commit()
	return
}
func (r *courtRepository) GetCourt(req *common.RequestContext, i *dto.RequestEntity) (
	res dto.ResultEntity[dto.CourtDto], err error) {
	courtQ := r.Q.Court
	courtID := r.convertor.StrtoInt(i.ID)
	err =  r.Q.WithContext(req.Ctx).Court.Select(
		courtQ.ID,courtQ.UUID,courtQ.Name, courtQ.Status, courtQ.CreatedAt,
	).Where(
		courtQ.CompanyID.Eq(req.ActiveCompany.ID),
		r.Q.Court.ID.Eq(courtID),
	).Scan(&res.Entity)
	if err != nil {
		return res, err
	}
	return
}
func (r *courtRepository) GetCourts(req *common.RequestContext, d dto.CourtsRequest) (
	res []dto.CourtDto, err error) {
		var (
			generateSQL strings.Builder
		)
		builder := r.Q.WithContext(req.Ctx).Court
		queryData := r.convertor.GenerateQueryMap(d)
		params := r.courtListQuery(req, queryData, &generateSQL)
		err = builder.UnderlyingDB().Raw(generateSQL.String(), params...).Scan(&res).Error
		return
}


func (m *courtRepository) courtListQuery(req *common.RequestContext, d map[string]string,
	generateSQL *strings.Builder) (params []interface{}) {
	var (
		whereSQL strings.Builder
	)
	generateSQL.WriteString(`SELECT 
		e.id,e.uuid,e.created_at,e.name,e.status
		from r_courts as e `)
	whereSQL.WriteString(` e.deleted_at is null and e.company_id = ?`)
	params = append(params, req.ActiveCompany.ID)
	columnFilters := []string{
		"id",
		"name",
		"created_at",
		"status",
	}
	m.query.FilterBuilder(&whereSQL, &params, d, columnFilters...)
	helper.JoinWhereBuilder(generateSQL, whereSQL)
	m.query.OrderAndLimitBuilder(generateSQL, d)
	return
}


func (r *courtRepository) GetFilterOptions(lng string) []dto.FilterOptionDto {
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