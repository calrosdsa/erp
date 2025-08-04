package stockconfig

import (
	"erp/api/common"
	"erp/gen/db/model"
)

type ItemStrategy interface {
	CreateItem(req *common.RequestContext, item *model.Item, itemPrice *model.ItemPrice) (err error)
}


type DefaultItemStrategy struct{}

func (s DefaultItemStrategy) CreateItem(req *common.RequestContext,item *model.Item, itemPrice *model.ItemPrice) (err error) {
	return nil
}
