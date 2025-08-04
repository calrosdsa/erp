package dto

type (
	RequestCashFlowStatement struct {
		AuthParams
		FromDate   string `query:"from_date" required:"true"`
		ToDate     string `query:"to_date" required:"true"`
		Currency   string `query:"currency" required:"false"`
		TimeUnit   string `query:"time_unit" required:"false"`
		CostCenter string `query:"cost_center_id" required:"false"`
		ProjectID    string `query:"project_id" required:"false"`
	}
	RequestFinancialStatement struct {
		AuthParams
		FromDate string `query:"from_date" required:"true"`
		ToDate   string `query:"to_date" required:"true"`
		Currency string `query:"currency" required:"false"`
		TimeUnit string `query:"time_unit" required:"false"`

		CostCenterID string `query:"cost_center_id" required:"false"`
		ProjectID      string `query:"project_id" required:"false"`
	}

	RequestProfitAndLossStatemnt struct {
		AuthParams
		FromDate string `query:"from_date" required:"true"`
		ToDate   string `query:"to_date" required:"true"`
		Currency string `query:"currency" required:"false"`
		TimeUnit string `query:"time_unit" required:"true"`
	}

	ProfitAndLossEntryDto struct {
		AccountType string `json:"account_type"`
		Name        string `json:"account_name"`
		PostingDate string `json:"posting_date"`
		Balance     int    `json:"balance"`
	}

	CashFlowEntryDto struct {
		AccountType     string `json:"account_type"`
		AccountName     string `json:"account_name"`
		CashFlowSection string `json:"cash_flow_section"`
		Amount          int    `json:"amount"`
	}
	BalanceSheetEntryDto struct {
		AccountType     string `json:"account_type"`
		AccountName     string `json:"account_name"`
		AccountRootType string `json:"account_root_type"`
		Debit           int    `json:"debit"`
		Credit          int    `json:"credit"`
	}
)
