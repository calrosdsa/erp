package activity

import (
	"context"
	"erp/internal/domain"
	"erp/pkg/di"
	"erp/pkg/system"
	activity_event "erp/project/core/activity/handler/event"
	activity_mcp "erp/project/core/activity/handler/mcp"
	activity_rest "erp/project/core/activity/handler/rest"
	activity_repo "erp/project/core/activity/repository"
	activity_ucase "erp/project/core/activity/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct{}

func (m Module) Startup(ctx context.Context, svc system.Service) error {
	return Root(ctx, svc)
}
func Root(ctx context.Context, svc system.Service) error {
	activityRepo := activity_repo.NewActivityRepository(svc.DBConn(), svc.Helpers())
	activityEventRepo := activity_repo.NewActivityEventRepo(svc.DBConn())
	activityUseCase := activity_ucase.NewActivityUseCase(
		svc.Logger(), activityRepo,svc.EventBus(),svc.Container(),
	)
	activity_event.NewActiviityEvnetHandler(svc.EventBus(),activityEventRepo,svc.Logger())
	activity_rest.NewActivityHandler(svc.HumaApi(), svc.Helpers(),
		svc.PermissionService(), huma.Middlewares{svc.Middlewares().Authenticate}, activityUseCase)
	activity_mcp.NewActivityMcp(svc.Mcp(), svc.Helpers(), svc.PermissionService(), activityUseCase)


	svc.Container().AddSingleton(domain.ActivityUseCase,func(c di.Container) (any, error) {
		return activityUseCase,nil
	})
	return nil
}
