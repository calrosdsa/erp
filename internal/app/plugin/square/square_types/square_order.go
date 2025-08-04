package squaretypes

import "time"

// type SquareOrderData struct {
// 	Type string `json:"type"`
// 	Data string `json:"data"`
// }

type SquareOrderResponse struct {
	PaymentLink struct {
		ID        string    `json:"id"`
		Version   int       `json:"version"`
		OrderID   string    `json:"order_id"`
		URL       string    `json:"url"`
		LongURL   string    `json:"long_url"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"payment_link"`
	RelatedResources struct {
		Orders []struct {
			ID         string `json:"id"`
			LocationID string `json:"location_id"`
			Source     struct {
				Name string `json:"name"`
			} `json:"source"`
			LineItems []struct {
				UID            string `json:"uid"`
				Name           string `json:"name"`
				Quantity       string `json:"quantity"`
				ItemType       string `json:"item_type"`
				BasePriceMoney Amount `json:"base_price_money"`
				VariationTotalPriceMoney Amount `json:"variation_total_price_money"`
				GrossSalesMoney Amount `json:"gross_sales_money"`
				TotalTaxMoney Amount `json:"total_tax_money"`
				TotalDiscountMoney Amount `json:"total_discount_money"`
				TotalMoney Amount `json:"total_money"`
				TotalServiceChargeMoney Amount `json:"total_service_charge_money"`
			} `json:"line_items"`
			Fulfillments []struct {
				UID   string `json:"uid"`
				Type  string `json:"type"`
				State string `json:"state"`
			} `json:"fulfillments"`
			NetAmounts struct {
				TotalMoney Amount `json:"total_money"`
				TaxMoney Amount `json:"tax_money"`
				DiscountMoney Amount `json:"discount_money"`
				TipMoney Amount `json:"tip_money"`
				ServiceChargeMoney Amount `json:"service_charge_money"`
			} `json:"net_amounts"`
			CreatedAt  time.Time `json:"created_at"`
			UpdatedAt  time.Time `json:"updated_at"`
			State      string    `json:"state"`
			Version    int       `json:"version"`
			TotalMoney Amount `json:"total_money"`
			TotalTaxMoney Amount `json:"total_tax_money"`
			TotalDiscountMoney Amount `json:"total_discount_money"`
			TotalTipMoney Amount `json:"total_tip_money"`
			TotalServiceChargeMoney Amount `json:"total_service_charge_money"`
			NetAmountDueMoney Amount `json:"net_amount_due_money"`
		} `json:"orders"`
	} `json:"related_resources"`
}

type LineItem struct {
	Name           string `json:"name"`
	Quantity       string    `json:"quantity"`
	BasePriceMoney Amount `json:"base_price_money"`
	Note string `json:"note,omitempty"`
}
