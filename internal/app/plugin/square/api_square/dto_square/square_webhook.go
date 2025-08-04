package dtosquare

import (
	squaretypes "erp/internal/app/plugin/square/square_types"
	"time"
)

type PaymentWeebhookRequest struct {
	Body struct {
		PaymentBody
	}
}
type PaymentBody struct {
	MerchantID string    `json:"merchant_id" required:"false"`
	Type       string    `json:"type" required:"false"`
	EventID    string    `json:"event_id" required:"false"`
	CreatedAt  time.Time `json:"created_at" required:"false"`
	Data       struct {
		Type   string `json:"type" required:"false"`
		ID     string `json:"id" required:"false"`
		Object struct {
			Payment struct {
				AmountMoney        squaretypes.Amount `json:"amount_money" required:"false"`
				ApplicationDetails struct {
					ApplicationID string `json:"application_id" required:"false"`
					SquareProduct string `json:"square_product" required:"false"`
				} `json:"application_details" required:"false"`
				ApprovedMoney squaretypes.Amount `json:"approved_money" required:"false"`
				Capabilities    []string `json:"capabilities" required:"false"`
				ExternalDetails ExternalDetails `json:"external_details" required:"false"`

				BuyerEmailAddress string `json:"buyer_email_address" required:"false"`
				CardDetails       struct {
					AvsStatus string `json:"avs_status" required:"false"`
					Card      struct {
						Bin         string `json:"bin" required:"false"`
						CardBrand   string `json:"card_brand" required:"false"`
						CardType    string `json:"card_type" required:"false"`
						ExpMonth    int    `json:"exp_month" required:"false"`
						ExpYear     int    `json:"exp_year" required:"false"`
						Fingerprint string `json:"fingerprint" required:"false"`
						Last4       string `json:"last_4" required:"false"`
						PrepaidType string `json:"prepaid_type" required:"false"`
					} `json:"card" required:"false"`
					CardPaymentTimeline struct {
						AuthorizedAt time.Time `json:"authorized_at" required:"false"`
						CapturedAt   time.Time `json:"captured_at" required:"false"`
					} `json:"card_payment_timeline" required:"false"`
					CvvStatus            string `json:"cvv_status" required:"false"`
					EntryMethod          string `json:"entry_method" required:"false"`
					StatementDescription string `json:"statement_description" required:"false"`
					Status               string `json:"status" required:"false"`
				} `json:"card_details" required:"false"`
				CreatedAt     time.Time          `json:"created_at" required:"false"`
				CustomerID    string             `json:"customer_id" required:"false"`
				DelayAction   string             `json:"delay_action" required:"false"`
				DelayDuration string             `json:"delay_duration" required:"false"`
				DelayedUntil  time.Time          `json:"delayed_until" required:"false"`
				ID            string             `json:"id" required:"false"`
				LocationID    string             `json:"location_id" required:"false"`
				OrderID       string             `json:"order_id" required:"false"`
				ProcessingFee []ProcessFee       `json:"processing_fee" required:"false"`
				ReceiptNumber string             `json:"receipt_number" required:"false"`
				ReceiptURL    string             `json:"receipt_url" required:"false"`
				SourceType    string             `json:"source_type" required:"false"`
				Status        string             `json:"status" required:"false"`
				TotalMoney    squaretypes.Amount `json:"total_money" required:"false"`
				UpdatedAt     time.Time          `json:"updated_at" required:"false"`
				Version       int                `json:"version" required:"false"`
			} `json:"payment" required:"false"`
		} `json:"object" required:"false"`
	} `json:"data" required:"false"`
}

type ExternalDetails struct {
	Source string `json:"source"`
	Type   string `json:"type"`
}

type ProcessFee struct {
	AmountMoney squaretypes.Amount `json:"amount_money"`
	EffectiveAt time.Time          `json:"effective_at"`
	Type        string             `json:"type"`
}
