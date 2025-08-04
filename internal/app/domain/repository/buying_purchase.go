package repository

import (
	"context"
	"erp/api/common"
	"erp/api/dto"
)

type PurchaseService interface {
	CreatePurchaseOrder(req *common.RequestContext, i *dto.CreatePurchaseOrderRequest) (
		dto.ResultEntity[dto.OrderDto], error)
}

type PurchaseRepository interface {
	CreatePurchaseOrder(ctx context.Context,req *common.RequestContext, i *dto.CreatePurchaseOrderRequest) (
		dto.ResultEntity[dto.OrderDto], error)
}