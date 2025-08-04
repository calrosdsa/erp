package selling

import (
	"erp/internal/app/service/helpers"
	"erp/internal/app/service/services"

	"github.com/danielgtaylor/huma/v2"
)

// optional code omitted

type AccountHandler struct {
	// salesOrderService *sellingservice.SalesOrderService
	sessionHelper *helpers.SessionHelper
}

func NewSaleOrderHandler(
	api *huma.API,
	services *services.Services,
	helpers *helpers.Helpers,
	base string,
	tags []string,
	commonMiddlewares huma.Middlewares,
) {

	// paths := NewSalesOrderPaths(base)
	// handler := AccountHandler{
	// 	sessionHelper:     helpers.Session,
	// 	// salesOrderService: services.SalesOrderService,
	// }
	// huma.Register(*api, huma.Operation{
	// 	OperationID:   "client-sales-orders",
	// 	Method:        http.MethodGet,
	// 	Path:          paths.ClientOrders,
	// 	Summary:       "Get client sales orders",
	// 	Tags:          tags,
	// 	DefaultStatus: http.StatusOK,
	// 	Middlewares:   commonMiddlewares,
	// }, handler.GetSalesOrders)

	// huma.Register(*api, huma.Operation{
	// 	OperationID:   "sale-order-detail",
	// 	Method:        http.MethodGet,
	// 	Path:          paths.DetailOrder,
	// 	Summary:       "Get sales order detail",
	// 	Tags:          tags,
	// 	DefaultStatus: http.StatusOK,
	// 	Middlewares:   commonMiddlewares,
	// }, handler.GetSalesOrderDetail)
}

// func (h *AccountHandler) GetSalesOrders(ctx context.Context, i *dto.RequestPaginationData) (
// 	*dto.PaginationResponse[dto.PaginationResult[[]entity.SalesOrder]], error,
// ) {
// 	req, err := h.sessionHelper.GetSession(ctx)
// 	if err != nil {
// 		return nil, huma.Error400BadRequest("Not Authorized", err)
// 	}
// 	h.sessionHelper.AppendPaginationParams(req, i)
// 	res, err := h.salesOrderService.GetSalesOrders(req, i)
// 	if err != nil {
// 		fmt.Println(err)
// 		return nil, huma.Error400BadRequest("Error", err)
// 	}
// 	response := &dto.PaginationResponse[dto.PaginationResult[[]entity.SalesOrder]]{}
// 	response.Body.PaginationResult = res
// 	return response, err
// }

// func (h *AccountHandler) GetSalesOrderDetail(ctx context.Context, i *dto.RequestSalesOrderDetail) (
// 	*dto.EntityResponse[dto.ResponseSalesOrderDetail], error,
// ) {
// 	req, err := h.sessionHelper.GetSession(ctx)
// 	if err != nil {
// 		return nil, huma.Error400BadRequest("Not Authorized", err)
// 	}
// 	res, err := h.salesOrderService.GetSalesOrderDetail(req, i)
// 	if err != nil {
// 		return nil, huma.Error400BadRequest("Error", err)
// 	}
// 	response := &dto.EntityResponse[dto.ResponseSalesOrderDetail]{}
// 	response.Body.Result = res
// 	return response, err
// }

// func (h AccountHandler) SignIn(ctx echo.Context) error {
// 	token ,err := h.jwtService.GenerateToken(common.Claims{
// 		ID: 1,
// 		Uuid: "5f4a8a86-0fd0-4c71-8615-f44a8475d595",
// 	})
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	time.Sleep(time.Duration(2)* time.Second)
// 	return ctx.JSON(http.StatusOK, SignInResponse{
// 		AccessToken: token,
// 	})
// }
