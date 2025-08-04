package domain

type PartyTypes string

const (
	PARTY_COMPANY = "company"

	PARTY_SUPPLIER_GROUP = "supplierGroup"
	PARTY_ITEM_GROUP     = "itemGroup"

	PARTY_WAREHOUSE = "warehouse"

	PARTY_STOCK_LEVEL    = "stockLevel"
	PARTY_ITEM           = "item"
	PARTY_ITEM_ATTRIBUTE = "itemAttribute"
	PARTY_ITEM_PRICE     = "itemPrice"
	PARTY_SUPPLIER       = "supplier"
	PARTY_PURCHASE_ORDER = "purchaseOrder"
	// PARTY_PURCHASE_ORDER = "purchase_order"

	PARTY_TAX = "tax"

	PARTY_CUSTOMER       = "customer"
	PARTY_CUSTOMER_GROUP = "customerGroup"

	PARTY_ADMIN    = "admin"
	PARTY_EMPLOYEE = "employee"
	PARTY_CLIENT   = "client"

	PARTY_ADDRESS    = "address"
	PARTY_CONTACT    = "contact"
	PURCHASE_INVOCIE = "purchaseInvoice"

	PARTY_RECEIPT = "receipt"
)
