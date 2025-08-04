package config

import (
	"fmt"
	"time"

	"erp/pkg/rpc"

	"github.com/spf13/viper"
)

type (
	WebConfig struct {
		Host string 
		Port string 
	}
	ApiConfig struct {
		JwtSecret string 
		Timeout time.Duration
	}
	PGConfig struct {
		Conn string `required:"true"`
		CryptoPass string 
	}
	NatsConfig struct {
		URL    string `required:"true"`
		Stream string `default:"mallbots"`
	}
	OtelConfig struct {
		ServiceName      string
		ExporterEndpoint string
	}
	ClientConfig struct {
		Url string
	}
	Email struct {
		Port int 
		Host string
		Password string
		User string
	}
	AppConfig struct {
		Client ClientConfig
		Email Email
		Environment     string
		LogLevel        string
		PG              PGConfig
		Nats            NatsConfig
		Rpc             rpc.RpcConfig
		Web             WebConfig
		Otel            OtelConfig
		ShutdownTimeout time.Duration
		Api     ApiConfig
	}
)

func (c WebConfig) Address() string {
	return fmt.Sprintf("%s%s", c.Host, c.Port)
}
func InitConfig() (cfg AppConfig, err error) {
	api := ApiConfig{
	}
	if api.Timeout = time.Duration(viper.GetInt("api.timeout")) * time.Second; api.Timeout == 0 {
		panic("No api timeout provided")
	}
	if api.JwtSecret = viper.GetString("api.jwtSecret"); api.JwtSecret == "" {
		panic("No api jwt secret provided")
	}
	otel := OtelConfig{
		ServiceName: viper.GetString("otel.service_name"),
	}
	otel.ExporterEndpoint = viper.GetString("otel.exporter_otlp")
	if otel.ExporterEndpoint == "" {
		panic("No exporter otlp provided in config files")
	}

	nats := NatsConfig{}

	if nats.URL = viper.GetString("nats.url"); nats.URL == "" {
		panic("No nats url provided")
	}
	if nats.Stream = viper.GetString("nats.stream"); nats.Stream == "" {
		panic("No nats stream provided")
	}
	dbName := viper.GetString("db.name")
	dbPass := viper.GetString("db.pass")
	dbUser := viper.GetString("db.user")
	dbPort := viper.GetInt("db.port")
	dbHost := viper.GetString("db.host")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", dbHost,
		dbUser, dbPass, dbName, dbPort)
	appConfig := AppConfig{
		Client: ClientConfig{
			Url: viper.GetString("client.url"),
		},
		Email: Email{
			User: viper.GetString("email.user"),
			Host: viper.GetString("email.host"),
			Port: viper.GetInt("email.port"),
			Password: viper.GetString("email.password"),
		},
		Nats:     nats,
		LogLevel: viper.GetString("log.level"),
		PG: PGConfig{
			Conn: dsn,
			CryptoPass: viper.GetString("db.crypto_pass"),
		},
		Otel:            otel,
		ShutdownTimeout: time.Duration(viper.GetInt("shutdown_timeout")) * time.Second,
		Rpc: rpc.RpcConfig{
			Host: "0.0.0.0",
			Port: ":7001",
		},
		Web: WebConfig{
			Host:viper.GetString("app.host"),
			Port: viper.GetString("app.port"),
		},
		
		Api: api,
	}
	return appConfig, nil
}
