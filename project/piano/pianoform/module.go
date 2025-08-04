package pianoform

import (
	"context"
	"erp/pkg/system"
	pianoform_rest "erp/project/piano/pianoform/internal/handler/rest"
	pianoform_repo "erp/project/piano/pianoform/internal/repository"
	pianoform_ucase "erp/project/piano/pianoform/internal/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct{}

func (m Module) Startup(ctx context.Context, svc system.Service) error {
	return Root(ctx, svc)
}

func Root(ctx context.Context, svc system.Service) error {
	pianoFormRepo := pianoform_repo.NewPianoReposotory(
		svc.DBConn(), svc.Helpers(),
	)
	pianoUcase := pianoform_ucase.NewPianoUseCase(
		svc.Logger(), pianoFormRepo,svc.CoreService(),
		svc.Helpers(),
	)
	pianoform_rest.NewPianoFormHandler(
		svc.HumaApi(), svc.Helpers(), svc.PermissionService(),
		huma.Middlewares{svc.Middlewares().Authenticate},
		pianoUcase,
	)
	return nil
}
