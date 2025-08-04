package terms_and_conditions

import (
	"context"
	"erp/pkg/system"
	terms_and_conditions_rest "erp/project/document/terms_and_conditions/handler/rest"
	terms_and_conditions_fsm "erp/project/document/terms_and_conditions/internal/pkg/fsm"
	terms_and_conditions_repo "erp/project/document/terms_and_conditions/repository"
	terms_and_conditions_ucase "erp/project/document/terms_and_conditions/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct{}

func (m *Module) Startup(ctx context.Context,svc system.Service)error {
	return Root(ctx,svc)
}

func Root(ctx context.Context,svc system.Service) error{
	fsm:= terms_and_conditions_fsm.NewFsm()
	termsAndConditionsRepo := terms_and_conditions_repo.NewRepository(
		svc.DBConn(),svc.Helpers(),
	)
	termsAndConditionsUcase := terms_and_conditions_ucase.NewUseCase(
		svc.Logger(),svc.CoreService(),svc.PermissionService(),termsAndConditionsRepo,fsm,
	)
	terms_and_conditions_rest.NewHandler(
		svc.HumaApi(),svc.Helpers(),termsAndConditionsUcase,huma.Middlewares{svc.Middlewares().Authenticate},
		svc.PermissionService(),
	)
	return nil
}