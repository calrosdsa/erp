package dtosquare

import "erp/api/dto"

type RequestSubscriptionCancel struct {
	dto.AuthParams
	Body struct {
		SubscriptionId string `json:"subscriptionId"`
	}
}