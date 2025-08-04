package payment_repo

import (
	"context"
	"erp/gen/proto"
	"erp/internal/domain/event"
	"fmt"
)

type PaymentEventRepository interface {
	OnInvoiceCancelled(ctx context.Context, payload event.StatusInvoiceEventData) error
}

type paymentEventRepo struct {
}

func NewPaymentEventRepository() PaymentEventRepository {
	return &paymentEventRepo{}
}

func (r *paymentEventRepo) OnInvoiceCancelled(ctx context.Context, payload event.StatusInvoiceEventData) (
	err error) {
	tx := payload.Tx
	//Get payment references
	referenceQ := tx.PartyReference
	partyQ := tx.Party
	references, err := tx.PartyReference.WithContext(ctx).Select(referenceQ.PartyID).
		Join(partyQ, partyQ.ID.EqCol(referenceQ.ReferenceID),
			partyQ.PartyTypeCode.In(
				proto.PartyType_saleInvoice.String(),
				proto.PartyType_purchaseInvoice.String(),
			)).
		Where(
			referenceQ.ReferenceID.Eq(payload.Invoice.ID),
		).Find()
	paymentIds := make([]int64, len(references))
	for i, reference := range references {
		paymentIds[i] = reference.PartyID
	}
	fmt.Println("PAYMENT IDS",paymentIds)
	_,err = tx.Payment.WithContext(ctx).Where(
		tx.Payment.ID.In(paymentIds...),
	).UpdateSimple(
		tx.Payment.Status.Value(proto.State_CANCELLED.String()),
	)
	return
}
