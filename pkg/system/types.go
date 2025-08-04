package system

import (
	"context"
	// "erp/gen/db/query"
	"erp/api/middlewares"
	"erp/internal/app/service/helpers"
	"erp/internal/domain/repository"
	"erp/pkg/bus"
	"erp/pkg/config"
	"erp/pkg/db"
	"erp/pkg/di"
	"erp/pkg/logger"
	"erp/pkg/waiter"
	"erp/pkg/ws"

	"github.com/danielgtaylor/huma/v2"
	"github.com/go-co-op/gocron/v2"
	"github.com/labstack/echo/v4"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
	"github.com/mark3labs/mcp-go/server"

	// "gorm.io/gorm"
)

type Service interface {
	Config() *config.AppConfig
	DBConn() db.Connection
	// DB() *gorm.DB
	// Q() *query.Query
	JS() nats.JetStreamContext
	Echo() *echo.Echo
	RPC() *grpc.Server
	Waiter() waiter.Waiter
	Logger() logger.Logger
	HumaApi() huma.API
	PermissionService() repository.PermissionService
	SessionService() repository.SessionService
	CoreService() repository.CoreService
	AccountingService() repository.AccountingService
	StockService() repository.StockService
	SettingService() repository.SettingService
	DocumentService() repository.DocumentService
	EventBus() bus.Bus

	Helpers() *helpers.Helpers

	Middlewares() middlewares.Middlewares

	Container() di.Container
	WsConn() ws.WsConn
	Scheduler() gocron.Scheduler

	Mcp() *server.MCPServer
}

type Module interface {
	Startup(context.Context, Service) error
}
