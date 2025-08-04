package invoice_rest

type InvoicePaths struct {
	Base string
	Type string
	Detail string
	Purchase string
	UpdateState string 
	Document string 
}

func NewInvoicePaths(base string)InvoicePaths{
	return InvoicePaths{
		Base: base,
		Type:   base + "/{party}",
		Detail: base + "/detail/{id}",
		Purchase: base + "/purchase/{id}",
		UpdateState: base + "/update-state",
		Document: base + "/export/document",
	}
}