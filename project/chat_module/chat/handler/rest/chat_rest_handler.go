package chat_rest

import (
	context "context"
	dto "erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	chat_ucase "erp/project/chat_module/chat/usecase"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type handler struct {
	sessionHelper helpers.SessionHelper
	locale        helpers.Locale
	errorHelper   helpers.ErrorHelper
	usecase       chat_ucase.ChatUseCase
}

func NewHandler(
	api huma.API,
	helpers *helpers.Helpers,
	middlewares huma.Middlewares,
	usecase chat_ucase.ChatUseCase,
) {
	base := domain.CHAT_BASE_ROUTE
	tags := []string{"Chat"}
	h := handler{
		sessionHelper: helpers.Session,
		locale:        helpers.Locale,
		errorHelper:   helpers.Error,
		usecase:       usecase,
	}

	huma.Register(api, huma.Operation{
		OperationID:   "chat",
		Method:        http.MethodGet,
		Summary:       "Chat",
		Tags:          tags,
		Path:          base,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetChats)
	huma.Register(api, huma.Operation{
		OperationID:   "chat-detail",
		Method:        http.MethodGet,
		Summary:       "Chat Detail",
		Tags:          tags,
		Path:          base + "/detail/{id}",
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetChat)

	huma.Register(api, huma.Operation{
		OperationID:   "update-member-last-read",
		Method:        http.MethodPut,
		Summary:       "Update Member Last Read",
		Tags:          tags,
		Path:          base + "/update-member-last-read/{id}",
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.UpdateMemberLastRead)

}

func (m *handler) UpdateMemberLastRead(ctx context.Context, d *dto.RequestEntity) (
	*dto.ResponseMessage, error) {
	req, err := m.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = m.usecase.UpdateMemberLastRead(req, *d)
	if err != nil {
		return nil, m.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.ResponseMessage
	response.Body.Message = m.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.SuccessfullyMessage"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil

}


func (m *handler) GetChat(ctx context.Context, d *dto.RequestEntity) (
	*dto.ResponseData[dto.ChatDetailDto], error) {
	req, err := m.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := m.usecase.GetChat(req, *d)
	if err != nil {
		return nil, m.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.ResponseData[dto.ChatDetailDto]
	response.Body.Result = res
	return &response, nil

}

func (m *handler) GetChats(ctx context.Context,d *dto.ChatsRequest) (*dto.ResponseData[[]dto.ChatDto], error) {
	req, err := m.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := m.usecase.GetChats(req,*d)
	if err != nil {
		return nil, m.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.ResponseData[[]dto.ChatDto]
	response.Body.Result = res
	return &response, nil

}
