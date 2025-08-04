package receipt_rest

type ReceiptPaths struct {
	Base   string
	Detail string
	Type string 
	UpdateState string 
	Document string 
}

func NewReceiptPaths(base string) ReceiptPaths {
	return ReceiptPaths{
		Base:   base,
		Detail: base + "/detail/{id}",
		Type: base + "/{party}",
		UpdateState: base + "/update-state",
		Document: base + "/export/document",
	}
}
