package repository

import (
	"erp/api/common"
	"erp/api/dto"
)

type InvoiceRepository interface {
	CreateInvoice(req *common.RequestContext, i dto.InvoiceBody) (dto.InvoiceDto, error)
	GetInvoices(req *common.RequestContext, i *dto.RequestPaginationData) (dto.PaginationResult[[]dto.InvoiceDto], error)
	GetInvoiceDetail(req *common.RequestContext, d *dto.RequestEntityWithParty) (
		dto.ResultEntity[dto.InvoiceDetailDto], error,
	)
	UpdateInvoiceState(req *common.RequestContext, id string, prevState, nexState string) (err error)
}

type InvoiceService interface {
	CreatePurchaseInvoice(req *common.RequestContext, i dto.InvoiceBody) (dto.InvoiceDto, error)
	GetPurchaseInvoices(req *common.RequestContext, i *dto.RequestPaginationData) (dto.PaginationResult[[]dto.InvoiceDto], error)
	GetPurchaseInvoiceDetail(req *common.RequestContext, d *dto.RequestEntityWithParty) (
		dto.ResultEntity[dto.InvoiceDetailDto], error,
	)
	UpdateInvoiceState(req *common.RequestContext, d *dto.UpdateStatusWithEvent) error
}
