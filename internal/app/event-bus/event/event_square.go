package event

import dtosquare "erp/internal/app/plugin/square/api_square/dto_square"

type SquarePaymenrCompleted struct {
	Body dtosquare.PaymentBody
}

const (
	SQUARE_PAYMENT_COMPLETED_EVENT = "square:payment:completed"
)

