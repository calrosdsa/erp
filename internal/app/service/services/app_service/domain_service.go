package appservice

import (
	"context"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/internal/app/connection"
	"erp/internal/app/service/helpers"
	"erp/pkg/logger"
	"time"
)

type DomainService struct {
	conn      *connection.Connection
	timeout   time.Duration
	emitLog   helpers.EmitLog
	convertor helpers.ConvertorHelper
}

func NewDomainService(
	conn *connection.Connection,
	timeout time.Duration,
	helpers *helpers.Helpers,
) *DomainService {
	return &DomainService{
		conn:      conn,
		timeout:   timeout,
		emitLog:   helpers.Logger.EmitLog("domain-service"),
		convertor: helpers.Convertor,
	}
}

func (s *DomainService) GetCurrency(ctx context.Context, code string) (model.Currency, error) {
	currency, err := s.conn.Q.Currency.WithContext(ctx).Where(
		s.conn.Q.Currency.Code.Eq(code),
	).First()
	if err != nil {
		return model.Currency{}, err
	}
	return *currency, err
}

func (s *DomainService) GetCurrencies(ctx context.Context, i *dto.RequestPaginationData) (
	dto.PaginationResult[[]dto.CurrencyDto], error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	var (
		err error
	)
	defer func() {
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("GetCurrencies"))
		}
	}()
	var response dto.PaginationResult[[]dto.CurrencyDto]
	limit, offset := s.convertor.ToPaginationParams(i.Page, i.Size)
	response.Total, err = s.conn.Q.Currency.WithContext(ctx).CountPaginate("", 0,
		i.Query, i.Enabled)
	if err != nil {
		return response, err
	}
	currencies, err := s.conn.Q.Currency.WithContext(ctx).Paginate("", 0, i.Query,
		limit, offset, i.OrderColumn, i.Order, i.Enabled)
	if err != nil {
		return response, err
	}
	currenciesDto := make([]dto.CurrencyDto, len(currencies))
	for i, currency := range currencies {
		currenciesDto[i] = dto.CurrencyDtoFromModel(&currency)
	}
	response.Results = currenciesDto

	return response, nil
}
