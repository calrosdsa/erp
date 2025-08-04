package rest_quotation

type QuotationPaths struct {
	Base         string
	Detail       string
	QuotationParty         string
	UpdateStatus string
}

func NewQuotationPaths(base string) QuotationPaths {
	return QuotationPaths{
		Base:         base,
		QuotationParty:         base + "/{party}",
		Detail:       base + "/detail/{id}",
		UpdateStatus: base + "/update-status",
	}
}
