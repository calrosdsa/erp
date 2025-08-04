package app

import (
	"context"
	_api "erp/api"
	"erp/internal/app/config"
	_config "erp/internal/app/config"
	"erp/internal/app/connection"
	eventbus "erp/internal/app/event-bus"
	"erp/internal/app/plugin"
	// "erp/internal/app/plugin/email"
	"erp/internal/app/repository"
	"erp/internal/app/service/helpers"
	"erp/internal/app/service/services"
	"erp/pkg/api"
	"erp/pkg/db"
	_logger "erp/pkg/logger"
	"erp/pkg/system"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"erp/api/common"

	"github.com/asaskevich/EventBus"

	"github.com/spf13/viper"
)

type Module struct{}

func (m Module) Startup(ctx context.Context, mono system.Service) (err error) {
	return Root(ctx, mono)
}

func Root(ctx context.Context, svc system.Service) (err error) {

	timeout := time.Duration(5) * time.Minute
	e := svc.Echo()
	e.Use(api.SetRequestContextWithTimeout(timeout))

	permifyEndpoint, err := getenvStr("PERMIFY_ENDPOINT")
	if err != nil {
		fmt.Println(err)
	}
	config := &config.AppConfig{
		ApiOptions: config.ApiOptions{
			Port:             getInt(os.Getenv("APP_PORT")),
			JwtSecret:        os.Getenv("JWT_SECRET"),
			JwtSecretFronend: os.Getenv("JWT_FRONTEND_SECRET"),
			TimeoutAPICall:   timeout,
			EchoServer:       e,
			Api:              svc.HumaApi(),
		},
		DbConfig: config.DbConfig{
			User:       viper.GetString("db.user"),
			Port:       viper.GetInt("db.port"),
			DbName:     viper.GetString("db.name"),
			Password:   viper.GetString("db.pass"),
			Host:       viper.GetString("db.host"),
			CryptoPass: viper.GetString("db.crypto_pass"),
			// CustomEntities: []interface{}{
			// 	&entitysquare.SquareObject{}, &entitysquare.SquareCustomer{}, &entitysquare.SquareSubscription{},
			// 	&entitysquare.SquareOrder{}},
		},
		Plugins: []config.PluginApp{
			{
				Name: config.PLUGIN_SQUARE,
			},
			{
				Name: config.PLUGIN_EMAIL,
			},
		},
		DefaultLanguage: string(common.LanguageCodeEN),
		ClientConfig: config.ClientConfig{
			Url: os.Getenv("CLIENT_HOST"),
		},
		PermifyAuthorization: config.PermifyAuthorization{
			Endpoint: permifyEndpoint,
		},
	}
	bus := EventBus.New()

	logger := _logger.New(
		"erp",
		"1.0",
	)

	configModule := _config.Init(config)

	helpers := helpers.Init(logger, &bus,svc.Config())

	connection := connection.NewDbConnection(configModule, helpers, &logger)

	dbConn := db.NewDbConnection(logger)

	plugins := plugin.InitPluginModule(connection, configModule.ConfigService, helpers, &bus)

	repositories := repository.NewRepositories(connection, dbConn, helpers)

	services := services.Init(connection, dbConn, configModule, helpers, logger, plugins, repositories)

	eventbus.NewEventBus(&bus, services, configModule.ConfigService, helpers, connection)

	// email.NewEmailModule(logger, configModule.ConfigService, &bus, connection, helpers, services)
	_api.Init(services, configModule, connection, helpers, plugins,svc.SessionService())
	return nil
}

func getenvStr(key string) (string, error) {
	v := os.Getenv(key)
	if v == "" {
		return v, config.ErrEnvVarEmpty
	}
	return v, nil
}

func getenvInt(key string) (int, error) {
	s, err := getenvStr(key)
	if err != nil {
		return 0, err
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func getInt(e string) int {
	val, err := strconv.Atoi(e)
	if err != nil {
		log.Fatal(err)
	}
	return val
}
