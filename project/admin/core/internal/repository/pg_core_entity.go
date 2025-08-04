package pg_core

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/internal/app/service/helpers"
	"erp/pkg/db"

	"gorm.io/gen"
	"gorm.io/gen/field"
)

type CoreEntityRepository interface {
	GetEntities(req *common.AdminRequestContext, d *dto.RequestPaginationData) (
		dto.PaginationResult[[]dto.EntityDto], error)
	GetEntity(req *common.AdminRequestContext, d *dto.RequestEntity) (
		dto.ResultEntity[dto.EntityDetailDto], error)
	CreateEntity(req *common.AdminRequestContext, d *dto.CreateEntityRequest) error
	AddEntityAction(req *common.AdminRequestContext, d *dto.AddEntityActionRequest) error
}

type coreEntityRepository struct {
	Q         *query.Query
	convertor helpers.ConvertorHelper
}

func NewCoreRepository(
	db db.Connection,
	helpers *helpers.Helpers,
) CoreEntityRepository {
	return &coreEntityRepository{
		Q:         db.GetQ(),
		convertor: helpers.Convertor,
	}
}
func (r *coreEntityRepository) AddEntityAction(req *common.AdminRequestContext, d *dto.AddEntityActionRequest) error {
	action := model.Action{}
	action.Name = d.Body.Name
	action.EntityID = d.Body.EntityID
	err := r.Q.Action.WithContext(req.Ctx).Save(&action)
	if err != nil {
		return err
	}
	return err
}

func (r *coreEntityRepository) GetEntity(req *common.AdminRequestContext, d *dto.RequestEntity) (
	res dto.ResultEntity[dto.EntityDetailDto], err error) {
	entityQ := r.Q.Entity
	actionQ := r.Q.Action
	entityID := r.convertor.StrtoInt(d.ID)
	err = entityQ.WithContext(req.Ctx).Where(
		r.Q.Entity.ID.Eq(entityID),
	).Select(
		entityQ.ID, entityQ.Name,
	).Scan(&res.Entity.Entity)
	if err != nil {
		return
	}
	err = actionQ.WithContext(req.Ctx).Select(
		actionQ.ID, actionQ.Name,
	).Where(actionQ.EntityID.Eq(res.Entity.Entity.ID)).Scan(&res.Entity.Actions)
	return
}

func (r *coreEntityRepository) CreateEntity(req *common.AdminRequestContext, d *dto.CreateEntityRequest) (err error) {
	tx := r.Q.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	entity := model.Entity{}
	entity.Name = d.Body.Name
	err = tx.Entity.WithContext(req.Ctx).Save(&entity)
	if err != nil {
		return
	}
	err = tx.Commit()
	return
}

func (r *coreEntityRepository) GetEntities(req *common.AdminRequestContext, d *dto.RequestPaginationData) (
	res dto.PaginationResult[[]dto.EntityDto], err error) {
	var (
		conds []gen.Condition
		order field.Expr
	)
	entityQ := r.Q.Entity
	builder := r.Q.WithContext(req.Ctx).Entity

	//ADDING CONDITIONS
	if d.Query != "" {
		conds = append(conds, entityQ.Name.Like("%"+d.Query+"%"))
	}

	limit, offset := r.convertor.ToPaginationParams(d.Page, d.Size)
	orderCol, ok := r.Q.Entity.GetFieldByName(d.OrderColumn) // maybe orderColStr == "id"
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
		entityQ.ID, entityQ.Name,entityQ.Href,
	).Where(conds...)
	total, err := builder.Count()
	if err != nil {
		return res, err
	}
	err = builder.Limit(limit).Offset(offset).Scan(&res.Results)
	res.Total = total
	return
}
