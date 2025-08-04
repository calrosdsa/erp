package domain

import (
	"context"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/app/service/services"
	pluginservice "erp/internal/app/service/services/plugin_service"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type PluginHandler struct {
	pluginService *pluginservice.PluginService
}

func NewHandler(
	api *huma.API,
	services *services.Services,
	helpers *helpers.Helpers,
	base string,
	tag string,
	middlewares huma.Middlewares,
){
	paths := NewConfigPaths(base)
	handler := PluginHandler{
		pluginService: services.PluginService,
	}
	huma.Register(*api,huma.Operation{
		OperationID:   "appConfig",
		Method:        http.MethodGet,
		Path:          paths.Plugins,
		Summary:       "Get available plugins",
		Tags:          []string{tag},
		DefaultStatus: http.StatusOK,
	},handler.GetPlugins)
}


func (h *PluginHandler)GetPlugins(ctx context.Context,i *struct{})(*dto.PluginsResponse,error){
	var response dto.PluginsResponse
	res := h.pluginService.GetPlugins()
	response.Body.Plugins = res
	return &response,nil
}
