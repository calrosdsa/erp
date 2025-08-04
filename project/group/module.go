package group

import (
	"context"
	"erp/pkg/system"
	rest_group "erp/project/group/internal/handler/rest"
	group_repo "erp/project/group/internal/repository"
	group_ucase "erp/project/group/internal/usecase"
	"fmt"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct{}

func (m *Module) Startup(ctx context.Context, svc system.Service) error {
	return Root(ctx, svc)
}

func Root(ctx context.Context, svc system.Service) error {
	groupRepo := group_repo.NewGroupRepository(svc.DBConn(), svc.Helpers())
	groupUcase := group_ucase.NewGroupUcaseCase(
		svc.Helpers(), svc.Logger(), groupRepo, svc.PermissionService(),
		svc.CoreService(),
	)
	rest_group.NewGroupHandler(
		svc.HumaApi(), svc.Helpers(), huma.Middlewares{svc.Middlewares().Authenticate},
		groupUcase, svc.PermissionService(),
	)
	fmt.Println("INIT GROUP MODULE...")
	return nil
}
