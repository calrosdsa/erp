package repository

import (
	"context"
	// "erp/gen/db/model"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
)

type AccountingService interface {
	DelTxnsByVoucherCode(ctx context.Context, tx *query.QueryTx, voucherCode string) (err error)
	GetCurrencyExchangeRate(ctx context.Context, tx *query.QueryTx, companyDefaults model.CompanyDefault,
		fromCurrency string, forSelling, forBuying bool) (int32, error)
	GetLedger(ctx context.Context, tx *query.Query, ledgerID int64, args ...interface{}) (res dto.LedgerDto, err error)
}
