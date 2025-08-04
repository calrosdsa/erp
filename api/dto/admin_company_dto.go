package dto

type (
	CreateCompanyAdminRequest struct {
		Body CompanyAdminData
	}
	CompanyAdminData struct {
		Name           string          `json:"name"`
		CompanyModules []CompanyModule `json:"company_modules"`
	}

	CreateUserAdminRequest struct {
		Body struct {
			CompanyID  int64  `json:"company_id"`
			Identifier string `json:"identifier"`
			FirstName  string `json:"first_name"`
			LastName   string `json:"last_name" required:"false"`
		}
	}

	CompanyModule struct {
		Name     string `json:"name"`
		Label    string `json:"label"`
		IconCode string `json:"icon_code"`
		IconName string `json:"icon_name"`
		Priority int32  `json:"priority"`
	}

	AddCompanyModules struct {
		Body struct {
			Modules []CompanyEntityDto `json:"modules"`
		}
	}
	CompanyEntityDto struct {
		EntityName string `json:"entity_name"`
		EntityID   int64  `json:"entity_id"`
		CompanyID  *int64 `json:"company_id"`
		Enabled    bool   `json:"enabled"`
	}
)
