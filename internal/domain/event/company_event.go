package event

import (
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
)

type CreatedCompanyEventData struct {
	Tx      *query.QueryTx
	Company model.Company
	CompanyDefaults model.CompanyDefault
	Body dto.CompanyAdminData
	IsRoot  bool
	LanguageCode string 
}
