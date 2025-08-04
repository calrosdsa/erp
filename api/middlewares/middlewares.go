package middlewares

import (
	// "net/http"

	"context"
	"erp/api/common"
	"erp/gen/db/model"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	"fmt"
	"net/http"
	"strings"

	"errors"

	"github.com/danielgtaylor/huma/v2"
	"github.com/labstack/echo/v4"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type Middlewares interface {
	Authenticate(ctx huma.Context, next func(huma.Context))
	AuthenticateAdmin(ctx huma.Context, next func(huma.Context))
	ValidateActiveCompany(ctx huma.Context, next func(huma.Context))
	AuthEcho() echo.MiddlewareFunc
	ToolMiddleware() server.ToolHandlerMiddleware
}



type middlewares struct {
	session    repository.SessionService
	jwtService helpers.JwtHelper
	api        huma.API
}

func NewMiddlewares(
	sessionService repository.SessionService,
	api huma.API,
	jwtHelper helpers.JwtHelper,
) Middlewares {
	return &middlewares{
		jwtService: jwtHelper,
		session:    sessionService,
		api:        api,
	}
}

func (m *middlewares) ValidateActiveCompany(ctx huma.Context, next func(huma.Context)) {

}

func (m *middlewares) AuthEcho() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()
			token := c.Request().Header.Get("Authorization")
			if token == "" {
				return nil
			}
			token = m.jwtService.GetToken(token)
			claims, err := m.jwtService.ExtractClaimsAdmin(token)
			if err != nil {
				return c.String(http.StatusUnauthorized, "No Authorized")
			}
			acceptLanguage := c.Request().Header.Get("Accept-Language")
			if acceptLanguage == "" {
				acceptLanguage = "en"
			}
			// }
			if len(acceptLanguage) > 2 {
				acceptLanguage = acceptLanguage[:2]
			}
			user, err := m.session.GetUser(ctx, claims.ID)
			if err != nil {
				return c.String(http.StatusUnauthorized, "No User Found")
			}

			session := &common.AdminRequestContext{
				LanguageCode: common.LanguageCode(acceptLanguage),
				User:         user,
				Ctx:          ctx,
			}
			nC := context.WithValue(c.Request().Context(),domain.SESSION_KEY, session)
			c.SetRequest(c.Request().WithContext(nC))
			return next(c)
		}
	}
}

func (m *middlewares) AuthenticateAdmin(ctx huma.Context, next func(huma.Context)) {
	token := ctx.Header("Authorization")
	if token == "" {
		huma.WriteErr(m.api, ctx, http.StatusUnauthorized, "No token provided")
		return
	}
	token = m.jwtService.GetToken(token)
	claims, err := m.jwtService.ExtractClaimsAdmin(token)
	if err != nil {
		huma.WriteErr(m.api, ctx, http.StatusUnauthorized, "No Authorized")
		return
	}
	acceptLanguage := ctx.Header("Accept-Language")
	if acceptLanguage == "" {
		acceptLanguage = "en"
	}
	// }
	if len(acceptLanguage) > 2 {
		acceptLanguage = acceptLanguage[:2]
	}
	user, err := m.session.GetUser(ctx.Context(), claims.ID)
	if err != nil {
		huma.WriteErr(m.api, ctx, http.StatusUnauthorized, "No User Found")
		return
	}

	session := &common.AdminRequestContext{
		LanguageCode: common.LanguageCode(acceptLanguage),
		User:         user,
		Ctx:          ctx.Context(),
	}
	ctx = huma.WithValue(ctx, domain.SESSION_KEY, session)
	next(ctx)
}

func (m *middlewares) Authenticate(ctx huma.Context, next func(huma.Context)) {
	token := ctx.Header("Authorization")
	if token == "" {
		huma.WriteErr(m.api, ctx, http.StatusUnauthorized, "No token provided")
		return
	}
	fmt.Println("TOKEN", token)
	token = m.jwtService.GetToken(token)
	_, err := m.jwtService.ExtractClaimsAdmin(token)
	if err != nil {
		fmt.Println("ERROR HERE 1", err)
		huma.WriteErr(m.api, ctx, http.StatusUnauthorized, "No Authorized")
		return
	}
	// userAccount, err := m.accountService.GetUserAccount(ctx.Context(), claims)
	// if err != nil {
	// 	fmt.Println("ERROR HERE 2")
	// 	huma.WriteErr(*m.api, ctx, http.StatusUnauthorized, "No Authorized")
	// 	return
	// }
	acceptLanguage := ctx.Header("Accept-Language")
	if acceptLanguage == "" {
		acceptLanguage = "en"
	}
	if len(acceptLanguage) > 2 {
		acceptLanguage = acceptLanguage[:2]
	}
	userRelationUuid := ctx.Header("Session-Uuid")
	var userRelation model.UserRelation
	if userRelationUuid != "" {
		userRelation, _ = m.session.GetUserRelation(ctx.Context(), userRelationUuid)
	}
	companyDefaults, err := m.session.GetCompanyDefaults(ctx.Context(), userRelation.Company.ID)
	if err != nil {
		huma.WriteErr(m.api, ctx, http.StatusUnauthorized, "No Authorized")
	}

	session := &common.RequestContext{
		LanguageCode:    common.LanguageCode(acceptLanguage),
		ActiveCompany:   userRelation.Company,
		Role:            userRelation.Role,
		User:            userRelation.User,
		Profile:         userRelation.Profile,
		SessionUuid:     userRelationUuid,
		CompanyDefaults: companyDefaults,
		Ctx:             ctx.Context(),
	}
	ctx = huma.WithValue(ctx, domain.SESSION_KEY, session)

	// Call the next middleware in the chain. This eventually calls the
	// operation handler as well.
	next(ctx)
}

func (m *middlewares) ToolMiddleware() server.ToolHandlerMiddleware {
	return func(thf server.ToolHandlerFunc) server.ToolHandlerFunc {
		return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			// token := c.Request().Header.Get("Authorization")
			token := request.Header.Get("Authorization")
			if token == "" {
				return nil, errors.New("no token provided")
			}
			fmt.Println("TOKEN TOOL", token)
			token = m.jwtService.GetToken(token)
			claims, err := m.jwtService.ExtractClaimsAdmin(token)
			if err != nil {
				return nil, errors.New("no Authorized")
			}
			acceptLanguage := request.Header.Get("Accept-Language")
			if acceptLanguage == "" {
				acceptLanguage = "en"
			}
			// }
			if len(acceptLanguage) > 2 {
				acceptLanguage = acceptLanguage[:2]
			}
			userRelation, err := m.session.GetUserRelationByUserID(ctx, claims.ID)
			if err != nil {
				return nil, errors.New("no Authorized")
			}
			companyDefaults, err := m.session.GetCompanyDefaults(ctx, userRelation.Company.ID)
			if err != nil {
				return nil, errors.New("no Authorized")
			}

			session := &common.RequestContext{
				LanguageCode:    common.LanguageCode(acceptLanguage),
				ActiveCompany:   userRelation.Company,
				Role:            userRelation.Role,
				User:            userRelation.User,
				Profile:         userRelation.Profile,
				SessionUuid:     userRelation.UUID,
				CompanyDefaults: companyDefaults,
				Ctx:             ctx,
			}
			nC := context.WithValue(ctx, domain.SESSION_KEY, session)
			return thf(nC, request)
		}
	}
}

func (m *middlewares) getCookieValue(cookieHeader, cookieName string) string {
	// Split the cookie header by "; " to get individual cookies
	cookies := strings.Split(cookieHeader, "; ")

	// Iterate over the cookies
	for _, cookie := range cookies {
		// Split each cookie by "=" to get name and value
		parts := strings.SplitN(cookie, "=", 2)
		if len(parts) == 2 {
			name := parts[0]
			value := parts[1]
			if name == cookieName {
				return value
			}
		}
	}
	return ""
}
