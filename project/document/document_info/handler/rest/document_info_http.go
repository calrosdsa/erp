package documentinfo_rest

import (
	"context"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	documentinfo_ucase "erp/project/document/document_info/usecase"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type documentInfoHandler struct {
	sessionHelper helpers.SessionHelper
	locale        helpers.Locale
	errorHelper   helpers.ErrorHelper
	ucase         documentinfo_ucase.DocumentInfoUseCase
	permission    repository.PermissionService
}

func NewDocumentInfoHandler(
	api huma.API,
	helpers *helpers.Helpers,
	middlewares huma.Middlewares,
	ucase documentinfo_ucase.DocumentInfoUseCase,
	permission repository.PermissionService,
) {
	base := domain.DOCUMENT_INFO_BASE_ROUTE
	tags := []string{"Document Info"}
	paths := NewPaths(base)
	h := documentInfoHandler{
		sessionHelper: helpers.Session,
		locale:        helpers.Locale,
		errorHelper:   helpers.Error,
		ucase:         ucase,
		permission:    permission,
	}
	huma.Register(api, huma.Operation{
		OperationID:   "edit-address-and-contact",
		Method:        http.MethodPut,
		Path:          paths.AddressAndContact,
		Summary:       "Edit Address & Contact",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.EditAddressAndContact)
	huma.Register(api, huma.Operation{
		OperationID:   "address-and-contact",
		Method:        http.MethodGet,
		Path:          paths.AddressAndContactDetail,
		Summary:       "Address & Contact",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetAddressAndContact)

	huma.Register(api, huma.Operation{
		OperationID:   "edit-doc-terms",
		Method:        http.MethodPut,
		Path:          paths.DocTerm,
		Summary:       "Edit Doc Terms",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.EditDocTerms)

	huma.Register(api, huma.Operation{
		OperationID:   "doc-terms",
		Method:        http.MethodGet,
		Path:          paths.DocTermDetail,
		Summary:       "Doc Terms",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetDocTerms)

	huma.Register(api, huma.Operation{
		OperationID:   "edit-doc-accounts",
		Method:        http.MethodPut,
		Path:          paths.DocAccounting,
		Summary:       "Edit Doc Accounts",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.EditDocAccounts)

	huma.Register(api, huma.Operation{
		OperationID:   "doc-accounting",
		Method:        http.MethodGet,
		Path:          paths.DocAccountingDetail,
		Summary:       "Doc Accounting",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetDocAccounts)
}

func (h *documentInfoHandler) EditDocAccounts(ctx context.Context, d *dto.EditDocAccountingRequest) (
	*dto.ResponseMessage, error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = h.ucase.EditDocAccounts(req, d.Body)
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

func (h *documentInfoHandler) GetDocAccounts(ctx context.Context, d *dto.RequestEntityWithParty) (
	*dto.ResponseData[dto.DocAccountingDto], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.ucase.GetDocAccounts(req, *d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.ResponseData[dto.DocAccountingDto]
	response.Body.Result = res
	return &response, nil
}

func (h *documentInfoHandler) EditDocTerms(ctx context.Context, d *dto.EditDocTermRequest) (
	*dto.ResponseMessage, error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = h.ucase.EditDocTerms(req, d.Body)
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

func (h *documentInfoHandler) GetDocTerms(ctx context.Context, d *dto.RequestEntityWithParty) (
	*dto.ResponseData[dto.DocTermsDto], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.ucase.GetDocTerms(req, *d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.ResponseData[dto.DocTermsDto]
	response.Body.Result = res
	return &response, nil
}

func (h *documentInfoHandler) GetAddressAndContact(ctx context.Context, d *dto.RequestEntityWithParty) (
	*dto.ResponseData[dto.AddressAndContactDto], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.ucase.GetAddressAndContact(req, *d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.ResponseData[dto.AddressAndContactDto]
	response.Body.Result = res
	return &response, nil
}

func (h *documentInfoHandler) EditAddressAndContact(ctx context.Context, d *dto.EditAddressAndContactRequest) (
	*dto.ResponseMessage, error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = h.ucase.EditAddressAndContact(req, d.Body)
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
