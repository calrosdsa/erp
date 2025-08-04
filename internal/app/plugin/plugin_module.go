package plugin

import (
	"erp/internal/app/config"
	"erp/internal/app/connection"
	"erp/internal/app/plugin/square"
	"erp/internal/app/service/helpers"

	"github.com/asaskevich/EventBus"
)

type PluginModule struct {
	conn          *connection.Connection
	configService *config.ConfigService
	helpers       *helpers.Helpers
	bus           *EventBus.Bus
}

func InitPluginModule(
	conn *connection.Connection,
	configService *config.ConfigService,
	helpers *helpers.Helpers,
	bus *EventBus.Bus,

) *PluginModule {
	pluginModule := &PluginModule{
		conn:          conn,
		configService: configService,
		helpers:       helpers,
		bus:           bus,
	}
	return pluginModule
}

func (p *PluginModule) InitHandlers(
	authenticateM config.AppMiddleware,
	validateCompanyM config.AppMiddleware,
) {
	square.InitHandler(p.conn, p.configService, p.bus, p.helpers, authenticateM, validateCompanyM)
}

func (p *PluginModule) GetPlugin(name string) *config.PluginApp {
	switch name {
	case config.PLUGIN_SQUARE:
		return square.Init(p.conn, p.configService, p.helpers)
	}
	return nil
}
