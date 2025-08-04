package accounting_service

import (
	"context"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	"erp/pkg/logger"

	"gorm.io/gen"
	"gorm.io/gen/field"
)

type accountingService struct {
	emitLog logger.EmitLog
}

func NewAccountingService(
	logger logger.Logger,

) repository.AccountingService {
	return &accountingService{}
}

func (s *accountingService) GetLedger(ctx context.Context,tx *query.Query,ledgerID int64,args ...interface{})(
	res dto.LedgerDto,err error){
	var (
		loadLedgerAccount bool
		columns []field.Expr
	)
	if len(args) > 0 {
		if val,ok := args[0].(bool);ok {
			loadLedgerAccount = val
		}
	}
	columns = append(columns, tx.Ledger.ID,tx.Ledger.AccountType,tx.Ledger.Name)
	builder := tx.Ledger.WithContext(ctx).Where(
		tx.Ledger.ID.Eq(ledgerID),
	)
	if loadLedgerAccount {
		builder = builder.Join(tx.LedgerAccount,tx.LedgerAccount.LedgerID.EqCol(tx.Ledger.ID))
		columns = append(columns, tx.LedgerAccount.Currency)
	}
	err = builder.Select(columns...).Limit(1).Scan(&res)
	return
}

func (s *accountingService) GetCurrencyExchangeRate(ctx context.Context, tx *query.QueryTx, companyDefaults model.CompanyDefault,
	fromCurrency string, forSelling, forBuying bool) (int32, error) {
	if companyDefaults.Currency == fromCurrency {
		return 1, nil
	}
	cExchangeQ := tx.CurrencyExchange
	var conds []gen.Condition
	if forBuying {
		conds = append(conds, cExchangeQ.ForBuying.Is(true))
	}
	if forSelling {
		conds = append(conds, cExchangeQ.ForSelling.Is(true))
	}

	currencyExchange, err := tx.WithContext(ctx).CurrencyExchange.
		Select(cExchangeQ.ExchangeRate).
		Where(
			cExchangeQ.ToCurrency.Eq(companyDefaults.Currency),
			cExchangeQ.FromCurrency.Eq(fromCurrency),
			cExchangeQ.CompanyID.Eq(companyDefaults.CompanyID),
		).First()
	if err != nil {
		return 0, domain.NO_CURRENCY_EXCHANGE_FOUND
	}
	return currencyExchange.ExchangeRate, nil
}

func (s *accountingService) DelTxnsByVoucherCode(ctx context.Context, tx *query.QueryTx, voucherCode string) (err error) {
	_, err = tx.TransactionLedger.WithContext(ctx).Where(
		tx.TransactionLedger.VoucherCode.Eq(voucherCode),
	).Delete()
	return
}
