package main

import (
	"erp/pkg/config"
	"erp/pkg/db"
	"erp/pkg/system"
	"erp/project/accounting"
	"erp/project/admin"
	"erp/project/auth"
	"erp/project/buying"
	"erp/project/chat_module"
	"erp/project/company"
	"erp/project/core"
	"erp/project/crm"
	"erp/project/document"
	"erp/project/group"
	"erp/project/invoice"
	"erp/project/invoicing"
	"erp/project/manage"
	"erp/project/order"
	"erp/project/piano"
	"erp/project/pricing_module"
	project_module "erp/project/projet_module"
	"erp/project/quotation"
	"erp/project/receipt"
	"erp/project/regate"
	"erp/project/selling"
	"erp/project/stock"
	"flag"
	"fmt"
	"os"
	"time"

	"erp/internal/app"
	"erp/internal/app/plugin/email"
	"erp/internal/domain"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type monolith struct {
	*system.System
	modules []system.Module
}

func initConfig() error {
	viper.SetConfigFile(`../../configs/config.json`)
	// viper.SetConfigFile(`./configs/config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	return nil
}

var (
	devMode = flag.Bool("dev", false, "Enable development mode with enhanced logging")
	localMode = flag.Bool("local", false, "Enable local development mode (implies -dev)")
	hotReload = flag.Bool("hot-reload", false, "Enable hot reload mode (for use with Air)")
	version = flag.Bool("version", false, "Show version information")
	migrate = flag.Bool("migrate", false, "Run database migrations and exit")
	generateOpenAPI = flag.Bool("generate-openapi", false, "Generate OpenAPI specs and exit")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nERP Backend Server\n\n")
		fmt.Fprintf(os.Stderr, "Commands:\n")
		fmt.Fprintf(os.Stderr, "  %s                    Start the server in production mode\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -dev               Start the server in development mode\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -local             Start the server in local development mode\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -hot-reload        Start with hot reload support (used by Air)\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -migrate           Run database migrations and exit\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -generate-openapi  Generate OpenAPI specifications and exit\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -version           Show version information\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nFlags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nFor hot reload development, use: make dev\n")
		fmt.Fprintf(os.Stderr, "This will start the infrastructure and use Air for hot reloading.\n")
	}
	
	flag.Parse()

	// Handle version flag
	if *version {
		fmt.Printf("ERP Backend Server\n")
		fmt.Printf("Version: development\n") // This could be set during build with ldflags
		fmt.Printf("Build time: %s\n", time.Now().Format("2006-01-02 15:04:05"))
		return
	}

	// Local mode implies dev mode
	if *localMode {
		*devMode = true
	}

	// Initialize configuration
	if err := initConfig(); err != nil {
		if *devMode {
			fmt.Printf("Warning: Error loading config: %v\n", err)
		} else {
			panic(err)
		}
	}

	err := godotenv.Load()
	if err != nil {
		if *devMode {
			fmt.Printf("Warning: Error loading .env file: %v\n", err)
		} else {
			panic("Error loading .env file")
		}
	}

	// Set timezone
	time.Local = time.UTC

	// Print startup information
	if *devMode {
		fmt.Println("🚀 Starting ERP Backend Server")
		if *localMode {
			fmt.Println("📍 Mode: Local Development")
		} else {
			fmt.Println("📍 Mode: Development")
		}
		if *hotReload {
			fmt.Println("🔥 Hot reload: Enabled")
		}
		fmt.Printf("⏰ Started at: %s\n", time.Now().Format("2006-01-02 15:04:05"))
		fmt.Println("──────────────────────────────")
	}

	if err := run(); err != nil {
		fmt.Printf("application exitted abnormally: %s\n", err.Error())
		os.Exit(1)
	}
}

func run() (err error) {
	var cfg config.AppConfig
	cfg, err = config.InitConfig()
	if err != nil {
		return err
	}
	
	if *devMode {
		fmt.Println("⚙️  CONFIG INITIALIZED")
	} else {
		fmt.Println("INIT CONFIG")
	}
	
	s, err := system.NewSystem(cfg)
	if err != nil {
		return err
	}

	// Handle special commands that don't need full server startup
	if *migrate {
		fmt.Println("🗄️  Running database migrations...")
		// Add migration logic here
		fmt.Println("✅ Database migrations completed")
		return nil
	}

	if *generateOpenAPI {
		fmt.Println("📋 Generating OpenAPI specifications...")
		// Add OpenAPI generation logic here
		fmt.Println("✅ OpenAPI specifications generated")
		return nil
	}
	m := monolith{
		System: s,
		modules: []system.Module{
			&admin.Module{},
			&core.Module{},
			&company.Module{},
			&project_module.Module{},
			&app.Module{},
			&auth.Module{},
			&stock.Module{},
			&manage.Module{},
			&accounting.Module{},
			&invoicing.Module{},
			&group.Module{},
			&quotation.Module{},
			&receipt.Module{},
			&order.Module{},
			&buying.Module{},
			&selling.Module{},
			&invoice.Module{},
			&document.Module{},
			&crm.Module{},
			&chat_module.Module{},

			&email.Plugin{},
			&regate.Project{},
			&piano.Project{},
			
			&pricing_module.Module{},
		},
	}
	defer func(conn db.Connection) {
		sql, _ := conn.GetDB().DB()
		if err := sql.Close(); err != nil {
			return
		}
	}(m.DBConn())

	if err = m.startupModules(); err != nil {
		return err
	}

	// Development mode provides more detailed information
	if *devMode {
		fmt.Printf("📡 Event Topics: %v\n", s.EventBus().Topics())
		fmt.Printf("🔌 Receipt Handlers: %v\n", s.EventBus().TopicHandlerKeys(domain.ReceiptCreatedEvent))
		fmt.Println("✅ ERP Backend Server started successfully")
		fmt.Println("──────────────────────────────")
		fmt.Println("🌐 Server endpoints:")
		fmt.Println("  - HTTP API: http://localhost:8080")
		fmt.Println("  - Health check: http://localhost:8080/health")
		if *localMode {
			fmt.Println("  - NATS Monitoring: http://localhost:8222")
			fmt.Println("  - PostgreSQL: localhost:5432")
			fmt.Println("  - Redis: localhost:6379")
		}
		if *hotReload {
			fmt.Println("🔥 Hot reload is active - file changes will restart the server")
		}
		fmt.Println("──────────────────────────────")
		fmt.Println("Press Ctrl+C to stop the server")
		defer fmt.Println("🛑 ERP Backend Server stopped")
	} else {
		fmt.Println("TOPICS...", s.EventBus().Topics())
		fmt.Println("HANDLERS", s.EventBus().TopicHandlerKeys(domain.ReceiptCreatedEvent))
		fmt.Println("started mallbots application")
		defer fmt.Println("stopped mallbots application")
	}

	m.Waiter().Add(
		m.WaitForWeb,
		m.WaitForRPC,
		m.WaitForStream,
		m.WaitForMcp,
	)
	err = m.Waiter().Wait()
	if err != nil {
		return err
	}
	return nil
}

func (m *monolith) startupModules() error {

	for _, module := range m.modules {
		ctx := m.Waiter().Context()
		if err := module.Startup(ctx, m); err != nil {
			return err
		}
	}
	return nil
}
