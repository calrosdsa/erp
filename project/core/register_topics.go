package core

import (
	"erp/internal/domain"
	"erp/pkg/bus"
)

func RegisterEventTopics(
	bus bus.Bus,
) {
	bus.RegisterTopics(
		//Activity
		domain.ActivityCreated,
		//User
		domain.UserCreatedEvent,
		//Auth
		domain.PasswordResetEvent,
		//Customer
		domain.EventCustomerCreated,
		domain.EventCustomerEdited,
		//Supplier
		domain.SupplierCreatedEvent,
		domain.SupplierEditedEvent,
		//Company
		domain.EventCompanyCreated,
		//Payment
		domain.PaymentCreatedEvent,
		domain.PaymentEditedEvent,
		domain.PaymentSubmittedEvent,
		domain.PaymentCancelledEvent,
		//ChargesTemplate
		domain.ChargesTemplateCreatedEvent,
		//Document
		domain.PaymentTermsTemplateCreatedEvent,
		domain.PaymentTermsTemplateEditedEvent,

		//Receipt
		domain.ReceiptCreatedEvent,
		domain.ReceiptSubmittedEvent,
		domain.ReceiptCancelledEvent,
		domain.ReceiptEditEvent,

		// Invoice
		domain.InvoiceSubmittedEvent,
		domain.InvoiceCreatedEvent,
		domain.InvoiceCancelledEvent,
		domain.InvoiceEditEvent,
		//Quotation
		domain.QuotationCreatedEvent,
		domain.QuotationSubmittedEvent,
		domain.QuotationEditEvent,

		//JournalEntry
		domain.JournalEntrySubmittedEvent,
		domain.JournalEntryCancelledEvent,

		//Cash outflow
		domain.CashOutflowCancelledEvent,
		domain.CashOutflowSubmittedEvent,
		domain.CashOutflowCreatedEvent,
		domain.CashOutflowEditedEvent,

		// Order
		domain.OrderCreatedEvent,
		domain.OrderEditEvent,
		//Stock
		domain.EventStockEntryCreated,
		domain.StockEntrySubmittedEvent,
		domain.StockEntryEditEvent,

		//Item
		domain.ItemCreatedEvent,
		domain.ItemEditedEvent,

		//DEAL
		domain.DealCreatedEvent,
		domain.DealEditedEvent,
		domain.PartyStageChange,
	)
}
