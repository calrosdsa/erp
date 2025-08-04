package tac_rest

import (
	"context"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	tac_usecase "erp/project/accounting/tax_and_charges/usecase"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type TacHandler struct {
	sessionHelper helpers.SessionHelper
	permission    repository.PermissionService
	errorHelper   helpers.ErrorHelper
	tacUseCase    tac_usecase.TacUseCase
	locale        helpers.Locale
}

func NewTacHandler(
	api huma.API,
	helpers *helpers.Helpers,
	middlewares huma.Middlewares,
	permission repository.PermissionService,
	tacUseCase tac_usecase.TacUseCase,
) {
	paths := NewTacPaths(domain.TAXES_AND_CHARGES)
	tag := []string{"Taxes and Charges"}
	handler := TacHandler{
		sessionHelper: helpers.Session,
		permission:    permission,
		tacUseCase:    tacUseCase,
		errorHelper:   helpers.Error,
		locale:        helpers.Locale,
	}
	huma.Register(api, huma.Operation{
		OperationID:   "edit-tax-and-charge",
		Method:        http.MethodPut,
		Path:          paths.Base,
		Summary:       "Edit tax and charge",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.EditTaxAndChargeLine)

	huma.Register(api, huma.Operation{
		OperationID:   "add-tax-and-charge",
		Method:        http.MethodPost,
		Path:          paths.Base,
		Summary:       "Add tax and charge",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.AddTaxAndChargeLine)

	huma.Register(api, huma.Operation{
		OperationID:   "tax-and-charges",
		Method:        http.MethodGet,
		Path:          paths.Base,
		Summary:       "Tax and charges",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetTacLines)

	huma.Register(api, huma.Operation{
		OperationID:   "delete-tax-and-charge",
		Method:        http.MethodDelete,
		Path:          paths.Base,
		Summary:       "Delete tax and charge",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.DeleteTaxAndChargeLine)
}

func (h *TacHandler) DeleteTaxAndChargeLine(ctx context.Context, i *dto.DeleteTaxLineRequest) (
	*dto.ResponseMessage, error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = h.tacUseCase.DeleteTaxAndChargeLine(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.ResponseMessage
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.DeletedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}

func (h *TacHandler) GetTacLines(ctx context.Context, i *dto.RequestTaxLines) (
	*dto.ResponseData[[]dto.TaxAndChargeLineDto], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.tacUseCase.GetTACLines(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.ResponseData[[]dto.TaxAndChargeLineDto]
	response.Body.Result = res
	return &response, nil
}

func (h *TacHandler) EditTaxAndChargeLine(ctx context.Context, i *dto.EditTaxLineRequest) (
	*dto.ResponseMessage, error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = h.tacUseCase.EditTaxAndChargeLine(req, i)
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

func (h *TacHandler) AddTaxAndChargeLine(ctx context.Context, i *dto.AddTaxLineRequest) (
	*dto.ResponseMessage, error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = h.tacUseCase.CreateTaxAndChargeLine(req, i)
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
