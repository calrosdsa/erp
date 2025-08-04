package dto

type (
	EditItemInventoryRequest struct {
		Body ItemInventoryFields
	}

	ItemInventoryFields struct {
		ItemID               int64   `json:"item_id" required:"false"`
		ShelfLifeInDays      *int32  `json:"shelf_life_in_days" required:"false"`
		WarrantyPeriodInDays *int32  `json:"warranty_period_in_days" required:"false"`
		HasSerialNo          *bool   `json:"has_serial_no" required:"false"`
		SerialNoTemplate     *string `json:"serial_no_template" required:"false"`
		WeightUomID          *int64  `json:"weight_uom_id" required:"false"`
		WeightPerUnit        *int32  `json:"weight_per_unit" required:"false"`
	}

	ItemInventoryDto struct {
		ShelfLifeInDays      *int32 `json:"shelf_life_in_days"`
		WarrantyPeriodInDays *int32 `json:"warranty_period_in_days"`

		HasSerialNo      *bool   `json:"has_serial_no"`
		SerialNoTemplate *string `json:"serial_no_template"`
		WeightUomID      *int64  `json:"weight_uom_id"`
		WeightUom        *string `json:"weight_uom"`
		WeightPerUnit    *int32  `json:"weight_per_unit"`
	}
)
