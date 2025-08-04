package admin

import (
	"context"
	"erp/pkg/system"
	auth_admin "erp/project/admin/auth"
	admin_company "erp/project/admin/company"
	admin_core "erp/project/admin/core"
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

	//Register topics event
	svc.EventBus().RegisterTopics()

	m := monolith{
		Service: svc,
		modules: []system.Module{
			&auth_admin.Module{},
			&admin_company.Module{},
			&admin_core.Module{},
		},
	}

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
