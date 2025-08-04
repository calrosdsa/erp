package document

import (
	"context"
	"erp/pkg/system"
	documentinfo "erp/project/document/document_info"
	payment_terms "erp/project/document/payment_term"
	"erp/project/document/payment_terms_template"
	"erp/project/document/terms_and_conditions"
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
			&terms_and_conditions.Module{},
			&payment_terms.Module{},
			&payment_terms_template.Module{},
			&documentinfo.Module{},
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
