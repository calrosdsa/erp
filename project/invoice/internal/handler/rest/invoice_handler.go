package invoice_rest

import (
	"context"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	invoice_ucase "erp/project/invoice/internal/usecase"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type InvoiceHandler struct {
	sessionHelper  helpers.SessionHelper
	invoiceUseCase invoice_ucase.InvoiceUseCase
	errorHelper    helpers.ErrorHelper
	permission     repository.PermissionService
	locale         helpers.Locale
}

func NewInvoiceHandler(
	api huma.API,
	helpers *helpers.Helpers,
	permission repository.PermissionService,
	middlewares huma.Middlewares,
	invoiceUseCase invoice_ucase.InvoiceUseCase,
) {
	tags := []string{"Invoice"}
	paths := NewInvoicePaths(domain.INVOICE_BASE_ROUTE)
	handler := InvoiceHandler{
		sessionHelper:  helpers.Session,
		invoiceUseCase: invoiceUseCase,
		errorHelper:    helpers.Error,
		permission:     permission,
		locale:         helpers.Locale,
	}
	huma.Register(api, huma.Operation{
		OperationID:   "create-invoice",
		Method:        http.MethodPost,
		Path:          paths.Base,
		Summary:       "Create Invoice",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.CreateInvoice)
	huma.Register(api, huma.Operation{
		OperationID:   "edit-invoice",
		Method:        http.MethodPut,
		Path:          paths.Base,
		Summary:       "Edit Invoice",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.EditInvoice)
	huma.Register(api, huma.Operation{
		OperationID:   "invoices",
		Method:        http.MethodGet,
		Path:          paths.Type,
		Summary:       "Get Invoices",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetInvoices)

	huma.Register(api, huma.Operation{
		OperationID:   "invoice",
		Method:        http.MethodGet,
		Path:          paths.Purchase,
		Summary:       "Get Invoice",
		Description:   "Get invoice",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetInvoice)

	huma.Register(api, huma.Operation{
		OperationID:   "update invoice state",
		Method:        http.MethodPut,
		Path:          paths.UpdateState,
		Summary:       "Update Invoice State",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.UpdateInvoiceState)

	huma.Register(api, huma.Operation{
		OperationID:   "export-invoice",
		Method:        http.MethodPost,
		Path:          paths.Document,
		Summary:       "Export Invoice",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.ExportInvoice)
}

func (h *InvoiceHandler) ExportInvoice(ctx context.Context, i *dto.ExportDocumentRequest) (*huma.StreamResponse, error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	bytes, err := h.invoiceUseCase.ExportInvoice(req, i.Body)
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

func (h *InvoiceHandler) UpdateInvoiceState(ctx context.Context, i *dto.UpdateStatusWithEvent) (
	*dto.ResponseMessage, error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = h.invoiceUseCase.UpdateInvoiceState(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToUpdate")
	}
	var response dto.ResponseMessage
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.UpdatedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}

func (h *InvoiceHandler) EditInvoice(ctx context.Context, d *dto.EditInvoiceRequest) (
	*dto.ResponseMessage, error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	err := h.invoiceUseCase.EditInvoice(req, d.Body)
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

func (h *InvoiceHandler) GetInvoice(ctx context.Context, i *dto.RequestEntityWithParty) (
	*dto.EntityResponse[dto.ResultEntity[dto.InvoiceDetailDto]], error,
) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.invoiceUseCase.GetInvoiceDetail(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	entity, err := h.invoiceUseCase.GetEntityInvoice(i.PartyType)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	actions := h.permission.GetActions(req.Ctx, entity.ID)
	var response dto.EntityResponse[dto.ResultEntity[dto.InvoiceDetailDto]]
	response.Body.Result = res
	response.Body.Actions = actions
	response.Body.AssociatedActions = h.getExtraActions(ctx, entity)
	return &response, nil
	// var
}

func (h *InvoiceHandler) GetInvoices(ctx context.Context, i *dto.RequestInvoices) (
	*dto.PaginationResponse[dto.PaginationResult[[]dto.InvoiceDto]], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	// h.sessionHelper.AppendPaginationParams(req, i)
	res, err := h.invoiceUseCase.GetInvoices(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	entity, err := h.invoiceUseCase.GetEntityInvoice(i.PartyType)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	actions := h.permission.GetActions(req.Ctx, entity.ID)
	response := &dto.PaginationResponse[dto.PaginationResult[[]dto.InvoiceDto]]{}
	response.Body.PaginationResult = res
	response.Body.Actions = actions
	return response, err
}

func (h *InvoiceHandler) CreateInvoice(ctx context.Context, i *dto.CreateInvoiceRequest) (
	*dto.ResponseData[dto.InvoiceDto], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error401Unauthorized("Not Authorized", err)
	}
	res, err := h.invoiceUseCase.CreateInvoice(req, i.Body)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToCreate")
	}
	var response dto.ResponseData[dto.InvoiceDto]
	response.Body.Result = res
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.CreatedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}

func (h *InvoiceHandler) getExtraActions(ctx context.Context, entity domain.EntityTemplate) map[int][]dto.ActionDto {
	var ids []int64
	switch entity {
	case domain.SALE_INVOICE:
		ids = append(ids, domain.SALE_ORDER.ID, domain.SALES_RECORD.ID)
	case domain.PURCHASE_INVOICE:
		ids = append(ids, domain.PURCHASE_ORDER.ID, domain.PURCHASE_RECORD.ID, domain.PURCHASE_RECEIPT.ID)
	}
	ids = append(ids, domain.SERIAL_NO.ID, domain.PAYMENT.ID, domain.GENERAL_LEDGER.ID,
		domain.STOCK_LEDGER.ID, domain.ADDRESS.ID, domain.CONTACT.ID,
		domain.PAYMENT_TERMS_TEMPLATE.ID,domain.TERMS_AND_CONDITIONS.ID,
		domain.LEDGER.ID,
	)
	r := h.permission.GetEntitiesActions(ctx,ids)	
	return r
}
