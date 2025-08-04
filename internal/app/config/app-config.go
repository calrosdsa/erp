package config

import (
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/labstack/echo/v4"
)

type AppConfig struct {
	ApiOptions      ApiOptions
	DbConfig        DbConfig
	Plugins         []PluginApp
	DefaultLanguage string
	EmailOptions    EmailOptions
	ClientConfig ClientConfig
	PermifyAuthorization PermifyAuthorization
}

type PermifyAuthorization struct {
	Endpoint string 
}

type ClientConfig struct {
	Url string
}


type DbConfig struct {
	User           string
	Password       string
	DbName         string
	Port           int
	Host           string
	CryptoPass     string
	CustomEntities []interface{}
}

type ApiOptions struct {
	Port             int
	JwtSecret        string
	JwtSecretFronend string
	TimeoutAPICall   time.Duration
	EchoServer       *echo.Echo
	Api              huma.API
	// Middlewares
}

type EmailOptions struct {
	Processor *ProcessorOptions
	Transport *TransportOptions
}

type TransportOptions struct {
	Host string
	Port int
	Auth Auth
}

type Auth struct {
	User string
	Pass string
}

type ProcessorOptions struct {
	QueueSize  int
	NumWorkers int
}

type AppMiddleware func(ctx huma.Context, next func(huma.Context))
