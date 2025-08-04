package email

import (
	"context"
	email_event "erp/internal/app/plugin/email/events"
	emailservice "erp/internal/app/plugin/email/service"
	"erp/pkg/system"
	"fmt"
)

type Plugin struct {}

func (m Plugin) Startup(ctx context.Context,svc system.Service)(error){
	return Root(ctx,svc)
}

func Root(ctx context.Context,svc system.Service)(error){
	fmt.Println("EMAIL PLUGIN...")

	emailService := emailservice.NewEmailService(svc.Config(),svc.Helpers())
	email_event.RegisterEmailEvents(
		svc.DBConn(),svc.Helpers(),svc.Config(),svc.EventBus(),
		svc.Logger(),emailService,
	)
	return nil
}

// type EmailModule struct {
// 	emailProcessor processor.EmailProcessor
// 	configService  *config.ConfigService
// 	logger         _logger.Logger
// }

// func NewEmailModule(
// 	logger _logger.Logger,
// 	configService *config.ConfigService,
// 	bus *EventBus.Bus,
// 	conn *connection.Connection,
// 	helpers *helpers.Helpers,
// 	services *services.Services,
// ) {
// 	busD := *bus

// 	em := &EmailModule{
// 		logger:        logger,
// 		configService: configService,
// 	}

// 	emailService := emailservice.NewEmailService(
// 		configService,helpers,
// 	)

// 	emailHandlers := emailhandler.NewEmailHandlers(
// 		configService, conn, helpers, emailService,services,
// 	)
// 	processorOptions := em.configService.GetEmailOptions().Processor

// 	em.emailProcessor = NewEmailProccessor(
// 		emailHandlers,
// 		Options.Logger(logger),
// 		Options.NumQueueSize(processorOptions.NumWorkers),
// 		Options.NumQueueSize(processorOptions.QueueSize),
// 	)

// 	busD.Subscribe(event.NOTIFICATION_EVENT, em.OnNotificationEvent)
// }

// func (m *EmailModule) OnNotificationEvent(d event.NotificationData) {
// 	m.emailProcessor.ProcessEmail(&d, processor.EmailOptions{})
// }

// func (m )
