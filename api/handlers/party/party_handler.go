package party

import (
	"context"
	"erp/api/dto"
	"erp/internal/app/domain/repository"
	"erp/internal/app/service/helpers"
	"erp/internal/app/service/services"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type PartyHandler struct {
	sessionHelper helpers.SessionHelper
	partyService  repository.PartyService
	errorHelper   helpers.ErrorHelper
	locale        helpers.Locale
}

func NewPartyHandler(
	api *huma.API,
	services *services.Services,
	helpers *helpers.Helpers,
	base string,
	tag []string,
	middlewares huma.Middlewares,
) {
	paths := NewPartyPaths(base)
	handler := PartyHandler{
		sessionHelper: helpers.Session,
		partyService:  services.PartyServices.PartyService,
		errorHelper:   helpers.Error,
		locale:        helpers.Locale,
	}
	huma.Register(*api, huma.Operation{
		OperationID:   "get party user types",
		Method:        http.MethodGet,
		Path:          paths.TypeUsers,
		Summary:       "Get Party User Types",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetUserPartyTypes)

	huma.Register(*api, huma.Operation{
		OperationID:   "get party by reference",
		Method:        http.MethodGet,
		Path:          paths.PartyByReferences,
		Summary:       "Get Party by reference",
		Description: "Retrieve parties by party type.",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetPartiesByReference)

	huma.Register(*api, huma.Operation{
		OperationID:   "get party references",
		Method:        http.MethodGet,
		Path:          paths.References,
		Summary:       "Get Party references",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetPartyReferences)
	huma.Register(*api, huma.Operation{
		OperationID:   "add party reference",
		Method:        http.MethodPost,
		Path:          paths.References,
		Summary:       "add party reference",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.AddPartyReferences)

	huma.Register(*api, huma.Operation{
		OperationID:   "get party references",
		Method:        http.MethodGet,
		Path:          paths.References,
		Summary:       "Get Party references",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetPartyReferences)

	huma.Register(*api, huma.Operation{
		OperationID:   "get party type references",
		Method:        http.MethodGet,
		Path:          paths.ReferencesType,
		Summary:       "Get Party type references",
		Description:   "Retrieve the allowed party types for reference",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetPartyTypeForReferences)
}

func (h *PartyHandler) GetPartyTypeForReferences(ctx context.Context, i *struct {
	dto.AuthParams
}) (
	*dto.ResponseData[[]dto.PartyTypeDto], error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	res := h.partyService.GetAllowedPartiesForReferences(req)
	var response dto.ResponseData[[]dto.PartyTypeDto]
	response.Body.Result = res
	return &response, nil
}

func (h *PartyHandler) AddPartyReferences(ctx context.Context, i *dto.RequestAddPartyReference) (
	*dto.ResponseMessage, error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	err := h.partyService.AddPartyReference(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToAdd")
	}
	var response dto.ResponseMessage
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.AddSuccess"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}

func (h *PartyHandler) GetPartyReferences(ctx context.Context, i *dto.RequestPaginationData) (
	*dto.PaginationResponse[dto.PaginationResult[[]dto.PartyReferenceDto]], error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	res, err := h.partyService.GetPartyReferences(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.PaginationResponse[dto.PaginationResult[[]dto.PartyReferenceDto]]
	response.Body.PaginationResult = res
	return &response, nil
}

func (h *PartyHandler) GetPartiesByReference(ctx context.Context, i *dto.RequestPartyReference) (
	*dto.EntityResponse[dto.ResultEntity[[]dto.PartyDto]], error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	res, err := h.partyService.GetPartiesByReference(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.EntityResponse[dto.ResultEntity[[]dto.PartyDto]]
	response.Body.Result = res
	return &response, nil
}

func (h *PartyHandler) GetUserPartyTypes(ctx context.Context, i *struct{ dto.AuthParams }) (
	*dto.EntityResponse[dto.ResultEntity[[]dto.PartyTypeDto]], error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	res := h.partyService.GetUserPartyTypes(req)
	var response dto.EntityResponse[dto.ResultEntity[[]dto.PartyTypeDto]]
	response.Body.Result = res
	return &response, nil
}
