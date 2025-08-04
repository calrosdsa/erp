package party

// import (
// 	"context"
// 	"erp/api/dto"
// 	"erp/internal/domain"
// 	"erp/internal/app/domain/repository"
// 	"erp/internal/app/service/helpers"
// 	"erp/internal/app/service/services"
// 	"erp/internal/app/service/services/account_service"
// 	"fmt"
// 	"net/http"

// 	"github.com/danielgtaylor/huma/v2"
// )

// type PartyAddressHandler struct {
// 	sessionHelper helpers.SessionHelper
// 	partyAddressService repository.PartyAddressService
// 	errorHelper helpers.ErrorHelper
// 	locale helpers.Locale
// 	roleService *account_service.RoleService
// }

// func NewPartyAddressHandler(
// 	api *huma.API,
// 	services *services.Services,
// 	helpers *helpers.Helpers,
// 	base string,
// 	tag []string,
// 	middlewares huma.Middlewares,
// ) {
// 	paths := NewPartyAddressPaths(base)
// 	handler := PartyAddressHandler{
// 		sessionHelper: helpers.Session,
// 		errorHelper: helpers.Error,
// 		partyAddressService: services.PartyServices.PartyAddress,
// 		locale:helpers.Locale,
// 		roleService: services.RoleService,
// 	}
// 	huma.Register(*api, huma.Operation{
// 		OperationID:   "create address",
// 		Method:        http.MethodPost,
// 		Path:          paths.Base,
// 		Summary:       "Create Address",
// 		Tags:          tag,
// 		DefaultStatus: http.StatusOK,
// 		Middlewares:middlewares,
// 	}, handler.CreatePartyAddress)

// 	huma.Register(*api, huma.Operation{
// 		OperationID:   "get addresses",
// 		Method:        http.MethodGet,
// 		Path:          paths.Base,
// 		Summary:       "Get Addresses",
// 		Tags:          tag,
// 		DefaultStatus: http.StatusOK,
// 		Middlewares:middlewares,
// 	}, handler.GetAddresses)

// 	huma.Register(*api, huma.Operation{
// 		OperationID:   "get address",
// 		Method:        http.MethodGet,
// 		Path:          paths.Detail,
// 		Summary:       "Get Address",
// 		Tags:          tag,
// 		DefaultStatus: http.StatusOK,
// 		Middlewares:middlewares,
// 	}, handler.GetAddress)

// 	huma.Register(*api, huma.Operation{
// 		OperationID:   "get address references",
// 		Method:        http.MethodGet,
// 		Path:          paths.References,
// 		Summary:       "Get Address References",
// 		Tags:          tag,
// 		DefaultStatus: http.StatusOK,
// 		Middlewares:middlewares,
// 	}, handler.GetAddressReferences)
// }

// func (h *PartyAddressHandler)GetAddressReferences(ctx context.Context,i *struct{
// 	dto.AuthParams
// })(
// 	*dto.EntityResponse[dto.ResultEntity[[]dto.PartyTypeDto]],error){
// 	req,_ := h.sessionHelper.GetSession(ctx)	
// 	res :=h.partyAddressService.GetAllowedPartiesForAddress(req)
// 	var response dto.EntityResponse[dto.ResultEntity[[]dto.PartyTypeDto]]
// 	response.Body.Result =res
// 	return &response,nil
// }
// func (h *PartyAddressHandler)GetAddress(ctx context.Context,i *dto.RequestEntity)(
// 	*dto.EntityResponse[dto.ResultEntity[dto.AddressDto]],error){
// 	req,_ := h.sessionHelper.GetSession(ctx)	
// 	res,err:=h.partyAddressService.GetAddress(req,i)
// 	if err != nil {
// 		return nil,h.errorHelper.HumaCustomError(string(req.LanguageCode),err)
// 	}
// 	actions := h.roleService.GetActions(req, domain.ADDRESS)
// 	var response dto.EntityResponse[dto.ResultEntity[dto.AddressDto]]
// 	response.Body.Result =res
// 	response.Body.Actions = actions
// 	return &response,nil
// }   
    
// func (h *PartyAddressHandler) GetAddresses(ctx context.Context, i *dto.RequestPaginationData) (
// 	*dto.PaginationResponse[dto.PaginationResult[[]dto.AddressDto]], error,
// ) {
// 	req, err := h.sessionHelper.GetSession(ctx)
// 	if err != nil {
// 		return nil, huma.Error400BadRequest("Not Authorized", err)
// 	}
// 	// h.sessionHelper.AppendPaginationParams(req, i)
// 	res, err := h.partyAddressService.GetAddresses(req, i)
// 	if err != nil {
// 		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode),err)
// 	}
// 	actions := h.roleService.GetActions(req, domain.ADDRESS)
// 	fmt.Println(actions)
// 	response := &dto.PaginationResponse[dto.PaginationResult[[]dto.AddressDto]]{}
// 	response.Body.PaginationResult = res
// 	response.Body.Actions = actions
// 	return response, err
// }


// func (h *PartyAddressHandler) CreatePartyAddress(ctx context.Context,i *dto.CreatePartyAddressRequest)(
// 	*dto.ResponseMessage,error){
// 	req,_:= h.sessionHelper.GetSession(ctx)	
// 	err := h.partyAddressService.CreatePartyAddress(req,i)
// 	if  err != nil {
// 		return nil,h.errorHelper.HumaCustomError(string(req.LanguageCode),err,"Error.FailToCreateAddress")
// 	}
// 	var response dto.ResponseMessage
// 	response.Body.Message = h.locale.MustLocalize(
// 		helpers.OptionsLocale.WithID("Message.CreateAddressSuccess"),
// 		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
// 	)
// 	return &response,nil
// }

// // func (h *PartyAddressHandler) GetUserPartyTypes(ctx context.Context,i *struct{dto.AuthParams})(
// // 	*dto.EntityResponse[dto.ResultEntity[[]dto.PartyTypeDto]],error){
// // 	req,_ := h.sessionHelper.GetSession(ctx)	
// // 	res := h.partyService.GetUserPartyTypes(req)
// // 	var response dto.EntityResponse[dto.ResultEntity[[]dto.PartyTypeDto]]
// // 	response.Body.Result = res
// // 	return &response,nil
// // }
