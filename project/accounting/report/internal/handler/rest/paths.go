package acct_report_rest

type AcctReportPaths struct {
	GeneralLedger           string
	AccountPayable          string
	AccountPayableSumary    string
	AccountReceivable       string
	AccountReceivableSumary string
	AccountBalance string 
}

func NewAcctReportPaths(base string) AcctReportPaths {
	return AcctReportPaths{
		GeneralLedger:           base + "/general",
		AccountPayable:          base + "/account-payable",
		AccountPayableSumary:    base + "/account-payable/sumary",
		AccountReceivable:       base + "/account-receivable",
		AccountReceivableSumary: base + "/account-receivable/sumary",
		AccountBalance: base + "/account-balance",
	}
}

type FinancialStatementPaths struct {
	ProfitAndLoss string
	CashFlow string
	BalanceSheet string
}

func NewFinancialStatementPaths(base string) FinancialStatementPaths {
	return FinancialStatementPaths{
		ProfitAndLoss: base + "/profit-and-loss",
		CashFlow: base + "/cash-flow",
		BalanceSheet: base + "/balance-sheet",
	}
}
