package setting_service

import (
	"context"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/internal/domain/repository"
	"erp/pkg/db"
	"erp/pkg/logger"
)

type settingService struct {
	emitLog logger.EmitLog
	Q       *query.Query
}

func NewSettingService(
	logger logger.Logger,
	conn db.Connection,
) repository.SettingService {
	return &settingService{
		Q:       conn.GetQ(),
		emitLog: logger.EmitLog("account-service"),
	}
}
// func(s *settingService) GetCurrencyExchange(ctx context.Context, fromCurrency, toCurrency string, 
// 	forSelling, forBuying bool,companyID int64)(err error) {
	
// 	return
// }

func(s *settingService) GetLedger(ctx context.Context,ledgerID int64)(res model.Ledger,err error){
	ledgerQ := s.Q.Ledger
	err = ledgerQ.WithContext(ctx).Where(ledgerQ.ID.Eq(ledgerID)).Select(
		ledgerQ.ID,ledgerQ.Name,
	).Scan(&res)
	return 
}

func (s *settingService) GetStockSettings(ctx context.Context, companyID int64) (res model.StockSetting, err error) {
	defer func() {
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("GetStockSettings"))
		}
	}()
	stockSetting, err := s.Q.StockSetting.WithContext(ctx).Where(
		s.Q.StockSetting.CompanyID.Eq(companyID),
	).First()
	if err != nil {
		return
	}
	return *stockSetting, err
}



func (s *settingService) GetAccountSettings(ctx context.Context, companyID int64) (res model.AccountSetting, err error) {
	defer func() {
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("GetAccountSettings"))
		}
	}()
	accountSetting, err := s.Q.AccountSetting.WithContext(ctx).Where(
		s.Q.AccountSetting.CompanyID.Eq(companyID),
	).First()
	if err != nil {
		return
	}
	return *accountSetting, err
}
