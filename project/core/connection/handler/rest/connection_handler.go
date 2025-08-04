package connection_rest

import (
	"context"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	connection_ucase "erp/project/core/connection/usecase"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type ConnectionHandler struct {
	sessionHelper   helpers.SessionHelper
	locale          helpers.Locale
	errorHelper     helpers.ErrorHelper
	connectionUcase connection_ucase.ConnectionUcase
	permission      repository.PermissionService
}

func NewConnectionHandler(
	api huma.API,
	helpers *helpers.Helpers,
	connectionUcase connection_ucase.ConnectionUcase,
	middlewares huma.Middlewares,
	permission repository.PermissionService,
) {
	base := domain.CONNECTION_ROUTE
	tags := []string{"Connection"}
	h := ConnectionHandler{
		sessionHelper:   helpers.Session,
		locale:          helpers.Locale,
		errorHelper:     helpers.Error,
		connectionUcase: connectionUcase,
		permission:      permission,
	}
	huma.Register(api, huma.Operation{
		OperationID:   "connections",
		Method:        http.MethodGet,
		Summary:       "",
		Tags:          tags,
		Path:          base + "/{id}",
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetConnectionsEntity)

}

func (h *ConnectionHandler) GetConnectionsEntity(ctx context.Context, d *dto.RequestEntity) (
	*dto.ResponseData[[]dto.ConnectionDto], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.connectionUcase.GetConnectionsEntity(req, *d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.ResponseData[[]dto.ConnectionDto]
	response.Body.Result = res
	return &response, nil
}
