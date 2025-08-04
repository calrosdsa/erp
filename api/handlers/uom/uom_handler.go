package uom

import (
	"context"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/app/service/services"
	uomservice "erp/internal/app/service/services/uom_service"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type PluginHandler struct {
	sessionHelper helpers.SessionHelper
	validator     *helpers.ValidatorHelper
	uomService *uomservice.UOMService
}

func NewHandler(
	api *huma.API,
	services *services.Services,
	helpers *helpers.Helpers,
	base string,
	tag []string,
	middlewares huma.Middlewares,
) {
	paths := NewPluginPaths(base)
	handler := PluginHandler{
		sessionHelper: helpers.Session,
		validator:     helpers.Validator,
		uomService: services.UOMService,
	}
	huma.Register(*api, huma.Operation{
		OperationID:   "uom",
		Method:        http.MethodGet,
		Path:          paths.Base,
		Summary:       "Retrieve UOMs (Units of Measure)",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:middlewares,
	}, handler.GetUnitOfMeasures)

}

func (h *PluginHandler) GetUnitOfMeasures(ctx context.Context, i *dto.UOMsRequest) (*dto.UOMsResponse, error) {
	var response dto.UOMsResponse
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error401Unauthorized("Not Authorized", err)
	}
	res,err := h.uomService.GetUnitOfMeasures(req,i)
	if err != nil {
		return nil,huma.NewError(http.StatusBadRequest,"Failed to retrieve UOM (Unit of Measure)")
	}
	response.Body.Results = res
	return &response, nil
}
