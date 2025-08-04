package batch_bundle_repo

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

type BatchBundleRepository interface {
	GetBatchBundle(req *common.RequestContext, d *dto.RequestEntity) (
		res dto.ResultEntity[dto.BatchBundleDto], err error)
	GetBatchBundles(req *common.RequestContext, d *dto.RequestPaginationData) (
		res dto.PaginationResult[[]dto.BatchBundleDto], err error)
	CreateBatchBundle(req *common.RequestContext, d *dto.CreateBatchBundleRequest) (res model.BatchBundle, err error)
}

type batchBundleRepository struct {
	Q         *query.Query
	convertor helpers.ConvertorHelper
	dbHelper  db.DbHelper
}

func NewBatchBundleRepository(
	db db.Connection,
	helpers *helpers.Helpers,
) BatchBundleRepository {
	return &batchBundleRepository{
		Q:         db.GetQ(),
		convertor: helpers.Convertor,
		dbHelper:  db.GetDbHelper(),
	}
}

func (r *batchBundleRepository) CreateBatchBundle(req *common.RequestContext, d *dto.CreateBatchBundleRequest,
	) (res model.BatchBundle, err error) {
	tx := r.Q.Begin()

	batchBundleID,err := tx.BatchBundle.InsertParty(proto.PartyType_batchBundle.String())
	if err != nil {
		return
	}
	res.ItemID = d.Body.ItemID
	res.VoucherCode = d.Body.VoucherCode
	res.VoucherType = d.Body.VoucherType
	res.PostingDate = d.Body.PostingDate
	res.PostingTime = d.Body.PostingTime
	res.WarehouseID = d.Body.WarehouseID
	res.CompanyID = req.ActiveCompany.ID
	res.ID = batchBundleID

	// err = 

	return
}

func (r *batchBundleRepository) GetBatchBundle(req *common.RequestContext, d *dto.RequestEntity) (
	res dto.ResultEntity[dto.BatchBundleDto], err error) {
	batchBundleQ := r.Q.BatchBundle
	warehouseQ := r.Q.WareHouse
	itemQ := r.Q.Item

	err = batchBundleQ.WithContext(req.Ctx).Select(
		batchBundleQ.ID, batchBundleQ.BatchBundleNo, batchBundleQ.VoucherType,
		batchBundleQ.CreatedAt,
		batchBundleQ.ItemID,itemQ.Name.As("item"), itemQ.Code.As("item_code"),
		warehouseQ.Name.As("warehouse"), warehouseQ.UUID.As("warehouse_uuid"),
	).
		Join(itemQ, itemQ.ID.EqCol(batchBundleQ.ItemID)).
		Join(warehouseQ, warehouseQ.ID.EqCol(batchBundleQ.WarehouseID)).
		Where(batchBundleQ.CompanyID.Eq(req.ActiveCompany.ID), batchBundleQ.BatchBundleNo.Eq(d.ID)).
		Scan(&res.Entity)
	return
}

func (r *batchBundleRepository) GetBatchBundles(req *common.RequestContext, d *dto.RequestPaginationData) (
	res dto.PaginationResult[[]dto.BatchBundleDto], err error) {
	var (
		conds []gen.Condition
		order field.Expr
	)
	batchBundleQ := r.Q.BatchBundle
	builder := r.Q.WithContext(req.Ctx).BatchBundle

	limit, offset := r.convertor.ToPaginationParams(d.Page, d.Size)
	orderCol, ok := r.Q.BatchBundle.GetFieldByName(d.OrderColumn) // maybe orderColStr == "id"
	if ok {
		if d.Order == "ASC" {
			order = orderCol.Asc()
		} else {
			order = orderCol.Desc()
		}
		builder = builder.Order(order)
	}
	//ADDING CONDITIONS
	conds = append(conds, batchBundleQ.CompanyID.Eq(req.ActiveCompany.ID))
	if d.Query != "" {
		conds = append(conds, batchBundleQ.BatchBundleNo.Like("%"+d.Query+"%"))
	}
	warehouseQ := r.Q.WareHouse
	itemQ := r.Q.Item
	builder = builder.Select(
		batchBundleQ.ID, batchBundleQ.BatchBundleNo, batchBundleQ.VoucherType,
		itemQ.Name.As("item"), itemQ.Code.As("item_code"),
		warehouseQ.Name.As("warehouse"), warehouseQ.UUID.As("warehouse_uuid"),
	).
		Join(itemQ, itemQ.ID.EqCol(batchBundleQ.ItemID)).
		Join(warehouseQ, warehouseQ.ID.EqCol(batchBundleQ.WarehouseID)).
		Where(conds...)
	total, err := builder.Count()
	if err != nil {
		return res, err
	}
	err = builder.Limit(limit).Offset(offset).Scan(&res.Results)
	res.Total = total
	return
}
