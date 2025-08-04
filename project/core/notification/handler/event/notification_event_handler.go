package notification_event

import (
	"context"
	"erp/internal/domain"
	"erp/internal/domain/event"
	"erp/pkg/bus"
	"erp/pkg/logger"
	notification_repo "erp/project/core/notification/repository"
)

type notificationEvent struct {
	repo    notification_repo.NotificationEventRepo
	emitLog logger.EmitLog
}

func NewNotificationEvent(
	bus bus.Bus,
	logger logger.Logger,
	repo notification_repo.NotificationEventRepo,
) {
	handler := notificationEvent{
		repo:    repo,
		emitLog: logger.EmitLog("notification-event"),
	}
	bus.RegisterHandler(domain.ActivityCreated, handler.OnActivityCreated())

}

func (s *notificationEvent) OnActivityCreated() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func() {
				if err != nil {
					s.emitLog.Err(err, logger.OptionsLog.WithMethod("OnActivityCreated"))
				}
			}()
			payload, ok := e.Data.(event.ActivityEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = s.repo.HandlerActivityNotification(ctx, payload)
			return
		},
		AbortOnError: true,
		Matcher:      domain.ActivityCreated,
	}
}
