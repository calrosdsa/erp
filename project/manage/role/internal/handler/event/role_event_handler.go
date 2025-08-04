package role_event

import (
	"context"
	"erp/internal/domain"
	"erp/internal/domain/event"
	"erp/pkg/bus"
	"erp/pkg/logger"
	role_repo "erp/project/manage/role/internal/repository"
	"fmt"
)

type roleEventHandler struct {
	emitLog       logger.EmitLog
	roleEventRepo role_repo.RoleEventRepository
}

func NewRoleEventHandler(
	logger logger.Logger,
	bus bus.Bus,
	roleEventRepo role_repo.RoleEventRepository,
) {
	h := roleEventHandler{
		roleEventRepo: roleEventRepo,
		emitLog:       logger.EmitLog("role-event-handler"),
	}
	bus.RegisterHandler(domain.EventCompanyCreated, h.OnCompanyCreated())
}

func (h *roleEventHandler) OnCompanyCreated() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			fmt.Println("CREATING COMPANY ROLE")
			defer func(){
				if err !=nil {
					h.emitLog.Err(err,logger.OptionsLog.WithMethod("OnCompanyCreated"))
				}
			}()
			payload, ok := e.Data.(event.CreatedCompanyEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.roleEventRepo.OnCompanyCreated(ctx, payload)
			if err != nil {
				return err
			}
			return nil
		},
		Matcher: domain.EventCompanyCreated,
		AbortOnError: true,
	}
}
