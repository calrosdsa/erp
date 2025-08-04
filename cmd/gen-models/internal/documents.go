package internal

import (
	"gorm.io/gen"
)

func DocumentModels(g *gen.Generator) []interface{} {
	invoice := g.GenerateModel("invoices")
	progressInvoice := g.GenerateModel("progress_invoices")
	receipt := g.GenerateModel("receipts")
	order := g.GenerateModel("orders")
	quotation := g.GenerateModel("quotations")
	taxAndChargeLine := g.GenerateModel("tax_and_charge_lines")
	progressOrder := g.GenerateModel("progress_orders")

	paymentTerm := g.GenerateModel("payment_terms")
	paymentTermsTemplate:= g.GenerateModel("payment_terms_templates")
	termsAndConditions := g.GenerateModel("terms_and_conditions")
	paymentTermsLine := g.GenerateModel("payment_terms_lines")

	addressAndContact := g.GenerateModel("address_and_contacts")
	docTerm := g.GenerateModel("doc_terms")
	docAccount := g.GenerateModel("doc_accounts")

	return []interface{}{
		order,
		invoice,
		progressInvoice,
		receipt,
		progressOrder,
		quotation,
		taxAndChargeLine,
		paymentTerm,
		paymentTermsTemplate,
		termsAndConditions,
		paymentTermsLine,
		addressAndContact,
		docTerm,
		docAccount,
	}
}
