package module_event

import (
	"context"
	"erp/internal/domain"
	"erp/internal/domain/event"
	"erp/pkg/bus"
	"erp/pkg/logger"
	module_repo "erp/project/core/module/repository"
)

type moduleEventHandler struct {
	emitLog logger.EmitLog
	moduleEventRepo module_repo.ModuleEventRepo
}

func NewModuleEventHandler(
	logger logger.Logger,
	moduleEventRepo module_repo.ModuleEventRepo,
	bus bus.Bus,
){
	h := moduleEventHandler{
		emitLog: logger.EmitLog("module-event"),
		moduleEventRepo: moduleEventRepo,
	}
	bus.RegisterHandler(domain.EventCompanyCreated,h.OnCompanyCreated())
}

func (h *moduleEventHandler) OnCompanyCreated() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) (err error) {
			defer func (){
				if err != nil {
					h.emitLog.Err(err,logger.OptionsLog.WithMethod("OnCompanyCreated"))
				}
			}()
			payload,ok := e.Data.(event.CreatedCompanyEventData)
			if !ok {
				return domain.FAIL_TYPE_ASSERTION
			}
			err = h.moduleEventRepo.CreateCompany(payload.Tx,ctx,payload.Company.ID,payload.Body.CompanyModules)
			return
		},
		AbortOnError: true,
		Matcher: domain.EventCompanyCreated,
	}
}