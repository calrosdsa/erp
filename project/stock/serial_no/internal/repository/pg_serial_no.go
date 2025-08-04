package serial_no_repo

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/query"
	"erp/internal/app/service/helpers"
	"erp/pkg/db"
	"strings"

	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gen/helper"
)

type SerialNoRepository interface {
	GetSerialNo(req *common.RequestContext, d *dto.RequestEntity) (
		res dto.ResultEntity[dto.SerialNoDto], err error)
	GetSerialNos(req *common.RequestContext, d *dto.RequestSerialNos) (
		res dto.PaginationResult[[]dto.SerialNoDto], err error)
	GetSerialNoTransactions(req *common.RequestContext, d *dto.RequestSerialNoTransactions) (
		res []dto.SerialNoTransactionDto, err error)
}

type serialNoRepository struct {
	Q         *query.Query
	convertor helpers.ConvertorHelper
	dbHelper  db.DbHelper
}

func NewSerialRepository(
	db db.Connection,
	helpers *helpers.Helpers,
) SerialNoRepository {
	return &serialNoRepository{
		Q:         db.GetQ(),
		convertor: helpers.Convertor,
		dbHelper:  db.GetDbHelper(),
	}
}
func (r *serialNoRepository) GetSerialNoTransactions(req *common.RequestContext, d *dto.RequestSerialNoTransactions) (
	res []dto.SerialNoTransactionDto, err error) {
	var (
		generateSQL strings.Builder
		whereSQL    strings.Builder
		params      []interface{}
	)
	generateSQL.WriteString(`
	SELECT 
    bb.posting_date,
    bb.posting_time,
    sn_tx.qty,
    sn_tx.status,
    bb.voucher_type,
    bb.voucher_code,
	bb.batch_bundle_no,
    sn.serial_no,
    sn.valuation_rate,
    w.name AS warehouse,
	w.id AS warehouse_id,
    w.uuid AS warehouse_uuid,
    it.name AS item,
    it.id AS item_id,
    it.code AS item_code
	FROM 
		serial_no_transactions AS sn_tx
	JOIN 
		serial_nos AS sn 
		ON sn.id = sn_tx.serial_no_id
	JOIN 
		batch_bundles AS bb 
		ON bb.id = sn_tx.batch_bundle_id
	JOIN 
		ware_houses AS w 
		ON w.id = bb.warehouse_id
	JOIN 
		items AS it 
		ON it.id = sn.item_id
	`)
	params = append(params, d.FromDate, d.ToDate, req.ActiveCompany.ID)
	whereSQL.WriteString(`sn_tx.deleted_at is null and bb.posting_date between ? and ? and bb.company_id = ?`)
	if d.VoucherCode != "" {
		params = append(params, d.VoucherCode)
		whereSQL.WriteString(` and bb.voucher_code = ?`)
	}
	if d.SerialNo != "" {
		params = append(params, d.SerialNo)
		whereSQL.WriteString(` and sn.serial_no = ?`)
	}
	if d.BatchBundleNo != "" {
		params = append(params, d.BatchBundleNo)
		whereSQL.WriteString(` and bb.batch_bundle_no = ?`)
	}
	if d.ItemID != "" {
		params = append(params, r.convertor.StrtoInt(d.ItemID))
		whereSQL.WriteString(` and bb.item_id = ?`)
	}
	if d.WarehouseID != "" {
		params = append(params, r.convertor.StrtoInt(d.WarehouseID))
		whereSQL.WriteString(` and bb.warehouse_id = ?`)
	}
	helper.JoinWhereBuilder(&generateSQL, whereSQL)
	generateSQL.WriteString(`ORDER BY sn_tx.id ASC`)
	err = r.Q.SerialNoTransaction.UnderlyingDB().Raw(generateSQL.String(), params...).Scan(&res).Error
	return
}

func (r *serialNoRepository) GetSerialNo(req *common.RequestContext, d *dto.RequestEntity) (
	res dto.ResultEntity[dto.SerialNoDto], err error) {
	serialNoQ := r.Q.SerialNo
	itemQ := r.Q.Item
	batchBundleQ := r.Q.BatchBundle
	warehouseQ := r.Q.WareHouse
	err = serialNoQ.WithContext(req.Ctx).Select(
		serialNoQ.ID, serialNoQ.SerialNo, serialNoQ.Status,serialNoQ.CreatedAt,
		serialNoQ.ItemID,itemQ.Name.As("item_name"), itemQ.Code.As("item_code"),
		warehouseQ.Name.As("warehouse"),warehouseQ.UUID.As("warehouse_uuid"),
	).
		Join(batchBundleQ, batchBundleQ.ID.EqCol(serialNoQ.BatchBundleID)).
		Join(itemQ, itemQ.ID.EqCol(serialNoQ.ItemID)).
		Join(warehouseQ, warehouseQ.ID.EqCol(batchBundleQ.WarehouseID)).
		Where(batchBundleQ.CompanyID.Eq(req.ActiveCompany.ID), serialNoQ.SerialNo.Eq(d.ID)).
		Scan(&res.Entity)
	return
}

func (r *serialNoRepository) GetSerialNos(req *common.RequestContext, d *dto.RequestSerialNos) (
	res dto.PaginationResult[[]dto.SerialNoDto], err error) {
	var (
		conds []gen.Condition
		order field.Expr
	)
	serialNoQ := r.Q.SerialNo
	batchBundleQ := r.Q.BatchBundle

	builder := r.Q.WithContext(req.Ctx).SerialNo

	limit, offset := r.convertor.ToPaginationParams(d.Page, d.Size)
	orderCol, ok := r.Q.SerialNo.GetFieldByName(d.OrderColumn) // maybe orderColStr == "id"
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
	if d.SerialNo != "" {
		conds = append(conds, serialNoQ.SerialNo.Like("%"+d.Query+"%"))
	}

	builder = builder.Select(
		serialNoQ.ID, serialNoQ.SerialNo, serialNoQ.Status,
	).
		Join(batchBundleQ, batchBundleQ.ID.EqCol(serialNoQ.BatchBundleID)).
		Where(conds...)
	total, err := builder.Count()
	if err != nil {
		return res, err
	}
	err = builder.Limit(limit).Offset(offset).Scan(&res.Results)
	res.Total = total
	return
}
