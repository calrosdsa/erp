package order_repo

import (
	"context"
	"erp/gen/db/query"
	"erp/gen/proto"
	"erp/internal/domain/event"
	"erp/pkg/db"

	"gorm.io/gen/field"
	"gorm.io/gorm"
)

type OrderEventRepository interface {
	OnReceiptSubmitted(ctx context.Context, payload event.StatusReceiptEventData) error
	OnReceiptCancelled(ctx context.Context, payload event.StatusReceiptEventData) error
	OnInvoiceSubmitted(ctx context.Context, payload event.StatusInvoiceEventData) error
	OnInvoiceCancelled(ctx context.Context, payload event.StatusInvoiceEventData) error
}

type orderEventRepository struct {
	conn db.Connection
}

func NewOrderEventRepository(
	conn db.Connection,
) OrderEventRepository {
	return &orderEventRepository{
		conn: conn,
	}
}

func (r *orderEventRepository) OnInvoiceCancelled(ctx context.Context, payload event.StatusInvoiceEventData) (err error) {
	tx := payload.Tx
	invoice := payload.Invoice
	orderQ := tx.Order
	if invoice.DocReferenceID == nil {
		return nil
	}
	//Remove or replace for a more  secure alternative
	order, err := tx.Order.WithContext(ctx).Select(orderQ.ID).Where(
		orderQ.ID.Eq(*invoice.DocReferenceID),
	).First()
	if err == gorm.ErrRecordNotFound {
		err = nil
		return err
	}
	if err != nil {
		return err
	}
	totalAmount := payload.LineItemsData.TotalAmount + payload.TaxLinesData.TotalAmount
	_, err = tx.ProgressOrder.WithContext(ctx).Where(
		tx.ProgressOrder.OrderID.Eq(order.ID),
	).UpdateSimple(tx.ProgressOrder.BilledAmount.Sub(totalAmount))
	if err != nil {
		return err
	}
	orderType := ""
	if payload.InvoicePartyType == proto.PartyType_purchaseInvoice.String() {
		orderType = proto.PartyType_purchaseOrder.String()
	} else if payload.InvoicePartyType == proto.PartyType_saleInvoice.String() {
		orderType = proto.PartyType_saleOrder.String()
	}
	err = r.updateOrderState(ctx, tx, order.ID, orderType)
	return err
}

func (r *orderEventRepository) OnInvoiceSubmitted(ctx context.Context, payload event.StatusInvoiceEventData) error {
	tx := payload.Tx
	invoice := payload.Invoice
	orderQ := tx.Order
	if invoice.DocReferenceID == nil {
		return nil
	}
	//Remove or replace for a more  secure alternative
	order, err := tx.Order.WithContext(ctx).Select(orderQ.ID).Where(
		orderQ.ID.Eq(*invoice.DocReferenceID),
	).First()
	if err == gorm.ErrRecordNotFound {
		err = nil
		return err
	}
	if err != nil {
		return err
	}
	totalAmount := payload.LineItemsData.TotalAmount + payload.TaxLinesData.TotalAmount

	var updateValues []field.AssignExpr

	updateValues = append(updateValues, tx.ProgressOrder.BilledAmount.Add(totalAmount))
	if invoice.UpdateStock {
		updateValues = append(updateValues, tx.ProgressOrder.ReceivedItems.Add(payload.LineItemsData.TotalQuantity))
	}
	_, err = tx.ProgressOrder.WithContext(ctx).Where(
		tx.ProgressOrder.OrderID.Eq(order.ID),
	).UpdateSimple(
		updateValues...,
	)
	if err != nil {
		return err
	}
	orderType := ""
	if payload.InvoicePartyType == proto.PartyType_purchaseInvoice.String() {
		orderType = proto.PartyType_purchaseOrder.String()
	} else if payload.InvoicePartyType == proto.PartyType_saleInvoice.String() {
		orderType = proto.PartyType_saleOrder.String()
	} 
	err = r.updateOrderState(ctx, tx, order.ID, orderType)
	return err
}

func (r *orderEventRepository) OnReceiptCancelled(ctx context.Context,
	payload event.StatusReceiptEventData) (err error) {
	tx := payload.Tx
	receipt := payload.Receipt
	if receipt.DocReferenceID == nil {
		return nil
	}
	orderQ := tx.Order
	order, err := tx.Order.WithContext(ctx).Select(orderQ.ID).Where(
		orderQ.ID.Eq(*receipt.DocReferenceID),
	).First()
	if err == gorm.ErrRecordNotFound {
		err = nil
		return err
	}
	if err != nil {
		return err
	}
	_, err = tx.ProgressOrder.WithContext(ctx).Where(
		tx.ProgressOrder.OrderID.Eq(order.ID),
	).UpdateSimple(tx.ProgressOrder.ReceivedItems.Sub(payload.LineItemsData.TotalQuantity))
	if err != nil {
		return err
	}
	orderType := ""
	if payload.ReceiptPartyType == proto.PartyType_purchaseReceipt.String() {
		orderType = proto.PartyType_purchaseOrder.String()
	} else if payload.ReceiptPartyType == proto.PartyType_deliveryNote.String() {
		orderType = proto.PartyType_saleOrder.String()
	}
	err = r.updateOrderState(ctx, tx, order.ID, orderType)
	return err
}

func (r *orderEventRepository) OnReceiptSubmitted(ctx context.Context, payload event.StatusReceiptEventData) error {
	tx := payload.Tx
	receipt := payload.Receipt
	if receipt.DocReferenceID == nil {
		return nil
	}
	orderQ := tx.Order
	order, err := tx.Order.WithContext(ctx).Select(orderQ.ID).Where(
		orderQ.ID.Eq(*receipt.DocReferenceID),
	).First()
	if err == gorm.ErrRecordNotFound {
		err = nil
		return err
	}
	if err != nil {
		return err
	}
	_, err = tx.ProgressOrder.WithContext(ctx).Where(
		tx.ProgressOrder.OrderID.Eq(order.ID),
	).UpdateSimple(tx.ProgressOrder.ReceivedItems.Add(payload.LineItemsData.TotalQuantity))
	if err != nil {
		return err
	}
	orderType := ""
	if payload.ReceiptPartyType == proto.PartyType_purchaseReceipt.String() {
		orderType = proto.PartyType_purchaseOrder.String()
	} else if payload.ReceiptPartyType == proto.PartyType_deliveryNote.String() {
		orderType = proto.PartyType_saleOrder.String()
	}
	err = r.updateOrderState(ctx, tx, order.ID, orderType)
	return err
}

func (r *orderEventRepository) updateOrderState(ctx context.Context, tx *query.QueryTx, orderID int64,
	partyType string) (err error) {
	progressOrder, err := tx.ProgressOrder.WithContext(ctx).Where(
		tx.ProgressOrder.OrderID.Eq(orderID),
	).First()
	if err != nil {
		return
	}
	orderState := ""
	// if progressOrder.TotalItems < progressOrder.ReceivedItems  {
	// }
	receivedItems := progressOrder.TotalItems <= progressOrder.ReceivedItems
	billedItems := progressOrder.TotalAmount <= progressOrder.BilledAmount
	if receivedItems {
		orderState = proto.State_TO_BILL.String()
	}
	if billedItems {
		if partyType == proto.PartyType_purchaseOrder.String() {
			orderState = proto.State_TO_RECEIVE.String()
		} else if partyType == proto.PartyType_saleOrder.String() {
			orderState = proto.State_TO_DELIVER.String()
		}
	}
	if billedItems && receivedItems {
		orderState = proto.State_COMPLETED.String()
	}
	if progressOrder.ReceivedItems == 0 && progressOrder.BilledAmount == 0 {
		if partyType == proto.PartyType_purchaseOrder.String() {
			orderState = proto.State_TO_RECEIVE_AND_BILL.String()
		} else if partyType == proto.PartyType_saleOrder.String() {
			orderState = proto.State_TO_DELIVER_AND_BILL.String()
		}
	}
	if orderState == "" {
		return nil
	}
	_, err = tx.Order.WithContext(ctx).Where(
		tx.Order.ID.Eq(orderID),
	).Update(tx.Order.Status, orderState)
	return err
}
