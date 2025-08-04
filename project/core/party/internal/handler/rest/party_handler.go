package party_rest

import (
	"context"
	"erp/api/dto"
	
	"erp/internal/app/service/helpers"
	
	"erp/internal/domain"
	party_ucase "erp/project/core/party/internal/usecase"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type PartyHandler struct {
	sessionHelper helpers.SessionHelper
	partyUseCase  party_ucase.PartyUseCase
	errorHelper   helpers.ErrorHelper
	locale        helpers.Locale
}

func NewPartyHandler(
	api huma.API,
	helpers *helpers.Helpers,
	middlewares huma.Middlewares,
	partyUseCase party_ucase.PartyUseCase,
) {
	paths := NewPartyPaths(domain.PARTY_BASE_ROUTE)
	tag := []string{"Party"}
	handler := PartyHandler{
		sessionHelper: helpers.Session,
		partyUseCase:  partyUseCase,
		errorHelper:   helpers.Error,
		locale:        helpers.Locale,
	}
	huma.Register(api, huma.Operation{
		OperationID:   "get party user types",
		Method:        http.MethodGet,
		Path:          paths.TypeUsers,
		Summary:       "Get Party User Types",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetUserPartyTypes)

	huma.Register(api, huma.Operation{
		OperationID:   "get party by reference",
		Method:        http.MethodGet,
		Path:          paths.PartyByReferences,
		Summary:       "Get Party by reference",
		Description: "Retrieve parties by party type.",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetPartiesByReference)

	huma.Register(api, huma.Operation{
		OperationID:   "get party references",
		Method:        http.MethodGet,
		Path:          paths.References,
		Summary:       "Get Party references",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetPartyReferences)
	huma.Register(api, huma.Operation{
		OperationID:   "add party reference",
		Method:        http.MethodPost,
		Path:          paths.References,
		Summary:       "add party reference",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.AddPartyReferences)

	huma.Register(api, huma.Operation{
		OperationID:   "get party references",
		Method:        http.MethodGet,
		Path:          paths.References,
		Summary:       "Get Party references",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetPartyReferences)

	huma.Register(api, huma.Operation{
		OperationID:   "get party type references",
		Method:        http.MethodGet,
		Path:          paths.ReferencesType,
		Summary:       "Get Party type references",
		Description:   "Retrieve the allowed party types for reference",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetPartyTypeForReferences)


	huma.Register(api, huma.Operation{
		OperationID:   "get party connections",
		Method:        http.MethodGet,
		Path:          paths.Connections,
		Summary:       "Get Party connections",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetPartyConnections)
}

func(h *PartyHandler) GetPartyConnections(ctx context.Context,i *dto.RequestEntityWithParty)(
 *dto.ResponseData[[]dto.PartyConnections],error){
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.partyUseCase.GetPartyConnections(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.ResponseData[[]dto.PartyConnections]
	response.Body.Result = res
	return &response, nil
}

func (h *PartyHandler) GetPartyTypeForReferences(ctx context.Context, i *struct {
	dto.AuthParams
}) (
	*dto.ResponseData[[]dto.PartyTypeDto], error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	res := h.partyUseCase.GetAllowedPartiesForReferences(req)
	var response dto.ResponseData[[]dto.PartyTypeDto]
	response.Body.Result = res
	return &response, nil
}

func (h *PartyHandler) AddPartyReferences(ctx context.Context, i *dto.RequestAddPartyReference) (
	*dto.ResponseMessage, error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	err := h.partyUseCase.AddPartyReference(req, i)
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
	res, err := h.partyUseCase.GetPartyReferences(req, i)
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
	res, err := h.partyUseCase.GetPartiesByReference(req, i)
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
	res := h.partyUseCase.GetUserPartyTypes(req)
	var response dto.EntityResponse[dto.ResultEntity[[]dto.PartyTypeDto]]
	response.Body.Result = res
	return &response, nil
}
