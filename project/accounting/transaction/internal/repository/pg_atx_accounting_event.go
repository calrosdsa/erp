package transaction_repo

import (
	"context"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/gen/proto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain/event"
	"erp/internal/domain/repository"
	"errors"
	"fmt"
	"time"
)

type AtxAccountignEventRepo interface {
	OnPaymentSubmitted(ctx context.Context, payload event.StatusPaymentEventData) (err error)
	OnPaymentCancelled(ctx context.Context, payload event.StatusPaymentEventData) (err error)
	OnJournalEntrySubmitted(ctx context.Context, payload event.StatusJournalEntryEventData) (err error)
	OnJournalEntryCancelled(ctx context.Context, payload event.StatusJournalEntryEventData) (err error)
	OnCashOutflowSubmitted(ctx context.Context,payload event.StatusCashOutflowEventData)(err error)
	OnCashOutflowCancelled(ctx context.Context,payload event.StatusCashOutflowEventData)(err error)
}

type atxAccountingEventRepo struct {
	transactionRepo TransactionRepository
	accounting      repository.AccountingService
	currency        helpers.CurrencyHelper
}

func NewAtxAccountingEventRepo(
	transactionRepo TransactionRepository,
	accounting repository.AccountingService,
	helpers *helpers.Helpers,
) AtxAccountignEventRepo {
	return &atxAccountingEventRepo{
		transactionRepo: transactionRepo,
		accounting:      accounting,
		currency:        helpers.Currency,
	}
}

func (r *atxAccountingEventRepo)OnCashOutflowSubmitted(ctx context.Context,payload event.StatusCashOutflowEventData)(err error){
	tx := payload.Tx
	cashOutflow := payload.CashOutflow
	companyDefaults := payload.CompanyDefaults
	var transactions []*model.TransactionLedger

	docAccounts,err := tx.DocAccount.Where(
		tx.DocAccount.DocID.Eq(cashOutflow.ID),
	).First()
	if err != nil {
		return
	}

	creditAccount,err := r.accounting.GetLedger(ctx,tx.Query,*docAccounts.CreditAccountID,true)
	if err != nil {
		return
	}
	debitAccount,err := r.accounting.GetLedger(ctx,tx.Query,*docAccounts.DebitAccountID,true)
	if err != nil {
		return
	}
	if debitAccount.Currency != creditAccount.Currency {
		return errors.New("currency accounts differ from each other")
	}
	exchangeRate, err := r.accounting.GetCurrencyExchangeRate(ctx, tx, companyDefaults, creditAccount.Currency, true, true)
	if err != nil {
		return fmt.Errorf("failed to get exchange rate for currency %s: %w", creditAccount.Currency, err)
	}
	var (
		totalAmount int64
		totalAmountPlusCharges int64
	)
	
	tacTransactions, totalTaxAmount, err := r.transactionRepo.ProcessTaxLines(ctx, tx, cashOutflow.ID, cashOutflow.PostingDate,
		cashOutflow.PostingTime, creditAccount.Currency, cashOutflow.Code,
		proto.PartyType_cashOutflow.String(), *cashOutflow.CashOutflowType, cashOutflow.CostCenterID,
		cashOutflow.ProjectID, payload.TaxLinesData.TaxLines, exchangeRate)
	if err != nil {
		return
	}
	totalAmount = r.currency.CurrencyExchange(cashOutflow.Amount, exchangeRate)

	totalAmountPlusCharges = totalAmount+ totalTaxAmount

	creditEntry := model.TransactionLedger{
		Ledger:         *docAccounts.CreditAccountID,
		Credit:         totalAmountPlusCharges,
		VoucherType:    proto.PartyType_cashOutflow.String(),
		VoucherCode:    cashOutflow.Code,
		PostingDate: cashOutflow.PostingDate,
		PostingTime: cashOutflow.PostingTime,
		VoucherSubtype: *cashOutflow.CashOutflowType,
		Currency: creditAccount.Currency,
		ProjectID: cashOutflow.ProjectID,
		CostCenterID: cashOutflow.CostCenterID,
	}
	debitEntry := model.TransactionLedger{
		Ledger:         *docAccounts.DebitAccountID,
		Debit:          totalAmount,
		VoucherType:    proto.PartyType_cashOutflow.String(),
		VoucherCode:    cashOutflow.Code,
		PostingDate: cashOutflow.PostingDate,
		PostingTime: cashOutflow.PostingTime,
		VoucherSubtype: *cashOutflow.CashOutflowType,
		PartyID: &cashOutflow.PartyID,
		Currency: debitAccount.Currency,
		ProjectID: cashOutflow.ProjectID,
		CostCenterID: cashOutflow.CostCenterID,
	}
	transactions = append(transactions,&creditEntry,&debitEntry)
	transactions = append(transactions, tacTransactions...)
	
	err = r.transactionRepo.SaveTransactionsBatch(ctx, tx, transactions)
	return
}
func (r *atxAccountingEventRepo)OnCashOutflowCancelled(ctx context.Context,payload event.StatusCashOutflowEventData)(err error){
	err = r.accounting.DelTxnsByVoucherCode(ctx, payload.Tx, payload.CashOutflow.Code)
	return
}

func (r *atxAccountingEventRepo) OnPaymentCancelled(ctx context.Context, payload event.StatusPaymentEventData) (err error) {
	err = r.accounting.DelTxnsByVoucherCode(ctx, payload.Tx, payload.Payment.Code)
	return
}

func (r *atxAccountingEventRepo) OnJournalEntryCancelled(ctx context.Context,
	payload event.StatusJournalEntryEventData) (err error) {
	err = r.accounting.DelTxnsByVoucherCode(ctx, payload.Tx, payload.JournalEntry.Code)
	return

}

func (r *atxAccountingEventRepo) OnJournalEntrySubmitted(ctx context.Context,
	payload event.StatusJournalEntryEventData) (err error) {
	tx := payload.Tx
	journalEntry := payload.JournalEntry
	fmt.Println("JOURNAL ENTRY TX ACCT", journalEntry)
	transactions := make([]*model.TransactionLedger, len(payload.Lines))
	for i, line := range payload.Lines {
		transaction := &model.TransactionLedger{}
		transaction.Credit = int64(line.Credit)
		transaction.Debit = int64(line.Debit)
		transaction.Ledger = line.LedgerID
		transaction.VoucherCode = journalEntry.Code
		transaction.VoucherType = proto.PartyType_journalEntry.String()
		transaction.VoucherSubtype = journalEntry.EntryType
		transaction.Currency = line.Currency
		transaction.ProjectID = line.ProjectID
		transaction.CostCenterID = line.CostCenterID
		transaction.PostingDate = journalEntry.PostingDate
		transactions[i] = transaction
	}
	err = tx.WithContext(ctx).TransactionLedger.CreateInBatches(transactions, len(transactions))
	return
}

func (r *atxAccountingEventRepo) OnPaymentSubmitted(ctx context.Context, payload event.StatusPaymentEventData) (err error) {
	tx := payload.Tx
	payment := payload.Payment
	companyDefaults := payload.CompanyDefaults
	paymentReferences, err := r.getPaymentReferences(ctx, tx, payment.ID)
	if err != nil {
		return
	}
	var transactions []*model.TransactionLedger
	accountFrom,err := r.accounting.GetLedger(ctx,tx.Query,payment.AccountPaidFromID,true)
	if err != nil {
		return
	}
	accountTo,err := r.accounting.GetLedger(ctx,tx.Query,payment.AccountPaidToID,true)
	if err != nil {
		return
	}
	if accountTo.Currency != accountFrom.Currency {
		return errors.New("currency accounts differ from each other")
	}
	fmt.Println("CURRENCY ACCOUNTS",accountFrom.Currency,accountTo.Currency)
	exchangeRate, err := r.accounting.GetCurrencyExchangeRate(ctx, tx, companyDefaults, accountFrom.Currency, true, true)
	if err != nil {
		return
	}
	var (
		totalAmount int64
		totalAmountPlusCharges int64
	)
	for _, ref := range paymentReferences {
		exchangeRateRef, err1 := r.accounting.GetCurrencyExchangeRate(ctx, tx, companyDefaults, ref.Currency, true, true)
		if err1 != nil {
			return err1
		}
		amountRef := r.currency.CurrencyExchange(ref.Allocated, exchangeRateRef)
		totalAmount += amountRef
	}
	// if not payment references
	if len(paymentReferences) == 0 {
		totalAmount += payment.Amount
		// err = r.AddAccountingEntries(ctx, tx, payment.Amount, payment,exchangeRate)
	}

	tacTransactions, totalTaxAmount, err := r.transactionRepo.ProcessTaxLines(ctx, tx, payment.ID, payment.PostingDate,
		time.Now().Format(time.TimeOnly), accountTo.Currency, payment.Code,
		proto.PartyType_payment.String(), payment.PaymentType, payment.CostCenterID,
		payment.ProjectID, payload.TaxLinesData.TaxLines, exchangeRate)
	if err != nil {
		return
	}
	totalAmount = r.currency.CurrencyExchange(totalAmount, exchangeRate)

	totalAmountPlusCharges = totalAmount+ totalTaxAmount

	paidFromEntry := model.TransactionLedger{
		Ledger:         payment.AccountPaidFromID,
		LedgerAgainst:  &payment.AccountPaidToID,
		Credit:         totalAmountPlusCharges,
		VoucherType:    proto.PartyType_payment.String(),
		VoucherCode:    payment.Code,
		VoucherSubtype: payment.PaymentType,
		PostingDate: payment.PostingDate,
		Currency: accountFrom.Currency,
		ProjectID: payment.ProjectID,
		CostCenterID: payment.CostCenterID,
	}
	paidToEntry := model.TransactionLedger{
		Ledger:         payment.AccountPaidToID,
		LedgerAgainst:  &payment.AccountPaidFromID,
		Debit:          totalAmount,
		VoucherType:    proto.PartyType_payment.String(),
		VoucherCode:    payment.Code,
		VoucherSubtype: payment.PaymentType,
		PostingDate: payment.PostingDate,
		Currency: accountTo.Currency,
		ProjectID: payment.ProjectID,
		CostCenterID: payment.CostCenterID,
	}
	switch payment.PaymentType {
	case proto.PaymentType_PAY.String():
		paidToEntry.PartyID = &payment.PartyID
	case proto.PaymentType_RECEIVE.String():
		paidFromEntry.PartyID = &payment.PartyID
	}
	transactions = append(transactions, &paidToEntry, &paidFromEntry,)
	transactions = append(transactions, tacTransactions...)
	
	err = r.transactionRepo.SaveTransactionsBatch(ctx, tx, transactions)

	return
}

func (r *atxAccountingEventRepo) AddAccountingEntries(ctx context.Context, tx *query.QueryTx,
	amount int64, payment model.Payment, exchangeRate int32) (err error) {
	amount = r.currency.CurrencyExchange(amount, exchangeRate)
	paidFromEntry := model.TransactionLedger{
		Ledger:         payment.AccountPaidFromID,
		LedgerAgainst:  &payment.AccountPaidToID,
		Credit:         amount,
		VoucherType:    proto.PartyType_payment.String(),
		VoucherCode:    payment.Code,
		VoucherSubtype: payment.PaymentType,
	}
	paidToEntry := model.TransactionLedger{
		Ledger:         payment.AccountPaidToID,
		LedgerAgainst:  &payment.AccountPaidFromID,
		Debit:          int64(amount),
		VoucherType:    proto.PartyType_payment.String(),
		VoucherCode:    payment.Code,
		VoucherSubtype: payment.PaymentType,
	}
	switch payment.PaymentType {
	case proto.PaymentType_PAY.String():
		paidToEntry.PartyID = &payment.PartyID
	case proto.PaymentType_RECEIVE.String():
		paidFromEntry.PartyID = &payment.PartyID
	}
	err = r.transactionRepo.SaveTransactionsBatch(ctx, tx, []*model.TransactionLedger{
		&paidToEntry, &paidFromEntry,
	})
	return
}

func (r *atxAccountingEventRepo) getPaymentReferences(ctx context.Context, tx *query.QueryTx, paymentID int64,
) (res []*model.PaymentReference, err error) {
	refQ := tx.PaymentReference
	res, err = refQ.WithContext(ctx).Where(
		refQ.PaymentID.Eq(paymentID),
	).Find()
	return
}
