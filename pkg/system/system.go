package system

import (
	"context"
	"erp/api/middlewares"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	accounting_service "erp/internal/service/accounting"
	"erp/internal/service/core"
	document_service "erp/internal/service/document"
	"erp/internal/service/permission"
	"erp/internal/service/session"
	setting_service "erp/internal/service/setting"
	stock_service "erp/internal/service/stock"
	"erp/pkg/bus"
	"erp/pkg/config"
	"erp/pkg/db"
	"erp/pkg/di"
	"erp/pkg/logger"
	"erp/pkg/waiter"
	"erp/pkg/ws"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humaecho"
	"github.com/labstack/echo/v4"
	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
	"github.com/stackus/errors"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/go-co-op/gocron/v2"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mark3labs/mcp-go/server"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type System struct {
	cfg *config.AppConfig
	// q *query.Query
	// db *gorm.DB
	nc                *nats.Conn
	js                nats.JetStreamContext
	echo              *echo.Echo
	rpc               *grpc.Server
	waiter            waiter.Waiter
	logger            logger.Logger
	tp                *sdktrace.TracerProvider
	humaApi           huma.API
	conn              db.Connection
	eventBus          bus.Bus
	permissionService repository.PermissionService
	sessionService    repository.SessionService
	coreService       repository.CoreService
	accountingService repository.AccountingService
	stockService      repository.StockService
	settingService    repository.SettingService
	documentService   repository.DocumentService
	helpers           *helpers.Helpers
	middlewares       middlewares.Middlewares
	di                di.Container
	wsConn            ws.WsConn
	scheduler         gocron.Scheduler
	mcpServer         *server.MCPServer
}

func NewSystem(cfg config.AppConfig) (*System, error) {
	s := &System{cfg: &cfg}
	s.initLogger()
	s.initWaiter()
	fmt.Println("INIT WRITTER")
	s.initDBConn()
	fmt.Println("INIT DB")
	// if err := s.initJS(); err != nil {
	// 	return nil, err
	// }
	fmt.Println("INIT JET STREAM")
	if err := s.initOpenTelemetry(); err != nil {
		return nil, err
	}
	fmt.Println("INIT OPEN TELEMETRY")
	s.initEcho()
	s.initHumaApi()
	s.initRpc()
	s.initScheduler()

	s.initHelpers()

	s.initEventBus()
	s.initContainer()

	s.initPermissionService()
	s.initSessionService()
	s.initCoreService()
	s.initAccountingService()
	s.initStockService()
	s.initSettingService()
	s.initDocumentService()

	s.initMiddliewares()
	s.initWsConn()
	s.initMcp()
	return s, nil
}

func (s *System) initMcp() {
	mcpServer := server.NewMCPServer(
		"Demo 🚀",
		"1.0.0",
		server.WithToolCapabilities(false),
		server.WithToolHandlerMiddleware(s.Middlewares().ToolMiddleware()),
	)
	
	
	s.mcpServer = mcpServer
}

func (s *System) Mcp() *server.MCPServer {
	return s.mcpServer
}

func (s *System) initContainer() {
	s.di = di.New()
	s.di.AddSingleton(domain.DbKey, func(c di.Container) (any, error) {
		return s.conn.GetQ(), nil
	})
}

func (s *System) Container() di.Container {
	return s.di
}

func (s *System) initMiddliewares() {
	s.middlewares = middlewares.NewMiddlewares(s.SessionService(), s.HumaApi(), s.Helpers().Jwt)
}

func (s *System) Middlewares() middlewares.Middlewares {
	return s.middlewares
}

func (s *System) Scheduler() gocron.Scheduler {
	return s.scheduler
}

func (s *System) initScheduler() {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		panic(err)
	}
	s.scheduler = scheduler
	go func() {
		s.scheduler.Start()
	}()
}

func (s *System) WsConn() ws.WsConn {
	return s.wsConn
}

func (s *System) initWsConn() {
	s.wsConn = ws.NewChatServer(s.echo, s.helpers, s.sessionService)
}

func (s *System) initHelpers() {
	s.helpers = helpers.InitHelpers(s.logger, s.cfg)
}

func (s *System) Helpers() *helpers.Helpers {
	return s.helpers
}

func (s *System) initEventBus() {
	s.eventBus = bus.NewBus()
}

func (s *System) EventBus() bus.Bus {
	return s.eventBus
}

func (s *System) initDocumentService() {
	s.documentService = document_service.NewDocumentService(
		s.conn,
		s.helpers.Currency,
		s.coreService,
	)
}
func (s *System) DocumentService() repository.DocumentService {
	return s.documentService
}

func (s *System) initSettingService() {
	s.settingService = setting_service.NewSettingService(s.logger, s.conn)
}
func (s *System) SettingService() repository.SettingService {
	return s.settingService
}

func (s *System) initCoreService() {
	s.coreService = core.NewCoreService(s.conn, s.logger, s.helpers)
}
func (s *System) CoreService() repository.CoreService {
	return s.coreService
}

func (s *System) initAccountingService() {
	s.accountingService = accounting_service.NewAccountingService(s.logger)
}
func (s *System) AccountingService() repository.AccountingService {
	return s.accountingService
}

func (s *System) initPermissionService() {
	s.permissionService = permission.NewPermissionService(s.conn, s.logger)
}

func (s *System) PermissionService() repository.PermissionService {
	return s.permissionService
}

func (s *System) initStockService() {
	s.stockService = stock_service.NewStockService(s.logger, s.conn)
}
func (s *System) StockService() repository.StockService {
	return s.stockService
}

func (s *System) initSessionService() {
	s.sessionService = session.New(
		s.conn, s.logger, s.cfg.ShutdownTimeout,
		s.helpers, s.cfg.PG,
	)
}

func (s *System) SessionService() repository.SessionService {
	return s.sessionService
}

func (s *System) initDBConn() {
	s.conn = db.NewDbConnection(s.logger)
}

func (s *System) DBConn() db.Connection {
	return s.conn
}

func (s *System) initEcho() {
	s.echo = echo.New()
	s.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))
	// s.mux.Use(middleware.Heartbeat("/liveness"))
	// s.mux.Method("GET", "/metrics", promhttp.Handler())
}
func (s *System) Echo() *echo.Echo {
	return s.echo
}

func (s *System) initOpenTelemetry() error {
	exporter, err := otlptracegrpc.New(context.Background())
	if err != nil {
		return err
	}

	s.tp = sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter))
	otel.SetTracerProvider(s.tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	s.waiter.Cleanup(func() {
		if err := s.tp.Shutdown(context.Background()); err != nil {
			// s.logger.Error().Err(err).Msg("ran into an issue shutting down the tracer provider")
			fmt.Println("ran into an issue shutting down the tracer provider")
		}
	})

	return nil
}

func (s *System) initRpc() {
	s.rpc = grpc.NewServer(
		grpc.StatsHandler(
			otelgrpc.NewServerHandler(),
		),
		grpc.ChainUnaryInterceptor(
			serverErrorUnaryInterceptor(),
		),
		// If there are streaming endpoints also add
		// grpc.StreamInterceptor(
		// 	otelgrpc.StreamServerInterceptor(),
		// ),
	)
	reflection.Register(s.rpc)
}

func (s *System) RPC() *grpc.Server {
	return s.rpc
}

func (s *System) initJS() (err error) {
	s.nc, err = nats.Connect(s.cfg.Nats.URL)
	if err != nil {
		return err
	}
	fmt.Println("NATS CONNECT")
	s.js, err = s.nc.JetStream()
	if err != nil {
		return err
	}
	fmt.Println("JET STREAM", s.cfg.Nats.Stream)

	_, err = s.js.AddStream(&nats.StreamConfig{
		Name:     s.cfg.Nats.Stream,
		Subjects: []string{fmt.Sprintf("%s.>", s.cfg.Nats.Stream)},
	})
	fmt.Println("Add stream", err)

	return err
}

func (s *System) JS() nats.JetStreamContext {
	return s.js
}

func (s *System) initLogger() {
	s.logger = logger.New("erp", "1.0")
	// s.logger = logger.NewLogger(logger.LogConfig{
	// 	Environment: s.cfg.Environment,
	// 	LogLevel:    logger.Level(s.cfg.LogLevel),
	// })
}

func (s *System) Logger() logger.Logger {
	return s.logger
}

func (s *System) Config() *config.AppConfig {
	return s.cfg
}

// func (s *System) initDB() (err error) {
// 	s.db, err = gorm.Open(postgres.Open(s.cfg.PG.Conn), &gorm.Config{})
// 	s.q = query.Use(s.db)
// 	return err
// }
// func (s *System) DB() *gorm.DB {
// 	return s.db
// }
// func (s *System) Q() *query.Query {
// 	return s.q
// }

func (s *System) initWaiter() {
	s.waiter = waiter.New(waiter.CatchSignals())
}

func (s *System) Waiter() waiter.Waiter {
	return s.waiter
}

func (s *System) initHumaApi() {
	humaConfig := huma.DefaultConfig("My API", "1.0.0")
	humaConfig.Servers = []*huma.Server{{URL: viper.GetString("app.url")}}
	s.humaApi = humaecho.New(s.echo, humaConfig)
}

func (s *System) HumaApi() huma.API {
	return s.humaApi
}

func (s *System) WaitForMcp(ctx context.Context) error {
	httpServer := server.NewStreamableHTTPServer(s.mcpServer)
	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		fmt.Printf("mcp server started; listening at http://localhost%s\n", "8080")
		defer fmt.Println("mcp server shutdown")
		if err := httpServer.Start(":8080"); err != nil {
			return err
		}
		return nil
	})
	group.Go(func() error {
		<-gCtx.Done()
		fmt.Println("mcp server to be shutdown")
		ctx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
		defer cancel()
		if err := httpServer.Shutdown(ctx); err != nil {
			return err
		}
		return nil
	})
	return group.Wait()
}

func (s *System) WaitForWeb(ctx context.Context) error {
	webServer := &http.Server{
		Addr:    s.cfg.Web.Address(),
		Handler: s.echo,
	}

	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		fmt.Printf("web server started; listening at http://localhost%s\n", s.cfg.Web.Port)
		defer fmt.Println("web server shutdown")
		// if err := s.echo.Start(fmt.Sprintf("%s",s.cfg.Web.Port)); err != http.ErrServerClosed {
		// 	return err
		// }
		if err := webServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	})
	group.Go(func() error {
		<-gCtx.Done()
		fmt.Println("web server to be shutdown")
		ctx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
		defer cancel()
		if err := webServer.Shutdown(ctx); err != nil {
			return err
		}
		return nil
	})

	return group.Wait()
}

func (s *System) WaitForRPC(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.cfg.Rpc.Address())
	if err != nil {
		return err
	}

	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		fmt.Println("rpc server started")
		defer fmt.Println("rpc server shutdown")
		if err := s.RPC().Serve(listener); err != nil && err != grpc.ErrServerStopped {
			return err
		}
		return nil
	})
	group.Go(func() error {
		<-gCtx.Done()
		fmt.Println("rpc server to be shutdown")
		stopped := make(chan struct{})
		go func() {
			s.RPC().GracefulStop()
			close(stopped)
		}()
		timeout := time.NewTimer(s.cfg.ShutdownTimeout)
		select {
		case <-timeout.C:
			// Force it to stop
			s.RPC().Stop()
			return fmt.Errorf("rpc server failed to stop gracefully")
		case <-stopped:
			return nil
		}
	})

	return group.Wait()
}

func (s *System) WaitForStream(ctx context.Context) error {
	// 	closed := make(chan struct{})
	// 	s.nc.SetClosedHandler(func(*nats.Conn) {
	// 		close(closed)
	// 	})
	// 	// group, gCtx := errgroup.WithContext(ctx)
	// 	// group.Go(func() error {
	// 	// 	fmt.Println("message stream started")
	// 	// 	defer fmt.Println("message stream stopped")
	// 	// 	<-closed
	// 	// 	return nil
	// 	// })
	// 	// group.Go(func() error {
	// 	// 	<-gCtx.Done()
	// 	// 	return s.nc.Drain()
	// 	// })
	// 	return group.Wait()
	return nil
}

func serverErrorUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		resp, err = handler(ctx, req)
		return resp, errors.SendGRPCError(err)
	}
}
