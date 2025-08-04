package warehouse_repo

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

type WarehouseRepository interface {
	GetWareHouseDetail(req *common.RequestContext, i *dto.RequestEntity) (
		res dto.ResultEntity[dto.WareHouseDto], err error)
	GetWareHouses(req *common.RequestContext, d *dto.RequestPaginationData) (
		res dto.PaginationResult[[]dto.WareHouseDto], err error)
	CreateWareHouse(req *common.RequestContext, i *dto.CreateWareHouseRequest) (err error)
	GetWarehouseTreeView(req *common.RequestContext) (
		res []dto.TreeEntryDto, err error)
	EditWarehouse(req *common.RequestContext, d *dto.EditWarehouseRequest) (err error)
}

type warehouseRepository struct {
	Q         *query.Query
	convertor helpers.ConvertorHelper
}

func NewWareHouseReposotory(
	conn db.Connection,
	helpers *helpers.Helpers,
) *warehouseRepository {
	return &warehouseRepository{
		Q:         conn.GetQ(),
		convertor: helpers.Convertor,
	}
}

func (r *warehouseRepository) EditWarehouse(req *common.RequestContext, d *dto.EditWarehouseRequest) (
	err error) {
	tx := r.Q.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	warehouseQ := tx.WareHouse
	var columns []field.AssignExpr

	columns = append(columns, warehouseQ.Name.Value(d.Body.Name))
	if d.Body.ParentID != 0 {
		columns = append(columns, warehouseQ.ParentID.Value(d.Body.ParentID))
	}
	_, err = tx.WareHouse.WithContext(req.Ctx).Where(
		warehouseQ.ID.Eq(d.Body.ID), warehouseQ.CompanyID.Eq(req.ActiveCompany.ID),
	).UpdateSimple(columns...)
	if err != nil {
		return
	}
	err = tx.WareHouse.InsertActivity(d.Body.ID, req.Profile.ID, proto.ActivityType_EDIT.String(), nil)
	if err != nil {
		return
	}
	err = tx.Commit()
	return
}

func (r *warehouseRepository) GetWarehouseTreeView(req *common.RequestContext) (
	res []dto.TreeEntryDto, err error) {
	query := `WITH RECURSIVE data_cte AS (
		SELECT 
			id, 
			uuid,
			parent_id,
			is_group,
			name
		FROM ware_houses
		WHERE parent_id IS NULL and company_id = ?
		UNION ALL 
		SELECT 
			l.id, 
			l.uuid,
			l.parent_id,
			l.is_group,
			l.name
		FROM ware_houses l
		INNER JOIN data_cte d
			ON l.parent_id = d.id  
	)
	SELECT * FROM data_cte;`
	err = r.Q.WareHouse.UnderlyingDB().WithContext(req.Ctx).Raw(query, req.ActiveCompany.ID).Scan(&res).Error
	return
}

func (s *warehouseRepository) GetWareHouseDetail(req *common.RequestContext, i *dto.RequestEntity) (
	res dto.ResultEntity[dto.WareHouseDto], err error) {
	id := s.convertor.StrtoInt(i.ID)
	warehouse, err := s.Q.WareHouse.WithContext(req.Ctx).Where(
		s.Q.WareHouse.ID.Eq(id),
		s.Q.WareHouse.CompanyID.Eq(req.ActiveCompany.ID),
	).First()
	if err != nil {
		return
	}
	res.Entity = dto.WarehouseDtoFromModel(warehouse)
	return res, err
}

func (r *warehouseRepository) GetWareHouses(req *common.RequestContext, d *dto.RequestPaginationData) (
	res dto.PaginationResult[[]dto.WareHouseDto], err error) {
	var (
		conds []gen.Condition
		order field.Expr
	)
	warehouseQ := r.Q.WareHouse
	builder := r.Q.WithContext(req.Ctx).WareHouse

	//ADDING CONDITIONS
	conds = append(conds, warehouseQ.CompanyID.Eq(req.ActiveCompany.ID))

	if d.IsGroup != "" {
		conds = append(conds, warehouseQ.IsGroup.Is(r.convertor.StrToBool(d.IsGroup)))
	}

	limit, offset := r.convertor.ToPaginationParams(d.Page, d.Size)
	orderCol, ok := r.Q.WareHouse.GetFieldByName(d.OrderColumn) // maybe orderColStr == "id"
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
		warehouseQ.ID, warehouseQ.Name, warehouseQ.IsGroup, warehouseQ.CreatedAt, warehouseQ.UUID,
	).
		Where(conds...)
	total, err := builder.Count()
	if err != nil {
		return res, err
	}
	err = builder.Limit(limit).Offset(offset).Scan(&res.Results)
	res.Total = total
	return
}

func (s *warehouseRepository) CreateWareHouse(req *common.RequestContext, i *dto.CreateWareHouseRequest) (err error) {
	var (
		wareHouse model.WareHouse
	)
	tx := s.Q.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	partyId, err := tx.WareHouse.InsertParty(domain.PARTY_WAREHOUSE)
	if err != nil {
		return err
	}
	wareHouse.ID = partyId
	wareHouse.Name = i.Body.Name
	wareHouse.CompanyID = req.ActiveCompany.ID
	if i.Body.ParentID != 0 {
		wareHouse.ParentID = &i.Body.ParentID
	}
	wareHouse.IsGroup = i.Body.IsGroup
	// wareHouse.Enabled = i.Body.Enabled
	err = tx.WareHouse.WithContext(req.Ctx).Save(&wareHouse)
	if err != nil {
		return
	}
	err = tx.Commit()
	return err
}
