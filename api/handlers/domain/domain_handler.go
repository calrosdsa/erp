package domain

import (
	"context"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/app/service/services"
	appservice "erp/internal/app/service/services/app_service"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type DomainHandler struct {
	domainService *appservice.DomainService
	errorHelper helpers.ErrorHelper
	sessionHelper helpers.SessionHelper
}

func NewDomainHandler(
	api *huma.API,
	services *services.Services,
	helpers *helpers.Helpers,
	base string,
	tags []string,
	middlewares huma.Middlewares,
){
	paths := NewDomainPaths(base)
	handler := DomainHandler{
		domainService: services.DomainService,
		errorHelper: helpers.Error,
		sessionHelper: helpers.Session,
	}
	huma.Register(*api,huma.Operation{
		OperationID:   "get currencies",
		Method:        http.MethodGet,
		Path:          paths.Currency,
		Summary:       "Get currencies",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares: middlewares,
	},handler.GetCurrencies)
}


func (h *DomainHandler)GetCurrencies(ctx context.Context,i *dto.RequestPaginationData)(
	*dto.PaginationResponse[dto.PaginationResult[[]dto.CurrencyDto]],error){
	req,_:= h.sessionHelper.GetSession(ctx)	
	res,err := h.domainService.GetCurrencies(ctx,i)
	if err != nil {
		return nil,h.errorHelper.HumaCustomError(string(req.LanguageCode),err)
	}
	var response dto.PaginationResponse[dto.PaginationResult[[]dto.CurrencyDto]]
	response.Body.PaginationResult = res
	return &response,nil
}
