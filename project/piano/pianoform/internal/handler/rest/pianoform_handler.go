package pianoform_rest

import (
	"context"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/internal/app/service/helpers"
	"erp/internal/domain/repository"
	pianoform_ucase "erp/project/piano/pianoform/internal/usecase"
	"fmt"
	"net/http"
	"github.com/danielgtaylor/huma/v2"
)

type PianoFormHandler struct {
	sessionHelper  helpers.SessionHelper
	locale         helpers.Locale
	errorHelper    helpers.ErrorHelper
	permission     repository.PermissionService
	pianoFormUcase pianoform_ucase.PianoFormUseCase
}

func NewPianoFormHandler(
	api huma.API,
	helpers *helpers.Helpers,
	permission repository.PermissionService,
	middlewares huma.Middlewares,
	pianoFormUcase pianoform_ucase.PianoFormUseCase,
) {
	paths := newPianoPaths("/pianoForms")
	tag := []string{"Piano Form"}
	handler := PianoFormHandler{
		sessionHelper:  helpers.Session,
		locale:         helpers.Locale,
		errorHelper:    helpers.Error,
		permission:     permission,
		pianoFormUcase: pianoFormUcase,
	}

	huma.Register(api, huma.Operation{
		OperationID:   "create piano form",
		Method:        http.MethodPost,
		Path:          paths.Base,
		Summary:       "Create Piano Form",
		Tags:          tag,
		DefaultStatus: http.StatusCreated,
	}, handler.CreatePianoForm)

	huma.Register(api, huma.Operation{
		OperationID:   "get piano forms",
		Method:        http.MethodGet,
		Path:          paths.Base,
		Summary:       "Get Piano Forms",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetPianoForms)
	huma.Register(api, huma.Operation{
		OperationID:   "get piano form",
		Method:        http.MethodGet,
		Path:          paths.Detail,
		Summary:       "Get Piano Form",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.GetPianoForm)

	huma.Register(api, huma.Operation{
		OperationID:   "export data",
		Method:        http.MethodPost,
		Path:          paths.Export,
		Summary:       "export data",
		Tags:          tag,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, handler.ExportData)
}

func (h *PianoFormHandler) ExportData(ctx context.Context, i *dto.PianoExportRequest) (*huma.StreamResponse, error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	bytes,err := h.pianoFormUcase.ExportData(req,i)
	fmt.Println("START STREAM RESPONSE",err)
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

func (h *PianoFormHandler) GetPianoForm(ctx context.Context, i *dto.RequestEntity) (
	*dto.EntityResponse[dto.ResultEntity[*model.PianoForm]], error,
) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	// h.sessionHelper.AppendPaginationParams(req, i)
	res, err := h.pianoFormUcase.GetPianoForm(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}

	var response dto.EntityResponse[dto.ResultEntity[*model.PianoForm]]
	response.Body.Result = res
	return &response, err
}

func (h *PianoFormHandler) GetPianoForms(ctx context.Context, i *dto.RequestPaginationData) (
	*dto.PaginationResponse[dto.PaginationResult[[]*model.PianoForm]], error,
) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	// h.sessionHelper.AppendPaginationParams(req, i)
	res, err := h.pianoFormUcase.GetPianoForms(req, i)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	// actions := h.permission.GetActions(req.Ctx, regate_domain.COURT.ID)

	var response dto.PaginationResponse[dto.PaginationResult[[]*model.PianoForm]]
	response.Body.PaginationResult = res
	// response.Body.Actions = actions
	return &response, err
}

func (h *PianoFormHandler) CreatePianoForm(ctx context.Context, i *dto.CreatePianoFormRequest) (
	*dto.ResponseMessage, error) {
	acceptLanguage := h.sessionHelper.ParseAcceptLanguage(i.AcceptLanguage)
	err := h.pianoFormUcase.CreatePianoForm(ctx, i.Body)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(acceptLanguage, err, "Error.FailToCreate")
	}
	var response dto.ResponseMessage
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.CreatedSuccessfully"),
		helpers.OptionsLocale.WithLang(acceptLanguage),
	)
	return &response, nil
}
