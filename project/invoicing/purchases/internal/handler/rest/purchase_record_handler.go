package purchase_record_rest

import (
	"context"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	purchase_record_ucase "erp/project/invoicing/purchases/internal/usecase"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type PurchaseRecordHandler struct {
	sessionHelper       helpers.SessionHelper
	locale              helpers.Locale
	errorHelper         helpers.ErrorHelper
	purchaseRecordUcase  purchase_record_ucase.PurchaseRecordUcase
	permission          repository.PermissionService
}

func NewHandler(
	api huma.API,
	helpers *helpers.Helpers,
	purchaseRecordUcase  purchase_record_ucase.PurchaseRecordUcase,
	middlewares huma.Middlewares,
	permission repository.PermissionService,
) {
	base := domain.PURCHASE_RECORD_ROUTE
	tags := []string{"Purchase Record"}
	path := NewPaths(base)
	h := PurchaseRecordHandler{
		sessionHelper:       helpers.Session,
		locale:              helpers.Locale,
		errorHelper:         helpers.Error,
		purchaseRecordUcase: purchaseRecordUcase,
		permission:          permission,
	}
	huma.Register(api, huma.Operation{
		OperationID:   "purchase-records",
		Method:        http.MethodGet,
		Summary:       "Purchase Records",
		Tags:          tags,
		Path:          path.Base,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetPurchaseRecords)
	huma.Register(api, huma.Operation{
		OperationID:   "purchase-record",
		Method:        http.MethodGet,
		Summary:       "Purchase Record",
		Tags:          tags,
		Path:          path.Detail,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetPurchaseRecord)

	huma.Register(api, huma.Operation{
		OperationID:   "create-purchase-record",
		Method:        http.MethodPost,
		Summary:       "Create Purchase Record",
		Tags:          tags,
		Path:          path.Base,
		DefaultStatus: http.StatusCreated,
		Middlewares:   middlewares,
	}, h.CreatePurchaseRecord)

	huma.Register(api, huma.Operation{
		OperationID:   "edit-purchase-record",
		Method:        http.MethodPut,
		Path:          path.Base,
		Summary:       "Edit Purchase Record",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.EditPurchaseRecord)

	huma.Register(api, huma.Operation{
		OperationID:   "update-status-purchase-record",
		Method:        http.MethodPut,
		Path:          path.UpdateStatus,
		Summary:       "Update Status Purchase Record",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.UpdateStatus)

	huma.Register(api, huma.Operation{
		OperationID:   "purchase-record-export",
		Method:        http.MethodPost,
		Path:          path.Export,
		Summary:       "Purchase Record Export",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.ExportData)

	huma.Register(api,huma.Operation{
		OperationID: "export-purchase-record",
		Method: http.MethodPost,
		Path: path.Document,
		Summary: "Export Purchase Record",
		Tags: tags,
		DefaultStatus: http.StatusOK,
		Middlewares: middlewares,
	},h.ExportPurchaseRecord)
}

func (h *PurchaseRecordHandler) ExportPurchaseRecord(ctx context.Context, i *dto.ExportDocumentRequest) (*huma.StreamResponse, error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	bytes,err := h.purchaseRecordUcase.ExportPurchaseRecord(req,i.Body)
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

func (h *PurchaseRecordHandler) ExportData(ctx context.Context, i *dto.ExportDataRequest) (*huma.StreamResponse, error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	bytes,err := h.purchaseRecordUcase.ExportData(req,i)
	// fmt.Println("START STREAM RESPONSE",err)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	return &huma.StreamResponse{
		Body: func(ctx huma.Context) {
			writer := ctx.BodyWriter()
			writer.Write(bytes.Bytes())
		},
	}, nil
}

func (h *PurchaseRecordHandler) UpdateStatus(ctx context.Context, d *dto.UpdateStatusWithEvent) (
	*dto.ResponseMessage, error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	err := h.purchaseRecordUcase.UpdateStatus(req, d)
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


func (h *PurchaseRecordHandler) EditPurchaseRecord(ctx context.Context, d *dto.EditPurchaseRecordRequest) (
	*dto.ResponseMessage, error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	err := h.purchaseRecordUcase.EditPurchaseRecord(req, d)
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

func (h *PurchaseRecordHandler) CreatePurchaseRecord(ctx context.Context, d *dto.CreatePurchaseRecordRequest) (
	*dto.ResponseData[dto.PurchaseRecordDto], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.purchaseRecordUcase.CreatePurchaseRecord(req, d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToCreate")
	}
	var response dto.ResponseData[dto.PurchaseRecordDto]
	response.Body.Result = res
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.CreatedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}

func (h *PurchaseRecordHandler) GetPurchaseRecord(ctx context.Context, d *dto.RequestEntity) (
	*dto.EntityResponse[dto.ResultEntity[dto.PurchaseRecordDto]], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.purchaseRecordUcase.GetPurchaseRecord(req, d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.EntityResponse[dto.ResultEntity[dto.PurchaseRecordDto]]
	response.Body.Result = res
	response.Body.Actions = h.permission.GetActions(req.Ctx, domain.PURCHASE_RECORD.ID)
	return &response, nil
}

func (h *PurchaseRecordHandler) GetPurchaseRecords(ctx context.Context, d *dto.PurchaseRecordsRequest) (
	*dto.PaginationResponse[dto.PaginationResult[[]dto.PurchaseRecordDto]], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.purchaseRecordUcase.GetPurchaseRecords(req, d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.PaginationResponse[dto.PaginationResult[[]dto.PurchaseRecordDto]]
	response.Body.PaginationResult = res
	response.Body.Actions = h.permission.GetActions(req.Ctx, domain.PURCHASE_RECORD.ID)
	return &response, nil
}
