package price_list_event

import (
	"context"
	"erp/internal/domain"
	"erp/internal/domain/event"
	"erp/pkg/bus"
	"erp/pkg/logger"
	price_list_repo "erp/project/stock/price_list/repository"
)

type priceListEventHandler struct {
	emitLog logger.EmitLog
	priceListRepoEvent price_list_repo.PriceListEventRepo
}

func NewPriceListEventHandler(
	logger logger.Logger,
	priceListRepoEvent price_list_repo.PriceListEventRepo,
	bus bus.Bus,
){
	h := priceListEventHandler{
		emitLog: logger.EmitLog("price-list-event-handler"),
		priceListRepoEvent: priceListRepoEvent,
	}
	bus.RegisterHandler(domain.EventCompanyCreated,h.OnCompanyCreated())
}

func (h *priceListEventHandler) OnCompanyCreated() bus.Handler  {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func () {
				if err != nil {
					h.emitLog.Err(err,logger.OptionsLog.WithMethod("OnCompanyCreated"))
				}
			}()
			payload,ok := e.Data.(event.CreatedCompanyEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.priceListRepoEvent.CreateDefaultPriceList(ctx,payload.Tx,payload.CompanyDefaults)
			return
		},
		AbortOnError: true,
		Matcher: domain.EventCompanyCreated,
	}
}