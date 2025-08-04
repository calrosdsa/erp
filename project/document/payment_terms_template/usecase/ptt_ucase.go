package payment_terms_t_ucase

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/query"
	"erp/internal/domain"
	"erp/internal/domain/event"
	"erp/internal/domain/repository"
	"erp/pkg/bus"
	"erp/pkg/di"
	"erp/pkg/fsm"
	"erp/pkg/logger"
	payment_terms_t_repo "erp/project/document/payment_terms_template/repository"
	"fmt"
)

type PaymentTermsTemplateUcase interface {
	GetPaymentTermsTemplates(req *common.RequestContext, d dto.PaymentTermsTemplateRequest) (
		res dto.ResponseDataList[[]dto.PaymentTermsTemplateDto], err error)
	GetPaymentTermsTemplateDetail(req *common.RequestContext, d dto.RequestEntity) (
		res dto.ResultEntity[dto.PaymentTermsTemplateDto], err error)
	CreatePaymentTermsTemplate(req *common.RequestContext, d dto.PaymentTermsTemplateData) (res dto.PaymentTermsTemplateDto, err error)
	EditPaymentTermsTemplate(req *common.RequestContext, d dto.PaymentTermsTemplateData) (err error)
	UpdateStatus(req *common.RequestContext, d dto.UpdateStatusWithEvent) (err error)
	Greet(name string) (string, error)
}
type paymentTermsUcase struct {
	emitLog                  logger.EmitLog
	core                     repository.CoreService
	permission               repository.PermissionService
	fsm                      fsm.FsmState
	paymentTermsTemplateRepo payment_terms_t_repo.PaymentTermsTemplateRepo
	bus                      bus.Bus
	c                        di.Container
}

func NewUseCase(
	logger logger.Logger,
	core repository.CoreService,
	permission repository.PermissionService,
	paymentTermsTemplateRepo payment_terms_t_repo.PaymentTermsTemplateRepo,
	fsm fsm.FsmState,
	bus bus.Bus,
	c di.Container,
) PaymentTermsTemplateUcase {
	return &paymentTermsUcase{
		emitLog:                  logger.EmitLog("payment-terms-ucase"),
		core:                     core,
		permission:               permission,
		paymentTermsTemplateRepo: paymentTermsTemplateRepo,
		fsm:                      fsm,
		bus: bus,
		c: c,
	}
}

func (u *paymentTermsUcase) Greet(name string) (string, error) {
	return fmt.Sprintf("Hello, %s",name),nil
}

func (u *paymentTermsUcase) GetPaymentTermsTemplates(req *common.RequestContext, d dto.PaymentTermsTemplateRequest) (
	res dto.ResponseDataList[[]dto.PaymentTermsTemplateDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetPaymentTermsTemplates"))
		}
	}()
	if err := u.permission.CheckPermissionEntity(req.Ctx, req, domain.PAYMENT_TERMS_TEMPLATE, domain.VIEW); err != nil {
		return res, err
	}
	res.Body.Result, err = u.paymentTermsTemplateRepo.GetPaymentTermsTemplates(req, d)
	if err != nil {
		return
	}
	res.Body.FilterOptions = u.paymentTermsTemplateRepo.GetFilterOptions(string(req.LanguageCode))
	return
}

func (u *paymentTermsUcase) GetPaymentTermsTemplateDetail(req *common.RequestContext, d dto.RequestEntity) (
	res dto.ResultEntity[dto.PaymentTermsTemplateDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetPaymentTermsTemplateDetail"))
		}
	}()
	if err := u.permission.CheckPermissionEntity(req.Ctx, req, domain.PAYMENT_TERMS_TEMPLATE, domain.VIEW); err != nil {
		return res, err
	}
	res.Entity, err = u.paymentTermsTemplateRepo.GetPaymentTermsTemplateDetail(req, d)
	if err != nil {
		return
	}
	res.Activities = u.core.GerActivitiesByPartyID(req, res.Entity.ID)
	return
}

func (u *paymentTermsUcase) CreatePaymentTermsTemplate(req *common.RequestContext, d dto.PaymentTermsTemplateData) (
	res dto.PaymentTermsTemplateDto, err error) {
	ctx := u.c.Scoped(req.Ctx)
	tx := di.Get(ctx, domain.DatabaseTransactionKey).(*query.QueryTx)
	defer func(tx *query.QueryTx) {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("CreatePaymentTermsTemplate"))
		}
		err = u.closeTx(tx, err)
	}(tx)
	if err := u.permission.CheckPermissionEntity(req.Ctx, req, domain.PAYMENT_TERMS_TEMPLATE, domain.CREATE); err != nil {
		return res, err
	}
	paymentTermsTemplate, err := u.paymentTermsTemplateRepo.CreatePaymentTermsTemplate(tx, req, d)
	if err != nil {
		return
	}
	err = u.bus.Emit(req.Ctx,domain.PaymentTermsTemplateCreatedEvent,event.PaymentTermsTemplateEventData{
		Tx: tx,
		PaymentTermsTemplateID:paymentTermsTemplate.ID,
		Body: d,
	})
	res = dto.PaymentTermsTemplateFromModel(paymentTermsTemplate)
	return
}
func (u *paymentTermsUcase) EditPaymentTermsTemplate(req *common.RequestContext, d dto.PaymentTermsTemplateData) (
	err error) {
	ctx := u.c.Scoped(req.Ctx)
	tx := di.Get(ctx, domain.DatabaseTransactionKey).(*query.QueryTx)
	defer func(tx *query.QueryTx) {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("EditPaymentTermsTemplate"))
		}
		err = u.closeTx(tx, err)
	}(tx)
	if err := u.permission.CheckPermissionEntity(req.Ctx, req, domain.PAYMENT_TERMS_TEMPLATE, domain.EDIT); err != nil {
		return err
	}
	err = u.paymentTermsTemplateRepo.EditPaymentTermsTemplate(tx, req, d)
	if err != nil {
		return
	}
	err = u.bus.Emit(req.Ctx,domain.PaymentTermsTemplateEditedEvent,event.PaymentTermsTemplateEventData{
		Tx: tx,
		PaymentTermsTemplateID:d.ID,
		Body: d,
	})
	return
}

func (u *paymentTermsUcase) UpdateStatus(req *common.RequestContext, d dto.UpdateStatusWithEvent) (err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("UpdateState"))
		}
	}()
	if err := u.permission.CheckPermissionEntity(req.Ctx, req, domain.PAYMENT_TERMS_TEMPLATE, domain.EDIT); err != nil {
		return err
	}
	nextState, err := u.fsm.NextState(d.Body.CurrentState, d.Body.Events)
	if err != nil {
		return err
	}
	err = u.paymentTermsTemplateRepo.UpdateStatus(req, d, nextState)
	return
}

func (u *paymentTermsUcase) closeTx(tx *query.QueryTx, err error) error {
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
