package sellingconfig

import (
	"erp/api/common"
	"erp/internal/app/entity"
)

type SalesOrderStrategy interface {
	GetSalesOrderDetail(req *common.RequestContext,salesOrder *entity.SalesOrder,lines []entity.SalesItemLine) (res string, err error)
	CreateSalesOrder(req *common.RequestContext,salesOrder *entity.SalesOrder,lines []entity.SalesItemLine)(res string,err error)
}

type DefaultSalesOrderStrategy struct{}

func (s *DefaultSalesOrderStrategy) GetSalesOrderDetail(req *common.RequestContext,salesOrder *entity.SalesOrder,
	lines []entity.SalesItemLine) (res string, err error) {
	return
}

func (s *DefaultSalesOrderStrategy) CreateSalesOrder(req *common.RequestContext,o *entity.SalesOrder,lines []entity.SalesItemLine)(res  string,err error){
	return
}
