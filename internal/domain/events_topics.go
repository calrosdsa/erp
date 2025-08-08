package domain

const (
	//Company
	EventCompanyCreated = "company.created"
	//Activity
	ActivityCreated = "activity.created"
	//USER
	UserCreatedEvent = "user.created"
	//Auth
	PasswordResetEvent = "account.password.reset"
	//Customer
	EventCustomerCreated = "customer.created"
	EventCustomerEdited = "customer.edited"
	
	//Supplier
	SupplierCreatedEvent = "supplier.created"
	SupplierEditedEvent = "supplier.edited"
	//Payment
	PaymentCreatedEvent   = "payment.created"
	PaymentEditedEvent    = "payment.edited"
	PaymentSubmittedEvent = "payment.submit"
	PaymentCancelledEvent = "payment.cancelled"
	//ChargesTemplate
	ChargesTemplateCreatedEvent = "chargestemplate.created"

	//JournalEntry
	JournalEntrySubmittedEvent = "journalentry.submit"
	JournalEntryCancelledEvent = "journalentry.cancelled"

	//Document
	PaymentTermsTemplateCreatedEvent = "paymentTermsTemplate.created"
	PaymentTermsTemplateEditedEvent  = "paymentTermsTemplate.edited"

	//Invoice
	InvoiceCreatedEvent   = "invoice.created"
	InvoiceSubmittedEvent = "invoice.submit"
	InvoiceCancelledEvent = "invoice.cancelled"
	InvoiceEditEvent      = "invoice.edit"

	//Quotation
	QuotationSubmittedEvent = "quotation.submit"
	QuotationCreatedEvent   = "quotation.created"
	QuotationEditEvent      = "quotation.edit"

	//Order
	OrderCreatedEvent   = "order.created"
	OrderSubmittedEvent = "order.submit"
	OrderEditEvent      = "order.edit"
	//Stock
	EventStockEntryCreated   = "stockentry.created"
	StockEntrySubmittedEvent = "stockentry.submit"
	StockEntryEditEvent      = "stockEntry.edit"

	//ItemLine
	ItemLineEditedEvent = "itemline.edited"
	ItemLineDelteEvent  = "itemline.delete"
	//Item
	ItemCreatedEvent = "item.created"
	ItemEditedEvent = "item.edited"

	//Receipt
	ReceiptCreatedEvent   = "receipt.created"
	ReceiptSubmittedEvent = "receipt.submit"
	ReceiptCancelledEvent = "receipt.cancelled"
	ReceiptEditEvent      = "receipt.edit"

	//CASH OUTFLOW
	CashOutflowCreatedEvent   = "cashOutflow.created"
	CashOutflowSubmittedEvent = "cashOutflow.submitted"
	CashOutflowCancelledEvent = "cashOutflow.cancelled"
	CashOutflowEditedEvent    = "cashOutflow.edited"

	//Deal
	DealCreatedEvent = "deal.created"
	DealEditedEvent  = "deal.edited"
	//Task
	TaskCreatedEvent = "task.created"
	TaskEditedEvent  = "task.edited"
	//Party
	PartyStageChange = "party.state_change"
)
