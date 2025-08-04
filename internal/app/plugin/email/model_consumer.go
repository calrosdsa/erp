package email

import "erp/internal/app/event-bus/event"

type ProccessEmail func(payload *event.NotificationData)

func ChainnedProccessEmail(emailProcessors ...ProccessEmail) ProccessEmail {
	return func(payload *event.NotificationData) {
		for _, processor := range emailProcessors {
			processor(payload)
		}
	}
}
