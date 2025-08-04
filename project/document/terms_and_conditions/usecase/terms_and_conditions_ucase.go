package terms_and_conditions_ucase

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	"erp/pkg/fsm"
	"erp/pkg/logger"
	terms_and_conditions_repo "erp/project/document/terms_and_conditions/repository"
)

type TermsAndConditionsUcase interface {
	GetTermsAndConditions(req *common.RequestContext, d *dto.TermsAndConditionsRequest) (
		res dto.ResponseDataList[[]dto.TermsAndConditionsDto], err error)
	GetTermsAndConditionsDetial(req *common.RequestContext, d *dto.RequestEntity) (
		res dto.ResultEntity[dto.TermsAndConditionsDto], err error)
	CreateTermsAndConditions(req *common.RequestContext, d dto.TermsAndConditionsData) (res dto.TermsAndConditionsDto, err error)
	EditTermsAndConditions(req *common.RequestContext, d dto.TermsAndConditionsData) (err error)
	UpdateStatus(req *common.RequestContext, d *dto.UpdateStatusWithEvent) (err error)
}
type termsAndConditionsUcase struct {
	emitLog                logger.EmitLog
	core                   repository.CoreService
	permission             repository.PermissionService
	fsm                    fsm.FsmState
	termsAndConditionsRepo terms_and_conditions_repo.TermsAndConditionsRepository
}

func NewUseCase(
	logger logger.Logger,
	core repository.CoreService,
	permission repository.PermissionService,
	termsAndConditionsRepo terms_and_conditions_repo.TermsAndConditionsRepository,
	fsm fsm.FsmState,
) TermsAndConditionsUcase {
	return &termsAndConditionsUcase{
		emitLog:                logger.EmitLog("terms-and-conditions-ucase"),
		core:                   core,
		permission:             permission,
		termsAndConditionsRepo: termsAndConditionsRepo,
		fsm:                    fsm,
	}
}

func (u *termsAndConditionsUcase) GetTermsAndConditions(req *common.RequestContext, d *dto.TermsAndConditionsRequest) (
	res dto.ResponseDataList[[]dto.TermsAndConditionsDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetTermsAndConditions"))
		}
	}()
	if allow := u.permission.CheckPermission(req.Ctx, req, domain.TERMS_AND_CONDITIONS, domain.VIEW); !allow {
		return res, domain.ACTION_NOT_ALLOWED
	}
	res.Body.Result, err = u.termsAndConditionsRepo.GetTermsAndConditions(req, d)
	if err != nil {
		return
	}
	res.Body.FilterOptions = u.termsAndConditionsRepo.GetFilterOptions(string(req.LanguageCode))
	return
}
func (u *termsAndConditionsUcase) GetTermsAndConditionsDetial(req *common.RequestContext, d *dto.RequestEntity) (
	res dto.ResultEntity[dto.TermsAndConditionsDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetTermsAndConditionsDetial"))
		}
	}()
	if allow := u.permission.CheckPermission(req.Ctx, req, domain.TERMS_AND_CONDITIONS, domain.VIEW); !allow {
		return res, domain.ACTION_NOT_ALLOWED
	}
	res.Entity, err = u.termsAndConditionsRepo.GetTermsAndConditionsDetial(req, d)
	if err != nil {
		return
	}
	res.Activities = u.core.GerActivitiesByPartyID(req, res.Entity.ID)

	return
}

func (u *termsAndConditionsUcase) CreateTermsAndConditions(req *common.RequestContext, d dto.TermsAndConditionsData) (
	res dto.TermsAndConditionsDto, err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("CreateTermsAndConditions"))
		}
	}()
	if allow := u.permission.CheckPermission(req.Ctx, req, domain.TERMS_AND_CONDITIONS, domain.CREATE); !allow {
		return res, domain.ACTION_NOT_ALLOWED
	}
	termsAndConditions, err := u.termsAndConditionsRepo.CreateTermsAndConditions(req, d)
	if err != nil {
		return
	}
	res = dto.TermsAndConditionFromModel(termsAndConditions)
	return
}
func (u *termsAndConditionsUcase) EditTermsAndConditions(req *common.RequestContext, d dto.TermsAndConditionsData) (
	err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("EditTermsAndConditions"))
		}
	}()
	if allow := u.permission.CheckPermission(req.Ctx, req, domain.TERMS_AND_CONDITIONS, domain.EDIT); !allow {
		return domain.ACTION_NOT_ALLOWED
	}
	err = u.termsAndConditionsRepo.EditTermsAndConditions(req, d)
	return
}

func (u *termsAndConditionsUcase) UpdateStatus(req *common.RequestContext, d *dto.UpdateStatusWithEvent) (err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("UpdateState"))
		}
	}()
	if allow := u.permission.CheckPermission(req.Ctx, req, domain.TERMS_AND_CONDITIONS, domain.EDIT); !allow {
		return domain.ACTION_NOT_ALLOWED
	}
	nextState, err := u.fsm.NextState(d.Body.CurrentState, d.Body.Events)
	if err != nil {
		return err
	}
	err = u.termsAndConditionsRepo.UpdateStatus(req, d, nextState)
	return
}
