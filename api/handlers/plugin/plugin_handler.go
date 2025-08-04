package plugin

import (
	"context"
	"erp/api/common"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/app/service/services"
	pluginservice "erp/internal/app/service/services/plugin_service"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type PluginHandler struct {
	pluginService *pluginservice.PluginService
	sessionHelper helpers.SessionHelper
	validator     *helpers.ValidatorHelper
}

func NewHandler(
	api *huma.API,
	services *services.Services,
	helpers *helpers.Helpers,
	base string,
	tag string,
	authMiddleware common.Middleware,
	validateCompany common.Middleware,
) {
	paths := NewPluginPaths(base)
	handler := PluginHandler{
		pluginService: services.PluginService,
		sessionHelper: helpers.Session,
		validator:     helpers.Validator,
	}
	huma.Register(*api, huma.Operation{
		OperationID:   "plugins",
		Method:        http.MethodGet,
		Path:          paths.Base,
		Summary:       "Get available plugins",
		Tags:          []string{tag},
		DefaultStatus: http.StatusOK,
	}, handler.GetPlugins)

	huma.Register(*api, huma.Operation{
		OperationID:   "add-plugin",
		Method:        http.MethodPost,
		Path:          paths.Base,
		Summary:       "Add plugin",
		Tags:          []string{tag},
		DefaultStatus: http.StatusOK,
		Middlewares:   huma.Middlewares{authMiddleware, validateCompany},
	}, handler.AddPlugin)

	huma.Register(*api, huma.Operation{
		OperationID:   "get-plugin",
		Method:        http.MethodGet,
		Path:          paths.Detail,
		Summary:       "Get plugin from company",
		Tags:          []string{tag},
		DefaultStatus: http.StatusOK,
		Middlewares:   huma.Middlewares{authMiddleware, validateCompany},
	}, handler.GetPlugin)

	huma.Register(*api, huma.Operation{
		OperationID:   "update-plugin-credentials",
		Method:        http.MethodPut,
		Path:          paths.Detail,
		Summary:       "Update plugin credentials",
		Tags:          []string{tag},
		DefaultStatus: http.StatusOK,
		Middlewares:   huma.Middlewares{authMiddleware, validateCompany},
	}, handler.UpdatePluginCredentials)
}

func (h *PluginHandler) GetPlugins(ctx context.Context, i *struct{}) (*dto.PluginsResponse, error) {
	var response dto.PluginsResponse
	res := h.pluginService.GetPlugins()
	response.Body.Plugins = res
	return &response, nil
}

func (h *PluginHandler) GetPlugin(ctx context.Context, i *dto.PluginDetailRequest) (*dto.PluginDetailResponse, error) {
	var response dto.PluginDetailResponse
	req, _ := h.sessionHelper.GetSession(ctx)
	res, err := h.pluginService.GetPlugin(req, i)
	if err != nil {
		return nil, huma.Error400BadRequest("Not plugin found", err)
	}
	response.Body.CompanyPlugin = res
	return &response, nil
}

func (h *PluginHandler) AddPlugin(ctx context.Context, i *dto.AddPluginRequest) (*dto.AddPluginResponse, error) {
	var response dto.AddPluginResponse
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error401Unauthorized("Not Authorized", err)
	}
	res, err := h.pluginService.AddPlugin(req, i)
	if err != nil {
		return nil, huma.Error400BadRequest("Fail to add plugin", err)
	}
	response.Body.CompanyPlugin = res
	return &response, nil
}

func (h *PluginHandler) UpdatePluginCredentials(ctx context.Context, i *dto.UpdateCredentialsPluginRequest) (*dto.ResponseMessage, error) {
	var response dto.ResponseMessage
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error401Unauthorized("Not Authorized", err)
	}
	err = h.pluginService.UpdatePluginCredentials(req, i)
	if err != nil {
		return nil, huma.Error400BadRequest("fail.to_update_plugin", err)
	}
	response.Body.Message = "successfully.updated_credentials"
	return &response, nil
}
