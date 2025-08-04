package chat_message_rest

import (
	context "context"
	dto "erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	chat_message_ucase "erp/project/chat_module/message/usecase"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type handler struct {
	sessionHelper helpers.SessionHelper
	locale        helpers.Locale
	errorHelper   helpers.ErrorHelper
	usecase       chat_message_ucase.ChatMessageUcase
}

func NewHandler(
	api huma.API,
	helpers *helpers.Helpers,
	middlewares huma.Middlewares,
	usecase chat_message_ucase.ChatMessageUcase,
) {
	base := domain.CHAT_MESSAGE_BASE_ROUTE
	tags := []string{"Chat Message"}
	h := handler{
		sessionHelper: helpers.Session,
		locale:        helpers.Locale,
		errorHelper:   helpers.Error,
		usecase:       usecase,
	}

	huma.Register(api, huma.Operation{
		OperationID:   "create-chat-message",
		Method:        http.MethodPost,
		Summary:       "Chat Message",
		Tags:          tags,
		Path:          base,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.CreateMessage)
	huma.Register(api, huma.Operation{
		OperationID:   "message",
		Method:        http.MethodGet,
		Summary:       "Chat Messages",
		Tags:          tags,
		Path:          base,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetMessages)

}
func (m *handler) CreateMessage(ctx context.Context, d *dto.ChatMessageDataRequest) (
	*dto.ResponseData[dto.ChatMessageDto], error) {
	req, err := m.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := m.usecase.CreateMessage(req, d.Body)
	if err != nil {
		return nil, m.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.ResponseData[dto.ChatMessageDto]
	response.Body.Result = res
	return &response, nil
}

func (m *handler) GetMessages(ctx context.Context, d *dto.ChatMessagesRequest) (
	*dto.ResponseData[[]dto.ChatMessageDto], error) {
	req, err := m.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := m.usecase.GetMessages(req, *d)
	if err != nil {
		return nil, m.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.ResponseData[[]dto.ChatMessageDto]
	response.Body.Result = res
	return &response, nil
}

// func (m *handler) GetChats(ctx context.Context,d *dto.ChatsRequest) (*dto.ResponseData[[]dto.ChatDto], error) {
// 	req, err := m.sessionHelper.GetSession(ctx)
// 	if err != nil {
// 		return nil, huma.Error400BadRequest("Not Authorized", err)
// 	}
// 	res, err := m.usecase.GetChats(req,*d)
// 	if err != nil {
// 		return nil, m.errorHelper.HumaCustomError(string(req.LanguageCode), err)
// 	}
// 	var response dto.ResponseData[[]dto.ChatDto]
// 	response.Body.Result = res
// 	return &response, nil

// }
