package cash_outflow_ucase

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
	cash_outflow_pdf "erp/project/accounting/cash_outflow/pkg/pdf"
	cash_outflow_repo "erp/project/accounting/cash_outflow/repository"
	"fmt"
)

type CashOutflowUseCase interface {
	Create(req *common.RequestContext, d dto.CashOutflowData) (res dto.CashOutflowDto, err error)
	UpdateStatus(req *common.RequestContext, d dto.UpdateStatusWithEvent) (err error)
	GetCashOutflow(req *common.RequestContext, d dto.RequestEntity) (
		res dto.ResultEntity[dto.CashOutflowDto], err error)
	GetCashOutflows(req *common.RequestContext, d dto.CashOutflowsRequest) (
		res dto.ResponseDataList[[]dto.CashOutflowDto], err error)
	EditCashOutflow(req *common.RequestContext, d dto.CashOutflowData) (err error)
	ExportCashOutflow(req *common.RequestContext, i dto.ExportDocumentData) (res []byte, err error)
}

type cashOutflowUcase struct {
	repo       cash_outflow_repo.CashOutflowRepository
	emitLog    logger.EmitLog
	permission repository.PermissionService
	fsm        fsm.FsmState
	bus        bus.Bus
	c          di.Container
	core       repository.CoreService
	document   repository.DocumentService
	pdfGenerator cash_outflow_pdf.CashOutflowPdf
}

func NewCashOutflowUseCase(
	logger logger.Logger,
	permission repository.PermissionService,
	repo cash_outflow_repo.CashOutflowRepository,
	fsm fsm.FsmState,
	bus bus.Bus,
	c di.Container,
	core repository.CoreService,
	document repository.DocumentService,
	pdfGenerator cash_outflow_pdf.CashOutflowPdf,
) CashOutflowUseCase {
	cashOutflowUcase := cashOutflowUcase{
		repo:       repo,
		emitLog:    logger.EmitLog("cash-outflow-usecase"),
		permission: permission,
		fsm:        fsm,
		c:          c,
		bus:        bus,
		core:       core,
		document:   document,
		pdfGenerator: pdfGenerator,
	}
	return &cashOutflowUcase
}



func(u *cashOutflowUcase) ExportCashOutflow(req *common.RequestContext, i dto.ExportDocumentData) (res []byte, err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("ExportCashOutflow"))
		}
	}()
	cashOutflow,err := u.GetCashOutflow(req,dto.RequestEntity{
		ID: i.ID,
	})
	if err != nil {
		return
	}
	res,err  =u.pdfGenerator.GenerateCashOutflowDocument(req,cashOutflow.Entity)
	return
}

func (u *cashOutflowUcase) UpdateStatus(req *common.RequestContext, d dto.UpdateStatusWithEvent) (err error) {
	ctx := u.c.Scoped(req.Ctx)
	tx := di.Get(ctx, domain.DbKey).(*query.Query).Begin()
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("UpdateStatus"))
		}
		err = u.closeTx(tx, err)
	}()
	if err := u.permission.CheckPermissionEntity(req.Ctx, req, domain.CASH_OUTFLOW, domain.EDIT); err != nil {
		return domain.ACTION_NOT_ALLOWED
	}
	nextState, err := u.fsm.NextState(d.Body.CurrentState, d.Body.Events)
	if err != nil {
		return err
	}
	cashOutflow, err := u.repo.UpdateStatus(tx, req, d, nextState)
	if err != nil {
		return
	}
	taxLinesData, err := u.document.GetTaxLines(tx.Query, req.Ctx, cashOutflow.ID)
	if err != nil {
		return
	}

	evtData := event.StatusCashOutflowEventData{
		Tx:              tx,
		CashOutflow:     *cashOutflow,
		TaxLinesData:    taxLinesData,
		CompanyDefaults: req.CompanyDefaults,
	}
	switch nextState {
	case proto.State_SUBMITTED.String():
		err = u.bus.Emit(req.Ctx, domain.CashOutflowSubmittedEvent, evtData)
	case proto.State_CANCELLED.String():
		err = u.bus.Emit(req.Ctx, domain.CashOutflowCancelledEvent, evtData)
	}
	return
}

func (u *cashOutflowUcase) EditCashOutflow(req *common.RequestContext, d dto.CashOutflowData) (err error) {
	ctx := u.c.Scoped(req.Ctx)
	tx := di.Get(ctx, domain.DbKey).(*query.Query).Begin()
	defer func(tx *query.QueryTx) {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("EditCashOutflow"))
		}
		err = u.closeTx(tx, err)
	}(tx)
	if err := u.permission.CheckPermissionEntity(req.Ctx, req, domain.CASH_OUTFLOW, domain.EDIT); err != nil {
		return domain.ACTION_NOT_ALLOWED
	}
	if err = u.repo.Edit(tx, req, d); err != nil {
		return
	}
	err = u.bus.Emit(req.Ctx, domain.CashOutflowEditedEvent, event.CashOutflowEventData{
		Tx:   tx,
		Data: d,
	})
	return
}

func (u *cashOutflowUcase) GetCashOutflow(req *common.RequestContext, d dto.RequestEntity) (
	res dto.ResultEntity[dto.CashOutflowDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetCashOutflow"))
		}
	}()
	if err := u.permission.CheckPermissionEntity(req.Ctx, req, domain.CASH_OUTFLOW, domain.VIEW); err != nil {
		return res, domain.ACTION_NOT_ALLOWED
	}
	res.Entity, err = u.repo.GetCashOutflow(req, d)
	if err != nil {
		return
	}
	res.Activities = u.core.GerActivitiesByPartyID(req, res.Entity.ID)
	return
}

func (u *cashOutflowUcase) GetCashOutflows(req *common.RequestContext, d dto.CashOutflowsRequest) (
	res dto.ResponseDataList[[]dto.CashOutflowDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetCashOutflows"))
		}
	}()
	if err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.CASH_OUTFLOW, domain.VIEW); err != nil {
		return res, domain.ACTION_NOT_ALLOWED
	}

	res.Body.Result, err = u.repo.GetCashOutflows(req, d)
	if err != nil {
		return
	}
	res.Body.FilterOptions = u.repo.GetFilterOptions()
	return
}

func (u *cashOutflowUcase) Create(req *common.RequestContext, d dto.CashOutflowData) (res dto.CashOutflowDto, err error) {
	ctx := u.c.Scoped(req.Ctx)
	tx := di.Get(ctx, domain.DbKey).(*query.Query).Begin()
	defer func(tx *query.QueryTx) {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("Create"))
		}
		err = u.closeTx(tx, err)
		fmt.Println("TRANSACTION ERROR", err)
	}(tx)
	if err := u.permission.CheckPermissionEntity(req.Ctx, req, domain.CASH_OUTFLOW, domain.CREATE); err != nil {
		return res, domain.ACTION_NOT_ALLOWED
	}
	cashOutflow, err := u.repo.Create(tx, req, d)
	if err != nil {
		return
	}
	res = dto.CashOutFlowModel(cashOutflow)
	err = u.bus.Emit(req.Ctx, domain.CashOutflowCreatedEvent, event.CashOutflowEventData{
		Tx:          tx,
		Data:        d,
		CashOutflow: cashOutflow,
	})
	return
}

func (u *cashOutflowUcase) closeTx(tx *query.QueryTx, err error) error {
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
