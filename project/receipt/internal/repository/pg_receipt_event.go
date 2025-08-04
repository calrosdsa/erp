package receipt_repo

import (
	"context"
	"erp/gen/db/query"
	"erp/gen/proto"
	"erp/internal/domain/event"
	"fmt"
)

type ReceiptEventRepository interface {
	OnInvoiceSubmitted(ctx context.Context, payload event.StatusInvoiceEventData) error
}

type receiptEventRepository struct {
}

func NewRecieptEventRepo() ReceiptEventRepository {
	return &receiptEventRepository{}
}

func (r *receiptEventRepository) OnInvoiceSubmitted(ctx context.Context,
	payload event.StatusInvoiceEventData) (err error) {
	invoice := payload.Invoice
	tx := payload.Tx

	fmt.Println("OnInvoiceSubmitted...")
	itemLineQ := tx.ItemLine
	invoicedItemLineQ := tx.InvoicedItemLine
	var invoiceItemLines []InvoiceItemLine
	err = itemLineQ.WithContext(ctx).
		Select(itemLineQ.ID, itemLineQ.Quantity, invoicedItemLineQ.ReceiptQuantity, itemLineQ.ItemLineReferenceID).
		Join(invoicedItemLineQ, invoicedItemLineQ.ItemLine.EqCol(itemLineQ.ID)).
		Where(
			itemLineQ.PartyID.Eq(invoice.ID),
		).Scan(&invoiceItemLines)
	fmt.Println("INVOICEITEM LINES", invoiceItemLines)

	receiptIdsToUpdated := make(map[int64]struct{})
	for _, invoiceItemLine := range invoiceItemLines {
		err = r.updateBilledQuantity(ctx, tx, &invoiceItemLine, receiptIdsToUpdated)
		if err != nil {
			return
		}
	}
	fmt.Println("RECEIPT IDS", receiptIdsToUpdated)
	for k, _ := range receiptIdsToUpdated {
		err = r.checkIfReciptIsBilled(ctx, tx, k)
		if err != nil {
			return
		}
	}
	return
}

func (r *receiptEventRepository) checkIfReciptIsBilled(ctx context.Context, tx *query.QueryTx,
	receiptID int64) (err error) {
	itemLineQ := tx.ItemLine
	itemLineReceiptQ := tx.ItemLineReceipt
	var receiptItemLines []ReceiptItemLine
	err = itemLineQ.WithContext(ctx).
		Select(itemLineQ.ID, itemLineQ.Quantity, itemLineReceiptQ.BilledQuantity, itemLineQ.ItemLineReferenceID).
		Join(itemLineReceiptQ, itemLineReceiptQ.ItemLine.EqCol(itemLineQ.ID)).
		Where(
			itemLineQ.PartyID.Eq(receiptID),
		).Scan(&receiptItemLines)
	if err != nil {
		return
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
			tx.Receipt.ID.Eq(receiptID),
		).UpdateSimple(tx.Receipt.Status.Value(proto.State_COMPLETED.String()))
		if err != nil {
			return
		}
	}

	return
}

func (r *receiptEventRepository) updateBilledQuantity(ctx context.Context, tx *query.QueryTx,
	invoiceItemLine *InvoiceItemLine, receiptIds map[int64]struct{}) (err error) {
	var receiptItemLines []ReceiptItemLine
	itemLineReceiptQ := tx.ItemLineReceipt
	receiptQ := tx.Receipt
	itemLineQ := tx.ItemLine
	err = itemLineQ.WithContext(ctx).
		Select(itemLineQ.ID, itemLineQ.Quantity, itemLineReceiptQ.BilledQuantity,
			receiptQ.ID.As("receipt_id")).
		Join(itemLineReceiptQ, itemLineReceiptQ.ItemLine.EqCol(itemLineQ.ID)).
		Join(receiptQ, receiptQ.ID.EqCol(itemLineQ.PartyID)).
		Where(
			itemLineQ.ItemLineReferenceID.Eq(invoiceItemLine.ItemLineReference),
			itemLineQ.Type.Eq(proto.ItemLineType_ITEM_LINE_RECEIPT.String()),
			receiptQ.Status.NotIn(proto.State_CANCELLED.String(), proto.State_DRAFT.String(),
				proto.State_COMPLETED.String()),
		).Scan(&receiptItemLines)
	if err != nil {
		return
	}
	if len(receiptItemLines) == 0 {
		return nil
	}
	totalQuantity := invoiceItemLine.Quantity - invoiceItemLine.ReceiptQuantity

	for _, receiptItemLine := range receiptItemLines {
		var billedQuantity int32
		billedQuantity = receiptItemLine.BilledQuantity
		allowedQnt := receiptItemLine.Quantity - receiptItemLine.BilledQuantity
		if allowedQnt > 0 && totalQuantity > 0 {
			billedQuantity = r.min(receiptItemLine.Quantity, totalQuantity+receiptItemLine.BilledQuantity)
			totalQuantity -= allowedQnt
		}
		if billedQuantity != receiptItemLine.BilledQuantity {
			_, err = itemLineReceiptQ.WithContext(ctx).Where(
				itemLineReceiptQ.ItemLine.Eq(receiptItemLine.ID),
			).UpdateSimple(itemLineReceiptQ.BilledQuantity.Value(billedQuantity))
			if err != nil {
				return err
			}
			_, ok := receiptIds[receiptItemLine.ReceiptID]
			if !ok {
				receiptIds[receiptItemLine.ReceiptID] = struct{}{}
			}
		}
		if totalQuantity <= 0 {
			break
		}
	}

	var receiptQuantity int32
	if totalQuantity <= 0 {
		receiptQuantity = invoiceItemLine.Quantity
	} else {
		receiptQuantity = invoiceItemLine.Quantity - totalQuantity
	}

	invoiceItemLine.ReceiptQuantity = receiptQuantity
	_, err = tx.InvoicedItemLine.WithContext(ctx).Where(
		tx.InvoicedItemLine.ItemLine.Eq(invoiceItemLine.ID),
	).UpdateSimple(tx.InvoicedItemLine.ReceiptQuantity.Value(receiptQuantity))
	if err != nil {
		return err
	}
	return

}

func (r *receiptEventRepository) min(x, y int32) int32 {
	if x < y {
		return x
	} else {
		return y
	}
}

type ReceiptItemLine struct {
	ID                int32
	Quantity          int32
	BilledQuantity    int32
	ItemLineReference int32
	ReceiptID         int64
}

type InvoiceItemLine struct {
	ID                int32
	Quantity          int32
	ItemLineReference int32
	ReceiptQuantity   int32
}

func (r *receiptEventRepository) getReference(ctx context.Context) {

}
