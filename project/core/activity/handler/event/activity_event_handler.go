package activity_event

import (
	"context"
	"erp/internal/domain"
	"erp/internal/domain/event"
	"erp/pkg/bus"
	"erp/pkg/logger"
	activity_repo "erp/project/core/activity/repository"
)

type activityEvent struct {
	repo    activity_repo.ActivityEventRepository
	emitLog logger.EmitLog
}

func NewActiviityEvnetHandler(
	bus bus.Bus,
	repo activity_repo.ActivityEventRepository,
	logger logger.Logger,
) {
	handler := activityEvent{
		repo:    repo,
		emitLog: logger.EmitLog("activity-event"),
	}
	bus.RegisterHandler(domain.PartyStageChange, handler.OnDealStageChange())
}

func (h *activityEvent) OnDealStageChange() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func() {
				if err != nil {
					h.emitLog.Err(err, logger.OptionsLog.WithMethod("OnDealStageChange"))
				}
			}()
			payload, ok := e.Data.(event.ChangeStageEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.repo.OnStagePartyChange(ctx, payload)
			return
		},
		AbortOnError: true,
		Matcher:      domain.PartyStageChange,
	}
}
