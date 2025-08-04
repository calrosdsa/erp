package acct_report_ucase

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	"erp/pkg/logger"
	acct_report_repo "erp/project/accounting/report/internal/repository"
)

type AcctReportUseCase interface {
	GetGeneralLedgerReport(req *common.RequestContext, i *dto.RequestGeneralLedger) (
		dto.GeneralLedgerData, error)
	GetAccountPayable(req *common.RequestContext, i *dto.RequestAccountPayable) (
		[]dto.AccountPayableEntryDto, error)
	GetAccountPayableSumary(req *common.RequestContext, i *dto.RequestAccountPayable) (
		[]dto.SumaryEntryDto, error)
	GetAccountReceivableSumary(req *common.RequestContext, i *dto.RequestAccountReceivable) (
		res []dto.SumaryEntryDto, err error)
	GetAccountReceivable(req *common.RequestContext, i *dto.RequestAccountReceivable) (
		res []dto.AccountReceivableEntryDto, err error)
	GetAccountBalance(req *common.RequestContext,d *dto.RequestAccountBalance)(res dto.GeneralLedgerOpening,err error)
}

type acctReportUseCase struct {
	emitLog        logger.EmitLog
	permission     repository.PermissionService
	acctReportRepo acct_report_repo.AcctReportRepository
}

func NewAcctReportUseCase(
	logger logger.Logger,
	permission repository.PermissionService,
	acctReportRepo acct_report_repo.AcctReportRepository,
) AcctReportUseCase {
	return &acctReportUseCase{
		emitLog:        logger.EmitLog("acct-report-usecase"),
		permission:     permission,
		acctReportRepo: acctReportRepo,
	}
}
func(u *acctReportUseCase)GetAccountBalance(req *common.RequestContext,d *dto.RequestAccountBalance)(res dto.GeneralLedgerOpening,err error) {
	defer func(){
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetAccountBalance"))
		}
	}()
	err = u.permission.CheckPermissionEntity(req.Ctx,req,domain.LEDGER,domain.VIEW)
	if err != nil {
		return
	}
	res, err = u.acctReportRepo.GetAccountBalance(req, d)
	return
}

func (u *acctReportUseCase) GetAccountReceivableSumary(req *common.RequestContext, i *dto.RequestAccountReceivable) (
	res []dto.SumaryEntryDto, err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetAccountReceivableSumary"))
		}
	}()
	err = u.permission.CheckPermissionEntity(req.Ctx,req,domain.ACCOUNT_RECEIVABLE,domain.VIEW)
	if err != nil {
		return
	}
	res, err = u.acctReportRepo.GetAccountReceivableSumary(req, i)
	return
}
func (u *acctReportUseCase) GetAccountReceivable(req *common.RequestContext, i *dto.RequestAccountReceivable) (
	res []dto.AccountReceivableEntryDto, err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetAccountReceivable"))
		}
	}()
	err = u.permission.CheckPermissionEntity(req.Ctx,req,domain.ACCOUNT_RECEIVABLE,domain.VIEW)
	if err != nil {
		return
	}
	res, err = u.acctReportRepo.GetAccountReceivable(req, i)
	return res, err
}

func (u *acctReportUseCase) GetAccountPayableSumary(req *common.RequestContext, i *dto.RequestAccountPayable) (
	res []dto.SumaryEntryDto, err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetAccountPayableSumary"))
		}
	}()
	err = u.permission.CheckPermissionEntity(req.Ctx,req,domain.ACCOUNT_PAYABLE,domain.VIEW)
	if err != nil {
		return
	}
	res, err = u.acctReportRepo.GetAccountPayableSumary(req, i)
	return res, err
}

func (u *acctReportUseCase) GetAccountPayable(req *common.RequestContext, i *dto.RequestAccountPayable) (
	res []dto.AccountPayableEntryDto, err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetAccountPayable"))
		}
	}()
	err = u.permission.CheckPermissionEntity(req.Ctx,req,domain.ACCOUNT_PAYABLE,domain.VIEW)
	if err != nil {
		return
	}
	res, err = u.acctReportRepo.GetAccountPayable(req, i)
	return res, err
}

func (u *acctReportUseCase) GetGeneralLedgerReport(req *common.RequestContext, i *dto.RequestGeneralLedger) (
	res dto.GeneralLedgerData, err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetGeneralLedgerReport"))
		}
	}()
	err = u.permission.CheckPermissionEntity(req.Ctx,req,domain.ACCOUNT_PAYABLE,domain.VIEW)
	if err != nil {
		return
	}
	res, err = u.acctReportRepo.GetGeneralLedgerReport(req, i)
	return res, err
}
