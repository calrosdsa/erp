package email_event

import (
	emailservice "erp/internal/app/plugin/email/service"
	"erp/internal/app/service/helpers"
	"erp/pkg/bus"
	"erp/pkg/config"
	"erp/pkg/db"
	"erp/pkg/logger"
)

func RegisterEmailEvents(
	conn db.Connection,
	helpers *helpers.Helpers,
	appConfig *config.AppConfig,
	bus bus.Bus,
	logger logger.Logger,
	emailService emailservice.EmailService,
) {
	NewUserEventHandler(conn,helpers,appConfig,bus,logger,emailService)
	NewAccountEventHandler(conn,helpers,appConfig,bus,logger,emailService)
}