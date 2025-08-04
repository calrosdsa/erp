package ledger_rest

type LedgerPaths struct {
	Base string 
	Detail string
	GeneralLedger string
	TreeView string
}
func NewLedgerPaths(base string) LedgerPaths{
	return LedgerPaths{
		Base:base,
		Detail: base + "/detail/{id}",
		GeneralLedger: base + "/general",
		TreeView:base + "/view/tree",
	}
}