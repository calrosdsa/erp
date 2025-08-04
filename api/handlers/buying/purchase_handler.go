package buying

import (
	"context"
	"erp/api/dto"
	"erp/internal/app/domain/repository"
	"erp/internal/app/service/helpers"
	"erp/internal/app/service/services"
	"erp/internal/app/service/services/account_service"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type PurchaseHandler struct {
	sessionHelper   helpers.SessionHelper
	errorHelper     helpers.ErrorHelper
	locale          helpers.Locale
	roleService     *account_service.RoleService
	purchaseService repository.PurchaseService
}

func NewPurchaseHandler(
	api *huma.API,
	services *services.Services,
	helpers *helpers.Helpers,
	base string,
	tag []string,
	middlewares huma.Middlewares,
) {
	paths := NewPurchasePaths(base)
	handler := PurchaseHandler{
		sessionHelper:   helpers.Session,
		errorHelper:     helpers.Error,
		purchaseService: services.Buying.PurchaseService,
		locale:          helpers.Locale,
		roleService:     services.RoleService,
	}
	// huma.Register(*api, huma.Operation{
	// 	OperationID:   "get supplier",
	// 	Method:        http.MethodGet,
	// 	Path:          paths.Detail,
	// 	Summary:       "Retrieve supplier",
	// 	Tags:          tag,
	// 	DefaultStatus: http.StatusOK,
	// 	Middlewares:   middlewares,
	// }, handler.GetSupplier)
	// huma.Register(*api, huma.Operation{
	// 	OperationID:   "get suppliers",
	// 	Method:        http.MethodGet,
	// 	Path:          paths.Base,
	// 	Summary:       "Retrieve suppliers",
	// 	Tags:          tag,
	// 	DefaultStatus: http.StatusOK,
	// 	Middlewares:   middlewares,
	// }, handler.GetSuppliers)

	huma.Register(*api, huma.Operation{
		OperationID:   "create purchase order",
		Method:        http.MethodPost,
		Path:          paths.Base,
		Summary:       "Create purchase order",
		Tags:          tag,
		DefaultStatus: http.StatusCreated,
		Middlewares:   middlewares,
	}, handler.CreatePurchaseOrder)

}

// func (h *PurchaseHandler)GetSupplier(ctx context.Context,i *dto.RequestEntity)(
// 	*dto.EntityResponse[dto.ResultEntity[dto.SupplierDto]],error){
// 	req,_ := h.sessionHelper.GetSession(ctx)
// 	res,err:=h.purchaseService.GetSupplier(req,i)
// 	if err != nil {
// 		return nil,h.errorHelper.HumaCustomError(string(req.LanguageCode),err)
// 	}
// 	actions := h.roleService.GetActions(req, domain.SUPPLIER)
// 	var response dto.EntityResponse[dto.ResultEntity[dto.SuplierDto]]
// 	response.Body.Result =res
// 	response.Body.Actions = actions
// 	return &response,nil
// }

func (h *PurchaseHandler) CreatePurchaseOrder(ctx context.Context, i *dto.CreatePurchaseOrderRequest) (
	*dto.EntityResponse[dto.ResultEntity[dto.OrderDto]], error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	res,err := h.purchaseService.CreatePurchaseOrder(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToCreatePurchaseOrder")
	}
	var response dto.EntityResponse[dto.ResultEntity[dto.OrderDto]]
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.CreatePurchaseOrderSuccess"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	response.Body.Result = res
	return &response, nil
}

// func (h *PurchaseHandler) GetSuppliers(ctx context.Context, i *dto.RequestPaginationData) (
// 	*dto.PaginationResponse[dto.PaginationResult[[]dto.SupplierDto]], error,
// ) {
// 	req, err := h.sessionHelper.GetSession(ctx)
// 	if err != nil {
// 		return nil, huma.Error400BadRequest("Not Authorized", err)
// 	}
// 	// h.sessionHelper.AppendPaginationParams(req, i)
// 	res, err := h.purchaseService.GetSuppliers(req, i)
// 	if err != nil {
// 		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode),err)
// 	}
// 	actions := h.roleService.GetActions(req, domain.SUPPLIER)

// 	response := &dto.PaginationResponse[dto.PaginationResult[[]dto.SupplierDto]]{}
// 	response.Body.PaginationResult = res
// 	response.Body.Actions = actions
// 	return response, err
// }
