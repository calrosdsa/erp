package deal_rest

import (
	context "context"
	dto "erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	deal_ucase "erp/project/crm/deal/usecase"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type handler struct {
	sessionHelper helpers.SessionHelper
	locale        helpers.Locale
	errorHelper   helpers.ErrorHelper
	permission    repository.PermissionService
	usecase       deal_ucase.DealUseCase
}

func NewHandler(
	api huma.API,
	helpers *helpers.Helpers,
	middlewares huma.Middlewares,
	permission repository.PermissionService,
	usecase deal_ucase.DealUseCase,
) {
	base := domain.DEAL_BASE_ROUTE
	tags := []string{"Deal"}
	h := handler{
		sessionHelper: helpers.Session,
		locale:        helpers.Locale,
		errorHelper:   helpers.Error,
		usecase:       usecase,
		permission:    permission,
	}

	huma.Register(api, huma.Operation{
		OperationID:   "deal",
		Method:        http.MethodGet,
		Summary:       "Deal",
		Description:   "Retrieve a paginated list of deals with optional filtering and sorting capabilities. Returns deal summary information including status, value, and key dates.",
		Tags:          tags,
		Path:          base,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetDeals)
	huma.Register(api, huma.Operation{
		OperationID:   "deal-detail",
		Method:        http.MethodGet,
		Summary:       "Deal Detail",
		Description:   "Retrieve comprehensive details for a specific deal by ID. Returns full deal information including contact details, pricing, timeline, and associated actions available to the current user.",
		Tags:          tags,
		Path:          base + "/detail/{id}",
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.GetDeal)

	huma.Register(api, huma.Operation{
		OperationID:   "create-deal",
		Method:        http.MethodPost,
		Summary:       "Create Deal",
		Description:   "Create a new deal record with the provided deal data. Validates input data, assigns initial status, and returns the created deal with success message.",
		Tags:          tags,
		Path:          base,
		DefaultStatus: http.StatusCreated,
		Middlewares:   middlewares,
	}, h.CreateDeal)

	huma.Register(api, huma.Operation{
		OperationID:   "edit-deal",
		Method:        http.MethodPut,
		Path:          base,
		Summary:       "Edit Deal",
		Description:   "Update an existing deal record with new data. Validates the provided changes and updates the deal information while maintaining audit trail.",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.EditDeal)

	huma.Register(api, huma.Operation{
		OperationID:   "deal-transition",
		Method:        http.MethodPut,
		Path:          base + "/transition",
		Summary:       "Deal Transition",
		Description:   "Transition a deal through its lifecycle states (e.g., from prospect to qualified, won, or lost). Validates state transition rules and updates deal status accordingly.",
		Tags:          tags,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
	}, h.DealTransition)

}

func (m *handler) DealTransition(ctx context.Context, d *dto.EntityTransitionRequest) (*dto.ResponseMessage, error) {
	req, err := m.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = m.usecase.DealTransition(req, d.Body)
	if err != nil {
		return nil, m.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToEdit")
	}
	var response dto.ResponseMessage
	response.Body.Message = m.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.EditedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil
}

func (m *handler) CreateDeal(ctx context.Context, d *dto.DealDataRequest) (*dto.ResponseData[dto.DealDto], error) {
	req, err := m.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := m.usecase.CreateDeal(req, d.Body)
	if err != nil {
		return nil, m.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToCreate")
	}
	var response dto.ResponseData[dto.DealDto]
	response.Body.Result = res
	response.Body.Message = m.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.CreatedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil

}

func (m *handler) EditDeal(ctx context.Context, d *dto.DealDataRequest) (*dto.ResponseMessage, error) {
	req, err := m.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	err = m.usecase.EditDeal(req, d.Body)
	if err != nil {
		return nil, m.errorHelper.HumaCustomError(string(req.LanguageCode), err, "Error.FailToEdit")
	}
	var response dto.ResponseMessage
	response.Body.Message = m.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.EditedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)
	return &response, nil

}

func (m *handler) GetDeal(ctx context.Context, d *dto.RequestEntity) (
	*dto.EntityResponse[dto.ResultEntity[dto.DealDetailDto]], error) {
	req, err := m.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := m.usecase.GetDeal(req, *d)
	if err != nil {
		return nil, m.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	var response dto.EntityResponse[dto.ResultEntity[dto.DealDetailDto]]
	response.Body.Result = res
	response.Body.Actions = m.permission.GetActions(ctx, domain.DEAL.ID)
	response.Body.AssociatedActions = m.getActions(ctx)
	return &response, nil

}

func (m *handler) GetDeals(ctx context.Context, d *dto.DealsRequest) (*dto.ResponseDataList[[]dto.DealDto], error) {
	req, err := m.sessionHelper.GetSession(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Not Authorized", err)
	}
	res, err := m.usecase.GetDeals(req, *d)
	if err != nil {
		return nil, m.errorHelper.HumaCustomError(string(req.LanguageCode), err)
	}
	// res.Body.AssociatedActions = m.getActions(ctx)
	res.Body.Actions = m.permission.GetActions(ctx, domain.DEAL.ID)
	return &res, nil

}

func (h *handler) getActions(ctx context.Context) map[int][]dto.ActionDto {
	var ids []int64
	ids = append(ids, domain.CONTACT.ID)
	ids = append(ids, domain.PRICING.ID)
	ids = append(ids, domain.QUOTATION.ID)
	ids = append(ids, domain.SUPPLIER_QUOTATION.ID)

	r := h.permission.GetEntitiesActions(ctx, ids)

	return r
}


