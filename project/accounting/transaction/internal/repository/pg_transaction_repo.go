package transaction_repo

import (
	"context"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/gen/proto"
	"erp/internal/app/service/helpers"
	"time"
)

type TransactionRepository interface {
	GetStockSettings(ctx context.Context, tx *query.QueryTx, companyID int64) (model.StockSetting, error)
	GetAccountSettings(ctx context.Context, tx *query.QueryTx, companyID int64) (model.AccountSetting, error)
	SaveTransaction(ctx context.Context, tx *query.QueryTx, model model.TransactionLedger) error
	SaveTransactionsBatch(ctx context.Context, tx *query.QueryTx, model []*model.TransactionLedger) error
	ProcessTaxLines(ctx context.Context, tx *query.QueryTx, docPartyID int64,
		postingDate time.Time, postingTime string, currency,
		voucherCode, voucherType, voucherSubtype string, costCenerID, projectID *int64,
		taxLines []dto.TaxLineWhitAccountHead,exchangeRate int32) (
		res []*model.TransactionLedger, total int64, err error)
}

type transactionRepository struct {
	currency helpers.CurrencyHelper
}

func NewTransactionRepository(
	helpers *helpers.Helpers,
) TransactionRepository {
	return &transactionRepository{
		currency:helpers.Currency,
	}
}



func (r *transactionRepository) SaveTransactionsBatch(ctx context.Context, tx *query.QueryTx,
	d []*model.TransactionLedger) (err error) {
	err = tx.TransactionLedger.WithContext(ctx).CreateInBatches(d, len(d))
	return
}

func (r *transactionRepository) SaveTransaction(ctx context.Context, tx *query.QueryTx,
	d model.TransactionLedger) (err error) {
	err = tx.TransactionLedger.WithContext(ctx).Save(&d)
	return
}

func (r *transactionRepository) GetStockSettings(ctx context.Context, tx *query.QueryTx, companyID int64,
) (res model.StockSetting, err error) {
	stockSetting, err := tx.StockSetting.WithContext(ctx).Where(
		tx.StockSetting.CompanyID.Eq(companyID),
	).First()
	if err != nil {
		return
	}
	return *stockSetting, err
}
func (r *transactionRepository) GetAccountSettings(ctx context.Context, tx *query.QueryTx, companyID int64,
) (res model.AccountSetting,
	err error) {
	accountSetting, err := tx.AccountSetting.WithContext(ctx).Where(
		tx.AccountSetting.CompanyID.Eq(companyID),
	).First()
	if err != nil {
		return
	}
	return *accountSetting, err
}

func (r *transactionRepository) ProcessTaxLines(ctx context.Context, tx *query.QueryTx, docPartyID int64,
	postingDate time.Time, postingTime string,
	currency,voucherCode, voucherType, voucherSubtype string, costCenerID, projectID *int64,
	taxLines []dto.TaxLineWhitAccountHead,exchangeRate int32,
) (res []*model.TransactionLedger, total int64, err error) {
	
	taxAndChargeTxs := make([]*model.TransactionLedger, len(taxLines))

	for i, taxLine := range taxLines {
		amount := r.currency.CurrencyExchange(taxLine.Amount, exchangeRate)
		m := &model.TransactionLedger{}
		m.PostingDate = postingDate
		m.PostingTime = postingTime
		m.Currency = currency
		m.Ledger = taxLine.AccountHead
		m.VoucherCode = voucherCode
		m.VoucherType = voucherType
		m.VoucherSubtype = voucherSubtype
		m.CostCenterID = costCenerID
		m.ProjectID = projectID
		switch taxLine.AcctHeadRootType {
		case proto.AccountType_LIABILITY.String():
			if taxLine.IsOffsetAccount {
				m.Debit = amount
			} else {
				m.Credit = amount
			}
		case proto.AccountType_EXPENSE.String():
			if taxLine.IsOffsetAccount {
				m.Credit = amount
			} else {
				m.Debit = amount
			}
		case proto.AccountType_REVENUE.String():
			if taxLine.IsOffsetAccount {
				m.Debit = amount
			} else {
				m.Credit = amount
			}
		}


		total += amount
		taxAndChargeTxs[i] = m
	}
	return taxAndChargeTxs, total, err
}
