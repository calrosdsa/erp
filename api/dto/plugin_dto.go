package dto

import (
	"erp/internal/app/config"
	"erp/internal/app/entity"
)

type PluginsResponse struct {
	Body struct {
		Plugins []config.PluginApp `json:"plugins"`
	}
}

type AddPluginRequest struct {
	AuthParams
	// ActiveCompanyHeader
	Body struct {
		Plugin string `json:"plugin"`
	}
}

type AddPluginResponse struct {
	Body struct {
		CompanyPlugin entity.CompanyPlugins `json:"company_plugin"`
	}
}

type PluginDetailRequest struct {
	AuthParams
	// ActiveCompanyHeader
	Plugin string `path:"plugin"`
}

type PluginDetailResponse struct {
	Body struct {
		CompanyPlugin entity.CompanyPlugins `json:"company_plugin"`
	}
}

type UpdateCredentialsPluginRequest struct {
	AuthParams
	Body struct {
		Credentials string `json:"credentials" required:"true" `
	}
	Plugin string `path:"plugin"`

}
