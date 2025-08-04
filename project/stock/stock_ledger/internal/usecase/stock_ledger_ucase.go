package stock_ledger_ucase

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/pkg/logger"
	stock_ledger_repo "erp/project/stock/stock_ledger/internal/repository"
)

type StockLedgerUseCase interface {
	GetStockLedgerReport(req *common.RequestContext, i *dto.RequestStockLedger) (
		res []dto.StockLedgerEntryDto, err error)
	GetStockBalanceReport(req *common.RequestContext, i *dto.RequestStockBalance) (
		res []dto.StockBalanceEntryDto, err error)
}

type stockLedgerUseCase struct {
	stockLedgerRepo stock_ledger_repo.StockLedgerTxRepository
	emitLog         logger.EmitLog
}

func NewStockLedgerUseCase(
	logger logger.Logger,
	stockLedgerRepo stock_ledger_repo.StockLedgerTxRepository,
) StockLedgerUseCase {
	return &stockLedgerUseCase{
		emitLog:         logger.EmitLog("stock-ledger-usecase"),
		stockLedgerRepo: stockLedgerRepo,
	}
}

func (u *stockLedgerUseCase) GetStockBalanceReport(req *common.RequestContext, i *dto.RequestStockBalance) (
	res []dto.StockBalanceEntryDto, err error) {
	defer func ()  {
		if err != nil{
			u.emitLog.Err(err,logger.OptionsLog.WithMethod("GetStockBalanceReport"))			
		}
	} ()
	res,err = u.stockLedgerRepo.GetStockBalanceReport(req,i)
	return 
}

func (u *stockLedgerUseCase) GetStockLedgerReport(req *common.RequestContext, i *dto.RequestStockLedger) (
	res []dto.StockLedgerEntryDto, err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetStockLedgerReport"))
		}
	}()
	res, err = u.stockLedgerRepo.GetStockLedgerReport(req, i)
	return
}
