package dto

import "erp/internal/app/config"

type AppConfigResponse struct {
	Body struct {
		Plugins []config.PluginApp
	}
}