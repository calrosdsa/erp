package documentinfo_rest

type Paths struct {
	AddressAndContact string
	AddressAndContactDetail string
	DocTerm string
	DocTermDetail string
	DocAccounting string
	DocAccountingDetail string
}

func NewPaths(base string) Paths {
	return Paths{
		AddressAndContact: base + "/address-and-contact",
		AddressAndContactDetail: base + "/address-and-contact/{id}",
		DocTerm: base + "/doc-term",
		DocTermDetail: base + "/doc-term/{id}",
		DocAccounting: base + "/doc-accounting",
		DocAccountingDetail: base + "/doc-accounting/{id}",
	}
}
