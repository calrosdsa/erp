package item_repo

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
	"gorm.io/gorm"
)

type ItemRepository interface {
	GetItemDetail(req *common.RequestContext, d *dto.RequestEntity) (
		res dto.ResultEntity[dto.ItemDetailDto], err error)
	GetItems(req *common.RequestContext, i *dto.RequestPaginationData) (
		res dto.PaginationResult[[]dto.ItemDto], err error)
	CreateItem(tx *query.QueryTx, req *common.RequestContext, d dto.ItemData) (dto.ItemDto, error)
	UpdateItem(req *common.RequestContext, i *dto.UpdateItemRequest) error
	EditItem(tx *query.QueryTx, req *common.RequestContext, d dto.ItemData) (err error)
	UpdateStatus(req *common.RequestContext, d dto.UpdateStatusWithEvent, nextState string) (err error)
}

type itemRepository struct {
	currencyHelper helpers.CurrencyHelper
	Q              *query.Query
	DB             *gorm.DB
	conn           db.Connection
	convertor      helpers.ConvertorHelper
}

func NewItemRepository(
	conn db.Connection,
	helpers *helpers.Helpers,
) ItemRepository {
	return &itemRepository{
		currencyHelper: helpers.Currency,
		Q:              conn.GetQ(),
		DB:             conn.GetDB(),
		conn:           conn,
		convertor:      helpers.Convertor,
	}
}
func (r *itemRepository) UpdateStatus(req *common.RequestContext, d dto.UpdateStatusWithEvent, nextState string) (
	err error) {
	e := r.Q.Item
	id := r.convertor.StrtoInt(d.Body.PartyID)
	_, err = r.Q.Item.WithContext(req.Ctx).Where(
		e.CompanyID.Eq(req.ActiveCompany.ID),
		e.Status.Eq(d.Body.CurrentState),
		e.ID.Eq(id),
	).UpdateSimple(
		e.Status.Value(nextState),
	)
	return
}
func (r *itemRepository) EditItem(tx *query.QueryTx, req *common.RequestContext, d dto.ItemData) (err error) {
	data, err := r.convertor.DataMap(d.Fields)
	if err != nil {
		return
	}
	err = tx.Item.UnderlyingDB().WithContext(req.Ctx).Model(
		&model.Item{ID: d.ID},
	).Updates(data).Error
	if err != nil {
		return
	}

	err = tx.Item.InsertActivity(d.ID, req.Profile.ID, proto.ActivityType_EDIT.String(), nil)
	if err != nil {
		return
	}
	return
}

func (r *itemRepository) GetItemDetail(req *common.RequestContext, d *dto.RequestEntity) (
	res dto.ResultEntity[dto.ItemDetailDto], err error) {
	id := r.convertor.StrtoInt(d.ID)
	itemQ := r.Q.Item
	groupQ := r.Q.Group
	uomQ := r.Q.UnitOfMeasure
	uomTranslateQ := r.Q.UnitOfMeasureTranslation
	uomTQ := r.Q.UnitOfMeasureTranslation.As("w_uom")
	itemInventoryQ := r.Q.ItemInventorySetting
	err = itemQ.WithContext(req.Ctx).Select(
		itemQ.ID, itemQ.UUID, itemQ.Name, itemQ.Code, itemQ.CreatedAt,
		itemQ.ItemType, itemQ.Description, itemQ.MaintainStock, itemQ.Status,
		groupQ.ID.As("group_id"), groupQ.Name.As("group_name"), groupQ.UUID.As("group_uuid"),
		uomQ.ID.As("uom_id"), uomTranslateQ.Name.As("uom_name"), uomQ.Code.As("uom_code"),

		itemInventoryQ.ShelfLifeInDays, itemInventoryQ.WarrantyPeriodInDays,
		itemInventoryQ.WeightPerUnit, itemInventoryQ.HasSerialNo,
		itemInventoryQ.WeightUomID, uomTQ.Name.As("weight_uom"), itemInventoryQ.SerialNoTemplate,
	).
		LeftJoin(itemInventoryQ, itemInventoryQ.ItemID.EqCol(itemQ.ID)).
		LeftJoin(uomTQ, uomTQ.BaseID.EqCol(itemInventoryQ.WeightUomID), uomTQ.LanguageCode.Eq(string(req.LanguageCode))).
		LeftJoin(groupQ, groupQ.ID.EqCol(itemQ.GroupID)).
		LeftJoin(uomQ, uomQ.ID.EqCol(itemQ.UnitOfMeasureID)).
		LeftJoin(uomTranslateQ, uomTranslateQ.BaseID.EqCol(itemQ.UnitOfMeasureID),
			uomTranslateQ.LanguageCode.Eq(string(req.LanguageCode))).
		Where(
			r.Q.Item.ID.Eq(id),
			r.Q.Item.CompanyID.Eq(req.ActiveCompany.ID),
		).Scan(&res.Entity)
	return res, err
}

func (r *itemRepository) GetItems(req *common.RequestContext, i *dto.RequestPaginationData) (
	res dto.PaginationResult[[]dto.ItemDto], err error) {
	var (
		conds []gen.Condition
		order field.Expr
	)
	itemQ := r.Q.Item
	builder := r.Q.WithContext(req.Ctx).Item

	//ADDING CONDITIONS
	conds = append(conds, itemQ.CompanyID.Eq(req.ActiveCompany.ID))
	if i.Query != "" {
		conds = append(conds, itemQ.Name.Like("%"+i.Query+"%"))
	}

	limit, offset := r.convertor.ToPaginationParams(i.Page, i.Size)
	orderCol, ok := r.Q.Item.GetFieldByName(i.OrderColumn) // maybe orderColStr == "id"
	if ok {
		if i.Order == "ASC" {
			order = orderCol.Asc()
		} else {
			order = orderCol.Desc()
		}
		builder = builder.Order(order)
		// User doesn't contains orderColStr
	}
	// uomTranslateQ := r.Q.UnitOfMeasureTranslation
	builder = builder.Select(
		itemQ.ID, itemQ.UUID, itemQ.Name, itemQ.Code, itemQ.CreatedAt,
		itemQ.Status,
		itemQ.ItemType, itemQ.UnitOfMeasureID.As("uom_id"),
		// uomTranslateQ.Name.As("uom"),
	).
		// LeftJoin(uomTranslateQ, uomTranslateQ.BaseID.EqCol(itemQ.UnitOfMeasureID)).
		Where(conds...)
	total, err := builder.Count()
	if err != nil {
		return res, err
	}
	err = builder.Limit(limit).Offset(offset).Scan(&res.Results)
	res.Total = total
	return
}

func (r *itemRepository) CreateItem(tx *query.QueryTx, req *common.RequestContext, d dto.ItemData) (
	res dto.ItemDto, err error) {
	var item model.Item
	id, err := tx.Address.InsertParty(proto.PartyType_item.String())
	if err != nil {
		return
	}
	fields := d.Fields
	item.ID = id
	item.CompanyID = req.ActiveCompany.ID

	if err = r.convertor.CopyStructData(fields, &item); err != nil {
		return
	}

	err = tx.WithContext(req.Ctx).Item.Save(&item)
	if err != nil {
		return
	}
	if err = tx.Item.InsertActivity(id, req.Profile.ID, proto.ActivityType_CREATE.String(), nil); err != nil {
		return
	}
	res = dto.ItemDtoFromModel(&item)
	return res, err
}

func (s *itemRepository) UpdateItem(req *common.RequestContext, i *dto.UpdateItemRequest) error {
	var (
		err error
	)
	// item.ID = i.Body.Entity.ID
	_, err = s.Q.Item.WithContext(req.Ctx).Where(
		s.Q.Item.UUID.Eq(i.FilterID),
	).Updates(model.Item{
		ItemType: i.Body.ItemType,
		Name:     i.Body.Name,
	})
	if err != nil {
		return err
	}

	return err
}
