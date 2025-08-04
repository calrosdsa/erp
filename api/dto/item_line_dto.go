package dto

type (
	ProductListDataRequest struct {
		Body ProductListData
	}

	ProductListData struct {
		PartyID   int64          `json:"party_id"`
		PartyType string         `json:"party_type"`
		Lines     []LineItemData `json:"lines" required:"true"`
	}

	EditLineItemRequest struct {
		AuthParams
		Body struct {
			ID               int32                 `json:"id" required:"true"`
			TotalAmountItems float64               `json:"total_amount_items" required:"false"`
			TotalAmount      float64               `json:"total_amount" required:"false"`
			TotalItems       int32                 `json:"total_items"`
			DocPartyID       int64                 `json:"doc_party_id" required:"false"`
			DocPartyType     string                `json:"doc_party_type" required:"false"`
			UpdateStock      bool                  `json:"update_stock" required:"false"`
			LineItemData     LineItemData          `json:"line_item_data" required:"true"`
			Charges          []TaxAndChargeLineDto `json:"charges" required:"false"`
		}
	}
	DeleteLineItemRequest struct {
		Body struct {
			ID               int32                 `json:"id" required:"true"`
			TotalAmountItems float64               `json:"total_amount_items" required:"false"`
			TotalAmount      float64               `json:"total_amount" required:"false"`
			TotalItems       int32                 `json:"total_items"`
			DocPartyID       int64                 `json:"doc_party_id" required:"false"`
			DocPartyType     string                `json:"doc_party_type" required:"false"`
			UpdateStock      bool                  `json:"update_stock" required:"false"`
			Charges          []TaxAndChargeLineDto `json:"charges" required:"false"`
		}
	}
	AddLineItemRequest struct {
		Body struct {
			TotalAmountItems float64               `json:"total_amount_items" required:"false"`
			TotalAmount      float64               `json:"total_amount" required:"false"`
			TotalItems       int32                 `json:"total_items" required:"false"`
			DocPartyID       int64                 `json:"doc_party_id" required:"false"`
			DocPartyType     string                `json:"doc_party_type" required:"false"`
			UpdateStock      bool                  `json:"update_stock" required:"false"`
			LineItemData     LineItemData          `json:"line_item_data" required:"true"`
			Charges          []TaxAndChargeLineDto `json:"charges" required:"false"`
		}
	}

	//Use for item table
	LineItemData struct {
		ID              int32   `json:"id" required:"false"`
		ItemID          int64   `json:"item_id" required:"true"`
		Rate            float64 `json:"rate" required:"true"`
		Quantity        int32   `json:"quantity" required:"true"`
		UnitOfMeasureID int64   `json:"unit_of_measure_id" required:"true"`

		LineItemReceipt     LineItemReceiptData    `json:"line_receipt" required:"false"`
		DeliveryLineItem    DeliveryLineItemData   `json:"delivery_line_item" required:"false"`
		LineItemStockEntry  LineItemStockEntryData `json:"line_stock_entry" required:"false"`
		ItemLineReferenceID *int32                 `json:"item_line_reference_id" required:"false"`
		LineType            string                 `json:"line_type" required:"true"`
	}

	LineItemReceiptData struct {
		AcceptedQuantity int32 `json:"accepted_quantity" required:"true"`
		RejectedQuantity int32 `json:"rejected_quantity" required:"true"`

		AcceptedWarehouse int64 `json:"accepted_warehouse" required:"false"`
		RejectedWarehouse int64 `json:"rejected_warehouse" required:"false"`
	}
	DeliveryLineItemData struct {
		SourceWarehouse int64 `json:"source_warehouse" required:"true"`
	}

	LineItemStockEntryData struct {
		SourceWarehouse int64 `json:"source_warehouse" required:"false"`
		TargetWarehouse int64 `json:"target_warehouse" required:"false"`
	}

	LineItemDto struct {
		ID              int32  `json:"id"`
		Rate            int32  `json:"rate"`
		Quantity        int32  `json:"quantity"`
		LineType        string `json:"line_type"`
		ItemID          int64  `json:"item_id"`
		UnitOfMeasureID int64  `json:"unit_of_measure_id"`

		// ItemPriceID int64  `json:"item_price_id"`
		ItemName        string `json:"item_name"`
		ItemCode        string `json:"item_code"`
		ItemDescription string `json:"item_description"`
		Uom             string `json:"uom"`

		ItemLineReferenceID *int32 `json:"item_line_reference_id"`

		LineItemReceiptDto
		DeliveryLineItemDto
		// LineItemReceiptDto    `json:"line_receipt" required:"false"`
		// LineItemStockEntryDto `json:"line_stock_entry" required:"false"`
	}

	LineItem struct {
		ID       int32
		Rate     int64
		Quantity int32
		//Item Data
		ItemID          int64
		ItemName        string
		ItemDescription string
		UnitOfMeasureID int64
		MaintainStock   bool
		//Receipt
		AcceptedWarehouse int64
		AcceptedQuantity  int32

		//Delivery
		SourceWarehouseID int64
		//Stock entry
		// SourceWarehouseID int64
		TargetWarehouseID int64
	}

	LineItemsData struct {
		TotalAmount   int64
		TotalQuantity int32
		LineItems     []LineItem
	}

	DeliveryLineItemDto struct {
		SourceWarehouse   string `json:"source_warehouse"`
		SourceWarehouseID int64  `json:"source_warehouse_id"`
	}

	LineItemReceiptDto struct {
		AcceptedQuantity    int32  `json:"accepted_quantity"`
		RejectedQuantity    int32  `json:"rejected_quantity"`
		AcceptedWarehouse   string `json:"accepted_warehouse"`
		RejectedWarehouse   string `json:"rejected_warehouse" required:"false"`
		AcceptedWarehouseID int64  `json:"accepted_warehouse_id"`
		RejectedWarehouseID int64  `json:"rejected_warehouse_id" required:"false"`
	}

	LineItemStockEntryDto struct {
		SourceWarehouse string `json:"source_warehouse" required:"false"`
		TargetWarehouse string `json:"target_warehouse" required:"false"`
	}

	ItemLineDto struct {
		// ItemPrice ItemPriceDto `json:"item_price" required:"true"`
		ID          int    `json:"id"`
		Rate        int32  `json:"rate" required:"true"`
		Quantity    int32  `json:"quantity" required:"true"`
		Type        string `json:"line_type"`
		ItemPriceID int64  `json:"item_price_id"`

		//Item
		ItemName string `json:"item_name"`
		ItemCode string `json:"item_code"`
		ItemUUID string `json:"item_uuid"`
		//UOM
		Uom string `json:"uom"`

		//ItemPrice
		ItemPriceUUID string `json:"item_price_uuid"`
	}

	RequestLineItems struct {
		LineType    string `query:"line_type"`
		UpdateStock string `query:"update_stock" required:"false"`
		ID          string `query:"id"`
		PartyType   string `query:"party_type" required:"false"`
	}
)
