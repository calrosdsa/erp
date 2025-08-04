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

func init() {
	viper.SetConfigFile(`../../configs/config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

    // handle err
    time.Local = time.UTC
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
	fmt.Println("INIT CONFIG")
	s, err := system.NewSystem(cfg)
	if err != nil {
		return err
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

	fmt.Println("TOPICS...", s.EventBus().Topics())
	fmt.Println("HANDLERS", s.EventBus().TopicHandlerKeys(domain.ReceiptCreatedEvent))
	fmt.Println("started mallbots application")
	defer fmt.Println("stopped mallbots application")

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
