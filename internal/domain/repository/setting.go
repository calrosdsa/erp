package repository

import (
	"context"
	"erp/gen/db/model"

)

type SettingService interface {
	GetAccountSettings(ctx context.Context, companyID int64) (model.AccountSetting, error)
	GetStockSettings(ctx context.Context, companyID int64) (model.StockSetting, error)
	// GetCurrencyExchange(ctx context.Context, fromCurrency, toCurrency string, 
	// forSelling, forBuying bool,companyID int64)(int,error)
}
