package core

import (
	"context"
	"erp/internal/domain"
	"erp/pkg/di"
	"erp/pkg/system"
	"erp/project/core/activity"
	"erp/project/core/address"
	"erp/project/core/connection"
	"erp/project/core/contact"
	"erp/project/core/currency_exchange"
	"erp/project/core/module"
	"erp/project/core/notification"
	"erp/project/core/party"
	"erp/project/core/stage"
	"erp/project/core/workspace"
)

type Module struct{}

type monolith struct {
	system.Service
	modules []system.Module
}

func (m Module) Startup(ctx context.Context, svc system.Service) error {
	return Root(ctx, svc)
}
func Root(ctx context.Context, svc system.Service) error {
	m := monolith{
		Service: svc,
		modules: []system.Module{
			&notification.Module{},
			&activity.Module{},
			&party.Module{},
			&contact.Module{},
			&currency_exchange.Module{},
			&address.Module{},
			&module.Module{},
			&stage.Module{},
			&workspace.Module{},
			&connection.Module{},
		},
	}
	svc.Container().AddSingleton(domain.DatabaseTransactionKey, func(c di.Container) (any, error) {
		return svc.DBConn().GetQ().Begin(), nil
	})
	RegisterEventTopics(m.EventBus())
	return m.startupModules()
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
