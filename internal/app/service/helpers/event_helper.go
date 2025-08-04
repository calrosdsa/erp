package helpers

import "github.com/asaskevich/EventBus"

type EventHelper struct {
	bus EventBus.Bus
}

func NewEventBus(bus *EventBus.Bus)*EventHelper{
	return &EventHelper{
		bus: *bus,
	}
}

func (h *EventHelper)Publish(event string,payload interface{}){
	h.bus.Publish(event,payload)
}