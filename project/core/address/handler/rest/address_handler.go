package address_rest

import (
	"context"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	address_ucase "erp/project/core/address/usecase"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type addressHandler struct {
	sessionHelper helpers.SessionHelper
	ucase         address_ucase.AddressUseCase
	errorHelper   helpers.ErrorHelper
	locale        helpers.Locale
	permission    repository.PermissionService
}

func NewAddressHandler(
	api huma.API,
	helpers *helpers.Helpers,
	middlewares huma.Middlewares,
	ucase address_ucase.AddressUseCase,
	permission repository.PermissionService,
) {
	paths := NewPaths(domain.ADDRESS_BASE_ROUTE)
	tags := []string{"Address"}
	handler := addressHandler{
		sessionHelper: helpers.Session,
		errorHelper:   helpers.Error,
		ucase:         ucase,
		locale:        helpers.Locale,
		permission:    permission,
	}
	huma.Register(api, huma.Operation{
		OperationID:   "create-address",
		Method:        http.MethodPost,
		Path:          paths.Base,
		Summary:       "Create Address",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.CreateAddress)

	huma.Register(api, huma.Operation{
		OperationID:   "edit-address",
		Method:        http.MethodPut,
		Path:          paths.Base,
		Summary:       "Edit Address",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.EditAddress)

	huma.Register(api, huma.Operation{
		OperationID:   "get-addresses",
		Method:        http.MethodGet,
		Path:          paths.Base,
		Summary:       "Get Addresses",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetAddresses)

	huma.Register(api, huma.Operation{
		OperationID:   "get-address",
		Method:        http.MethodGet,
		Path:          paths.Detail,
		Summary:       "Get Address",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetAddress)

	huma.Register(api, huma.Operation{
		OperationID:   "get-address-references",
		Method:        http.MethodGet,
		Path:          paths.References,
		Summary:       "Get Address References",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetAddressReferences)

	huma.Register(api, huma.Operation{
		OperationID:   "update-status-address",
		Method:        http.MethodPut,
		Path:          paths.Base + "/update-status",
		Summary:       "Update Status Address",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.UpdateStatus)
}

func (m *addressHandler) UpdateStatus(ctx context.Context, d *dto.UpdateStatusWithEvent) (_a0 *dto.ResponseMessage, _a1 error) {

	req, err := m.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = m.ucase.UpdateStatus(req, *d)
	if err != nil {
		return nil, m.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.ErrorMsg")
	}
	var response dto.ResponseMessage
	response.Body.Message = m.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.UpdatedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil

}

func (h *addressHandler) GetAddressReferences(ctx context.Context, i *struct {
	dto.AuthParams
}) (
	*dto.EntityResponse[dto.ResultEntity[[]dto.PartyTypeDto]], error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	res := h.ucase.GetAllowedPartiesForAddress(req)
	var response dto.EntityResponse[dto.ResultEntity[[]dto.PartyTypeDto]]
	response.Body.Result = res
	return &response, nil
}
func (h *addressHandler) GetAddress(ctx context.Context, i *dto.RequestEntity) (
	*dto.EntityResponse[dto.ResultEntity[dto.AddressDto]], error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	res, err := h.ucase.GetAddress(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	actions := h.permission.GetActions(ctx, domain.ADDRESS.ID)
	var response dto.EntityResponse[dto.ResultEntity[dto.AddressDto]]
	response.Body.Result = res
	response.Body.Actions = actions
	return &response, nil
}

func (h *addressHandler) GetAddresses(ctx context.Context, i *dto.RequestAddresses) (
	*dto.ResponseDataList[[]dto.AddressDto], error,
) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	// h.sessionHelper.AppendPaginationParams(req, i)
	res, err := h.ucase.GetAddresses(req, *i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	actions := h.permission.GetActions(ctx, domain.ADDRESS.ID)
	fmt.Println(actions)
	res.Body.Actions = actions
	return &res, err
}

func (h *addressHandler) CreateAddress(ctx context.Context, i *dto.AddressDataRequest) (
	*dto.ResponseData[dto.AddressDto], error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	res, err := h.ucase.CreateAddress(req, i.Body)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.ErrorMsg")
	}
	var response dto.ResponseData[dto.AddressDto]
	response.Body.Result = res
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.SuccessfullyMessage"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}

func (h *addressHandler) EditAddress(ctx context.Context, i *dto.AddressDataRequest) (*dto.ResponseMessage, error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	err := h.ucase.EditAddress(req, i.Body)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.ErrorMsg")
	}
	var response dto.ResponseMessage
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.SuccessfullyMessage"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}

// func (h  *addressHandler) GetUserPartyTypes(ctx context.Context,i *struct{dto.AuthParams})(
// 	*dto.EntityResponse[dto.ResultEntity[[]dto.PartyTypeDto]],error){
// 	req,_ := h.sessionHelper.GetSession(ctx)
// 	res := h.partyService.GetUserPartyTypes(req)
// 	var response dto.EntityResponse[dto.ResultEntity[[]dto.PartyTypeDto]]
// 	response.Body.Result = res
// 	return &response,nil
// }
