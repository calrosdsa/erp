package config

type Entity struct {
	Name string
	ID   uint
}

type Entities struct {
	Company       Entity
	Item          Entity
	ItemPrice     Entity
	ItemGroup     Entity
	ItemStock     Entity
	ItemAttribute Entity
	Warehouse     Entity
	Tax           Entity
	PriceList     Entity
	Role          Entity
	User          Entity
}

func GetAllEntities() *Entities {
	return &Entities{
		Company:       Entity{Name: "Company", ID: 1},
		Item:          Entity{Name: "Item", ID: 2},
		ItemPrice:     Entity{Name: "Item Price", ID: 3},
		ItemGroup:     Entity{Name: "Item Group", ID: 4},
		ItemStock:     Entity{Name: "Item Stock", ID: 5},
		ItemAttribute: Entity{Name: "Item Attributes", ID: 6},
		Warehouse:     Entity{Name: "Warehouse", ID: 7},
		Tax:           Entity{Name: "Tax", ID: 8},
		PriceList:     Entity{Name: "Price List", ID: 9},
		Role:          Entity{Name: "Role", ID: 10},
		User:          Entity{Name: "User", ID: 11},
		
	}
}
