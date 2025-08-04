package repository

import (
	"erp/api/common"
	"erp/api/dto"
)

type CustomerRepository interface {
	CreateCustomer(req *common.RequestContext, i *dto.CustomerFields) error
	GetCustomers(req *common.RequestContext, i *dto.RequestPaginationData) (dto.PaginationResult[[]dto.CustomerDto], error)
	GetCustomerDetail(req *common.RequestContext, id *dto.RequestEntity) (dto.ResultEntity[dto.CustomerDto], error)
	GetCustomerTypes() []string
}

type ConvertorHelper interface {
	StrtoInt(s string) int64
	ToPaginationParams(page, size string) (limit, offset int)
}
