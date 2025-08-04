	package user_rest

	import (
		"context"
		"erp/api/dto"
		"erp/internal/app/service/helpers"
		"erp/internal/domain"
		"erp/internal/domain/repository"
		user_ucase "erp/project/manage/user/internal/usecase"
		"net/http"

		"github.com/danielgtaylor/huma/v2"
	)

	type UserHandler struct {
		sessionHelper  helpers.SessionHelper
		locale         helpers.Locale
		errorHelper    helpers.ErrorHelper
		userUseCase user_ucase.UserUseCase
		permission repository.PermissionService
	}

	func NewUserHandler(
		api huma.API,
		helpers *helpers.Helpers,
		permission repository.PermissionService,
		middlewares huma.Middlewares,
		userUseCase user_ucase.UserUseCase,
	) {
		paths := NewUserPaths(domain.USER_BASE_ROUTE)
		tag := []string{"User"}
		handler := UserHandler{
			sessionHelper:  helpers.Session,
			locale:         helpers.Locale,
			errorHelper:    helpers.Error,
			userUseCase: userUseCase,
			permission: permission,
		}
		
		huma.Register(api, huma.Operation{
			OperationID:   "create user",
			Method:        http.MethodPost,
			Path:          paths.Base,
			Summary:       "Create User",
			Tags:          tag,
			DefaultStatus: http.StatusCreated,
			Middlewares:   middlewares,
		}, handler.CreateUser)

		
	}


	func (h *UserHandler) CreateUser(ctx context.Context, i *dto.CreateUserRequest) (
		*dto.ResponseMessage, error) {
		req, _ := h.sessionHelper.GetSession(ctx)
		err := h.userUseCase.CreateUser(req, i)
		if err != nil {
			return nil, h.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToCreateUser")
		}
		var res dto.ResponseMessage
		res.Body.Message = h.locale.MustLocalize(
			helpers.OptionsLocale.WithID("Message.CreateUserSuccess"),
			helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
		)
		return &res, nil
	}

