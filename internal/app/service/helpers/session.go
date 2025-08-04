package helpers

import (
	"context"
	"erp/api/common"
	"erp/api/dto"
	"erp/internal/domain"
	"fmt"
)

type SessionHelper interface {
	GetSession(ctx context.Context) (*common.RequestContext, error)
	GetAdminSession(ctx context.Context) (*common.AdminRequestContext, error)
	AppendPaginationParams(req *common.RequestContext, i *dto.RequestPaginationData)
	AppendParam(req *common.RequestContext, name AllowedParams, param string)
	ParseAcceptLanguage(lang string) string
}

type sessionHelper struct {
	emitLog       EmitLog
}



func NewSessionHelper(
	logerr *LoggerHelper,
) SessionHelper {
	return &sessionHelper{
		emitLog:       logerr.EmitLog("session-helper"),
	}
}

// func (h *sessionHelper)CreateRequestContext(
// 	conn *connection.Connection,
// 	company *entity.Company,
// 	user *entity.User,
// 	)(*common.RequestContext,error){
// 	appConfig := h.configService.GetDefaultLanguage()
// 	var reqCtx common.RequestContext
// 	reqCtx.LanguageCode = common.LanguageCode(appConfig)
// 	ctx := context.Background()
// 	if company != nil {
// 		err := conn.Db.WithContext(ctx).Where(&company).First(&company).Error
// 		if err != nil {
// 			h.emitLog.Err(err,utils.OptionsLog.WithMethod("CreateRequestContextCompany"))
// 		}
// 		reqCtx.ActiveCompany = *company
// 	}
// 	if user != nil {
// 		err := conn.Db.WithContext(ctx).Where(&user).First(&user).Error
// 		if err != nil {
// 			h.emitLog.Err(err,utils.OptionsLog.WithMethod("CreateRequestContextCompany"))
// 		}
// 		reqCtx.Session.CurrentUser = *user
// 	}
// 	return &reqCtx,nil
// }

func (h *sessionHelper) ParseAcceptLanguage(lang string) string {
	if lang == "" {
		lang = domain.DEFAULT_LANGUAGE
	}
	if len(lang) > 2 {
		lang = lang[:2]
	}
	return lang
}

func (h *sessionHelper) GetSession(ctx context.Context) (*common.RequestContext, error) {
	v := ctx.Value(domain.SESSION_KEY)
	if v == nil {
		return nil, fmt.Errorf("no session in context")
	}

	if userProfile, ok := v.(*common.RequestContext); ok {
		return userProfile, nil
	}

	return nil, fmt.Errorf("request context is nil")
}

func (h *sessionHelper)GetAdminSession(ctx context.Context) (*common.AdminRequestContext, error) {
	v := ctx.Value(domain.SESSION_KEY)
	if v == nil {
		return nil, fmt.Errorf("no session in context")
	}

	if session, ok := v.(*common.AdminRequestContext); ok {
		return session, nil
	}

	return nil, fmt.Errorf("admin request context is nil")
}

func (h *sessionHelper) AppendPaginationParams(req *common.RequestContext, i *dto.RequestPaginationData) {
	params := make(map[string]string)
	params["page"] = i.Page
	params["size"] = i.Size
	params["order"] = i.Order
	params["column"] = i.OrderColumn
	req.Params = params
}

func (h *sessionHelper) AppendParam(req *common.RequestContext, name AllowedParams, param string) {
	params := req.Params
	params[string(name)] = param
	req.Params = params
}

type AllowedParams string

const (
	COLUMN_PARAM = "column"
	ORDER_PARAM  = "order"
)

// func (h *SessionHelper) AppendParams(req *common.RequestContext,i *dto.RequestPaginationData){
// 	params := make(map[string]string)
// 	params["page"] = i.Page
// 	params["size"] = i.Size
// 	req.Params = params
// }
