package currency_exchange_rest

import (
	"context"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	currency_exchange_ucase "erp/project/core/currency_exchange/internal/usecase"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type CurrencyExchangeHandler struct {
	sessionHelper       helpers.SessionHelper
	locale              helpers.Locale
	errorHelper         helpers.ErrorHelper
	currencyExchangeUcase currency_exchange_ucase.CurrencyExchangeUcase
	permission          repository.PermissionService
}

func NewHandler(
	api huma.API,
	helpers *helpers.Helpers,
	currencyExchangeUcase currency_exchange_ucase.CurrencyExchangeUcase,
	middlewares huma.Middlewares,
	permission repository.PermissionService,
) {
	base := domain.CURRENCY_EXCHANGE_ROUTE
	tags := []string{"Currency Exchange"}
	path := NewPaths(base)
	h := CurrencyExchangeHandler{
		sessionHelper:       helpers.Session,
		locale:              helpers.Locale,
		errorHelper:         helpers.Error,
		currencyExchangeUcase: currencyExchangeUcase,
		permission:          permission,
	}
	huma.Register(api, huma.Operation{
		OperationID:   "currency-exchanges",
		Method:        http.MethodGet,
		Summary:       "Currency Exchanges",
		Tags:          tags,
		Path:          path.Base,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetCurrencyExchanges)
	huma.Register(api, huma.Operation{
		OperationID:   "currency-exchange",
		Method:        http.MethodGet,
		Summary:       "Currency Exchange",
		Tags:          tags,
		Path:          path.Detail,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetCurrencyExchange)

	huma.Register(api, huma.Operation{
		OperationID:   "create-currency-exchange",
		Method:        http.MethodPost,
		Summary:       "Create Currency Exchange",
		Tags:          tags,
		Path:          path.Base,
		DefaultStatus: http.StatusCreated,
		Middlewares:   middlewares,
	}, h.CreateCurrencyExchange)

	huma.Register(api, huma.Operation{
		OperationID:   "edit-currency-exchange",
		Method:        http.MethodPut,
		Path:          path.Base,
		Summary:       "Edit Currency Exchange",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.EditCurrencyExchange)

	huma.Register(api, huma.Operation{
		OperationID:   "update-status-currency-exchange",
		Method:        http.MethodPut,
		Path:          path.UpdateStatus,
		Summary:       "Update Status Currency Exchange",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.UpdateStatus)
}

func (h *CurrencyExchangeHandler) UpdateStatus(ctx context.Context, d *dto.UpdateStatusWithEvent) (
	*dto.ResponseMessage, error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	err := h.currencyExchangeUcase.UpdateStatus(req, d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToUpdate")
	}
	var response dto.ResponseMessage
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.UpdatedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}


func (h *CurrencyExchangeHandler) EditCurrencyExchange(ctx context.Context, d *dto.EditCurrencyExchangeRequest) (
	*dto.ResponseMessage, error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	err := h.currencyExchangeUcase.EditCurrencyExchange(req, d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToEdit")
	}
	var response dto.ResponseMessage
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.EditedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}

func (h *CurrencyExchangeHandler) CreateCurrencyExchange(ctx context.Context, d *dto.CreateCurrencyExchangeRequest) (
	*dto.ResponseData[dto.CurrencyExchangeDto], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.currencyExchangeUcase.CreateCurrencyExchange(req, d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToCreate")
	}
	var response dto.ResponseData[dto.CurrencyExchangeDto]
	response.Body.Result = res
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.CreatedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}

func (h *CurrencyExchangeHandler) GetCurrencyExchange(ctx context.Context, d *dto.RequestEntity) (
	*dto.EntityResponse[dto.ResultEntity[dto.CurrencyExchangeDto]], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.currencyExchangeUcase.GetCurrencyExchange(req, d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.EntityResponse[dto.ResultEntity[dto.CurrencyExchangeDto]]
	response.Body.Result = res
	response.Body.Actions = h.permission.GetActions(req.Ctx, domain.CURRENCY_EXCHANGE.ID)
	return &response, nil
}

func (h *CurrencyExchangeHandler) GetCurrencyExchanges(ctx context.Context, d *dto.RequestPaginationData) (
	*dto.PaginationResponse[dto.PaginationResult[[]dto.CurrencyExchangeDto]], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.currencyExchangeUcase.GetCurrencyExchanges(req, d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.PaginationResponse[dto.PaginationResult[[]dto.CurrencyExchangeDto]]
	response.Body.PaginationResult = res
	response.Body.Actions = h.permission.GetActions(req.Ctx, domain.CURRENCY_EXCHANGE.ID)
	return &response, nil
}
