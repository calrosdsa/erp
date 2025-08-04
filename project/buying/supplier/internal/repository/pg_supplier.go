package supplier_repo

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/gen/proto"
	"erp/internal/app/service/helpers"
	"erp/pkg/db"

	"gorm.io/gen"
	"gorm.io/gen/field"
)

type SupplierRepository interface {
	GetSupplier(req *common.RequestContext, i *dto.RequestEntity) (
		dto.ResultEntity[dto.SupplierDto], error)
	GetSuppliers(req *common.RequestContext, i *dto.RequestPaginationData) (
		dto.PaginationResult[[]dto.SupplierDto], error)
	CreateSupplier(req *common.RequestContext, tx *query.QueryTx, i dto.SupplierData) (dto.SupplierDto, error)
	EditSupplier(tx *query.QueryTx,req *common.RequestContext, d dto.SupplierData) (err error)
	// EditSupplier(req *common.RequestContext, d *dto.EditSupplierRequest) (err error)
	UpdateStatus(req *common.RequestContext, d *dto.UpdateStatusWithEvent, nextState string) (err error)
}

type supplierRepository struct {
	Q         *query.Query
	convertor helpers.ConvertorHelper
}

func NewSupplierRepository(
	conn db.Connection,
	helpers *helpers.Helpers,
) SupplierRepository {
	return &supplierRepository{
		Q:         conn.GetQ(),
		convertor: helpers.Convertor,
	}
}

func (r *supplierRepository) UpdateStatus(req *common.RequestContext, d *dto.UpdateStatusWithEvent,
	nextState string) (err error) {
	supplierQ := r.Q.Supplier
	_, err = r.Q.Supplier.WithContext(req.Ctx).Where(
		supplierQ.CompanyID.Eq(req.ActiveCompany.ID),
		supplierQ.Status.Eq(d.Body.CurrentState),
		supplierQ.UUID.Eq(d.Body.PartyID),
	).UpdateSimple(
		supplierQ.Status.Value(nextState),
	)
	return
}

func (r *supplierRepository) GetSupplier(req *common.RequestContext, i *dto.RequestEntity) (
	res dto.ResultEntity[dto.SupplierDto], err error) {
	id := r.convertor.StrtoInt(i.ID)
	supplierQ := r.Q.Supplier
	groupQ := r.Q.Group
	err = supplierQ.WithContext(req.Ctx).Select(
		supplierQ.ID, supplierQ.UUID, supplierQ.Name, supplierQ.Status,
		supplierQ.GroupID, groupQ.Name.As("group"), groupQ.UUID.As("group_uuid"),
	).
		LeftJoin(groupQ, groupQ.ID.EqCol(supplierQ.GroupID)).
		Where(
			supplierQ.CompanyID.Eq(req.ActiveCompany.ID),
			supplierQ.ID.Eq(id),
		).Scan(&res.Entity)
	return res, err
}

func (r *supplierRepository) GetSuppliers(req *common.RequestContext, i *dto.RequestPaginationData) (
	dto.PaginationResult[[]dto.SupplierDto], error) {
	var (
		res   dto.PaginationResult[[]dto.SupplierDto]
		conds []gen.Condition
		order field.Expr
	)
	supplierQ := r.Q.Supplier
	builder := r.Q.WithContext(req.Ctx).Supplier

	//ADDING CONDITIONS
	conds = append(conds, supplierQ.CompanyID.Eq(req.ActiveCompany.ID))
	if i.Query != "" {
		conds = append(conds, supplierQ.Name.Like("%"+i.Query+"%"))
	}
	builder = builder.Where(conds...)

	total, err := builder.Count()
	if err != nil {
		return res, err
	}
	limit, offset := r.convertor.ToPaginationParams(i.Page, i.Size)
	orderCol, ok := r.Q.Supplier.GetFieldByName(i.OrderColumn) // maybe orderColStr == "id"
	if ok {
		if i.Order == "ASC" {
			order = orderCol.Asc()
		} else {
			order = orderCol.Desc()
		}
		builder = builder.Order(order)
	}

	err = builder.
		Select(supplierQ.ID, supplierQ.Name, supplierQ.UUID, supplierQ.CreatedAt, supplierQ.Status).
		Limit(limit).Offset(offset).Scan(&res.Results)

	res.Total = total
	return res, err
}



func (r *supplierRepository) CreateSupplier(req *common.RequestContext, tx *query.QueryTx,
	i dto.SupplierData) (res dto.SupplierDto, err error) {
	id, err := tx.Supplier.InsertParty(proto.PartyType_supplier.String())
	if err != nil {
		return
	}
	var supplier model.Supplier
	fields := i.Fields
	supplier.CompanyID = req.ActiveCompany.ID
	supplier.ID = id
	if err = r.convertor.CopyStructData(fields, &supplier); err != nil {
		return
	}

	err = tx.WithContext(req.Ctx).Supplier.Save(&supplier)
	if err != nil {
		return
	}
	if err = tx.Supplier.InsertActivity(id, req.Profile.ID, proto.ActivityType_CREATE.String(), nil); err != nil {
		return
	}
	res = dto.SupplierDtoFromModel(&supplier)
	return
}

func (r *supplierRepository)EditSupplier(tx *query.QueryTx,req *common.RequestContext, d dto.SupplierData) (err error) {
	data, err := r.convertor.DataMap(d.Fields)
	if err != nil {
		return
	}
	err = tx.Supplier.UnderlyingDB().WithContext(req.Ctx).Model(
		&model.Supplier{ID: d.ID},
	).Updates(data).Error
	if err != nil {
		return
	}
	err = tx.Supplier.InsertActivity(d.ID, req.Profile.ID, proto.ActivityType_EDIT.String(), nil)
	if err != nil {
		return
	}
	return
}
