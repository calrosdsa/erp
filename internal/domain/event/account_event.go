package event

import "erp/gen/db/model"

type PasswordResetEventData struct {
	LanguageCode string
	User         model.User
	Profile      model.Profile
	Company      model.Company
}
