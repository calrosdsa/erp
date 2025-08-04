package repository

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/internal/domain"
)

type OrderRepository interface {
	CreateOrder(req *common.RequestContext, i *dto.CreateOrderRequest) (
		dto.ResultEntity[dto.OrderDto], error,
	)
	GetOrder(req *common.RequestContext, i *dto.RequestEntityWithParty) (dto.ResultEntity[dto.OrderDto], error)
	GetOrders(req *common.RequestContext, d *dto.RequestPaginationPartyData) (
		dto.PaginationResult[[]dto.OrderDto], error)
}

type OrderService interface {
	CreateOrder(req *common.RequestContext, i *dto.CreateOrderRequest) (
		dto.ResultEntity[dto.OrderDto], error,
	)
	GetOrder(req *common.RequestContext, i *dto.RequestEntityWithParty) (dto.ResultEntity[dto.OrderDto], error)
	GetOrders(req *common.RequestContext, d *dto.RequestPaginationPartyData) (
		dto.PaginationResult[[]dto.OrderDto], error)
	GetEntityOrder(partyCode string) (domain.EntityTemplate, error)
}
