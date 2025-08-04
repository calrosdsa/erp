package client

import (
	"context"
	"erp/api/dto"
	"erp/internal/app/entity"
	"erp/internal/app/service/helpers"
	"erp/internal/app/service/services"
	clientservice "erp/internal/app/service/services/client_service"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type ClientHandler struct {
	clientService *clientservice.ClientService
	sessionHelper helpers.SessionHelper
	locale        helpers.Locale
}

func NewHandler(
	api *huma.API,
	services *services.Services,
	helpers *helpers.Helpers,
	base string,
	tag []string,
	middlewares huma.Middlewares,
) {
	paths := NewClientPaths(base)
	handler := ClientHandler{
		clientService: services.ClientService,
		sessionHelper: helpers.Session,
		locale:        helpers.Locale,
	}
	huma.Register(*api, huma.Operation{
		OperationID:   "get-client-profile",
		Summary:       "Get client Profile",
		Path:          paths.Base,
		Method:        http.MethodGet,
		Tags:          tag,
		Middlewares:   middlewares,
		DefaultStatus: http.StatusOK,
	}, handler.GetClientProfile)

	huma.Register(*api, huma.Operation{
		OperationID:   "update-client-profile",
		Summary:       "Update client Profile",
		Path:          paths.Base,
		Method:        http.MethodPut,
		Tags:          tag,
		Middlewares:   middlewares,
		DefaultStatus: http.StatusOK,
	}, handler.EditClient)
}

func (h *ClientHandler) EditClient(ctx context.Context, i *dto.EditClientRequest) (*dto.ResponseMessage, error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	err := h.clientService.EditClient(req, i)
	if err != nil {
		return nil, huma.Error400BadRequest(h.locale.MustLocalize(
			helpers.OptionsLocale.WithID("Error.FailToEdit"),
			helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
		), err)
	}
	var response dto.ResponseMessage
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.EditSuccess"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}

func (h *ClientHandler) GetClientProfile(ctx context.Context, i *dto.AuthParams) (
	*dto.EntityResponse[entity.Client], error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	res, err := h.clientService.GetClientProfile(req)
	if err != nil {
		return nil, huma.Error400BadRequest("Fail to get profile", err)
	}
	var response dto.EntityResponse[entity.Client]
	response.Body.Result = res
	return &response, nil
}
