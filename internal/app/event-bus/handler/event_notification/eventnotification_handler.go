package eventnotification

import (
	"erp/internal/app/event-bus/event"
	"fmt"

	"github.com/asaskevich/EventBus"
)

type EventNotification struct {
	// conn *connection.Connection
}

func NewEventNotification(bus *EventBus.Bus) {
	handler := EventNotification{}
	b := *bus
	b.Subscribe(event.NOTIFICATION_EVENT, handler.OnNotificationEvent)
}

func (h *EventNotification) OnNotificationEvent(d event.NotificationData) {
	fmt.Println(d)
}
