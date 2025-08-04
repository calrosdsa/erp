package auth_admin_repo

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/gen/proto"
	"erp/internal/app/service/helpers"
	"erp/pkg/db"
	"fmt"

	"gorm.io/gen"
	"gorm.io/gen/field"
)

type RoleTemplateRepository interface {
	GetRoleTemplates(req *common.AdminRequestContext, d *dto.RequestPaginationData) (
		dto.PaginationResult[[]dto.RoleTemplateDto], error)
	GetRoleTemplate(req *common.AdminRequestContext, d *dto.RequestEntity) (
		dto.ResultEntity[dto.RoleTemplateDto], error)
	CreateRoleTemplate(req *common.AdminRequestContext, d *dto.CreateRoleTemplateRequest) error
}

type roleTemplateRepository struct {
	Q         *query.Query
	convertor helpers.ConvertorHelper
}

func NewRoleTemplateRepo(
	db db.Connection,
	helpers *helpers.Helpers,
) RoleTemplateRepository {
	return &roleTemplateRepository{
		Q:         db.GetQ(),
		convertor: helpers.Convertor,
	}
}

func (r *roleTemplateRepository) GetRoleTemplates(req *common.AdminRequestContext, d *dto.RequestPaginationData) (
	res dto.PaginationResult[[]dto.RoleTemplateDto], err error) {
	var (
		conds []gen.Condition
		order field.Expr
	)
	roleTemplateQ := r.Q.RoleTemplate
	builder := r.Q.WithContext(req.Ctx).RoleTemplate

	if d.Query != "" {
		conds = append(conds, roleTemplateQ.Name.Like("%"+d.Query+"%"))
	}

	limit, offset := r.convertor.ToPaginationParams(d.Page, d.Size)
	orderCol, ok := r.Q.RoleTemplate.GetFieldByName(d.OrderColumn) // maybe orderColStr == "id"
	if ok {
		if d.Order == "ASC" {
			order = orderCol.Asc()
		} else {
			order = orderCol.Desc()
		}
		builder = builder.Order(order)
		// User doesn't contains orderColStr
	}

	builder = builder.Select(
		roleTemplateQ.ID, roleTemplateQ.Name, roleTemplateQ.CreatedAt,
	).Where(conds...)
	total, err := builder.Count()
	if err != nil {
		return res, err
	}
	err = builder.Limit(limit).Offset(offset).Scan(&res.Results)
	res.Total = total
	return
}

func (r *roleTemplateRepository) GetRoleTemplate(req *common.AdminRequestContext, d *dto.RequestEntity) (
	res dto.ResultEntity[dto.RoleTemplateDto], err error) {
	roleTemplateQ := r.Q.RoleTemplate
	id := r.convertor.StrtoInt(d.ID)
	err = r.Q.RoleTemplate.WithContext(req.Ctx).Select(
		roleTemplateQ.Name, roleTemplateQ.ID, roleTemplateQ.CreatedAt,
	).Where(roleTemplateQ.ID.Eq(id)).Scan(&res.Entity)
	return
}

func (r *roleTemplateRepository) CreateRoleTemplate(req *common.AdminRequestContext,
	d *dto.CreateRoleTemplateRequest) (err error) {
	tx := r.Q.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	
	partyID, err := tx.RoleTemplate.InsertParty(proto.PartyAdminType_roleTemplate.String())
	if err != nil {
		return
	}
	fmt.Println("PARTY ID",partyID)
	roleTemplate := model.RoleTemplate{}
	roleTemplate.Name = d.Body.Name
	roleTemplate.ID = partyID
	err = tx.RoleTemplate.WithContext(req.Ctx).Save(&roleTemplate)
	if err != nil {
		return
	}
	err = tx.Commit()
	return
}
