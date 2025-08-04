package acct_master

import (
	"context"
	"erp/pkg/system"
)

type Module struct{}

func (m *Module) Startup(ctx context.Context, svc system.Service) (err error) {
	return Root(ctx, svc)
}

func Root(ctx context.Context, svc system.Service) error {
	return nil
}
