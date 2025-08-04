package sales_record_rest

import (
	"context"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	sales_record_ucase "erp/project/invoicing/sales/internal/usecase"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type SalesRecordHandler struct {
	sessionHelper       helpers.SessionHelper
	locale              helpers.Locale
	errorHelper         helpers.ErrorHelper
	salesRecordUcase sales_record_ucase.SalesRecordUcase
	permission          repository.PermissionService
}

func NewHandler(
	api huma.API,
	helpers *helpers.Helpers,
	salesRecordUcase sales_record_ucase.SalesRecordUcase,
	middlewares huma.Middlewares,
	permission repository.PermissionService,
) {
	base := domain.SALES_RECORD_ROUTE
	tags := []string{"Sales Record"}
	paths := NewPaths(base)
	h := SalesRecordHandler{
		sessionHelper:       helpers.Session,
		locale:              helpers.Locale,
		errorHelper:         helpers.Error,
		salesRecordUcase: salesRecordUcase,
		permission:          permission,
	}
	huma.Register(api, huma.Operation{
		OperationID:   "sales-records",
		Method:        http.MethodGet,
		Summary:       "Sales Records",
		Tags:          tags,
		Path:          paths.Base,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetSalesRecords)
	huma.Register(api, huma.Operation{
		OperationID:   "sales-record",
		Method:        http.MethodGet,
		Summary:       "Sales Record",
		Tags:          tags,
		Path:          paths.Detail,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetSalesRecord)

	huma.Register(api, huma.Operation{
		OperationID:   "create-sales-record",
		Method:        http.MethodPost,
		Summary:       "Create Sales Record",
		Tags:          tags,
		Path:          paths.Base,
		DefaultStatus: http.StatusCreated,
		Middlewares:   middlewares,
	}, h.CreateSalesRecord)

	huma.Register(api, huma.Operation{
		OperationID:   "edit-sales-record",
		Method:        http.MethodPut,
		Path:          paths.Base,
		Summary:       "Edit Sales Record",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.EditSalesRecord)

	huma.Register(api, huma.Operation{
		OperationID:   "update-status-sales-record",
		Method:        http.MethodPut,
		Path:          paths.UpdateStatus,
		Summary:       "Update Status Sales Record",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.UpdateStatus)

	huma.Register(api, huma.Operation{
		OperationID:   "sales-record-export",
		Method:        http.MethodPost,
		Path:          paths.Export,
		Summary:       "Sales Record Export",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.ExportData)
}

func (h *SalesRecordHandler) ExportData(ctx context.Context, i *dto.ExportDataRequest) (*huma.StreamResponse, error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	bytes,err := h.salesRecordUcase.ExportData(req,i)
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

func (h *SalesRecordHandler) UpdateStatus(ctx context.Context, d *dto.UpdateStatusWithEvent) (
	*dto.ResponseMessage, error) {
	req, err:= h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = h.salesRecordUcase.UpdateStatus(req, d)
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


func (h *SalesRecordHandler) EditSalesRecord(ctx context.Context, d *dto.EditSalesRecordRequest) (
	*dto.ResponseMessage, error) {
	req, _ := h.sessionHelper.GetSession(ctx)
	err := h.salesRecordUcase.EditSalesRecord(req, d)
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

func (h *SalesRecordHandler) CreateSalesRecord(ctx context.Context, d *dto.CreateSalesRecordRequest) (
	*dto.ResponseData[dto.SalesRecordDto], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.salesRecordUcase.CreateSalesRecord(req, d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToCreate")
	}
	var response dto.ResponseData[dto.SalesRecordDto]
	response.Body.Result = res
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.CreatedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}

func (h *SalesRecordHandler) GetSalesRecord(ctx context.Context, d *dto.RequestEntity) (
	*dto.EntityResponse[dto.ResultEntity[dto.SalesRecordDto]], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.salesRecordUcase.GetSalesRecord(req, d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.EntityResponse[dto.ResultEntity[dto.SalesRecordDto]]
	response.Body.Result = res
	response.Body.Actions = h.permission.GetActions(req.Ctx, domain.SALES_RECORD.ID)
	return &response, nil
}

func (h *SalesRecordHandler) GetSalesRecords(ctx context.Context, d *dto.SalesRecordsRequest) (
	*dto.PaginationResponse[dto.PaginationResult[[]dto.SalesRecordDto]], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.salesRecordUcase.GetSalesRecords(req, d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.PaginationResponse[dto.PaginationResult[[]dto.SalesRecordDto]]
	response.Body.PaginationResult = res
	response.Body.Actions = h.permission.GetActions(req.Ctx, domain.SALES_RECORD.ID)
	return &response, nil
}
