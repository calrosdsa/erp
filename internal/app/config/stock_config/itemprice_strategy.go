package stockconfig

import (
	"context"
	"erp/api/common"
	"erp/gen/db/model"
)

type ItemPriceStrategy interface {
	GetItemPrice(ctx context.Context, req *common.RequestContext, itemPrice *model.ItemPrice) (string, error)
}

type DefaultItemPriceStrategy struct{}

func (d *DefaultItemPriceStrategy) GetItemPrice(ctx context.Context, req *common.RequestContext, itemPrice *model.ItemPrice) (
	res string, err error) {
	return
}
