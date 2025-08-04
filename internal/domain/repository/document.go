package repository

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/query"
	"context"
	// "erp/gen/db/model"
)

type DocumentService interface {
	GetTaxLines(q *query.Query, ctx context.Context, id int64) (dto.TaxLinesData, error)
	GetLineItems(q *query.Query, req *common.RequestContext, id int64, opts ...OptionStock) (dto.LineItemsData, error)

	GenerateOrderDocumentPdf(req *common.RequestContext,d dto.OrderDto,orderPartyType string) ([]byte, error)
	GenerateInvoiceDocumentPdf(req *common.RequestContext,d dto.InvoiceDto,docPartyType string) ([]byte, error)
	GenerateReceiptDocumentPdf(req *common.RequestContext,d dto.ReceiptDto,docPartyType string) ([]byte, error)
}



type options struct {
	LoadItemInLine       bool
	LoadReceiptLineItem  bool
	LoadLineStockEntry   bool
	LoadDeliveryLineItem bool
}

type OptionStock func(o *options)

var OptionsStock options

func (_ *options) WithLoadLineStockEntry(d bool) OptionStock {
	return func(o *options) {
		o.LoadLineStockEntry = d
	}
}

func (_ *options) WithLoadItemInLine(d bool) OptionStock {
	return func(o *options) {
		o.LoadItemInLine = d
	}
}
func (_ *options) WithLoadReceiptLineItem(d bool) OptionStock {
	return func(o *options) {
		o.LoadReceiptLineItem = d
	}
}
func (_ *options) WithLoadDeliveryLineItem(d bool) OptionStock {
	return func(o *options) {
		o.LoadDeliveryLineItem = d
	}
}

func (o *options) GetLoadItemInLine() bool {
	return o.LoadItemInLine
}

func (_ *options) Apply(opts ...OptionStock) options {
	options := options{}
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(&options)
	}
	return options
}
