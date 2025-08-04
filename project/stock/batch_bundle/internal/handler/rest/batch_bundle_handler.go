package rest_batch_bundle

import (
	"context"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	batch_bundle_ucase "erp/project/stock/batch_bundle/internal/usecase"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type BatchBundleHandler struct {
	sessionHelper    helpers.SessionHelper
	locale           helpers.Locale
	errorHelper      helpers.ErrorHelper
	batchBundleUcase batch_bundle_ucase.BatchBundleUseCase
	permission       repository.PermissionService
}

func NewBatchBundleHandler(
	api huma.API,
	helpers *helpers.Helpers,
	batchBundleUcase batch_bundle_ucase.BatchBundleUseCase,
	middlewares huma.Middlewares,
	permission repository.PermissionService,
) {
	base := domain.BATCH_BUNDLE_BASE_ROUTE
	tags := []string{"Batch Bundle"}
	path := NewBatchBundlePaths(base)
	h := BatchBundleHandler{
		sessionHelper:    helpers.Session,
		locale:           helpers.Locale,
		errorHelper:      helpers.Error,
		batchBundleUcase: batchBundleUcase,
		permission:       permission,
	}
	huma.Register(api, huma.Operation{
		OperationID:   "batch-bundles",
		Method:        http.MethodGet,
		Summary:       "Bath Bundles",
		Tags:          tags,
		Path:          path.Base,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetBatchBundles)
	huma.Register(api, huma.Operation{
		OperationID:   "batch-bundle",
		Method:        http.MethodGet,
		Summary:       "Batch Bundle",
		Tags:          tags,
		Path:          path.Detail,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetBatchBundle)

}

func (h *BatchBundleHandler) GetBatchBundle(ctx context.Context, d *dto.RequestEntity) (
	*dto.EntityResponse[dto.ResultEntity[dto.BatchBundleDto]], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.batchBundleUcase.GetBatchBundle(req, d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.EntityResponse[dto.ResultEntity[dto.BatchBundleDto]]
	response.Body.Result = res
	response.Body.Actions = h.permission.GetActions(ctx, domain.BATCH_BUNDLE.ID)
	return &response, nil
}

func (h *BatchBundleHandler) GetBatchBundles(ctx context.Context, d *dto.RequestPaginationData) (
	*dto.PaginationResponse[dto.PaginationResult[[]dto.BatchBundleDto]], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.batchBundleUcase.GetBatchBundles(req, d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.PaginationResponse[dto.PaginationResult[[]dto.BatchBundleDto]]
	response.Body.PaginationResult = res
	response.Body.Actions = h.permission.GetActions(ctx, domain.BATCH_BUNDLE.ID)
	return &response, nil
}
