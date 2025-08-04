package rest_group

import (
	"context"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	group_ucase "erp/project/group/internal/usecase"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type GroupHandler struct {
	sessionHelper helpers.SessionHelper
	errorHelper   helpers.ErrorHelper
	locale        helpers.Locale
	permission   repository.PermissionService
	groupUcase  group_ucase.GroupUseCase
}

func NewGroupHandler(
	api huma.API,
	helpers *helpers.Helpers,
	middlewares huma.Middlewares,
	groupUseCase group_ucase.GroupUseCase,
	permission repository.PermissionService,
) {
	paths := NewGroupPaths(domain.GROUP_BASE_ROUTE)
	tags := []string{"Group"}
	handler := GroupHandler{
		sessionHelper: helpers.Session,
		errorHelper:   helpers.Error,
		groupUcase:  groupUseCase,
		locale:        helpers.Locale,
		permission:   permission,
	}
	huma.Register(api, huma.Operation{
		OperationID:   "get groups by party code",
		Method:        http.MethodGet,
		Path:          paths.Type,
		Summary:       "Retrieve groups by party code",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetGroups)
	huma.Register(api, huma.Operation{
		OperationID:   "get group",
		Method:        http.MethodGet,
		Path:          paths.Detail,
		Summary:       "Retrieve group",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetGroup)

	huma.Register(api, huma.Operation{
		OperationID:   "create group",
		Method:        http.MethodPost,
		Path:          paths.Base,
		Summary:       "Create Group",
		Tags:          tags,
		DefaultStatus: http.StatusCreated,
		Middlewares:   middlewares,
	}, handler.CreateGroup)

	huma.Register(api, huma.Operation{
		OperationID:   "get group descendents",
		Method:        http.MethodGet,
		Path:          paths.Decendents,
		Summary:       "Get group descendents",
		Tags:          tags,
		DefaultStatus: http.StatusCreated,
		Middlewares:   middlewares,
	}, handler.GetGroupDescendents)

	huma.Register(api, huma.Operation{
		OperationID:   "edit-group",
		Method:        http.MethodPut,
		Path:          paths.Base,
		Summary:       "Edit Group",
		Tags:          tags,
		DefaultStatus: http.StatusCreated,
		Middlewares:   middlewares,
	}, handler.EditGroup)
}



func (h *GroupHandler) EditGroup(ctx context.Context, d *dto.EditGroupRequest) (*dto.ResponseMessage, error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = h.groupUcase.EditGroup(req, d)
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


func (h *GroupHandler) GetGroupDescendents(ctx context.Context, i *dto.RequestEntityWithParty) (
	*dto.EntityResponse[dto.ResultEntity[[]dto.GroupHierarchyDto]], error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	res, err := h.groupUcase.GetGroupDescendents(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.EntityResponse[dto.ResultEntity[[]dto.GroupHierarchyDto]]
	response.Body.Result = res
	return &response, nil
}

func (h *GroupHandler) GetGroup(ctx context.Context, i *dto.RequestEntityWithParty) (
	*dto.EntityResponse[dto.ResultEntity[dto.GroupDto]], error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	res, err := h.groupUcase.GetGroup(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	entity, err := h.groupUcase.GetEntityGroup(i.PartyType)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	actions := h.permission.GetActions(req.Ctx, entity.ID)
	var response dto.EntityResponse[dto.ResultEntity[dto.GroupDto]]
	response.Body.Result = res
	response.Body.Actions = actions
	return &response, nil
}

func (h *GroupHandler) CreateGroup(ctx context.Context, i *dto.CreateGroupRequest) (
	*dto.ResponseMessage, error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	err := h.groupUcase.CreateGroup(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToCreateGroup")
	}
	var response dto.ResponseMessage
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.CreateGroupSuccess"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}

func (h *GroupHandler) GetGroups(ctx context.Context, i *dto.RequestPaginationPartyData) (
	*dto.PaginationResponse[dto.PaginationResult[[]dto.GroupDto]], error,
) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	// h.sessionHelper.AppendPaginationParams(req, i)
	res, err := h.groupUcase.GetGroups(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	entity, err := h.groupUcase.GetEntityGroup(i.PartyType)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	actions := h.permission.GetActions(req.Ctx, entity.ID)

	response := &dto.PaginationResponse[dto.PaginationResult[[]dto.GroupDto]]{}
	response.Body.PaginationResult = res
	response.Body.Actions = actions
	return response, err
}
