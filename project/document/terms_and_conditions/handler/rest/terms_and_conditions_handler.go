package terms_and_conditions_rest

import (
	"context"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	terms_and_conditions_ucase "erp/project/document/terms_and_conditions/usecase"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type termsAndConditionsHandler struct {
	sessionHelper       helpers.SessionHelper
	locale              helpers.Locale
	errorHelper         helpers.ErrorHelper
	termsAndConditionsUcase terms_and_conditions_ucase.TermsAndConditionsUcase
	permission          repository.PermissionService
}

func NewHandler(
	api huma.API,
	helpers *helpers.Helpers,
	termsAndConditionsUcase terms_and_conditions_ucase.TermsAndConditionsUcase,
	middlewares huma.Middlewares,
	permission repository.PermissionService,
) {
	base := domain.TERMS_AND_CONDITIONS_ROUTE
	tags := []string{"Terms and Conditions"}
	path := NewPaths(base)
	h := termsAndConditionsHandler{
		sessionHelper:       helpers.Session,
		locale:              helpers.Locale,
		errorHelper:         helpers.Error,
		termsAndConditionsUcase: termsAndConditionsUcase,
		permission:          permission,
	}
	huma.Register(api, huma.Operation{
		OperationID:   "terms-and-conditions",
		Method:        http.MethodGet,
		Summary:       "Terms & Conditions",
		Tags:          tags,
		Path:          path.Base,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetTermsAndConditions)
	huma.Register(api, huma.Operation{
		OperationID:   "terms-and-conditions-details",
		Method:        http.MethodGet,
		Summary:       "Terms & Conditions Detail",
		Tags:          tags,
		Path:          path.Detail,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetTermsAndConditionsDetail)

	huma.Register(api, huma.Operation{
		OperationID:   "create-terms-and-conditions",
		Method:        http.MethodPost,
		Summary:       "Create Terms & Conditions",
		Tags:          tags,
		Path:          path.Base,
		DefaultStatus: http.StatusCreated,
		Middlewares:   middlewares,
	}, h.CreateTermsAndConditions)

	huma.Register(api, huma.Operation{
		OperationID:   "edit-terms-and-conditions",
		Method:        http.MethodPut,
		Path:          path.Base,
		Summary:       "Edit Terms & Conditions",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.EditTermsAndConditions)

	huma.Register(api, huma.Operation{
		OperationID:   "update-status-terms-and-conditions",
		Method:        http.MethodPut,
		Path:          path.UpdateStatus,
		Summary:       "Update Status Terms & Conditions",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.UpdateStatus)
}

func (h *termsAndConditionsHandler) UpdateStatus(ctx context.Context, d *dto.UpdateStatusWithEvent) (
	*dto.ResponseMessage, error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	err := h.termsAndConditionsUcase.UpdateStatus(req, d)
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


func (h *termsAndConditionsHandler) EditTermsAndConditions(ctx context.Context, d *dto.TermsAndConditionsDataRequest) (
	*dto.ResponseMessage, error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	err := h.termsAndConditionsUcase.EditTermsAndConditions(req, d.Body)
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

func (h *termsAndConditionsHandler) CreateTermsAndConditions(ctx context.Context, d *dto.TermsAndConditionsDataRequest) (
	*dto.ResponseData[dto.TermsAndConditionsDto], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.termsAndConditionsUcase.CreateTermsAndConditions(req, d.Body)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToCreate")
	}
	var response dto.ResponseData[dto.TermsAndConditionsDto]
	response.Body.Result = res
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.CreatedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}

func (h *termsAndConditionsHandler) GetTermsAndConditionsDetail(ctx context.Context, d *dto.RequestEntity) (
	*dto.EntityResponse[dto.ResultEntity[dto.TermsAndConditionsDto]], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.termsAndConditionsUcase.GetTermsAndConditionsDetial(req, d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.EntityResponse[dto.ResultEntity[dto.TermsAndConditionsDto]]
	response.Body.Result = res
	response.Body.Actions = h.permission.GetActions(req.Ctx, domain.TERMS_AND_CONDITIONS.ID)
	return &response, nil
}

func (h *termsAndConditionsHandler) GetTermsAndConditions(ctx context.Context, d *dto.TermsAndConditionsRequest) (
	*dto.ResponseDataList[[]dto.TermsAndConditionsDto], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.termsAndConditionsUcase.GetTermsAndConditions(req, d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	res.Body.Actions = h.permission.GetActions(req.Ctx, domain.TERMS_AND_CONDITIONS.ID)
	return &res, nil
}
