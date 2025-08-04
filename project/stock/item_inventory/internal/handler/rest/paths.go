package item_inventory_rest 

type Paths struct {
	InventorySettingsDetail string
	InventorySettings string
}

func NewPaths(base string)Paths{
	return Paths{
		InventorySettingsDetail: base + "/inventory-setting/{id}",
		InventorySettings: base + "/inventory-setting",
	}
}