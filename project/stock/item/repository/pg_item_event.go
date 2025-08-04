package item_repo

import (
	"context"
	"erp/internal/domain/event"
	// "github.com/samber/lo"
)

type ItemEventRepository interface {
	OnPurchaseOrderCreated(ctx context.Context, d event.OrderEventData) (err error)
}

type itemEventRepository struct {
}

func NewItemEventRepository() ItemRepository {
	return &itemRepository{}
}

func (r *itemEventRepository) OnPurchaseOrderCreated(ctx context.Context, d event.OrderEventData) (err error) {

	return
}

func (r *itemEventRepository) proccessReferencesToPurchaseOrder(ctx context.Context, d event.OrderEventData,
	orderID int64) (
	err error) {
	
	return
}
