package payment_terms_ucase

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	"erp/pkg/fsm"
	"erp/pkg/logger"
	payment_terms_repo "erp/project/document/payment_term/repository"
)

type PaymentTermsUcase interface {
	GetPaymentTerms(req *common.RequestContext, d *dto.PaymentTermsRequest) (
		res dto.ResponseDataList[[]dto.PaymentTermsDto], err error)
	GetPaymentTermsDetail(req *common.RequestContext, d *dto.RequestEntity) (
		res dto.ResultEntity[dto.PaymentTermsDto], err error)
	CreatePaymentTerms(req *common.RequestContext, d dto.PaymentTermsData) (res dto.PaymentTermsDto, err error)
	EditPaymentTerms(req *common.RequestContext, d dto.PaymentTermsData) (err error)
	UpdateStatus(req *common.RequestContext, d *dto.UpdateStatusWithEvent) (err error)

	GetPaymentTermLines(req *common.RequestContext, d dto.RequestEntity) (res []dto.PaymentTermsLineDto, err error)
}
type paymentTermsUcase struct {
	emitLog          logger.EmitLog
	core             repository.CoreService
	permission       repository.PermissionService
	fsm              fsm.FsmState
	paymentTermsRepo payment_terms_repo.PaymentTermsRepo
	paymentTermLineRepo payment_terms_repo.PaymentTermsLineRepo
}

func NewUseCase(
	logger logger.Logger,
	core repository.CoreService,
	permission repository.PermissionService,
	paymentTermsRepo payment_terms_repo.PaymentTermsRepo,
	fsm fsm.FsmState,
	paymentTermLineRepo payment_terms_repo.PaymentTermsLineRepo,
) PaymentTermsUcase {
	return &paymentTermsUcase{
		emitLog:          logger.EmitLog("payment-terms-ucase"),
		core:             core,
		permission:       permission,
		paymentTermsRepo: paymentTermsRepo,
		fsm:              fsm,
		paymentTermLineRepo: paymentTermLineRepo,
	}
}

func (u *paymentTermsUcase) GetPaymentTermLines(req *common.RequestContext, d dto.RequestEntity) (res []dto.PaymentTermsLineDto, err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetPaymentTermLines"))
		}
	}()
	res,err = u.paymentTermLineRepo.GetPaymentTermLines(req,d)
	return 
}

func (u *paymentTermsUcase) GetPaymentTerms(req *common.RequestContext, d *dto.PaymentTermsRequest) (
	res dto.ResponseDataList[[]dto.PaymentTermsDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetPaymentTerms"))
		}
	}()
	if err := u.permission.CheckPermissionEntity(req.Ctx, req, domain.PAYMENT_TERMS, domain.VIEW); err != nil  {
		return res, err
	}
	res.Body.Result, err = u.paymentTermsRepo.GetPaymentTerms(req, d)
	if err != nil {
		return
	}
	res.Body.FilterOptions = u.paymentTermsRepo.GetFilterOptions(string(req.LanguageCode))
	return
}
func (u *paymentTermsUcase) GetPaymentTermsDetail(req *common.RequestContext, d *dto.RequestEntity) (
	res dto.ResultEntity[dto.PaymentTermsDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetPaymentTermsDetail"))
		}
	}()
	if err := u.permission.CheckPermissionEntity(req.Ctx, req, domain.PAYMENT_TERMS, domain.VIEW);err != nil  {
		return res, err
	}
	res.Entity, err = u.paymentTermsRepo.GetPaymentTermsDetail(req, d)
	if err != nil {
		return
	}
	res.Activities = u.core.GerActivitiesByPartyID(req, res.Entity.ID)

	return
}

func (u *paymentTermsUcase) CreatePaymentTerms(req *common.RequestContext, d dto.PaymentTermsData) (
	res dto.PaymentTermsDto, err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("CreatePaymentTerms"))
		}
	}()
	if err := u.permission.CheckPermissionEntity(req.Ctx, req, domain.PAYMENT_TERMS, domain.CREATE); err != nil {
		return res, err
	}
	paymentTerms, err := u.paymentTermsRepo.CreatePaymentTerms(req, d)
	if err != nil {
		return
	}
	res = dto.PaymentTermsFromModel(paymentTerms)
	return
}
func (u *paymentTermsUcase) EditPaymentTerms(req *common.RequestContext, d dto.PaymentTermsData) (
	err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("EditPaymentTerms"))
		}
	}()
	if err := u.permission.CheckPermissionEntity(req.Ctx, req, domain.PAYMENT_TERMS, domain.EDIT); err != nil {
		return err
	}
	err = u.paymentTermsRepo.EditPaymentTerms(req, d)
	return
}

func (u *paymentTermsUcase) UpdateStatus(req *common.RequestContext, d *dto.UpdateStatusWithEvent) (err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("UpdateState"))
		}
	}()
	
	
	return
}
