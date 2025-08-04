package order_rest

import (
	"context"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	order_usecase "erp/project/order/usecase"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type OrderHandler struct {
	sessionHelper helpers.SessionHelper
	errorHelper   helpers.ErrorHelper
	locale        helpers.Locale
	permission    repository.PermissionService
	orderUseCase  order_usecase.OrderUseCase
}

func NewOrderHandler(
	api huma.API,
	helpers *helpers.Helpers,
	middlewares huma.Middlewares,
	permission repository.PermissionService,
	orderUseCase order_usecase.OrderUseCase,
) {
	// base :=
	paths := NewOrderPaths(domain.ORDER_BASE_ROUTE)
	tags := []string{"Order"}
	handler := OrderHandler{
		sessionHelper: helpers.Session,
		errorHelper:   helpers.Error,
		orderUseCase:  orderUseCase,
		locale:        helpers.Locale,
		permission:    permission,
	}
	huma.Register(api, huma.Operation{
		OperationID:   "orders",
		Method:        http.MethodGet,
		Path:          paths.Type,
		Summary:       "Retrieve orders",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetOrders)

	huma.Register(api, huma.Operation{
		OperationID:   "order",
		Method:        http.MethodGet,
		Path:          paths.Detail,
		Summary:       "Retrieve order",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetOrder)

	huma.Register(api, huma.Operation{
		OperationID:   "create-order",
		Method:        http.MethodPost,
		Path:          paths.Base,
		Summary:       "Create Order",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.CreateOrder)

	huma.Register(api, huma.Operation{
		OperationID:   "edit-order",
		Method:        http.MethodPut,
		Path:          paths.Base,
		Summary:       "Edit Order",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.EditOrder)

	huma.Register(api, huma.Operation{
		OperationID:   "update-order-status",
		Method:        http.MethodPut,
		Path:          paths.UpdateStatus,
		Summary:       "Update Order Status",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.UpdateOrderStatus)

	huma.Register(api, huma.Operation{
		OperationID:   "export-order",
		Method:        http.MethodPost,
		Path:          paths.Document,
		Summary:       "Export Order",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.ExportOrder)
}

func (h *OrderHandler) ExportOrder(ctx context.Context, i *dto.ExportDocumentRequest) (*huma.StreamResponse, error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	bytes, err := h.orderUseCase.ExportOrder(req, i.Body)
	// fmt.Println("START STREAM RESPONSE",err)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	return &huma.StreamResponse{
		Body: func(ctx huma.Context) {
			writer := ctx.BodyWriter()
			writer.Write(bytes)
		},
	}, nil
}

func (h *OrderHandler) UpdateOrderStatus(ctx context.Context, i *dto.UpdateStatusWithEvent) (
	*dto.ResponseMessage, error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = h.orderUseCase.UpdateOrderStatus(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToUpdateOrder")
		// return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.ResponseMessage
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.UpdateOrderSuccess"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}

func (h *OrderHandler) EditOrder(ctx context.Context, d *dto.EditOrderRequest) (
	*dto.ResponseMessage, error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	err := h.orderUseCase.EditOrder(req, d.Body)
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

func (h *OrderHandler) CreateOrder(ctx context.Context, i *dto.CreateOrderRequest) (
	*dto.EntityResponse[dto.ResultEntity[dto.OrderDto]], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.orderUseCase.CreateOrder(req, i.Body)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToCreateOrder")
		// return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.EntityResponse[dto.ResultEntity[dto.OrderDto]]
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.CreatePurchaseOrderSuccess"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	response.Body.Result = res
	return &response, nil
}

func (h *OrderHandler) GetOrder(ctx context.Context, i *dto.RequestEntityWithParty) (
	*dto.EntityResponse[dto.ResultEntity[dto.OrderDetailDto]], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.orderUseCase.GetOrder(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	entity, err := h.orderUseCase.GetEntityOrder(i.PartyType)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	actions := h.permission.GetActions(req.Ctx, entity.ID)
	var response dto.EntityResponse[dto.ResultEntity[dto.OrderDetailDto]]
	response.Body.Result = res
	response.Body.Actions = actions
	response.Body.AssociatedActions = h.getExtraActions(ctx, entity)
	return &response, nil
}

func (h *OrderHandler) GetOrders(ctx context.Context, i *dto.RequestOrders) (
	*dto.PaginationResponse[dto.PaginationResult[[]dto.OrderDto]], error,
) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	// h.sessionHelper.AppendPaginationParams(req, i)
	res, err := h.orderUseCase.GetOrders(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	entity, err := h.orderUseCase.GetEntityOrder(i.PartyType)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	actions := h.permission.GetActions(req.Ctx, entity.ID)

	response := &dto.PaginationResponse[dto.PaginationResult[[]dto.OrderDto]]{}
	response.Body.PaginationResult = res
	response.Body.Actions = actions
	return response, err
}

func (h *OrderHandler) getExtraActions(ctx context.Context, entity domain.EntityTemplate) map[int][]dto.ActionDto {
	var ids []int64
	switch entity {
	case domain.PURCHASE_ORDER:
		ids = append(ids, domain.PURCHASE_INVOICE.ID, domain.PURCHASE_RECEIPT.ID)
	case domain.SALE_ORDER:
		ids = append(ids, domain.SALE_INVOICE.ID, domain.DELIVERY_NOTE.ID)

	}
	ids = append(ids, domain.PAYMENT.ID, domain.ADDRESS.ID, domain.CONTACT.ID, 
		domain.PAYMENT_TERMS_TEMPLATE.ID, domain.TERMS_AND_CONDITIONS.ID,domain.LEDGER.ID)

	r := h.permission.GetEntitiesActions(ctx, ids)
	return r
}
