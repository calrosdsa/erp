package eventbus

import (
	"erp/internal/app/config"
	"erp/internal/app/connection"
	// eventnotification "erp/internal/app/event-bus/handler/event_notification"
	// eventsquarehandler "erp/internal/app/event-bus/handler/event_square_handler"
	"erp/internal/app/service/helpers"
	"erp/internal/app/service/services"
	"fmt"

	evbus "github.com/asaskevich/EventBus"
)

type EventBus struct {
	bus evbus.Bus
}

func NewEventBus(bus *evbus.Bus, services *services.Services, configService *config.ConfigService,
	helpers *helpers.Helpers, conn *connection.Connection) *EventBus {
	busD := *bus

	busD.Subscribe("main:calculator", calculator)
	// r := &Event{}
	// busD.Publish("main:calculator",r.Calculate(1,2))

	// eventsquarehandler.NewEventSquareHandler(conn, bus, services, configService, helpers)
	// eventnotification.NewEventNotification(&busD)
	return &EventBus{
		bus: busD,
	}
}

func (e *EventBus) Publish(topic string, handler interface{}) {
	e.bus.Publish(topic, handler)
}

func calculator(a int, b int) {
	fmt.Println("calculator", a, b)
}
