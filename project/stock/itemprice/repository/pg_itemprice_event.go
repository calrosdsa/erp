package itemprice_repo

import (
	"context"

	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/event"
)

type ItemPriceEventRepo interface {
	OnItemCreated(ctx context.Context, payload event.ItemCreatedEventData) (err error)
	OnItemEdited(ctx context.Context, payload event.ItemCreatedEventData) (err error)
}

type itemPriceEventRepo struct {
	currency      helpers.CurrencyHelper
	itemPriceRepo ItemPriceRepository
}

func NewItemPriceEventRepo(
	helpers *helpers.Helpers,
	itemPriceRepo ItemPriceRepository,
) ItemPriceEventRepo {
	return &itemPriceEventRepo{
		currency:      helpers.Currency,
		itemPriceRepo: itemPriceRepo,
	}
}

func (r *itemPriceEventRepo) OnItemEdited(ctx context.Context, payload event.ItemCreatedEventData) (err error) {
	tx := payload.Tx
	b := payload.Body
	for _, itemPriceData := range b.ItemPrices {
		switch itemPriceData.Action {
			case string(domain.CREATE):
				itemPriceData.Fields.ItemID = payload.Item.ID
				_, err = r.itemPriceRepo.CreateItemPriceTx(tx, &payload.Req, itemPriceData)
				if err != nil {
					return
				}
			case string(domain.EDIT):
				itemPriceData.Fields.ItemID = payload.Body.ID
				err = r.itemPriceRepo.EditItemPriceTx(tx, &payload.Req, itemPriceData)
				if err != nil {
					return
				}
			case string(domain.DELETE):
				err = r.itemPriceRepo.DeleteItemPriceTx(tx, &payload.Req, itemPriceData.ID)
		}
	}
	return
}

func (r *itemPriceEventRepo) OnItemCreated(ctx context.Context, payload event.ItemCreatedEventData) (err error) {
	tx := payload.Tx
	b := payload.Body
	for _, itemPriceData := range b.ItemPrices {
		switch itemPriceData.Action {
		case string(domain.CREATE):
			itemPriceData.Fields.ItemID = payload.Item.ID
			_, err = r.itemPriceRepo.CreateItemPriceTx(tx, &payload.Req, itemPriceData)
			if err != nil {
				return
			}
		case string(domain.EDIT):
			itemPriceData.Fields.ItemID = payload.Item.ID
			err = r.itemPriceRepo.EditItemPriceTx(tx, &payload.Req, itemPriceData)
			if err != nil {
				return
			}
		case string(domain.DELETE):
			err = r.itemPriceRepo.DeleteItemPriceTx(tx, &payload.Req, itemPriceData.ID)
		}
	}
	return
}
