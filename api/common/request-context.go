package common

import (
	"context"
	"erp/gen/db/model"

	"github.com/golang-jwt/jwt/v5"
)

var DEFAULT_REQ_CONTEXT = RequestContext{
	LanguageCode: LanguageCodeES,
	Ctx: context.TODO(),
}

type RequestContext struct {
	// CurrencyCode CurrencyCode
	LanguageCode  LanguageCode
	Ctx           context.Context
	Params        map[string]string
	ActiveCompany model.Company
	CompanyDefaults model.CompanyDefault
	Role          model.Role
	Profile       model.Profile
	User          model.User
	SessionUuid   string
}

func(r *RequestContext) SetContext(ctx context.Context)  {
	r.Ctx = ctx
}

type AdminRequestContext struct {
	LanguageCode LanguageCode
	Ctx          context.Context
	User         model.User
}

type Session struct {
	CurrentUser model.User
	AccessToken string `json:"access_token"`
	CompanyUuid string `json:"companyUuud"`
	Role        string `json:"role"`
	//Can be clientID or ...
	//Base on the role
	ProfileSessionUuid string `json:"userSessionUuid"`
}

func (r *RequestContext) GetClientID() uint {
	// clients := r.Session.CurrentUser.Clients
	// for _, client := range clients {
	// 	if client.Uuid == r.Session.ProfileSessionUuid {
	// 		return client.ID
	// 	}
	// }
	return 0
}

// type UserSession struct {
// 	AccessToken string `json:"access_token"`
// 	Locale      string `json:"locale"`
// 	CompanyUuid string `json:"companyUuud"`
// 	Role        string `json:"role"`
// 	//Can be clientID or ...
// 	//Base on the role
// 	UserSessionUuid    string    `json:"userSessionUuid"`
// 	}

type Claims struct {
	jwt.RegisteredClaims
	ID   int64  `json:"id"`
	Uuid string `json:"uuid"`
}

type ErrorMessage struct {
	Message string `json:"message"`
}

func GetSession(ctx context.Context) *RequestContext {
	v := ctx.Value("session")
	if v == nil {
		return nil
	}

	if userProfile, ok := v.(*RequestContext); ok {
		return userProfile
	}

	return nil
}
