package domain

type EntityTemplate struct {
	Name string
	ID   int64
}

var (
	COMPANY                   = EntityTemplate{Name: "Company", ID: 1}
	ITEM                      = EntityTemplate{Name: "Item", ID: 2}
	ITEM_PRICE                = EntityTemplate{Name: "Item-Price", ID: 3}
	ITEM_GROUP                = EntityTemplate{Name: "Item-Group", ID: 4}
	ITEM_STOCK                = EntityTemplate{Name: "Item-Stock", ID: 5}
	ITEM_ATTRIBUTE            = EntityTemplate{Name: "Item-Attributes", ID: 6}
	WAREHOUSE                 = EntityTemplate{Name: "Warehouse", ID: 7}
	TAX                       = EntityTemplate{Name: "Tax", ID: 8}
	PRICE_LIST                = EntityTemplate{Name: "Price-List", ID: 9}
	ROLE                      = EntityTemplate{Name: "Role", ID: 10}
	USER                      = EntityTemplate{Name: "User", ID: 11}
	SUPPLIER                  = EntityTemplate{Name: "Supplier", ID: 12}
	PURCHASE_ORDER            = EntityTemplate{Name: "Purchase-Order", ID: 13}
	CUSTOMER                  = EntityTemplate{Name: "Customer", ID: 14}
	ADDRESS                   = EntityTemplate{Name: "Address", ID: 15}
	CONTACT                   = EntityTemplate{Name: "Contact", ID: 16}
	PURCHASE_INVOICE          = EntityTemplate{Name: "Purchase-Invoice", ID: 17}
	PAYMENT                   = EntityTemplate{Name: "Payment", ID: 18}
	LEDGER                    = EntityTemplate{Name: "Ledger", ID: 19}
	PURCHASE_RECEIPT          = EntityTemplate{Name: "Purchase-Receipt", ID: 20}
	SALE_ORDER                = EntityTemplate{Name: "Sale-Order", ID: 24}
	SALE_INVOICE              = EntityTemplate{Name: "Sale-Invoice", ID: 25}
	PIANO_FORMS               = EntityTemplate{Name: "Piano-Forms", ID: 26}
	REGATE_CHART              = EntityTemplate{Name: "Regate-Chart", ID: 27}
	DELIVERY_NOTE             = EntityTemplate{Name: "Delivery-Note", ID: 28}
	JOURNAL_ENTRY             = EntityTemplate{Name: "Journal-Entry", ID: 29}
	COST_CENTER               = EntityTemplate{Name: "Cost-Center", ID: 30}
	PROJECT                   = EntityTemplate{Name: "Project", ID: 31}
	STOCK_ENTRY               = EntityTemplate{Name: "Stock-Entry", ID: 32}
	GENERAL_LEDGER            = EntityTemplate{Name: "General-Ledger", ID: 33}
	ACCOUNT_RECEIVABLE        = EntityTemplate{Name: "Account-Receivable", ID: 34}
	ACCOUNT_PAYABLE           = EntityTemplate{Name: "Account-Payable", ID: 35}
	FINANCIAL_STATEMENTS      = EntityTemplate{Name: "Financial-Statements", ID: 36}
	STOCK_SETTING             = EntityTemplate{Name: "Stock-Setting", ID: 37}
	SERIAL_NO                 = EntityTemplate{Name: "Serial-No", ID: 38}
	BATCH_BUNDLE              = EntityTemplate{Name: "Batch-Bundle", ID: 39}
	SUPPLIER_QUOTATION        = EntityTemplate{Name: "Supplier-Quotation", ID: 40}
	QUOTATION                 = EntityTemplate{Name: "Quotation", ID: 41}
	CHARGES_TEMPLATE          = EntityTemplate{Name: "Charges-Template", ID: 42}
	CURRENCY_EXCHANGE         = EntityTemplate{Name: "Currency-Exchange", ID: 43}
	PURCHASE_RECORD           = EntityTemplate{Name: "Purchase-Record", ID: 44}
	SALES_RECORD              = EntityTemplate{Name: "Sales-Record", ID: 45}
	STOCK_LEDGER              = EntityTemplate{Name: "Stock-Ledger", ID: 46}
	MODULE                    = EntityTemplate{Name: "Module", ID: 47}
	STOCK_BALANCE             = EntityTemplate{Name: "Stock-Balance", ID: 48}
	SERIALNO_RESUME           = EntityTemplate{Name: "SerialNo-Resume", ID: 49}
	INCOME_STATEMENT          = EntityTemplate{Name: "Income-Statement", ID: 50}
	CASH_FLOW                 = EntityTemplate{Name: "Cash-Flow", ID: 51}
	BALANCE_SHEET             = EntityTemplate{Name: "Balance-Sheet", ID: 52}
	ACCOUNT_RECEIVABLE_SUMARY = EntityTemplate{Name: "Account-Receivable-Sumary", ID: 53}
	ACCOUNT_PAYABLE_SUMARY    = EntityTemplate{Name: "Account-Payable-Sumary", ID: 54}
	SUPPLIER_GROUP            = EntityTemplate{Name: "Supplier-Group", ID: 55}
	CUSTOMER_GROUP            = EntityTemplate{Name: "Customer-Group", ID: 56}
	PRICING                   = EntityTemplate{Name: "Pricing", ID: 57}
	TERMS_AND_CONDITIONS      = EntityTemplate{Name: "Terms and Conditions", ID: 58}
	PAYMENT_TERMS             = EntityTemplate{Name: "Payment Terms", ID: 59}
	PAYMENT_TERMS_TEMPLATE    = EntityTemplate{Name: "Payment Terms Template", ID: 60}
	BANK                      = EntityTemplate{Name: "Payment Terms Template", ID: 61}
	BANK_ACCOUNT              = EntityTemplate{Name: "Payment Terms Template", ID: 62}
	CASH_OUTFLOW              = EntityTemplate{Name: "Cash Outflow", ID: 63}
	DEAL                      = EntityTemplate{Name: "Deal", ID: 64}
	STAGE                     = EntityTemplate{Name: "Stage", ID: 65}
	CHAT                      = EntityTemplate{Name: "Chat", ID: 66}
	BOOKING_SCHEDULE                      = EntityTemplate{Name: "Calendario Reservas", ID: 67}
	WORKSPACE = EntityTemplate{Name: "Workspace",ID: 68}
	// SELLER_GROUP   = EntityTemplate{Name: "Seller Group", ID: 13}
)
