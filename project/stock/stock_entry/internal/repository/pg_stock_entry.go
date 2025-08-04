package stock_entry_repo

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

type StockEntryRepository interface {
	GetStockEntry(req *common.RequestContext, d *dto.RequestEntity) (
		res dto.ResultEntity[dto.StockEntryDto], err error)
	CreateStockEntry(tx *query.QueryTx, req *common.RequestContext, d dto.StockEntryBody) (
		res model.StockEntry, err error)
	GetStockEntries(req *common.RequestContext, d *dto.RequestPaginationData) (
		res dto.PaginationResult[[]dto.StockEntryDto], err error)
	UpdateStockEntryStatus(req *common.RequestContext, tx *query.QueryTx,
		id, prevState, nextState string) (res model.StockEntry, err error)
	EditStockEntry(tx *query.QueryTx, req *common.RequestContext, d dto.StockEntryBody) (err error)
}

type stockEntryRepository struct {
	Q         *query.Query
	convertor helpers.ConvertorHelper
	dbHelper  db.DbHelper
}

func NewStockEntryRepository(
	db db.Connection,
	helpers *helpers.Helpers,
) StockEntryRepository {
	return &stockEntryRepository{
		Q:         db.GetQ(),
		convertor: helpers.Convertor,
		dbHelper:  db.GetDbHelper(),
	}
}
func (r *stockEntryRepository) EditStockEntry(tx *query.QueryTx, req *common.RequestContext, d dto.StockEntryBody) (
	err error) {
	if err = r.dbHelper.DeleteReferences(req.Ctx, tx, d.StockEntry.ID); err != nil {
		return
	}
	data, err := r.convertor.DataMap(d.StockEntry.Fields)
	if err != nil {
		return
	}
	err = tx.StockEntry.UnderlyingDB().WithContext(req.Ctx).Model(
		&model.Pricing{ID: d.StockEntry.ID},
	).Updates(data).Error
	if err != nil {
		return
	}

	err = tx.StockEntry.InsertActivity(d.StockEntry.ID, req.Profile.ID, proto.ActivityType_EDIT.String(), nil)
	if err != nil {
		return
	}
	var references []*int64
	references = append(references,
		d.StockEntry.Fields.CostCenterID,
		d.StockEntry.Fields.ProjectID,
	)
	if err = r.dbHelper.InsertReferences(req.Ctx, tx, d.StockEntry.ID, references); err != nil {
		return
	}
	return
}

func (r *stockEntryRepository) UpdateStockEntryStatus(req *common.RequestContext, tx *query.QueryTx,
	id, prevState, nextState string) (res model.StockEntry, err error) {
	stockEntryQ := tx.StockEntry
	_, err = tx.StockEntry.WithContext(req.Ctx).Where(
		stockEntryQ.CompanyID.Eq(req.ActiveCompany.ID),
		stockEntryQ.Status.Eq(prevState),
		stockEntryQ.Code.Eq(id),
	).UpdateSimple(stockEntryQ.Status.Value(nextState))
	if err != nil {
		return
	}
	stockEntry, err := tx.StockEntry.WithContext(req.Ctx).Where(
		stockEntryQ.CompanyID.Eq(req.ActiveCompany.ID),
		stockEntryQ.Code.Eq(id),
	).First()
	if err != nil {
		return
	}

	return *stockEntry, err
}

func (r *stockEntryRepository) GetStockEntry(req *common.RequestContext, d *dto.RequestEntity) (
	res dto.ResultEntity[dto.StockEntryDto], err error) {
	stockEntryQ := r.Q.StockEntry
	projectQ := r.Q.Project
	costCenterQ := r.Q.CostCenter
	sWarehouseQ := r.Q.WareHouse.As("source")
	tWarehouseQ := r.Q.WareHouse.As("target")

	err = stockEntryQ.WithContext(req.Ctx).Select(
		stockEntryQ.ID, stockEntryQ.Code, stockEntryQ.Status, stockEntryQ.EntryType,
		stockEntryQ.PostingDate, stockEntryQ.PostingTime, stockEntryQ.Tz,
		stockEntryQ.Currency,
		projectQ.Name.As("project"), projectQ.ID.As("project_id"), projectQ.UUID.As("project_uuid"),
		costCenterQ.Name.As("cost_center"), costCenterQ.ID.As("cost_center_id"), costCenterQ.UUID.As("cost_center_uuid"),
		sWarehouseQ.Name.As("source_warehouse"), sWarehouseQ.ID.As("source_warehouse_id"), 
		sWarehouseQ.UUID.As("source_warehouse_uuid"),
		tWarehouseQ.Name.As("target_warehouse"), tWarehouseQ.ID.As("target_warehouse_id"), 
		tWarehouseQ.UUID.As("target_warehouse_uuid"),
	).
		LeftJoin(projectQ, projectQ.ID.EqCol(stockEntryQ.ProjectID)).
		LeftJoin(costCenterQ, costCenterQ.ID.EqCol(stockEntryQ.CostCenterID)).
		LeftJoin(sWarehouseQ,sWarehouseQ.ID.EqCol(stockEntryQ.SourceWarehouseID)).
		LeftJoin(tWarehouseQ,tWarehouseQ.ID.EqCol(stockEntryQ.TargetWarehouseID)).
		Where(stockEntryQ.CompanyID.Eq(req.ActiveCompany.ID), stockEntryQ.Code.Eq(d.ID)).
		Scan(&res.Entity)
	return
}

func (r *stockEntryRepository) CreateStockEntry(tx *query.QueryTx, req *common.RequestContext, d dto.StockEntryBody) (
	res model.StockEntry, err error) {
	res.Code = r.dbHelper.GenerateCode(r.Q.StockEntry.UnderlyingDB(), model.StockEntry{}, req.ActiveCompany.ID)
	res.CompanyID = req.ActiveCompany.ID
	err = r.dbHelper.ValidateName(tx.StockEntry.UnderlyingDB(), &res)
	if err != nil {
		return res, domain.ERROR_NAME_TAKEN
	}
	stockEntryID, err := tx.StockEntry.InsertParty(proto.PartyType_stockEntry.String())
	if err != nil {
		return
	}
	res.ID = stockEntryID
	res.Status = proto.State_DRAFT.String()
	fields := d.StockEntry.Fields
	if err = r.convertor.CopyStructData(fields, &res); err != nil {
		return
	}

	err = tx.StockEntry.WithContext(req.Ctx).Save(&res)
	if err != nil {
		return
	}
	err = tx.StockEntry.InsertActivity(res.ID, req.Profile.ID, proto.ActivityType_CREATE.String(), nil)
	if err != nil {
		return
	}

	var references []*int64
	references = append(references,
		d.StockEntry.Fields.CostCenterID,
		d.StockEntry.Fields.ProjectID,
	)
	if err = r.dbHelper.InsertReferences(req.Ctx, tx, d.StockEntry.ID, references); err != nil {
		return
	}

	return
}

func (r *stockEntryRepository) GetStockEntries(req *common.RequestContext, d *dto.RequestPaginationData) (
	res dto.PaginationResult[[]dto.StockEntryDto], err error) {
	var (
		conds []gen.Condition
		order field.Expr
	)
	stockEntryQ := r.Q.StockEntry
	builder := r.Q.WithContext(req.Ctx).StockEntry

	limit, offset := r.convertor.ToPaginationParams(d.Page, d.Size)
	orderCol, ok := r.Q.StockEntry.GetFieldByName(d.OrderColumn) // maybe orderColStr == "id"
	if ok {
		if d.Order == "ASC" {
			order = orderCol.Asc()
		} else {
			order = orderCol.Desc()
		}
		builder = builder.Order(order)
	}
	//ADDING CONDITIONS
	conds = append(conds, stockEntryQ.CompanyID.Eq(req.ActiveCompany.ID))
	if d.Query != "" {
		conds = append(conds, stockEntryQ.Code.Like("%"+d.Query+"%"))
	}

	builder = builder.Select(
		stockEntryQ.ID, stockEntryQ.Code, stockEntryQ.Status,
		stockEntryQ.EntryType, stockEntryQ.PostingDate,
	).Where(conds...)
	total, err := builder.Count()
	if err != nil {
		return res, err
	}
	err = builder.Limit(limit).Offset(offset).Scan(&res.Results)
	res.Total = total
	return
}
