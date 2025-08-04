package itemprice_rest

import (
	"context"
	"erp/api/dto"
	"erp/gen/proto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	itemprice_ucase "erp/project/stock/itemprice/usecase"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type ItemPriceHandler struct {
	itemPriceUcase itemprice_ucase.ItemPriceUseCase
	sessionHelper  helpers.SessionHelper
	locale         helpers.Locale
	errorHelper    helpers.ErrorHelper
	permission     repository.PermissionService
}

func NewItemPriceHandler(
	api huma.API,
	helpers *helpers.Helpers,
	middlewares huma.Middlewares,
	permission repository.PermissionService,
	itemPriceUcase itemprice_ucase.ItemPriceUseCase,
) {
	paths := NewItemPricePaths(domain.ITEMPRICE_BASE_ROUTE)
	tags := []string{"Item Price"}
	handler := ItemPriceHandler{
		itemPriceUcase: itemPriceUcase,
		sessionHelper:  helpers.Session,
		locale:         helpers.Locale,
		errorHelper:    helpers.Error,
		permission:     permission,
	}

	huma.Register(api, huma.Operation{
		OperationID:   "get-item-price-detail",
		Method:        http.MethodGet,
		Path:          paths.Detail,
		Summary:       "Get Item Price",
		Description:   "Retrieve Item Price Datail",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetItemPrice)

	huma.Register(api, huma.Operation{
		OperationID:   "get-item-prices",
		Method:        http.MethodGet,
		Path:          paths.Base,
		Summary:       "Get Item Prices",
		Description:   "Retrieve Item Prices",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetItemPrices)

	huma.Register(api, huma.Operation{
		OperationID:   "get-list-by-item",
		Method:        http.MethodGet,
		Path:          paths.Item,
		Summary:       "Get List by item",
		Description:   "Retrieve Item Prices by Item Code",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetItemPricesByItemCode)

	huma.Register(api, huma.Operation{
		OperationID:   "create-item-price",
		Method:        http.MethodPost,
		Path:          paths.Base,
		Summary:       "Create item price",
		Tags:          tags,
		DefaultStatus: http.StatusCreated,
		Middlewares:   middlewares,
	}, handler.CreateItemPrice)

	huma.Register(api, huma.Operation{
		OperationID:   "get-item-prices-for-order",
		Method:        http.MethodGet,
		Path:          paths.Order,
		Summary:       "Get Item Price For Order",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetItemPricesForOrder)

	huma.Register(api, huma.Operation{
		OperationID:   "associated-actions",
		Method:        http.MethodGet,
		Path:          paths.AssociatedActions,
		Summary:       "Get Associated Actions",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetItemPriceActions)

	huma.Register(api, huma.Operation{
		OperationID:   "edit-item-price",
		Method:        http.MethodPut,
		Path:          paths.Base,
		Summary:       "Edit item price",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.EditItemPrice)
}
func (h *ItemPriceHandler) EditItemPrice(ctx context.Context, d *dto.ItemPriceRequestData) (
	*dto.ResponseMessage, error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	err := h.itemPriceUcase.EditItemPrice(req, d.Body)
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

func (h *ItemPriceHandler) GetItemPriceActions(ctx context.Context, i *struct{}) (
	*dto.ResponseData[any], error,
) {
	_, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	var res dto.ResponseData[any]
	res.Body.AssociatedActions = h.getExtraActions(ctx)
	return &res, nil
}

func (h *ItemPriceHandler) GetItemPricesForOrder(ctx context.Context, i *dto.RequestItemPricesForOrder) (
	*dto.EntityResponse[dto.ResultEntity[[]dto.ItemPriceDto]], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.itemPriceUcase.GetItemPricesForOrders(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	actions := h.permission.GetActions(req.Ctx, domain.ITEM.ID)
	var response dto.EntityResponse[dto.ResultEntity[[]dto.ItemPriceDto]]
	response.Body.Result = res
	response.Body.Actions = actions
	return &response, nil
}

func (h *ItemPriceHandler) CreateItemPrice(ctx context.Context, i *dto.ItemPriceRequestData) (
	*dto.ResponseData[dto.ItemPriceDto], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.itemPriceUcase.CreateItemPrice(req, i.Body)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToCreate")
	}
	var response dto.ResponseData[dto.ItemPriceDto]
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.CreatedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	response.Body.Result = res
	return &response, nil
}

func (h *ItemPriceHandler) GetItemPrice(ctx context.Context, i *dto.RequestEntity) (
	*dto.EntityResponse[dto.ResultEntity[dto.ItemPriceDto]], error,
) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.itemPriceUcase.GetItemPrice(req, i)
	if err != nil {
		return nil, huma.Error400BadRequest("Error", err)
	}
	actions := h.permission.GetActions(req.Ctx, domain.ITEM_PRICE.ID)
	response := &dto.EntityResponse[dto.ResultEntity[dto.ItemPriceDto]]{}
	response.Body.Result = res
	response.Body.Actions = actions
	response.Body.AssociatedActions = h.getExtraActions(req.Ctx)
	return response, err
}

func (h *ItemPriceHandler) GetItemPrices(ctx context.Context, i *dto.ItemPricesRequest) (
	*dto.ResponseDataList[[]dto.ItemPriceDto], error,
) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}

	res, err := h.itemPriceUcase.GetItemPrices(req, *i)
	if err != nil {
		return nil, huma.Error400BadRequest("Error", err)
	}
	actions := h.permission.GetActions(req.Ctx, domain.ITEM_PRICE.ID)
	res.Body.Actions = actions
	return &res, err
}

func (h *ItemPriceHandler) GetItemPricesByItemCode(ctx context.Context, i *dto.RequestItemPriceByCode) (
	*dto.PaginationResponse[dto.PaginationResult[[]dto.ItemPriceDto]], error,
) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	paginationParams := dto.RequestPaginationData{
		PaginationParams: i.PaginationParams,
	}
	h.sessionHelper.AppendPaginationParams(req, &paginationParams)
	res, err := h.itemPriceUcase.GetItemPricesByItemCode(req, i)
	if err != nil {
		return nil, huma.Error400BadRequest("Error", err)
	}
	actions := h.permission.GetActions(req.Ctx, domain.ITEM_PRICE.ID)
	response := &dto.PaginationResponse[dto.PaginationResult[[]dto.ItemPriceDto]]{}
	response.Body.PaginationResult = res
	response.Body.Actions = actions
	return response, err
}

// func (h *ItemPriceHandler) CreateItemPrice(ctx context.Context, i *dto.UpsertItemPriceRequest) (*dto.ResponseMessage, error) {
// 	req, _ := h.sessionHelper.GetSession(ctx)
// 	var response dto.ResponseMessage
// 	err := h.itemPriceUcase.UpsertItemPrice(req, i)
// 	// h.locale.MustLocalize()
// 	if err != nil {
// 		return nil, huma.Error400BadRequest("Fail to upsert item price")
// 	}
// 	return &response, err
// }

func (h *ItemPriceHandler) UpdateItemPrice(ctx context.Context, i *dto.UpsertItemPriceRequest) (*dto.ResponseMessage, error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	var response dto.ResponseMessage
	err := h.itemPriceUcase.UpsertItemPrice(req, i)
	// h.locale.MustLocalize()
	if err != nil {
		return nil, huma.Error400BadRequest("Fail to upsert item price")
	}
	return &response, err
}

func (h *ItemPriceHandler) getExtraActions(ctx context.Context) map[int][]dto.ActionDto {
	r := make(map[int][]dto.ActionDto)
	r[int(proto.PartyType_tax)] = h.permission.GetActions(ctx, domain.TAX.ID)
	r[int(proto.PartyType_priceList)] = h.permission.GetActions(ctx, domain.PRICE_LIST.ID)
	return r
}
