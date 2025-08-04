package project_repo

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/gen/proto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/pkg/db"

	"gorm.io/gen"
	"gorm.io/gen/field"
)

type ProjectRepository interface {
	GetProject(req *common.RequestContext, d *dto.RequestEntity) (
		res dto.ResultEntity[dto.ProjectDto], err error)
	CreateProject(req *common.RequestContext, d *dto.CreateProjectRequest) (
		res dto.ProjectDto, err error)
	GetProjects(req *common.RequestContext, d *dto.RequestPaginationData) (
		res dto.PaginationResult[[]dto.ProjectDto], err error)
}

type costCenterRepository struct {
	Q         *query.Query
	convertor helpers.ConvertorHelper
	dbHelper  db.DbHelper
}

func NewProjectRepository(
	db db.Connection,
	helpers *helpers.Helpers,
) ProjectRepository {
	return &costCenterRepository{
		Q:         db.GetQ(),
		convertor: helpers.Convertor,
		dbHelper:  db.GetDbHelper(),
	}
}

func (r *costCenterRepository) GetProject(req *common.RequestContext, d *dto.RequestEntity) (
	res dto.ResultEntity[dto.ProjectDto], err error) {
	id := r.convertor.StrtoInt(d.ID)
	projectQ := r.Q.Project
	err = projectQ.WithContext(req.Ctx).Select(
		projectQ.ID, projectQ.Name, projectQ.Status,
	).
		Where(projectQ.CompanyID.Eq(req.ActiveCompany.ID), projectQ.ID.Eq(id)).
		Scan(&res.Entity)
	return
}

func (r *costCenterRepository) CreateProject(req *common.RequestContext, d *dto.CreateProjectRequest) (
	res dto.ProjectDto, err error) {
	tx := r.Q.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	project := model.Project{}
	project.Name = d.Body.Name
	project.CompanyID = req.ActiveCompany.ID
	err = r.dbHelper.ValidateName(tx.Project.UnderlyingDB(), &project)
	if err != nil {
		return res, domain.ERROR_NAME_TAKEN
	}
	projectID, err := tx.Project.InsertParty(proto.PartyType_project.String())
	if err != nil {
		return
	}
	project.ID = projectID
	project.Status = d.Body.Status
	err = tx.Project.WithContext(req.Ctx).Save(&project)
	if err != nil {
		return
	}
	err = tx.Commit()
	if err != nil {
		return
	}
	res = dto.ProjectDtoFromModel(&project)
	return
}

func (r *costCenterRepository) GetProjects(req *common.RequestContext, d *dto.RequestPaginationData) (
	res dto.PaginationResult[[]dto.ProjectDto], err error) {
	var (
		conds []gen.Condition
		order field.Expr
	)
	projectQ := r.Q.Project
	builder := r.Q.WithContext(req.Ctx).Project

	limit, offset := r.convertor.ToPaginationParams(d.Page, d.Size)
	orderCol, ok := r.Q.Project.GetFieldByName(d.OrderColumn) // maybe orderColStr == "id"
	if ok {
		if d.Order == "ASC" {
			order = orderCol.Asc()
		} else {
			order = orderCol.Desc()
		}
		builder = builder.Order(order)
	}
	//ADDING CONDITIONS
	conds = append(conds, projectQ.CompanyID.Eq(req.ActiveCompany.ID))
	if d.Query != "" {
		conds = append(conds, projectQ.Name.Like("%"+d.Query+"%"))
	}

	builder = builder.Select(
		projectQ.ID, projectQ.Name, projectQ.Status,
	).Where(conds...)
	total, err := builder.Count()
	if err != nil {
		return res, err
	}
	err = builder.Limit(limit).Offset(offset).Scan(&res.Results)
	res.Total = total
	return
}
