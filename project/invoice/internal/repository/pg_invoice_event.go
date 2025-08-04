package invoice_repo

import (
	"context"
	"erp/gen/db/query"
	"erp/gen/proto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain/event"
	"fmt"

	"gorm.io/gen/field"
	"gorm.io/gorm"
)

type InvoiceEventRepository interface {
	OnReceiptSubmitted(ctx context.Context, payload event.StatusReceiptEventData) error
	OnPaymentSubmitted(ctx context.Context, payload event.StatusPaymentEventData) error
	OnPaymentCancelled(ctx context.Context, payload event.StatusPaymentEventData) error
}
type invoiceEventRepository struct {
	util helpers.Util
}

func NewInvoiceEventRepo(
	helpers *helpers.Helpers,
) InvoiceEventRepository {
	return &invoiceEventRepository{
		util: helpers.Util,
	}
}

func (r *invoiceEventRepository) OnPaymentCancelled(ctx context.Context,
	payload event.StatusPaymentEventData) error {
	tx := payload.Tx
	for _, paymentRef := range payload.References {
		progressInvoice, err := tx.ProgressInvoice.WithContext(ctx).Where(
			tx.ProgressInvoice.InvoiceID.Eq(paymentRef.PartyID),
		).First()
		if err == gorm.ErrRecordNotFound {
			continue
		}
		if err != nil {
			return err
		}
		totalPaidAmount := progressInvoice.PaidAmount - paymentRef.Allocated
		progressInvoice.PaidAmount = totalPaidAmount
		_, err = tx.ProgressInvoice.WithContext(ctx).Where(
			tx.ProgressInvoice.InvoiceID.Eq(paymentRef.PartyID),
		).UpdateSimple(
			tx.ProgressInvoice.PaidAmount.Value(totalPaidAmount),
		)
		if err != nil {
			return err
		}
		invoiceStatus := ""
		if totalPaidAmount >= progressInvoice.TotalAmount {
			invoiceStatus = proto.State_PAID.String()
		} else if totalPaidAmount > 0 {
			invoiceStatus = proto.State_PARTIALLY_PAID.String()
		} else {
			invoiceStatus = proto.State_UNPAID.String()
		}
		if invoiceStatus != "" {
			_, err = tx.Invoice.WithContext(ctx).Where(
				tx.Invoice.ID.Eq(paymentRef.PartyID),
			).UpdateSimple(tx.Invoice.Status.Value(invoiceStatus))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *invoiceEventRepository) OnPaymentSubmitted(ctx context.Context, payload event.StatusPaymentEventData) error {
	fmt.Println("On payment created ...")
	tx := payload.Tx
	for _, paymentRef := range payload.References {
		progressInvoice, err := tx.ProgressInvoice.WithContext(ctx).Where(
			tx.ProgressInvoice.InvoiceID.Eq(paymentRef.PartyID),
		).First()
		if err == gorm.ErrRecordNotFound {
			continue
		}
		if err != nil {
			return err
		}
		totalPaidAmount := progressInvoice.PaidAmount + paymentRef.Allocated
		progressInvoice.PaidAmount = totalPaidAmount
		_, err = tx.ProgressInvoice.WithContext(ctx).Where(
			tx.ProgressInvoice.InvoiceID.Eq(paymentRef.PartyID),
		).UpdateSimple(
			tx.ProgressInvoice.PaidAmount.Value(totalPaidAmount),
		)
		if err != nil {
			return err
		}
		fmt.Println("TOTALS", totalPaidAmount, progressInvoice.TotalAmount)
		invoiceStatus := ""
		if totalPaidAmount >= progressInvoice.TotalAmount && progressInvoice.TotalAmount != 0 {
			invoiceStatus = proto.State_PAID.String()
		} else if totalPaidAmount > 0 {
			invoiceStatus = proto.State_PARTIALLY_PAID.String()
		}
		if invoiceStatus != "" {
			_, err = tx.Invoice.WithContext(ctx).Where(
				tx.Invoice.ID.Eq(paymentRef.PartyID),
			).UpdateSimple(tx.Invoice.Status.Value(invoiceStatus))
			if err != nil {
				return err
			}
		}
	}
	return nil
}


// func (r *invoiceEventRepository) updatePaidAmount(ctx context.Context,tx *query.QueryTx,)(error){

// }

func (r *invoiceEventRepository) getLineItemReceipts(tx *query.QueryTx, ctx context.Context, docID int64,
	docPartyType string) (
	res []ItemLineReceipt, err error) {
	var columns []field.Expr
	itemLineQ := tx.ItemLine
	fmt.Println("DOCUNENT ID",docID)
	columns = append(columns,itemLineQ.ID,itemLineQ.Quantity,itemLineQ.ItemLineReferenceID)
	builder := itemLineQ.WithContext(ctx)
	if docPartyType == proto.PartyType_purchaseReceipt.String(){
		itemLineReceiptQ := tx.ItemLineReceipt
		columns = append(columns, itemLineReceiptQ.BilledQuantity)
		builder = builder.Join(itemLineReceiptQ, itemLineReceiptQ.ItemLine.EqCol(itemLineQ.ID))
	}
	if docPartyType == proto.PartyType_deliveryNote.String(){
		deliveryLineItemQ := tx.DeliveryLineItem
		columns = append(columns, deliveryLineItemQ.BilledQuantity)
		builder = builder.Join(deliveryLineItemQ, deliveryLineItemQ.ItemLineID.EqCol(itemLineQ.ID))
	}
	
	err = builder.Select(columns...).Where(
		itemLineQ.PartyID.Eq(docID),
	).Scan(&res)
	if err != nil {
		return
	}
	return
}

func (r *invoiceEventRepository) OnReceiptSubmitted(ctx context.Context,
	payload event.StatusReceiptEventData) (err error) {
	receipt := payload.Receipt
	tx := payload.Tx
	//Get invoice id from receipt reference
	fmt.Println("OnReceiptSubmitted...")

	receiptItemLines, err := r.getLineItemReceipts(tx, ctx, receipt.ID,payload.ReceiptPartyType)
	if err != nil {
		return
	}
	fmt.Println("RECEIPT ITEM LINES", receiptItemLines)

	for i, receiptItemLine := range receiptItemLines {
		err = r.updateBilledQuantity(ctx, tx, &receiptItemLine)
		if err != nil {
			return
		}
		receiptItemLines[i] = receiptItemLine
	}

	var receiveItemsCount int32
	var billedItemsCount int32
	for _, receiptItemLine := range receiptItemLines {
		receiveItemsCount += receiptItemLine.Quantity
		billedItemsCount += receiptItemLine.BilledQuantity
	}
	fmt.Println("COUNT", receiptItemLines, billedItemsCount)
	//Update to confirm if all items listed on the receipt were billed through an invoice.
	if billedItemsCount >= receiveItemsCount {
		_, err = tx.Receipt.WithContext(ctx).Where(
			tx.Receipt.ID.Eq(receipt.ID),
		).UpdateSimple(tx.Receipt.Status.Value(proto.State_PAID.String()))
		if err != nil {
			return
		}
	}

	return
}

type ItemLineReceipt struct {
	ID                int32
	Quantity          int32
	BilledQuantity    int32
	ItemLineReferenceID int32
}

type InvoiceItemLine struct {
	ID              int32
	Quantity        int32
	ReceiptQuantity int32

	//For update paid amount
	Rate       int32
	PaidAmount int32
}

func (r *invoiceEventRepository) updateBilledQuantity(ctx context.Context, tx *query.QueryTx, receiptItemLine *ItemLineReceipt) error {
	var invoicedItemLines []InvoiceItemLine
	invoicedItemLineQ := tx.InvoicedItemLine
	invoiceQ := tx.Invoice
	itemLineQ := tx.ItemLine

	// Query for invoiced item lines
	err := itemLineQ.WithContext(ctx).
		Select(itemLineQ.ID, itemLineQ.Quantity, invoicedItemLineQ.ReceiptQuantity).
		Join(invoicedItemLineQ, invoicedItemLineQ.ItemLine.EqCol(itemLineQ.ID)).
		Join(invoiceQ, invoiceQ.ID.EqCol(itemLineQ.PartyID)).
		Where(
			itemLineQ.ItemLineReferenceID.Eq(receiptItemLine.ItemLineReferenceID),
			itemLineQ.Type.Eq(proto.ItemLineType_ITEM_LINE_INVOICE.String()),
			invoiceQ.Status.NotIn(proto.State_CANCELLED.String(), proto.State_DRAFT.String()),
		).Scan(&invoicedItemLines)

	if err != nil {
		return err
	}

	// If no invoiced item lines are found, there's nothing to update
	if len(invoicedItemLines) == 0 {
		return nil
	}

	// Calculate the total quantity to be billed
	totalQuantity := receiptItemLine.Quantity - receiptItemLine.BilledQuantity

	// Iterate over invoiced item lines
	for _, invoicedItemLine := range invoicedItemLines {
		receiptQuantity := invoicedItemLine.ReceiptQuantity
		allowedQuantity := invoicedItemLine.Quantity - invoicedItemLine.ReceiptQuantity

		// Only bill if there is available quantity and the total quantity to bill is greater than zero
		if allowedQuantity > 0 && totalQuantity > 0 {
			// Calculate the new receipt quantity (the quantity to be billed)
			receiptQuantity = r.util.Min(invoicedItemLine.Quantity, totalQuantity+invoicedItemLine.ReceiptQuantity)
			totalQuantity -= allowedQuantity // Reduce the total quantity to bill
		}

		// If the billed quantity has changed, update it in the database
		if receiptQuantity != invoicedItemLine.ReceiptQuantity {
			_, err = invoicedItemLineQ.WithContext(ctx).Where(
				invoicedItemLineQ.ItemLine.Eq(invoicedItemLine.ID),
			).UpdateSimple(invoicedItemLineQ.ReceiptQuantity.Value(receiptQuantity))
			if err != nil {
				return err
			}
		}

		// If there's no more quantity to bill, exit the loop
		if totalQuantity <= 0 {
			break
		}
	}

	// Calculate the final billed quantity (adjust for any remaining quantity)
	var billedQuantity int32
	if totalQuantity <= 0 {
		billedQuantity = receiptItemLine.Quantity
	} else {
		billedQuantity = receiptItemLine.Quantity - totalQuantity
	}

	// Update the BilledQuantity field in ItemLineReceipt
	receiptItemLine.BilledQuantity = billedQuantity
	_, err = tx.ItemLineReceipt.WithContext(ctx).Where(
		tx.ItemLineReceipt.ItemLine.Eq(receiptItemLine.ID),
	).UpdateSimple(tx.ItemLineReceipt.BilledQuantity.Value(billedQuantity))

	return err
}
