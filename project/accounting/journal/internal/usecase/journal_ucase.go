package journal_ucase

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/query"
	"erp/gen/proto"
	"erp/internal/domain"
	"erp/internal/domain/event"
	"erp/internal/domain/repository"
	"erp/pkg/bus"
	"erp/pkg/di"
	"erp/pkg/fsm"
	"erp/pkg/logger"
	journal_repo "erp/project/accounting/journal/internal/repository"
	"fmt"
)

type JournalUseCase interface {
	CreateJournalEntry(req *common.RequestContext, d dto.JournalEntryData) (dto dto.JournalEntryDto, err error)
	EditJournalEntry(req *common.RequestContext, d dto.JournalEntryData) (err error)
	GetJournalEntry(req *common.RequestContext, d *dto.RequestEntity) (
		res dto.ResultEntity[dto.JournalEntryDetailDto], err error)
	GetJournalEntries(req *common.RequestContext, d *dto.RequestPaginationData) (
		res dto.PaginationResult[[]dto.JournalEntryDto], err error)
	UpdateStatus(req *common.RequestContext, d dto.UpdateStatusWithEvent) (err error)
}

type journalUseCase struct {
	emitLog     logger.EmitLog
	journalRepo journal_repo.JournalRepository
	permission  repository.PermissionService
	core        repository.CoreService
	fsm         fsm.FsmState
	c           di.Container
	bus         bus.Bus
}

func NewJournalUseCase(
	logger logger.Logger,
	permission repository.PermissionService,
	core repository.CoreService,
	journalRepo journal_repo.JournalRepository,
	c di.Container,
	bus bus.Bus,
	fsm fsm.FsmState,
) JournalUseCase {
	return &journalUseCase{
		emitLog:     logger.EmitLog("journal-usecase"),
		permission:  permission,
		core:        core,
		journalRepo: journalRepo,
		c:           c,
		bus:         bus,
		fsm:         fsm,
	}
}

func (u *journalUseCase) UpdateStatus(req *common.RequestContext, d dto.UpdateStatusWithEvent) (err error) {
	ctx := u.c.Scoped(req.Ctx)
	tx := di.Get(ctx, domain.DatabaseTransactionKey).(*query.QueryTx)
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("UpstateOrderState"))
		}
		err = u.closeTx(tx, err)
		fmt.Println("TRANSACTION ERROR", err)
	}()
	err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.JOURNAL_ENTRY, domain.EDIT)
	nextState, err := u.fsm.NextState(d.Body.CurrentState, d.Body.Events, d.Body.PartyType)
	if err != nil {
		return err
	}
	journalEntry, lines, err := u.journalRepo.UpdateStatus(req, tx,d, nextState)
	if err != nil {
		return
	}
	payload := event.StatusJournalEntryEventData{
		Tx:           tx,
		JournalEntry: journalEntry,
		Lines:        lines,
	}
	switch nextState {
	case proto.State_SUBMITTED.String():
		err = u.bus.Emit(req.Ctx, domain.JournalEntrySubmittedEvent, payload)
	case proto.State_CANCELLED.String():
		err = u.bus.Emit(req.Ctx, domain.JournalEntryCancelledEvent, payload)
	}
	return
}

func (u *journalUseCase) CreateJournalEntry(req *common.RequestContext, d dto.JournalEntryData) (res dto.JournalEntryDto, err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("CreateJournalEntry"))
		}
	}()
	if err := u.permission.CheckPermissionEntity(req.Ctx, req, domain.JOURNAL_ENTRY, domain.CREATE); err != nil {
		return res, err
	}
	res, err = u.journalRepo.CreateJournalEntry(req, d)
	if err != nil {
		return
	}
	return
}

func (u *journalUseCase) EditJournalEntry(req *common.RequestContext, d dto.JournalEntryData) (err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("EditJournalEntry"))
		}
	}()
	if err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.JOURNAL_ENTRY, domain.CREATE); err != nil {
		return
	}
	err = u.journalRepo.EditJournalEntry(req, d)
	if err != nil {
		return
	}
	return
}

func (u *journalUseCase) GetJournalEntry(req *common.RequestContext, d *dto.RequestEntity) (
	res dto.ResultEntity[dto.JournalEntryDetailDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetJournalEntry"))
		}
	}()
	if allow := u.permission.CheckPermission(req.Ctx, req, domain.JOURNAL_ENTRY, domain.VIEW); !allow {
		return res, domain.ACTION_NOT_ALLOWED
	}
	res, err = u.journalRepo.GetJournalEntry(req, d)
	if err != nil {
		return
	}
	res.Activities = u.core.GerActivitiesByPartyID(req, res.Entity.JournalEntry.ID)
	return
}
func (u *journalUseCase) GetJournalEntries(req *common.RequestContext, d *dto.RequestPaginationData) (
	res dto.PaginationResult[[]dto.JournalEntryDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetJournalEntries"))
		}
	}()
	if allow := u.permission.CheckPermission(req.Ctx, req, domain.JOURNAL_ENTRY, domain.VIEW); !allow {
		return res, domain.ACTION_NOT_ALLOWED
	}
	res, err = u.journalRepo.GetJournalEntries(req, d)
	return
}

func (s *journalUseCase) closeTx(tx *query.QueryTx, err error) error {
	if p := recover(); p != nil {
		_ = tx.Rollback()
		panic(p)
	} else if err != nil {
		_ = tx.Rollback()
		return err
	} else {
		return tx.Commit()
	}
}
