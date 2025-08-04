package emailhandler

import (
	"erp/internal/app/config"
	"erp/internal/app/connection"
	"erp/internal/app/event-bus/event"
	emailuser "erp/internal/app/plugin/email/emailhandler/email_user"
	emailservice "erp/internal/app/plugin/email/service"
	"erp/internal/app/service/helpers"
	"erp/internal/app/service/services"
)

type EmailHandlers struct {
	UserEmailHandler *emailuser.UserEmailHandler
}

func NewEmailHandlers(
	configService *config.ConfigService,
	conn *connection.Connection,
	helpers *helpers.Helpers,
	emailService emailservice.EmailService,
	services *services.Services,
) *EmailHandlers {
	return &EmailHandlers{
		UserEmailHandler: emailuser.NewEmailUserHandler(configService, conn, helpers, emailService,services),
	}
}

func (h *EmailHandlers) SendEmail(payload *event.NotificationData) {
	switch payload.NotificationEventType {
	case event.NOTIFY_NEW_CLIENT:
		h.UserEmailHandler.SendUserCredentials(payload)
		return
	}
}
