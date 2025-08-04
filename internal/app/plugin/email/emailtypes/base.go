package emailtypes

import "erp/gen/db/model"

type EmailBuilderOpts struct {
	LanguageCode string
	Company *model.Company
}