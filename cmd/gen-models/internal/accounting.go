package internal

import (
	"gorm.io/gen"
)

func AccountingModels(g *gen.Generator) []interface{} {
	// transactionLedger := g.GenerateModel("transaction_ledgers")
	transaction := g.GenerateModel("transaction_ledgers")
	ledger := g.GenerateModel("ledgers")

	payment := g.GenerateModel("payments")
	paymentReferences := g.GenerateModel("payment_references")
	cashOutflow := g.GenerateModel("cash_outflows")
	journalEntry := g.GenerateModel("journal_entries")
	journalEntryLine := g.GenerateModel("journal_entry_lines")

	ledgerAccount := g.GenerateModel("ledger_accounts")
	accountSettings := g.GenerateModel("account_settings")
	costCenter := g.GenerateModel("cost_centers")
	chargesTemplate := g.GenerateModel("charges_template")

	bank := g.GenerateModel("banks")
	bankAccount := g.GenerateModel("bank_accounts")



	return []interface{}{
		transaction,
		ledger,
		ledgerAccount,

		payment,
		paymentReferences,
		journalEntry,
		journalEntryLine,
		cashOutflow,

		accountSettings,
		costCenter,
		chargesTemplate,
		bank,
		bankAccount,
	}
}
