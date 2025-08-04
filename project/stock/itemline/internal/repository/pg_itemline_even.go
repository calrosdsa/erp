package itemline_repo

import (
	"context"
	"erp/api/dto"
	"erp/gen/db/query"
	"erp/gen/proto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain/event"
	"erp/pkg/db"
)

type ItemLineEventRepo interface {
	CreateReceiptItemLines(ctx context.Context, i event.ReceiptEventData) (err error)
	CreateInvoiceItemLines(ctx context.Context, i event.InvoiceEventData) (err error)
	CreateOrderItemLines(ctx context.Context, i event.OrderEventData) (err error)
	CreateStockEntryItemLines(ctx context.Context, i event.StockEntryEventData) (err error)
	CreateQuotationLineItems(ctx context.Context, i event.QuotationEventData) (err error)

	EditQuotationLineItems(ctx context.Context, i event.QuotationEventData) (err error)
	EditOrderLineItems(ctx context.Context, i event.OrderEventData) (err error)
	EditReceiptLineItems(ctx context.Context, i event.ReceiptEventData) (err error)
	EditInvoiceLineItems(ctx context.Context, i event.InvoiceEventData) (err error)
	EditStockEntryLineItems(ctx context.Context, i event.StockEntryEventData) (err error)

}

type itemLineEventRepo struct {
	conn         db.Connection
	Q            *query.Query
	currency     helpers.CurrencyHelper
	convertor    helpers.ConvertorHelper
	lineItemRepo ItemLineRepository
}

func NewItemLineEventRepo(
	conn db.Connection,
	helpers *helpers.Helpers,
	lineItemRepo ItemLineRepository,
) ItemLineEventRepo {
	return &itemLineEventRepo{
		conn:         conn,
		Q:            conn.GetQ(),
		currency:     helpers.Currency,
		convertor:    helpers.Convertor,
		lineItemRepo: lineItemRepo,
	}
}

func (r *itemLineEventRepo) EditQuotationLineItems(ctx context.Context, i event.QuotationEventData) (
	err error) {
	tx := i.Tx
	err = r.deleteItemLines(ctx, tx, i.Body.Quotation.ID)
	if err != nil {
		return
	}
	err = r.createItemLines(ctx, tx, i.Body.CreateItemLines,
		i.Body.Quotation.ID, proto.ItemLineType_QUOTATION_LINE_ITEM.String(), "")
	return
}

func (r *itemLineEventRepo) EditOrderLineItems(ctx context.Context, i event.OrderEventData) (
	err error) {
	tx := i.Tx
	err = r.deleteItemLines(ctx, tx, i.Body.Order.ID)
	if err != nil {
		return
	}
	err = r.createItemLines(ctx, tx, i.Body.CreateItemLines,
		i.Body.Order.ID, proto.ItemLineType_ITEM_LINE_ORDER.String(), "")
	return
}

func (r *itemLineEventRepo) EditReceiptLineItems(ctx context.Context, i event.ReceiptEventData) (
	err error) {
	tx := i.Tx
	err = r.deleteItemLines(ctx, tx, i.Body.Receipt.ID)
	if err != nil {
		return
	}
	err = r.createItemLines(ctx, tx, i.Body.CreateItemLines,
		i.Body.Receipt.ID, proto.ItemLineType_ITEM_LINE_RECEIPT.String(), "")
	return
}

func (r *itemLineEventRepo) EditInvoiceLineItems(ctx context.Context, i event.InvoiceEventData) (
	err error) {
	tx := i.Tx
	err = r.deleteItemLines(ctx, tx, i.Body.Invoice.ID)
	if err != nil {
		return
	}
	err = r.createItemLines(ctx, tx, i.Body.CreateItemLines,
		i.Body.Invoice.ID, proto.ItemLineType_ITEM_LINE_INVOICE.String(), i.Body.Invoice.InvoicePartyType,
	i.Body.Invoice.Fields.UpdateStock)
	return
}

func (r *itemLineEventRepo) EditStockEntryLineItems(ctx context.Context, i event.StockEntryEventData) (
	err error) {
	tx := i.Tx
	err = r.deleteItemLines(ctx, tx, i.StockEntryBody.StockEntry.ID)
	if err != nil {
		return
	}
	err = r.createItemLines(ctx, tx, i.StockEntryBody.CreateItemLines,
		i.StockEntryBody.StockEntry.ID, proto.ItemLineType_ITEM_LINE_STOCK_ENTRY.String(), "")
	return
}

func (r *itemLineEventRepo) CreateQuotationLineItems(ctx context.Context, i event.QuotationEventData) (
	err error) {
	tx := i.Tx
	err = r.createItemLines(ctx, tx, i.Body.CreateItemLines,
		i.Quotation.ID, proto.ItemLineType_QUOTATION_LINE_ITEM.String(), "")
	return
}

func (r *itemLineEventRepo) CreateOrderItemLines(ctx context.Context, i event.OrderEventData) (
	err error) {
	tx := i.Tx
	err = r.createItemLines(ctx, tx, i.Body.CreateItemLines,
		i.Order.ID, proto.ItemLineType_ITEM_LINE_ORDER.String(), "")
	return
}

func (r *itemLineEventRepo) CreateInvoiceItemLines(ctx context.Context, i event.InvoiceEventData) (
	err error) {
	tx := i.Tx
	err = r.createItemLines(ctx, tx, i.Body.CreateItemLines,
		i.Invoice.ID, proto.ItemLineType_ITEM_LINE_INVOICE.String(), i.InvoicePartyType,i.Invoice.UpdateStock)
	return
}

func (r *itemLineEventRepo) CreateReceiptItemLines(ctx context.Context, i event.ReceiptEventData) (
	err error) {
	tx := i.Tx
	switch i.ReceiptPartyType {
	case proto.PartyType_purchaseReceipt.String():
		err = r.createItemLines(ctx, tx, i.Body.CreateItemLines,
			i.Receipt.ID, proto.ItemLineType_ITEM_LINE_RECEIPT.String(), i.ReceiptPartyType)
	case proto.PartyType_deliveryNote.String():
		err = r.createItemLines(ctx, tx, i.Body.CreateItemLines,
			i.Receipt.ID, proto.ItemLineType_DELIVERY_LINE_ITEM.String(), i.ReceiptPartyType)
	}
	return
}

func (r *itemLineEventRepo) CreateStockEntryItemLines(ctx context.Context, i event.StockEntryEventData) (
	err error) {
	tx := i.Tx
	err = r.createItemLines(ctx, tx, i.StockEntryBody.CreateItemLines,
		i.StockEntry.ID, proto.ItemLineType_ITEM_LINE_STOCK_ENTRY.String(), "")
	return
}

func (r *itemLineEventRepo) deleteItemLines(ctx context.Context, tx *query.QueryTx, docPartyID int64) (err error) {
	err = r.lineItemRepo.DeleteLineItems(ctx, tx, docPartyID)
	return
}

func (r *itemLineEventRepo) createItemLines(ctx context.Context, tx *query.QueryTx, d dto.CreateItemLines,
	docPartyID int64, lineType string, documentType string,args ...interface{}) (err error) {
	var (
		updateStock bool
	)
	if len(args) == 1 {
		if val,ok := args[0].(bool);ok {
			updateStock = val
		}
	}
	// invoicedItemLines := make([]*model.InvoicedItemLine,0)
	for _, line := range d.Lines {
		err = r.lineItemRepo.CreateLineItem(ctx, tx, line, docPartyID, documentType,updateStock)
		if err != nil {
			return err
		}
	}
	return nil
}
