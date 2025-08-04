package clientconfig

import (
	"erp/api/common"
	"erp/internal/app/entity"
)

type ClientStrategy interface {
	CreateCustomer(req *common.RequestContext, d *entity.Client,metadata string) (err error)
}

type DefaultItemStrategy struct{}

func (s DefaultItemStrategy) CreateCustomer(req *common.RequestContext, d *entity.Client,metadata string) (err error) {
	return nil
}
