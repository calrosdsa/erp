package cost_center_repo

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/gen/proto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/pkg/db"
	"fmt"

	"gorm.io/gen"
	"gorm.io/gen/field"
)

type CostCenterRepository interface {
	GetCostCenter(req *common.RequestContext, d *dto.RequestEntity) (
		res dto.ResultEntity[dto.CostCenterDto], err error)
	CreateCostCenter(req *common.RequestContext, d *dto.CreateCostCenterRequet) (
		res dto.CostCenterDto, err error)
	GetCostCenters(req *common.RequestContext, d *dto.RequestPaginationData) (
		res dto.PaginationResult[[]dto.CostCenterDto], err error)
}

type costCenterRepository struct {
	Q         *query.Query
	convertor helpers.ConvertorHelper
	dbHelper  db.DbHelper
}

func NewCostCenterRepo(
	db db.Connection,
	helpers *helpers.Helpers,
) CostCenterRepository {
	return &costCenterRepository{
		Q:         db.GetQ(),
		convertor: helpers.Convertor,
		dbHelper:  db.GetDbHelper(),
	}
}

func (r *costCenterRepository) GetCostCenter(req *common.RequestContext, d *dto.RequestEntity) (
	res dto.ResultEntity[dto.CostCenterDto], err error) {
	id := r.convertor.StrtoInt(d.ID)
	costCenterQ := r.Q.CostCenter
	err = costCenterQ.WithContext(req.Ctx).Select(
		costCenterQ.ID, costCenterQ.Name, costCenterQ.Status,
	).
		Where(costCenterQ.CompanyID.Eq(req.ActiveCompany.ID), costCenterQ.ID.Eq(id)).
		Scan(&res.Entity)
	return
}

func (r *costCenterRepository) CreateCostCenter(req *common.RequestContext, d *dto.CreateCostCenterRequet) (
	res dto.CostCenterDto, err error) {
	tx := r.Q.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	costCenter := model.CostCenter{}
	costCenter.Name = d.Body.Name
	costCenter.CompanyID = req.ActiveCompany.ID
	err = r.dbHelper.ValidateName(tx.CostCenter.UnderlyingDB(), &costCenter)
	if err != nil {
		return res, domain.ERROR_NAME_TAKEN
	}
	costCenterID, err := tx.CostCenter.InsertParty(proto.PartyType_costCenter.String())
	if err != nil {
		return
	}
	costCenter.ID = costCenterID
	costCenter.Status = proto.State_ENABLED.String()
	err = tx.CostCenter.WithContext(req.Ctx).Save(&costCenter)
	if err != nil {
		return
	}
	err = tx.Commit()
	if err != nil {
		return
	}
	res = dto.CostCenterDtoFromModel(&costCenter)
	return
}

func (r *costCenterRepository) GetCostCenters(req *common.RequestContext, d *dto.RequestPaginationData) (
	res dto.PaginationResult[[]dto.CostCenterDto], err error) {
	var (
		conds []gen.Condition
		order field.Expr
	)
	costCenterQ := r.Q.CostCenter
	builder := r.Q.WithContext(req.Ctx).CostCenter

	limit, offset := r.convertor.ToPaginationParams(d.Page, d.Size)
	orderCol, ok := r.Q.CostCenter.GetFieldByName(d.OrderColumn) // maybe orderColStr == "id"
	if ok {
		if d.Order == "ASC" {
			order = orderCol.Asc()
		} else {
			order = orderCol.Desc()
		}
		builder = builder.Order(order)
	}
	//ADDING CONDITIONS
	conds = append(conds, costCenterQ.CompanyID.Eq(req.ActiveCompany.ID))
	fmt.Println("QUERY", d.Query)
	if d.Query != "" {
		conds = append(conds, costCenterQ.Name.Like("%"+d.Query+"%"))
	}

	builder = builder.Select(
		costCenterQ.ID, costCenterQ.Name, costCenterQ.Status, costCenterQ.CreatedAt,
	).Where(conds...)
	total, err := builder.Count()
	if err != nil {
		return res, err
	}
	err = builder.Limit(limit).Offset(offset).Scan(&res.Results)
	res.Total = total
	return
}
