package auth

import (
	"context"
	"erp/pkg/system"
	account_module "erp/project/auth/account"
)

type Module struct{}

type monolith struct {
	system.Service
	modules []system.Module
}


func (m Module) Startup(ctx context.Context,svc system.Service) error {
	return Root(ctx,svc)
}
func Root(ctx context.Context,svc system.Service) error{
	m := monolith{
		Service: svc,
		modules: []system.Module{
			&account_module.Module{},
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
