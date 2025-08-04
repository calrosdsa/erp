package payment_rest

type PaymentPaths struct {
	Base              string
	Parties           string
	Detail            string
	PaymentAccounts string
	AssociatedActions string
	UpdateState string
	Document string 
}

func NewPaymentPaths(base string) PaymentPaths {
	return PaymentPaths{
		Base:    base,
		Parties: base + "/parties",
		Detail:  base + "/detail/{id}",
		AssociatedActions: base + "/associated-actions",
		UpdateState: base + "/update-state",
		PaymentAccounts: base + "/payment-accounts",
		Document: base + "/export/document",
	}
}
