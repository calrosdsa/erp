package rest_journal

import (
	"context"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	journal_ucase "erp/project/accounting/journal/internal/usecase"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type JournalHandler struct {
	sessionHelper  helpers.SessionHelper
	locale         helpers.Locale
	errorHelper    helpers.ErrorHelper
	journalUseCase journal_ucase.JournalUseCase
	permission     repository.PermissionService
}

func NewJournalHandler(
	api huma.API,
	helpers *helpers.Helpers,
	journalUseCase journal_ucase.JournalUseCase,
	middlewares huma.Middlewares,
	permission repository.PermissionService,
) {
	base := domain.JOURNAL_BASE_ROUTE
	tags := []string{"Journal Entry"}
	path := NewJournalPaths(base)
	h := JournalHandler{
		sessionHelper:  helpers.Session,
		locale:         helpers.Locale,
		errorHelper:    helpers.Error,
		journalUseCase: journalUseCase,
		permission:     permission,
	}
	huma.Register(api, huma.Operation{
		OperationID:   "journal-entries",
		Method:        http.MethodGet,
		Summary:       "Journal Entries",
		Tags:          tags,
		Path:          path.Base,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetJournalEntries)
	huma.Register(api, huma.Operation{
		OperationID:   "journal-entry",
		Method:        http.MethodGet,
		Summary:       "Journal Entry",
		Tags:          tags,
		Path:          path.Detail,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetJournalEntry)

	huma.Register(api, huma.Operation{
		OperationID:   "create-journal-entry",
		Method:        http.MethodPost,
		Summary:       "Create Journal Entry",
		Tags:          tags,
		Path:          path.Base,
		DefaultStatus: http.StatusCreated,
		Middlewares:   middlewares,
	}, h.CreateJournalEntry)

	huma.Register(api, huma.Operation{
		OperationID:   "edit-journal-entry",
		Method:        http.MethodPut,
		Summary:       "Edit Journal Entry",
		Tags:          tags,
		Path:          path.Base,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.EditJournalEntry)

	huma.Register(api, huma.Operation{
		OperationID:   "create-journal-entry",
		Method:        http.MethodPost,
		Summary:       "Create Journal Entry",
		Tags:          tags,
		Path:          path.Base,
		DefaultStatus: http.StatusCreated,
		Middlewares:   middlewares,
	}, h.CreateJournalEntry)

	huma.Register(api, huma.Operation{
		OperationID:   "update-journal-entry-status",
		Method:        http.MethodPut,
		Summary:       "Update Journal Entry Status",
		Path:          path.UpdateStatus,
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.UpdateStatus)
}
func (h *JournalHandler) UpdateStatus(ctx context.Context, d *dto.UpdateStatusWithEvent) (
	*dto.ResponseMessage, error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = h.journalUseCase.UpdateStatus(req, *d)
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

func (h *JournalHandler) CreateJournalEntry(ctx context.Context, d *dto.JournalEntryRequestData) (
	*dto.ResponseData[dto.JournalEntryDto], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.journalUseCase.CreateJournalEntry(req, d.Body)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.ErrorMsg")
	}
	var response dto.ResponseData[dto.JournalEntryDto]
	response.Body.Result = res
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.CreatedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}

func (h *JournalHandler) EditJournalEntry(ctx context.Context, d *dto.JournalEntryRequestData) (
	*dto.ResponseData[dto.JournalEntryDto], error) {
		fmt.Println("JOURNAL ENTRY EDIT",d.Body)
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = h.journalUseCase.EditJournalEntry(req, d.Body)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.ErrorMsg")
	}
	var response dto.ResponseData[dto.JournalEntryDto]
	response.Body.Message = h.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.SuccessfullyMessage"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}

func (h *JournalHandler) GetJournalEntry(ctx context.Context, d *dto.RequestEntity) (
	*dto.EntityResponse[dto.ResultEntity[dto.JournalEntryDetailDto]], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.journalUseCase.GetJournalEntry(req, d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.EntityResponse[dto.ResultEntity[dto.JournalEntryDetailDto]]
	response.Body.Result = res
	response.Body.Actions = h.permission.GetActions(ctx, domain.JOURNAL_ENTRY.ID)
	return &response, nil
}

func (h *JournalHandler) GetJournalEntries(ctx context.Context, d *dto.RequestPaginationData) (
	*dto.PaginationResponse[dto.PaginationResult[[]dto.JournalEntryDto]], error) {
	req, err := h.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := h.journalUseCase.GetJournalEntries(req, d)
	if err != nil {
		return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.PaginationResponse[dto.PaginationResult[[]dto.JournalEntryDto]]
	response.Body.PaginationResult = res
	response.Body.Actions = h.permission.GetActions(ctx, domain.JOURNAL_ENTRY.ID)
	return &response, nil
}
