package domain

type StockMovementType string 

const (
	ADJUSTMENT StockMovementType = "ADJUSTMENT"
	ALLOCATION StockMovementType = "ALLOCATION"
	CANCELLATION StockMovementType = "CANCELLATION"
	RELEASE StockMovementType = "RELEASE"
	RETURN StockMovementType = "RETURN"
	SALE StockMovementType = "SALE"
)


type ItemType string

const (
	ITEM_TYPE          = "ITEM"
	ITEM_VARIANT_TYPE  = "ITEM_VARIANT"
	ITEM_TEMPLATE_TYPE = "ITEM_TEMPLATE"
)
