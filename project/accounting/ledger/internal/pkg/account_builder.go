package pkg_ledger

import (
	"context"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/gen/proto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain/event"
	"fmt"
)

type AccountBuilder struct {
	locale helpers.Locale
}

func NewAccountBuilder(
	locale helpers.Locale,
) *AccountBuilder {
	return &AccountBuilder{
		locale: locale,
	}
}

func (r *AccountBuilder) CreateChartofAccount(ctx context.Context, tx *query.QueryTx, payload event.CreatedCompanyEventData) (err error) {
	// lang := payload.LanguageCode
	lang := "es"
	asset := r.BuildChartAssetAccounts(lang)
	liability := r.BuildChartLiabilitiesAccounts(lang)
	income := r.BuildChartIncomeAccounts(lang)
	expense := r.BuildChartExpenseAccounts(lang)
	accounts, err := r.InsertAccountTree(ctx, tx, []Ledger{asset, liability, income, expense}, payload.Company.ID, []string{
		proto.AccountType_ASSET.String(), proto.AccountType_LIABILITY.String(),
		proto.AccountType_REVENUE.String(), proto.AccountType_EXPENSE.String(),
	})
	if err != nil {
		return
	}
	fmt.Println("ACCOUNT", len(accounts), accounts)
	err = tx.Ledger.WithContext(ctx).CreateInBatches(accounts, len(accounts))
	if err != nil {
		return
	}
	err = r.CompanyDefaultAccounts(ctx, tx, payload.Company.ID, accounts)
	return
}

func (r *AccountBuilder) CompanyDefaultAccounts(ctx context.Context, tx *query.QueryTx,
	companyID int64, accounts []*model.Ledger) (err error) {
	accountSetting := model.AccountSetting{}
	stockAccounts := model.StockSetting{}
	for _, account := range accounts {
		if account.AccountType != nil {

			switch *account.AccountType {
			case proto.AccountType_BANK.String():
				accountSetting.BankAccount = account.ID
			case proto.AccountType_CASH.String():
				accountSetting.CashAccunt = account.ID
			case proto.AccountType_PAYABLE.String():
				accountSetting.PayableAccount = account.ID
			case proto.AccountType_COST_OF_GOODS_SOLD.String():
				accountSetting.CostOfGoodSoldAccount = account.ID
			case proto.AccountType_RECEIVABLE.String():
				accountSetting.ReceivableAccount = account.ID
			case proto.AccountType_SALES_REVENUE.String():
				accountSetting.IncomeAccount = account.ID
			case proto.AccountType_INVENTORY.String():
				stockAccounts.InventoryAccount = account.ID
			case proto.AccountType_ASSET_RECEIVED_BUT_NOT_BILLED.String():
				stockAccounts.StockReceivedButNotBilled = account.ID
			case proto.AccountType_STOCK_ADJUSTMENT.String():
				stockAccounts.StockAdjustment = account.ID
			}
		}

	}
	stockAccounts.CompanyID = companyID
	accountSetting.CompanyID = companyID
	err = tx.WithContext(ctx).AccountSetting.Save(&accountSetting)
	if err != nil {
		return
	}
	err = tx.WithContext(ctx).StockSetting.Save(&stockAccounts)
	if err != nil {
		return
	}
	return
}

func (r *AccountBuilder) InsertAccountTree(ctx context.Context, tx *query.QueryTx, d []Ledger,
	companyID int64, rootAccounts []string) ([]*model.Ledger, error) {
	var accounts []*model.Ledger

	var insertAccount func(account Ledger, rooAccount string) error
	insertAccount = func(account Ledger, rootAccount string) (err error) {
		partyID, err := tx.Ledger.WithContext(ctx).InsertParty(proto.PartyType_ledger.String())
		if err != nil {
			return
		}
		ledger := &model.Ledger{}
		ledger.ID = partyID
		ledger.AccountRootType = rootAccount
		ledger.CompanyID = companyID
		if account.AccountType != "" {
			ledger.AccountType = &account.AccountType
		}
		ledger.ReportType = account.ReportType
		ledger.CashFlowSection = account.CashFlowSection
		ledger.Name = account.Name
		ledger.IsGroup = account.IsGroup
		ledger.LedgerParent = account.ParentID
		accounts = append(accounts, ledger)
		if len(account.Childs) == 0 {
			return nil
		}
		for _, l := range account.Childs {
			l.ParentID = &ledger.ID
			err = insertAccount(*l, rootAccount)
			if err != nil {
				return
			}
		}
		return
	}
	for i, ledger := range d {
		err := insertAccount(ledger, rootAccounts[i])
		if err != nil {
			return accounts, err
		}
	}
	return accounts, nil
}

func (r *AccountBuilder) BuildChartAssetAccounts(lang string) Ledger {
	t := r.locale.Translate(lang)
	balanceSheet := proto.FinacialReport_BALANCE_SHEET.String()
	cashOperating := proto.CashFlowSection_OPERATING.String()
	asset := Ledger{
		Name:    t("Account.Asset"),
		IsGroup: true,
		Childs: []*Ledger{
			{
				Name:    t("Account.CurrentAsset"),
				IsGroup: true,
				Childs: []*Ledger{
					{
						Name:    t("Account.CashInHand"),
						IsGroup: true,
						Childs: []*Ledger{
							{
								Name:        t("Account.Cash"),
								AccountType: proto.AccountType_CASH.String(),
								ReportType:  &balanceSheet,
							},
						},
					},
					{
						Name:    t("Account.CashAtBank"),
						IsGroup: true,
						Childs: []*Ledger{
							{
								Name:        t("Account.Bank"),
								AccountType: proto.AccountType_BANK.String(),
								ReportType:  &balanceSheet,
							},
						},
					},
					{
						Name:    t("Account.AccountsReceivable"),
						IsGroup: true,
						Childs: []*Ledger{
							{
								Name:            t("Account.Debtors"),
								AccountType:     proto.AccountType_RECEIVABLE.String(),
								ReportType:      &balanceSheet,
								CashFlowSection: &cashOperating,
							},
						},
					},
					{
						Name:    t("Account.StockAssets"),
						IsGroup: true,
						Childs: []*Ledger{
							{
								Name:            t("Account.Inventory"),
								AccountType:     proto.AccountType_INVENTORY.String(),
								ReportType:      &balanceSheet,
								CashFlowSection: &cashOperating,
							},
						},
					},
				},
			},
		},
	}
	return asset
}
func (r *AccountBuilder) BuildChartLiabilitiesAccounts(lang string) Ledger {
	t := r.locale.Translate(lang)
	balanceSheet := proto.FinacialReport_BALANCE_SHEET.String()
	operating := proto.CashFlowSection_OPERATING.String()
	asset := Ledger{
		Name:    t("Account.Liabilities"),
		IsGroup: true,
		Childs: []*Ledger{
			{
				Name:    t("Account.CurrentLiabilities"),
				IsGroup: true,
				Childs: []*Ledger{
					{
						Name:    t("Account.AccountsPayable"),
						IsGroup: true,
						Childs: []*Ledger{
							{
								Name:        t("Account.Creditors"),
								AccountType: proto.AccountType_PAYABLE.String(),
								ReportType: &balanceSheet,
								CashFlowSection: &operating,
							},
						},
					},
					{
						Name:    t("Account.StockLiabilities"),
						IsGroup: true,
						Childs: []*Ledger{
							{
								Name:        t("Account.StockReceivedButNotBilled"),
								AccountType: proto.AccountType_ASSET_RECEIVED_BUT_NOT_BILLED.String(),
								ReportType: &balanceSheet,
							},
						},
					},
				},
			},
		},
	}
	return asset
}

func (r *AccountBuilder) BuildChartIncomeAccounts(lang string) Ledger {
	t := r.locale.Translate(lang)
	reportType := proto.FinacialReport_PROFIT_AND_LOSS.String()
	asset := Ledger{
		Name:    t("Account.Income"),
		IsGroup: true,
		Childs: []*Ledger{
			{
				Name:    t("Account.DirectIncome"),
				IsGroup: true,
				Childs: []*Ledger{
					{
						Name:        t("Account.Sales"),
						AccountType: proto.AccountType_SALES_REVENUE.String(),
						ReportType: &reportType,
					},
				},
			},
		},
	}
	return asset
}

func (r *AccountBuilder) BuildChartExpenseAccounts(lang string) Ledger {
	t := r.locale.Translate(lang)
	reportType := proto.FinacialReport_PROFIT_AND_LOSS.String()
	asset := Ledger{
		Name:    t("Account.Expense"),
		IsGroup: true,
		Childs: []*Ledger{
			{
				Name:    t("Account.DirectExpense"),
				IsGroup: true,
				Childs: []*Ledger{
					{
						Name:    t("Account.StockExpenses"),
						IsGroup: true,
						Childs: []*Ledger{
							{
								Name:        t("Account.CostOfGoodsSold"),
								AccountType: proto.AccountType_COST_OF_GOODS_SOLD.String(),
								ReportType: &reportType,
							},
							{
								Name:        t("Account.StockAdjustment"),
								AccountType: proto.AccountType_STOCK_ADJUSTMENT.String(),
							},
						},
					},
				},
			},
			{
				Name:    t("Account.IndirectExpenses"),
				IsGroup: true,
				Childs: []*Ledger{
					{
						Name:        t("Account.RoundedOff"),
						AccountType: proto.AccountType_OPERATING_EXPENSES.String(),
					},
				},
			},
		},
	}
	return asset
}

type Ledger struct {
	Childs          []*Ledger
	ParentID        *int64
	Name            string
	AccountRootType string
	AccountType     string
	ReportType      *string
	CashFlowSection *string
	IsGroup         bool
}
