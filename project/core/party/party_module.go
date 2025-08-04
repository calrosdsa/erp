package party

import (
	"context"
	"erp/pkg/system"
	party_rest "erp/project/core/party/internal/handler/rest"
	party_repo "erp/project/core/party/internal/repository"
	party_ucase "erp/project/core/party/internal/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct{}

func (m Module) Startup(ctx context.Context, svc system.Service) error {
	return Root(ctx, svc)
}

func Root(ctx context.Context, svc system.Service) error {
	partyRepository := party_repo.NewPartyRepository(svc.DBConn(), svc.Helpers())
	partyUseCase := party_ucase.NewPartyUseCase(
		svc.Logger(),partyRepository,
	)
	party_rest.NewPartyHandler(
		svc.HumaApi(),svc.Helpers(),huma.Middlewares{svc.Middlewares().Authenticate},
		partyUseCase,
	)
	return nil
}
