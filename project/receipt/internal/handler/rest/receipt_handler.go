package receipt_rest

import (
	"context"

	// "erp/api/common"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	receipt_ucase "erp/project/receipt/internal/usecase"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type ReceiptHandler struct {
	receiptUseCase receipt_ucase.ReceiptUseCase
	sessionHelper  helpers.SessionHelper
	errorHelper    helpers.ErrorHelper
	locale         helpers.Locale
	permission     repository.PermissionService
}

func NewReceiptHandler(
	api huma.API,
	middlewares huma.Middlewares,
	helpers *helpers.Helpers,
	receiptUseCase receipt_ucase.ReceiptUseCase,
	permissionService repository.PermissionService,
) {
	tags := []string{"Receipt"}
	paths := NewReceiptPaths(domain.RECEIPT_BASE_ROUTE)
	handler := ReceiptHandler{
		receiptUseCase: receiptUseCase,
		sessionHelper:  helpers.Session,
		errorHelper:    helpers.Error,
		locale:         helpers.Locale,
		permission:     permissionService,
	}

	huma.Register(api, huma.Operation{
		OperationID:   "create-receipt",
		Method:        http.MethodPost,
		Path:          paths.Base,
		Summary:       "Create Receipt`",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.CreateReceipt)

	huma.Register(api, huma.Operation{
		OperationID:   "edit-receipt",
		Method:        http.MethodPut,
		Path:          paths.Base,
		Summary:       "Edit Receipt",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.EditReceipt)

	huma.Register(api, huma.Operation{
		OperationID:   "get-receipts",
		Method:        http.MethodGet,
		Path:          paths.Type,
		Summary:       "Get Receipts",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetReceipts)

	huma.Register(api, huma.Operation{
		OperationID:   "get-receipt-detail",
		Method:        http.MethodGet,
		Path:          paths.Detail,
		Summary:       "Get Receipt Detail",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetReceiptDetail)

	huma.Register(api, huma.Operation{
		OperationID:   "update-receipt-state",
		Method:        http.MethodPut,
		Path:          paths.UpdateState,
		Summary:       "Update Receipt State",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.UpdateReceiptState)

	huma.Register(api, huma.Operation{
		OperationID:   "export-receipt",
		Method:        http.MethodPost,
		Path:          paths.Document,
		Summary:       "Export Receipt",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.ExportReceipt)

}

func (h *ReceiptHandler) ExportReceipt(ctx context.Context, i *dto.ExportDocumentRequest) (*huma.StreamResponse, error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	bytes, err := h.receiptUseCase.ExportReceipt(req, i.Body)
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

func (h *ReceiptHandler) UpdateReceiptState(ctx context.Context, i *dto.UpdateStatusWithEvent) (
	*dto.ResponseMessage, error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = h.receiptUseCase.UpdateReceiptState(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToSubmit")
	}
	var response dto.ResponseMessage
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.SubmitSuccess"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}

func (h *ReceiptHandler) EditReceipt(ctx context.Context, d *dto.EditReceiptRequest) (
	*dto.ResponseMessage, error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	err := h.receiptUseCase.EditReceipt(req, d.Body)
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

func (h *ReceiptHandler) GetReceiptDetail(ctx context.Context, i *dto.RequestEntityWithParty) (
	*dto.EntityResponse[dto.ResultEntity[dto.ReceiptDetailDto]], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.receiptUseCase.GetReceiptDetail(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	entity, err := h.receiptUseCase.GetReceiptEntity(i.PartyType)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	actions := h.permission.GetActions(req.Ctx, entity.ID)
	var response dto.EntityResponse[dto.ResultEntity[dto.ReceiptDetailDto]]
	response.Body.Result = res
	response.Body.Actions = actions
	response.Body.AssociatedActions = h.getExtraActions(ctx, entity)
	return &response, nil
}
func (h *ReceiptHandler) GetReceipts(ctx context.Context, i *dto.RequestReceipts) (
	*dto.PaginationResponse[dto.PaginationResult[[]dto.ReceiptDto]], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.receiptUseCase.GetReceipts(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	entity, err := h.receiptUseCase.GetReceiptEntity(i.PartyType)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	actions := h.permission.GetActions(req.Ctx, entity.ID)
	var response dto.PaginationResponse[dto.PaginationResult[[]dto.ReceiptDto]]
	response.Body.Actions = actions
	response.Body.PaginationResult = res
	return &response, nil
}

func (h *ReceiptHandler) CreateReceipt(ctx context.Context, i *dto.CreateReceiptRequest) (
	*dto.ResponseData[dto.ReceiptDto], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.receiptUseCase.CreateReceipt(req, i.Body)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToCreateReceipt")
	}

	var response dto.ResponseData[dto.ReceiptDto]
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.CreateReceiptSuccess"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	response.Body.Result = res
	return &response, nil
}

func (h *ReceiptHandler) getExtraActions(ctx context.Context, entity domain.EntityTemplate) map[int][]dto.ActionDto {
	var ids []int64
	switch entity {
	case domain.PURCHASE_RECEIPT:
		ids = append(ids, domain.PURCHASE_ORDER.ID,domain.PURCHASE_INVOICE.ID)
	case domain.DELIVERY_NOTE:
		ids = append(ids, domain.SALE_ORDER.ID,domain.SALE_INVOICE.ID)
	}
	ids = append(ids, domain.SERIAL_NO.ID,domain.GENERAL_LEDGER.ID, domain.STOCK_LEDGER.ID,
	domain.ADDRESS.ID,domain.CONTACT.ID,
	domain.PAYMENT_TERMS_TEMPLATE.ID,domain.TERMS_AND_CONDITIONS.ID,domain.LEDGER.ID)
	
	r := h.permission.GetEntitiesActions(ctx,ids)
	return r
}
