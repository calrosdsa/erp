package acct_report_ucase

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/pkg/logger"
	acct_report_repo "erp/project/accounting/report/internal/repository"
)

type FinancialStatementUcase interface {
	ProfitAndLossStatement(req *common.RequestContext, d *dto.RequestFinancialStatement) (
		res []dto.ProfitAndLossEntryDto, err error)
	CashFlowStatement(req *common.RequestContext, d *dto.RequestFinancialStatement) (
		res []dto.CashFlowEntryDto, err error)
	BalanceSheetStatement(req *common.RequestContext, d *dto.RequestFinancialStatement) (
		res []dto.BalanceSheetEntryDto, err error)
}

type financialStatementUcase struct {
	emitLog                logger.EmitLog
	financialStatementRepo acct_report_repo.FinancialStatementRepo
}

func NewFinancialStmtUcase(
	logger logger.Logger,
	financialStatementRepo acct_report_repo.FinancialStatementRepo,
) FinancialStatementUcase {
	return &financialStatementUcase{
		emitLog:                logger.EmitLog("financial-statement-usecase"),
		financialStatementRepo: financialStatementRepo,
	}
}
func (r *financialStatementUcase) BalanceSheetStatement(req *common.RequestContext, d *dto.RequestFinancialStatement) (
	res []dto.BalanceSheetEntryDto, err error) {
	defer func() {
		if err != nil {
			r.emitLog.Err(err, logger.OptionsLog.WithMethod("BalanceSheetStatement"))
		}
	}()
	res, err = r.financialStatementRepo.BalanceSheetStatement(req, d)
	return
}
func (r *financialStatementUcase) CashFlowStatement(req *common.RequestContext, d *dto.RequestFinancialStatement) (
	res []dto.CashFlowEntryDto, err error) {		
	defer func() {
		if err != nil {
			r.emitLog.Err(err, logger.OptionsLog.WithMethod("CashFlowStatement"))
		}
	}()
	res, err = r.financialStatementRepo.CashFlowStatement(req, d)
	return
}

func (r *financialStatementUcase) ProfitAndLossStatement(req *common.RequestContext, d *dto.RequestFinancialStatement) (
	res []dto.ProfitAndLossEntryDto, err error) {
	defer func() {
		if err != nil {
			r.emitLog.Err(err, logger.OptionsLog.WithMethod("ProfitAndLossStatement"))
		}
	}()
	res, err = r.financialStatementRepo.ProfitAndLossStatement(req, d)
	return
}
