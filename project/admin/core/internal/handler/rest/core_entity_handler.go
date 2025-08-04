package rest_core

import (
	"context"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	core_ucase "erp/project/admin/core/internal/usecase"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type CoreEntityHandler struct {
	jwt               helpers.JwtHelper
	sessionHelper     helpers.SessionHelper
	locale            helpers.Locale
	errorHelper       helpers.ErrorHelper
	coreEntityUseCase core_ucase.CoreEntityUseCase
}

func NewCoreEntityHandler(
	api huma.API,
	helpers *helpers.Helpers,
	coreEntityUseCase core_ucase.CoreEntityUseCase,
	middlewares huma.Middlewares,
) {
	base := domain.A_CORE_ENTITY_BASE_ROUTE
	tags := []string{"Core Entity"}
	path := NewCorePaths(base)
	h := CoreEntityHandler{
		jwt:               helpers.Jwt,
		sessionHelper:     helpers.Session,
		locale:            helpers.Locale,
		errorHelper:       helpers.Error,
		coreEntityUseCase: coreEntityUseCase,
	}
	huma.Register(api, huma.Operation{
		OperationID:   "entities",
		Method:        http.MethodGet,
		Summary:       "Entities",
		Tags:          tags,
		Path:          path.Base,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetCoreEntities)

	huma.Register(api, huma.Operation{
		OperationID:   "create-entity",
		Method:        http.MethodPost,
		Summary:       "Create Entity",
		Tags:          tags,
		Path:          path.Base,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.CreateEntity)

	huma.Register(api, huma.Operation{
		OperationID:   "entity",
		Method:        http.MethodGet,
		Summary:       "Entity",
		Tags:          tags,
		Path:          path.Detail,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetEntity)

	huma.Register(api, huma.Operation{
		OperationID:   "add-entity-action",
		Method:        http.MethodPost,
		Summary:       "Add Entity Action",
		Tags:          tags,
		Path:          path.Action,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.AddEntityAction)
}
func (h *CoreEntityHandler) AddEntityAction(ctx context.Context, d *dto.AddEntityActionRequest) (
	*dto.ResponseMessage, error) {
	req, err := h.sessionHelper.GetAdminSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = h.coreEntityUseCase.AddEntityAction(req, d)
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

func (h *CoreEntityHandler) CreateEntity(ctx context.Context, d *dto.CreateEntityRequest) (
	*dto.ResponseMessage, error) {
	req, err := h.sessionHelper.GetAdminSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = h.coreEntityUseCase.CreateEntity(req, d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToCreate")
	}
	var response dto.ResponseMessage
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.CreatedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}

func (h *CoreEntityHandler) GetEntity(ctx context.Context, d *dto.RequestEntity) (
	*dto.EntityResponse[dto.ResultEntity[dto.EntityDetailDto]], error) {
	req, err := h.sessionHelper.GetAdminSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.coreEntityUseCase.GetEntity(req, d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.EntityResponse[dto.ResultEntity[dto.EntityDetailDto]]
	response.Body.Result = res
	return &response, nil
}

func (h *CoreEntityHandler) GetCoreEntities(ctx context.Context, d *dto.RequestPaginationData) (
	*dto.PaginationResponse[dto.PaginationResult[[]dto.EntityDto]], error) {
	req, err := h.sessionHelper.GetAdminSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.coreEntityUseCase.GetEntities(req, d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.PaginationResponse[dto.PaginationResult[[]dto.EntityDto]]
	response.Body.PaginationResult = res
	return &response, nil
}
