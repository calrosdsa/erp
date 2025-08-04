package profile

import (
	"context"
	"erp/pkg/system"
	profile_rest "erp/project/manage/profile/handler"
	profile_repo "erp/project/manage/profile/repository"
	profile_ucase "erp/project/manage/profile/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct {
}

func (m *Module) Startup(ctx context.Context,svc system.Service,)(error){
	return Root(ctx,svc)
}

func Root(ctx context.Context,svc system.Service)(error){
	profileRepo := profile_repo.NewProfileRepo(svc.DBConn(),svc.Helpers(),svc.SessionService())
	profileUcase := profile_ucase.NewProfileUseCase(svc.PermissionService(),
	svc.CoreService(),profileRepo,svc.Logger(),svc.EventBus(),svc.Container())	
	profile_rest.NewProfileHandler(svc.HumaApi(),svc.Helpers(),svc.PermissionService(),
	profileUcase,huma.Middlewares{svc.Middlewares().Authenticate})	
	return nil
}