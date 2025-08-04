package piano

import (
	"context"
	"erp/pkg/system"
	"erp/project/piano/pianoform"
)

type Project struct{}

type monolith struct {
	system.Service
	modules []system.Module
}


func (m Project) Startup(ctx context.Context,svc system.Service) error {
	return Root(ctx,svc)
}
func Root(ctx context.Context,svc system.Service) error{

	//Register topics event
	// svc.EventBus().RegisterTopics()

	m := monolith{
		Service: svc,
		modules: []system.Module{
			&pianoform.Module{},
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
