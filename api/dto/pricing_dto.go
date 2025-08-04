package dto

import "erp/gen/db/model"

type (
	RequestPricings struct {
		PaginationParams
		OptionalQueryParams
		ID   string `query:"id" required:"false"`
		Code string `query:"code" required:"false"`
	}

	CreatePricingRequest struct {
		Body PricingData
	}

	PricingDataRequest struct {
		Body struct {
			Pricing          PricingDto            `json:"pricing"`
			PricingLineItems []PricingLineItemData `json:"pricing_line_items"`
		}
	}
	GenerateQuotationRequest struct {
		Body struct {
			Pricing PricingDto `json:"pricing"`
			PricingData
		}
	}

	EditPricingRequest struct {
		Body PricingData
	}

	PricingData struct {
		ID               int64                  `json:"id" required:"false"`
		PricingFields PricingFields `json:"fields" required:"true"`
		
		PricingCharges   []PricingChargeData    `json:"pricing_charges"`		
		PricingLineItems []PricingLineItemData  `json:"pricing_line_items"`
	}

	PricingFields struct {
		CustomerID *int64 `json:"customer_id" required:"false"`
		ProjectID    *int64 `json:"project_id" required:"false"`
		CostCenterID *int64 `json:"cost_center_id" required:"false"`
	}

	PricingChargeData struct {
		Name string  `json:"name"`
		Rate float64 `json:"rate"`
	}

	PricingLineItemData struct {
		SupplierID  *int64  `json:"supplier_id" required:"false"`
		PartNumber  *string `json:"part_number" required:"false"`
		Description *string `json:"description" required:"false"`

		Quantity         *int32   `json:"quantity" required:"false"`
		PlUnit           *float64 `json:"pl_unit" required:"false"`
		FobUnit          *float64 `json:"fob_unit" required:"false"`
		Retention        *float64 `json:"retention" required:"false"`
		CostZf           *float64 `json:"cost_zf" required:"false"`
		CostAlm          *float64 `json:"cost_alm" required:"false"`
		Tva              *float64 `json:"tva" required:"false"`
		Cantidad         *int32   `json:"cantidad" required:"false"`
		PrecioUnitario   *float64 `json:"precio_unitario" required:"false"`
		PrecioTotal      *float64 `json:"precio_total" required:"false"`
		PrecioUnitarioTc *float64 `json:"precio_unitario_tc" required:"false"`
		PrecioTotalTc    *float64 `json:"precio_total_tc" required:"false"`
		FobTotal         *float64 `json:"fob_total" required:"false"`
		GplTotal         *float64 `json:"gpl_total" required:"false"`
		TvaTotal         *float64 `json:"tva_total" required:"false"`

		FobUnitFn          *string `json:"fob_unit_fn" required:"false"`
		RetentionFn        *string `json:"retention_fn" required:"false"`
		CostZfFn           *string `json:"cost_zf_fn" required:"false"`
		CostAlmFn          *string `json:"cost_alm_fn" required:"false"`
		TvaFn              *string `json:"tva_fn" required:"false"`
		CantidadFn         *string `json:"cantidad_fn" required:"false"`
		PrecioUnitarioFn   *string `json:"precio_unitario_fn" required:"false"`
		PrecioTotalFn      *string `json:"precio_total_fn" required:"false"`
		PrecioUnitarioTcFn *string `json:"precio_unitario_tc_fn" required:"false"`
		PrecioTotalTcFn    *string `json:"precio_total_tc_fn" required:"false"`
		FobTotalFn         *string `json:"fob_total_fn" required:"false"`
		GplTotalFn         *string `json:"gpl_total_fn" required:"false"`
		TvaTotalFn         *string `json:"tva_total_fn" required:"false"`
		IsTitle            *bool   `json:"is_title" required:"false"`
		Color              *string `json:"color" required:"false"`
	}

	PricingDetailDto struct {
		PricingDto       PricingDto           `json:"pricing"`
		PricingLineItems []PricingLineItemDto `json:"pricing_line_items"`
		PricingCharges   []PricingChargeDto   `json:"pricing_charges"`
	}
	PricingDto struct {
		ID     int64  `json:"id"`
		Code   string `json:"code"`
		Status string `json:"status"`

		CustomerID   *int64  `json:"customer_id"`
		Customer     *string `json:"customer"`
		CustomerUUID *string `json:"customer_uuid"`

		AccountingDimensionDto
	}

	PricingChargeDto struct {
		Name string `json:"name"`
		Rate int32  `json:"rate"`
	}

	PricingLineItemDto struct {
		SupplierID         int64  `json:"supplier_id"`
		Supplier           string `json:"supplier"`
		PartNumber         string `json:"part_number"`
		Description        string `json:"description"`
		Quantity           int32  `json:"quantity"`
		PlUnit             int32  `json:"pl_unit"`
		FobUnitFn          string `json:"fob_unit_fn"`
		RetentionFn        string `json:"retention_fn"`
		CostZfFn           string `json:"cost_zf_fn"`
		CostAlmFn          string `json:"cost_alm_fn"`
		TvaFn              string `json:"tva_fn"`
		CantidadFn         string `json:"cantidad_fn"`
		PrecioUnitarioFn   string `json:"precio_unitario_fn"`
		PrecioTotalFn      string `json:"precio_total_fn"`
		PrecioUnitarioTcFn string `json:"precio_unitario_tc_fn"`
		PrecioTotalTcFn    string `json:"precio_total_tc_fn"`
		FobTotalFn         string `json:"fob_total_fn"`
		GplTotalFn         string `json:"gpl_total_fn"`
		TvaTotalFn         string `json:"tva_total_fn"`
		IsTitle            bool   `json:"is_title"`
		Color              string `json:"color"`
		// fob_total_fn
		// precio_unitario_tc_fn
	}
)

func PricingDtoFromModel(m *model.Pricing) PricingDto {
	return PricingDto{
		ID:     m.ID,
		Code:   m.Code,
		Status: m.Status,
	}
}
