package domain

type EntityTemplate struct {
	Name string
	ID   int64
}

var (
	COMPANY          = EntityTemplate{Name: "Company", ID: 1}
	ITEM             = EntityTemplate{Name: "Item", ID: 2}
	ITEM_PRICE       = EntityTemplate{Name: "Item-Price", ID: 3}
	ITEM_GROUP       = EntityTemplate{Name: "Item-Group", ID: 4}
	ITEM_STOCK       = EntityTemplate{Name: "Item-Stock", ID: 5}
	ITEM_ATTRIBUTE   = EntityTemplate{Name: "Item-Attributes", ID: 6}
	WAREHOUSE        = EntityTemplate{Name: "Warehouse", ID: 7}
	TAX              = EntityTemplate{Name: "Tax", ID: 8}
	PRICE_LIST       = EntityTemplate{Name: "Price-List", ID: 9}
	ROLE             = EntityTemplate{Name: "Role", ID: 10}
	USER             = EntityTemplate{Name: "User", ID: 11}
	SUPPLIER         = EntityTemplate{Name: "Supplier", ID: 12}
	PURCHASE_ORDER   = EntityTemplate{Name: "Purchase-Order", ID: 14}
	CUSTOMER         = EntityTemplate{Name: "Customer", ID: 15}
	ADDRESS          = EntityTemplate{Name: "Address", ID: 16}
	CONTACT          = EntityTemplate{Name: "Contact", ID: 17}
	PURCHASE_INVOICE = EntityTemplate{Name: "Purchase-Invoice", ID: 18}
	PAYMENT          = EntityTemplate{Name: "Payment", ID: 19}
	// SELLER_GROUP   = EntityTemplate{Name: "Seller Group", ID: 13}
)
