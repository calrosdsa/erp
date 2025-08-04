package notification

import (
	"context"
	"erp/internal/domain"
	"erp/pkg/di"
	"erp/pkg/system"
	notification_event "erp/project/core/notification/handler/event"
	notification_rest "erp/project/core/notification/handler/http"
	notification_repo "erp/project/core/notification/repository"
	notification_ucase "erp/project/core/notification/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct{}

func (m *Module) Startup(ctx context.Context, svc system.Service) error {
	return Root(ctx, svc)
}

func Root(ctx context.Context, svc system.Service) error {
	notificationRepo := notification_repo.NewNotificationRepo(svc.DBConn())
	notificationEventRepo := notification_repo.NewNotificationEventRepo(svc.Helpers())
	notificationUcase := notification_ucase.NewNotificationUcase(svc.Logger(), notificationRepo,
		svc.Scheduler(), svc.WsConn())
	svc.Container().AddSingleton(domain.NotificationUseCase, func(c di.Container) (any, error) {
		return notificationUcase, nil
	})
	notification_event.NewNotificationEvent(svc.EventBus(), svc.Logger(), notificationEventRepo)
	notification_rest.NewNotificationHandler(svc.HumaApi(), svc.Helpers(),
		huma.Middlewares{svc.Middlewares().Authenticate}, notificationUcase)
	return nil
}
