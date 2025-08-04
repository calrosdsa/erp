package item_inventory_repo

import (
	"context"
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/event"
	"erp/pkg/db"
	"fmt"
)

type ItemInventoryRepo interface {
	GetItemInventory(req *common.RequestContext, d *dto.RequestEntity) (res dto.ItemInventoryDto, err error)
	EditItemInventory(req *common.RequestContext, d dto.ItemInventoryFields) (err error)
	OnCreateItemInventory(ctx context.Context, payload event.ItemCreatedEventData) (err error)
	OnEditInventory(ctx context.Context, payload event.ItemCreatedEventData) (err error)
}

type itemInventoryRepo struct {
	Q         *query.Query
	convertor helpers.ConvertorHelper
}

func NewItemInventoryRepo(
	conn db.Connection,
	helpers *helpers.Helpers,
) ItemInventoryRepo {
	return &itemInventoryRepo{
		Q:         conn.GetQ(),
		convertor: helpers.Convertor,
	}
}

func (r *itemInventoryRepo) OnEditInventory(ctx context.Context, payload event.ItemCreatedEventData) (err error) {
	err = r.editItemInventory(payload.Tx, ctx, payload.Body.ItemInventory)
	return
}

func (r *itemInventoryRepo) OnCreateItemInventory(ctx context.Context, payload event.ItemCreatedEventData) (err error) {
	if payload.Item == nil {
		return domain.NIL_POINTER
	}
	d := payload.Body.ItemInventory
	tx := payload.Tx
	itemInventory := model.ItemInventorySetting{
		ItemID:               payload.Item.ID,
		ShelfLifeInDays:      d.ShelfLifeInDays,
		WarrantyPeriodInDays: d.WarrantyPeriodInDays,
		WeightUomID:          d.WeightUomID,
		WeightPerUnit:        d.WeightPerUnit,
		HasSerialNo:          d.HasSerialNo,
		SerialNoTemplate:     d.SerialNoTemplate,
	}
	err = tx.ItemInventorySetting.UnderlyingDB().WithContext(ctx).Save(&itemInventory).Error
	return
}

func (r *itemInventoryRepo) EditItemInventory(req *common.RequestContext, d dto.ItemInventoryFields) (err error) {
	fmt.Println("BODY", d)
	tx := r.Q.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	err = r.editItemInventory(tx, req.Ctx, d)
	if err != nil {
		return err
	}
	return
}

func (r *itemInventoryRepo) editItemInventory(tx *query.QueryTx, ctx context.Context, fields dto.ItemInventoryFields) (err error) {
	data, err := r.convertor.DataMap(fields)
	if err != nil {
		return
	}
	err = tx.ItemInventorySetting.UnderlyingDB().WithContext(ctx).Model(
		&model.ItemInventorySetting{ItemID: fields.ItemID},
	).Updates(data).Error
	if err != nil {
		return
	}

	return
}

func (r *itemInventoryRepo) GetItemInventory(req *common.RequestContext, d *dto.RequestEntity) (
	res dto.ItemInventoryDto, err error) {
	uomTQ := r.Q.UnitOfMeasureTranslation
	itemInventoryQ := r.Q.ItemInventorySetting
	err = itemInventoryQ.WithContext(req.Ctx).Select(
		itemInventoryQ.ShelfLifeInDays, itemInventoryQ.WarrantyPeriodInDays,
		itemInventoryQ.WeightPerUnit, itemInventoryQ.HasSerialNo,
		uomTQ.Name.As("weight_uom"), itemInventoryQ.SerialNoTemplate,
	).
		LeftJoin(uomTQ, uomTQ.BaseID.EqCol(itemInventoryQ.WeightUomID), uomTQ.LanguageCode.Eq(string(req.LanguageCode))).
		Where(
			itemInventoryQ.ItemID.Eq(r.convertor.StrtoInt(d.ID)),
		).Scan(&res)
	return
}
