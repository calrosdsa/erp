package auth_admin_rest

import (
	"context"
	"erp/api/common"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	auth_admin_ucase "erp/project/admin/auth/internal/usecase"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type AuthAdminHandler struct {
	jwt     helpers.JwtHelper
	sessionHelper  helpers.SessionHelper
	locale         helpers.Locale
	errorHelper    helpers.ErrorHelper
	authAdminUcase auth_admin_ucase.AuthAdminUseCase
}

func NewAuthAdminHandler(
	api huma.API,
	helpers *helpers.Helpers,
	authAdminUcase auth_admin_ucase.AuthAdminUseCase,
) {
	base := domain.AUTH_ADMIN_BASE_ROUTE
	tags := []string{"Auth Admin"}
	path := NewAuthAdminPaths(base)
	h := AuthAdminHandler{
		jwt: helpers.Jwt,
		sessionHelper: helpers.Session,
		locale: helpers.Locale,
		errorHelper: helpers.Error,
		authAdminUcase: authAdminUcase,
	}
	huma.Register(api,huma.Operation{
		OperationID: "sign-in-admin",
		Method: http.MethodPost,
		Summary: "Sign in user admin",
		Tags: tags,
		Path: path.SignIn,
		DefaultStatus: http.StatusOK,
	},h.SignIn)
}

func (h *AuthAdminHandler) SignIn(ctx context.Context,d *dto.SignInRequest) (*dto.SignInResponse,error){
	res, err := h.authAdminUcase.SignIn(ctx, d)
	if err != nil {
		return nil, huma.Error400BadRequest("Invalid credentials", err)
	}
	token, err := h.jwt.GenerateToken(common.Claims{
		ID:   res.Body.U.ID,
		Uuid: res.Body.U.UUID,
	})
	if err != nil {
		return nil, huma.Error400BadRequest("Fial to generate token", err)
	}
	// session := common.GetSession(ctx)
	// fmt.Println(i.Body,session.LanguageCode)
	res.Body.AccessToken = token
	return &res, nil
}

