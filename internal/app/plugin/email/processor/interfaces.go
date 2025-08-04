package processor

import (
	"erp/internal/app/event-bus/event"
	"errors"
	"io"
)

var ErrBusy = errors.New("server busy")

type EmailOptions struct {
	EventType string
}

type EmailProcessor interface {
	ProcessEmail(payload *event.NotificationData, options EmailOptions) (bool, error)
	io.Closer
}
