package contact_rest

import (
	"context"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	contact_ucase "erp/project/core/contact/usecase"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type ContactHandler struct {
	sessionHelper  helpers.SessionHelper
	errorHelper    helpers.ErrorHelper
	locale         helpers.Locale
	permission     repository.PermissionService
	contactUseCase contact_ucase.ContactUseCase
}

func NewContactHandler(
	api huma.API,
	helpers *helpers.Helpers,
	middlewares huma.Middlewares,
	contactUseCase contact_ucase.ContactUseCase,
	permission repository.PermissionService,
) {
	base := domain.CONTACT_BASE_ROUTE
	paths := NewContactPaths(domain.CONTACT_BASE_ROUTE)
	tag := []string{"Contact"}
	handler := ContactHandler{
		sessionHelper:  helpers.Session,
		errorHelper:    helpers.Error,
		contactUseCase: contactUseCase,
		locale:         helpers.Locale,
		permission:     permission,
	}
	huma.Register(api, huma.Operation{
		OperationID:   "create-contact",
		Method:        http.MethodPost,
		Path:          paths.Base,
		Summary:       "Create Contact",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.CreateContact)

	huma.Register(api, huma.Operation{
		OperationID:   "contacts",
		Method:        http.MethodGet,
		Path:          paths.Base,
		Summary:       "Get Contacts",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetContacts)

	huma.Register(api, huma.Operation{
		OperationID:   "contact",
		Method:        http.MethodGet,
		Path:          paths.Detail,
		Summary:       "Get Contact",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetContact)
	
	huma.Register(api, huma.Operation{
		OperationID:   "edit-contact",
		Method:        http.MethodPut,
		Path:          paths.Base,
		Summary:       "Edit contact",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.EditCustomer)

	huma.Register(api, huma.Operation{
		OperationID:   "contact-bulk",
		Method:        http.MethodPost,
		Path:          base + "/bulk",
		Summary:       "Contact Bulk",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.ContactBulk)

}
func (h *ContactHandler) ContactBulk(ctx context.Context, d *dto.ContactBulkDataRequest) (
	*dto.ResponseMessage, error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	err := h.contactUseCase.ContactBulk(req, d.Body)
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

func (h *ContactHandler) EditCustomer(ctx context.Context, d *dto.ContactDataRequest) (
	*dto.ResponseMessage, error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	err := h.contactUseCase.EditContact(req, d.Body)
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

func (h *ContactHandler) GetContact(ctx context.Context, i *dto.RequestEntity) (
	*dto.EntityResponse[dto.ResultEntity[dto.ContactDto]], error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	res, err := h.contactUseCase.GetContact(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	actions := h.permission.GetActions(req.Ctx, domain.CONTACT.ID)
	var response dto.EntityResponse[dto.ResultEntity[dto.ContactDto]]
	response.Body.Result = res
	response.Body.Actions = actions
	return &response, nil
}

func (h *ContactHandler) GetContacts(ctx context.Context, i *dto.ContactsRequest) (
	*dto.ResponseDataList[[]dto.ContactDto], error,
) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	// h.sessionHelper.AppendPaginationParams(req, i)
	res, err := h.contactUseCase.GetContacts(req, *i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	actions := h.permission.GetActions(req.Ctx, domain.CONTACT.ID)
	fmt.Println(actions)
	var response dto.ResponseDataList[[]dto.ContactDto]
	response.Body.Result = res
	response.Body.Actions = actions
	return &response, err
}

func (h *ContactHandler) CreateContact(ctx context.Context, i *dto.ContactDataRequest) (
	*dto.ResponseData[dto.ContactDto], error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	res, err := h.contactUseCase.CreateContact(req, i.Body)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToCreateContact")
	}
	var response dto.ResponseData[dto.ContactDto]
	response.Body.Result = res
	return &response, nil
}

// func (h *ContactHandler) GetUserPartyTypes(ctx context.Context,i *struct{dto.AuthParams})(
// 	*dto.EntityResponse[dto.ResultEntity[[]dto.PartyTypeDto]],error){
// 	req,_ := h.sessionHelper.GetSession(ctx)
// 	res := h.partyService.GetUserPartyTypes(req)
// 	var response dto.EntityResponse[dto.ResultEntity[[]dto.PartyTypeDto]]
// 	response.Body.Result = res
// 	return &response,nil
// }
