package payment_ucase

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
	payment_pdf "erp/project/accounting/payment/internal/pkg/pdf"
	payment_repo "erp/project/accounting/payment/internal/repository"
	"fmt"
)

type PaymentUseCase interface {
	CreatePayment(req *common.RequestContext, d *dto.CreatePaymentRequest) (
		dto.PaymentDto, error)
	GetAllowedParties(req *common.RequestContext) []dto.PartyTypeDto
	GetPayments(req *common.RequestContext, d *dto.RequestPayments) (
		dto.PaginationResult[[]dto.PaymentDto], error,
	)
	GetPaymentDetail(req *common.RequestContext, d *dto.RequestEntity) (
		dto.ResultEntity[dto.PaymentDetailDto], error)
	UpdatePaymentState(req *common.RequestContext, i *dto.UpdateStatusWithEvent) (err error)
	GetPaymentAccounts(req *common.RequestContext) (res dto.PaymentAccountsDto, err error)
	EditPayment(req *common.RequestContext, d dto.PaymentBody) (err error)
	ExportPayment(req *common.RequestContext, i dto.ExportDocumentData) (res []byte, err error)
}

type paymentUseCase struct {
	paymentRepo payment_repo.PaymentRepository
	emitLog     logger.EmitLog
	permission  repository.PermissionService
	core repository.CoreService

	// setting           repository.SettingService
	bus bus.Bus
	c   di.Container
	fsm fsm.FsmState
	paymentPdf payment_pdf.PaymentPDF
}

func NewPaymentUseCase(
	paymentRepo payment_repo.PaymentRepository,
	logger logger.Logger,
	permission repository.PermissionService,
	// setting repository.SettingService,
	bus bus.Bus,
	c di.Container,
	fsm fsm.FsmState,
	paymentPdf payment_pdf.PaymentPDF,
	core repository.CoreService,
) PaymentUseCase {
	return &paymentUseCase{
		paymentRepo: paymentRepo,
		permission:  permission,
		emitLog:     logger.EmitLog("payment-usecase"),
		bus:         bus,
		c:           c,
		fsm:         fsm,
		core: core,
		paymentPdf: paymentPdf,
	}
}

func(u *paymentUseCase) ExportPayment(req *common.RequestContext, i dto.ExportDocumentData) (res []byte, err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("ExportPayment"))
		}
	}()
	paymentDetail,err := u.GetPaymentDetail(req,&dto.RequestEntity{
		ID: i.ID,
	})
	if err != nil {
		return
	}
	res,err  =u.paymentPdf.GeneratePaymentDocument(req,paymentDetail.Entity)
	return
}

func (u *paymentUseCase) EditPayment(req *common.RequestContext, d dto.PaymentBody) (err error) {
	ctx := u.c.Scoped(req.Ctx)
	tx := di.Get(ctx, domain.DatabaseTransactionKey).(*query.QueryTx)
	defer func(tx *query.QueryTx) {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("EditPayment"))
		}
		err = u.closeTx(tx, err)
	}(tx)
	err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.PAYMENT, domain.EDIT)
	if err != nil {
		return
	}
	err = u.paymentRepo.EditPayment(req, d)
	if err != nil {
		return
	}
	err = u.bus.Emit(req.Ctx, domain.PaymentEditedEvent, event.PaymentEventData{
		Tx:   tx,
		Body: d,
	})
	return
}

func (u *paymentUseCase) GetPaymentAccounts(req *common.RequestContext) (res dto.PaymentAccountsDto, err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetPaymentAccounts"))
		}
	}()
	if allow := u.permission.CheckPermission(req.Ctx, req, domain.PAYMENT, domain.VIEW); !allow {
		return
	}
	res, err = u.paymentRepo.GetPaymentAccounts(req)
	return
}

func (s *paymentUseCase) UpdatePaymentState(req *common.RequestContext, i *dto.UpdateStatusWithEvent) (err error) {
	ctx := s.c.Scoped(req.Ctx)
	tx := di.Get(ctx, domain.DatabaseTransactionKey).(*query.QueryTx)
	defer func() {
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("UpdatePaymentState"))
		}
		err = s.closeTx(tx, err)
		fmt.Println("TRANSACTION ERROR", err)
	}()
	if allow := s.permission.CheckPermission(req.Ctx, req, domain.PAYMENT, domain.EDIT); !allow {
		return domain.ACTION_NOT_ALLOWED
	}
	nextState, err := s.fsm.NextState(i.Body.CurrentState, i.Body.Events)
	if err != nil {
		return err
	}
	payment, paymentRef, err := s.paymentRepo.UpdatePaymentState(tx, req, i.Body.PartyID, i.Body.CurrentState, nextState)
	if err != nil {
		return
	}
	payload := event.StatusPaymentEventData{
		Tx:              tx,
		Payment:         *payment,
		References:      paymentRef,
		CompanyDefaults: req.CompanyDefaults,
	}
	switch nextState {
	case proto.State_SUBMITTED.String():
		err = s.bus.Emit(req.Ctx, domain.PaymentSubmittedEvent, payload)
	case proto.State_CANCELLED.String():
		err = s.bus.Emit(req.Ctx, domain.PaymentCancelledEvent, payload)
	}
	return
}

func (u *paymentUseCase) GetPayments(req *common.RequestContext, i *dto.RequestPayments) (
	res dto.PaginationResult[[]dto.PaymentDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetPayments"))
		}
	}()
	if allow := u.permission.CheckPermission(req.Ctx, req, domain.PAYMENT, domain.VIEW); !allow {
		return
	}

	res, err = u.paymentRepo.GetPayments(req, i)
	if err != nil {
		return
	}
	res.FilterOptions = u.paymentRepo.GetFilterOptions(string(req.LanguageCode))
	return
}

func (u *paymentUseCase) GetPaymentDetail(req *common.RequestContext, i *dto.RequestEntity) (
	res dto.ResultEntity[dto.PaymentDetailDto], err error,
) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetPaymentDetail"))
		}
	}()
	if allow := u.permission.CheckPermission(req.Ctx, req, domain.PAYMENT, domain.VIEW); !allow {
		return
	}
	res, err = u.paymentRepo.GetPaymentDetail(req, i)
	if err != nil {
		return
	}
	res.Activities = u.core.GerActivitiesByPartyID(req,res.Entity.ID)
	return
}

func (u *paymentUseCase) GetAllowedParties(req *common.RequestContext) []dto.PartyTypeDto {
	return u.paymentRepo.GetAllowedParties(req)
}

func (u *paymentUseCase) CreatePayment(req *common.RequestContext, d *dto.CreatePaymentRequest) (
	res dto.PaymentDto, err error) {
	ctx := u.c.Scoped(req.Ctx)
	tx := di.Get(ctx, domain.DatabaseTransactionKey).(*query.QueryTx)
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("CreatePayment"))
		}
		err = u.closeTx(tx, err)
		fmt.Println("TRANSACTION ERROR", err)
	}()
	if allow := u.permission.CheckPermission(req.Ctx, req, domain.PAYMENT, domain.CREATE); !allow {
		return res, domain.ACTION_NOT_ALLOWED
	}
	payment, err := u.paymentRepo.CreatePayment(req, tx, d)
	if err != nil {
		return
	}
	//Seding event to the subscribers
	u.bus.Emit(req.Ctx, domain.PaymentCreatedEvent, event.PaymentEventData{
		Body:    d.Body,
		Payment: payment,
		Tx:      tx,
	})
	res = dto.PaymentDtoFromModel(&payment)
	return res, err
}

func (s *paymentUseCase) closeTx(tx *query.QueryTx, err error) error {
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
