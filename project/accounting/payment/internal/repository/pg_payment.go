package payment_repo

import (
	"context"
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/gen/proto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	"erp/pkg/db"
	"strings"

	"gorm.io/gen/helper"
)

type PaymentRepository interface {
	GetPayments(req *common.RequestContext, d *dto.RequestPayments) (
		dto.PaginationResult[[]dto.PaymentDto], error,
	)
	GetPaymentDetail(req *common.RequestContext, d *dto.RequestEntity) (
		dto.ResultEntity[dto.PaymentDetailDto], error)
	CreatePayment(req *common.RequestContext, tx *query.QueryTx, d *dto.CreatePaymentRequest) (
		model.Payment, error)
	GetAllowedParties(req *common.RequestContext) []dto.PartyTypeDto
	UpdatePaymentState(tx *query.QueryTx, req *common.RequestContext, id string, prevState, nexState string) (
		*model.Payment, []*model.PaymentReference, error)
	GetPaymentAccounts(req *common.RequestContext) (res dto.PaymentAccountsDto, err error)
	EditPayment(req *common.RequestContext, d dto.PaymentBody) (err error)
	GetFilterOptions(lng string) []dto.FilterOptionDto
}

type paymentRepository struct {
	locale     helpers.Locale
	Q          *query.Query
	dbHelper   db.DbHelper
	convertor  helpers.ConvertorHelper
	currency   helpers.CurrencyHelper
	setting    repository.SettingService
	query      helpers.QueryHelper
	accounting repository.AccountingService
}

func NewPaymentRepository(
	conn db.Connection,
	helpers *helpers.Helpers,
	setting repository.SettingService,
	accounting repository.AccountingService,
) PaymentRepository {
	return &paymentRepository{
		locale:     helpers.Locale,
		Q:          conn.GetQ(),
		dbHelper:   conn.GetDbHelper(),
		convertor:  helpers.Convertor,
		currency:   helpers.Currency,
		setting:    setting,
		query:      helpers.Query,
		accounting: accounting,
	}
}
func (r *paymentRepository) GetPaymentAccounts(req *common.RequestContext) (
	res dto.PaymentAccountsDto, err error) {
	acctS, err := r.setting.GetAccountSettings(req.Ctx, req.ActiveCompany.ID)
	if err != nil {
		return
	}
	cashAcct, err := r.accounting.GetLedger(req.Ctx, r.Q, acctS.CashAccunt, true)
	if err != nil {
		return
	}
	payableAcct, err := r.accounting.GetLedger(req.Ctx, r.Q, acctS.PayableAccount, true)
	if err != nil {
		return
	}
	receivableAcct, err := r.accounting.GetLedger(req.Ctx, r.Q, acctS.ReceivableAccount, true)
	if err != nil {
		return
	}
	res.CashAcct = cashAcct.Name
	res.CashAcctID = cashAcct.ID
	res.CashAcctCurrency = cashAcct.Currency
	res.PayableAcct = payableAcct.Name
	res.PayableAcctID = payableAcct.ID
	res.PayableAcctCurrency = payableAcct.Currency
	res.ReceivableAcct = receivableAcct.Name
	res.ReceivableAcctID = receivableAcct.ID
	res.ReceivableAcctCurrency = receivableAcct.Currency
	return
}

func (r *paymentRepository) UpdatePaymentState(tx *query.QueryTx, req *common.RequestContext, id string, prevState, nexState string) (
	res *model.Payment, paymentReferences []*model.PaymentReference, err error) {
	paymentQ := tx.Payment
	_, err = tx.Payment.WithContext(req.Ctx).Where(
		paymentQ.CompanyID.Eq(req.ActiveCompany.ID),
		paymentQ.Status.Eq(prevState),
		paymentQ.Code.Eq(id),
	).UpdateSimple(paymentQ.Status.Value(nexState))
	if err != nil {
		return
	}
	res, err = tx.Payment.WithContext(req.Ctx).Where(
		paymentQ.Code.Eq(id),
		paymentQ.CompanyID.Eq(req.ActiveCompany.ID),
	).First()
	if err != nil {
		return
	}
	paymentReferences, err = tx.PaymentReference.WithContext(req.Ctx).Where(
		tx.PaymentReference.PaymentID.Eq(res.ID),
	).Find()
	if err != nil {
		return
	}
	return
}

func (r *paymentRepository) GetPaymentDetail(req *common.RequestContext, i *dto.RequestEntity) (
	res dto.ResultEntity[dto.PaymentDetailDto], err error) {
	paymentQ := r.Q.Payment

	// transactionLedgerQ := r.Q.TransactionLedger
	ledgerFromQ := r.Q.Ledger
	ledgerToQ := r.Q.Ledger.As("ledger_to")
	ledgerAFromQ := r.Q.LedgerAccount
	ledgerAToQ := r.Q.LedgerAccount.As("ledger_account_to")
	partyQ := r.Q.Party
	//parties assoiated with payments
	supplierQ := r.Q.Supplier
	customerQ := r.Q.Customer

	cBAQ := r.Q.BankAccount.As("company_bank_account")
	pBAQ := r.Q.BankAccount.As("party_bank_account")

	projectQ := r.Q.Project
	costCenterQ := r.Q.CostCenter

	err = paymentQ.WithContext(req.Ctx).Select(
		paymentQ.ID, paymentQ.Code, paymentQ.PostingDate, paymentQ.CreatedAt, paymentQ.Status,
		paymentQ.Amount, paymentQ.PaymentType, paymentQ.PartyID,
		partyQ.PartyTypeCode.As("party_type"),
		supplierQ.Name.As("party_name"), supplierQ.UUID.As("party_uuid"),
		customerQ.Name.As("party_name"), customerQ.UUID.As("party_uuid"),
		paymentQ.AccountPaidFromID.As("paid_from_id"), ledgerFromQ.Name.As("paid_from_name"),
		ledgerFromQ.UUID.As("paid_from_uuid"), ledgerAFromQ.Currency.As("paid_from_currency"),
		paymentQ.AccountPaidToID.As("paid_to_id"), ledgerToQ.Name.As("paid_to_name"),
		ledgerToQ.UUID.As("paid_to_uuid"), ledgerAToQ.Currency.As("paid_to_currency"),
		paymentQ.CompanyBankAccountID, cBAQ.AccountName.As("company_bank_account"), cBAQ.UUID.As("company_bank_account_uuid"),
		paymentQ.PartyBankAccountID, pBAQ.AccountName.As("party_bank_account"), pBAQ.UUID.As("party_bank_account_uuid"),
		paymentQ.ChequeReferenceDate,paymentQ.ChequeReferenceNo,
		paymentQ.ProjectID,projectQ.Name.As("project"),projectQ.UUID.As("project_uuid"),
		paymentQ.CostCenterID,costCenterQ.Name.As("cost_center"),costCenterQ.UUID.As("cost_center_uuid"),
	).
		Join(partyQ, partyQ.ID.EqCol(paymentQ.PartyID)).
		// Join(transactionLedgerQ, transactionLedgerQ.VoucherCode.EqCol(paymentQ.Code)).
		LeftJoin(cBAQ, cBAQ.ID.EqCol(paymentQ.CompanyBankAccountID)).
		LeftJoin(pBAQ, pBAQ.ID.EqCol(paymentQ.PartyBankAccountID)).
		LeftJoin(ledgerFromQ, ledgerFromQ.ID.EqCol(paymentQ.AccountPaidFromID)).
		LeftJoin(ledgerToQ, ledgerToQ.ID.EqCol(paymentQ.AccountPaidToID)).
		LeftJoin(ledgerAFromQ, ledgerAFromQ.LedgerID.EqCol(paymentQ.AccountPaidFromID)).
		LeftJoin(ledgerAToQ, ledgerAToQ.LedgerID.EqCol(paymentQ.AccountPaidToID)).
		LeftJoin(projectQ, projectQ.ID.EqCol(paymentQ.ProjectID)).
		LeftJoin(costCenterQ, costCenterQ.ID.EqCol(paymentQ.CostCenterID)).
		LeftJoin(supplierQ, partyQ.PartyTypeCode.Eq(proto.PartyType_supplier.String()), supplierQ.ID.EqCol(partyQ.ID)).
		LeftJoin(customerQ, partyQ.PartyTypeCode.Eq(proto.PartyType_customer.String()), customerQ.ID.EqCol(partyQ.ID)).
		// .Join(joinTable, joinExprs...)
		Where(
			paymentQ.CompanyID.Eq(req.ActiveCompany.ID), paymentQ.Code.Eq(i.ID),
		).Scan(&res.Entity)

	if err != nil {
		return
	}

	res.Entity.PaymentReferences, err = r.getPaymentReferences(req.Ctx, res.Entity.ID)
	return res, err
}

func (r *paymentRepository) getPaymentReferences(ctx context.Context, paymentID int64) (
	res []dto.PaymentReferenceDto, err error) {
	partyQ := r.Q.Party
	referenceQ := r.Q.PaymentReference
	err = referenceQ.WithContext(ctx).Select(
		referenceQ.PartyID, referenceQ.PartyCode, partyQ.PartyTypeCode.As("party_type"),
		referenceQ.Total, referenceQ.Outstanding, referenceQ.Allocated,
		referenceQ.Currency,
	).Join(partyQ, partyQ.ID.EqCol(referenceQ.PartyID)).
		Where(referenceQ.PaymentID.Eq(paymentID)).
		Scan(&res)
	return
}

func (r *paymentRepository) GetPayments(req *common.RequestContext, d *dto.RequestPayments) (
	res dto.PaginationResult[[]dto.PaymentDto], err error) {
	var (
		generateSQL strings.Builder
	)
	builder := r.Q.WithContext(req.Ctx).Invoice
	queryData := r.convertor.GenerateQueryMap(d)
	params := r.paymentsQuery(req, queryData, &generateSQL)
	err = builder.UnderlyingDB().Raw(generateSQL.String(), params...).Scan(&res.Results).Error
	return
}

func (r *paymentRepository) paymentsQuery(req *common.RequestContext, d map[string]string,
	generateSQL *strings.Builder) (params []interface{}) {
	var (
		whereSQL strings.Builder
	)
	generateSQL.WriteString(`SELECT 
		e.id,e.code,e.created_at,e.status,e.posting_date,e.payment_type,
		e.amount,
		p.party_type_code as party_type,coalesce(s.name,c.name) as party_name,
		coalesce(s.uuid,c.uuid) as party_uuid
		from payments as e
		join parties as p on p.id = e.party_id
		left join customers as c on c.id = e.party_id 
		left join suppliers as s on s.id = e.party_id 
	`)
	whereSQL.WriteString(` e.deleted_at is null and e.company_id = ? `)
	params = append(params, req.ActiveCompany.ID)

	columnFilters := []string{
		"id",
		"project_id",
		"cost_center_id",
		"status",
		"code",
		"party_id",
		"posting_date",
		"payment_type",
		"amount",
		"account_paid_to_id",
		"account_paid_from_id",
	}
	r.query.FilterBuilder(&whereSQL, &params, d, columnFilters...)

	r.query.ReferenceFilterBuilder(generateSQL, &whereSQL, &params, d, "invoice_id")

	helper.JoinWhereBuilder(generateSQL, whereSQL)

	r.query.OrderAndLimitBuilder(generateSQL, d)
	return
}

func (r *paymentRepository) EditPayment(req *common.RequestContext, d dto.PaymentBody) (
	err error) {
	tx := r.Q.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	var references []*int64

	data, err := r.convertor.DataMap(d.PaymentData.Fields)
	if err != nil {
		return
	}
	err = tx.Payment.UnderlyingDB().WithContext(req.Ctx).Model(
		&model.Payment{ID: d.PaymentData.ID},
	).Updates(data).Error
	if err != nil {
		return
	}

	err = tx.Payment.InsertActivity(d.PaymentData.ID, req.Profile.ID, proto.ActivityType_EDIT.String(), nil)
	if err != nil {
		return
	}

	_, err = tx.PaymentReference.WithContext(req.Ctx).Unscoped().Where(
		tx.PaymentReference.PaymentID.Eq(d.PaymentData.ID),
	).Delete()
	if err != nil {
		return
	}

	err = r.createPaymentReferences(req.Ctx, tx, d.PaymentData.ID, d.PaymentReferences, &references)
	if err != nil {
		return
	}
	references = append(references,
		&d.PaymentData.Fields.PartyID,
		d.PaymentData.Fields.CompanyBankAccountID,
		d.PaymentData.Fields.PartyBankAccountID,
		d.PaymentData.Fields.ProjectID,
		d.PaymentData.Fields.CostCenterID,
	)
	err = r.dbHelper.InsertReferences(req.Ctx, tx, d.PaymentData.ID, references, true)
	if err != nil {
		return
	}
	err = tx.Commit()
	return
}

func (r *paymentRepository) CreatePayment(req *common.RequestContext, tx *query.QueryTx, d *dto.CreatePaymentRequest) (
	payment model.Payment, err error) {
	var references []*int64
	paymentPartyID, err := tx.Payment.InsertParty(proto.PartyType_payment.String())
	if err != nil {
		return
	}
	code := r.dbHelper.GenerateCode(tx.Payment.UnderlyingDB(), model.Payment{}, req.ActiveCompany.ID)
	payment.ID = paymentPartyID
	payment.Code = code
	payment.CompanyID = req.ActiveCompany.ID
	payment.Status = proto.State_DRAFT.String()

	fields := d.Body.PaymentData.Fields
	if err = r.convertor.CopyStructData(fields, &payment); err != nil {
		return
	}

	err = tx.Payment.Save(&payment)
	if err != nil {
		return
	}
	err = r.createPaymentReferences(req.Ctx, tx, payment.ID, d.Body.PaymentReferences, &references)
	if err != nil {
		return
	}
	references = append(references,
		fields.ProjectID,
		fields.CostCenterID,
		fields.CompanyBankAccountID,
		fields.PartyBankAccountID,
		&fields.PartyID,
	)
	err = r.dbHelper.InsertReferences(req.Ctx, tx, payment.ID, references)
	return
}

func (r *paymentRepository) createPaymentReferences(ctx context.Context, tx *query.QueryTx,
	paymentID int64, d []dto.CreatePaymentReference, references *[]*int64) (err error) {

	paymentReferences := make([]*model.PaymentReference, len(d))

	// Loop through each reference and map it to the model.PaymentReference
	for i, reference := range d {
		paymentReference := &model.PaymentReference{
			Allocated:   r.currency.FloatToInt(reference.Allocated),
			Total:       r.currency.FloatToInt(reference.Total),
			Outstanding: r.currency.FloatToInt(reference.Outstanding),
			PartyCode:   reference.PartyCode,
			PartyID:     reference.PartyID,
			Currency:    reference.Currency,
			PaymentID:   paymentID,
		}
		paymentReferences[i] = paymentReference

		// Append PartyID to the references slice
		*references = append(*references, &reference.PartyID)
	}

	// Create the payment references in batches
	err = tx.PaymentReference.WithContext(ctx).CreateInBatches(paymentReferences, len(paymentReferences))
	if err != nil {
		return err
	}

	return nil
}

func (r *paymentRepository) GetAllowedParties(req *common.RequestContext) []dto.PartyTypeDto {
	paymentParties := []dto.PartyTypeDto{
		{Name: r.locale.MustLocalize(
			helpers.OptionsLocale.WithID("Party.customer"),
			helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
		), Code: domain.PARTY_CUSTOMER},
		{Name: r.locale.MustLocalize(
			helpers.OptionsLocale.WithID("Party.supplier"),
			helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
		), Code: domain.PARTY_SUPPLIER},
	}
	return paymentParties
}

func (r *paymentRepository) GetFilterOptions(lng string) []dto.FilterOptionDto {
	filterOptions := []dto.FilterOptionDto{}
	t := r.locale.Translate(lng)
	status := dto.FilterOptionDto{
		Param:     "status",
		Name:      "Estado",
		Type:      dto.FILTER_TYPE_STRING,
		Operators: dto.StringOperators,
		Options: []string{proto.State_DRAFT.String(), proto.State_SUBMITTED.String(),
			proto.State_CANCELLED.String()},
	}
	paymentType := dto.FilterOptionDto{
		Param:     "payment_type",
		Name:      "Tipo de Pago",
		Type:      dto.FILTER_TYPE_STRING,
		Operators: dto.StringOperators,
		Options:   []string{proto.PaymentType_PAY.String(), proto.PaymentType_RECEIVE.String()},
	}

	postingDate := dto.FilterOptionDto{
		Name:      "Fecha de Publicacion",
		Param:     "posting_date",
		Type:      dto.FILTER_TYPE_DATE,
		Operators: dto.DateOperators,
	}

	amount := dto.FilterOptionDto{
		Name:      "Monto",
		Param:     "amount",
		Type:      dto.FILTER_TYPE_NUMBER,
		Operators: dto.NumberOperators,
	}

	code := dto.FilterOptionDto{
		Name:      "ID",
		Param:     "code",
		Type:      dto.FILTER_TYPE_STRING,
		Operators: dto.StringOperators,
	}

	purchaseInvoice := dto.FilterOptionDto{
		Name:      t("Entity." + domain.PURCHASE_INVOICE.Name),
		Param:     "pi_id",
		Type:      dto.FILTER_TYPE_STRING,
		Operators: dto.StringOperators,
		PartyType: proto.PartyType_purchaseInvoice.String(),
	}

	saleInvoice := dto.FilterOptionDto{
		Name:      t("Entity." + domain.SALE_INVOICE.Name),
		Param:     "si_id",
		Type:      dto.FILTER_TYPE_STRING,
		Operators: dto.StringOperators,
		PartyType: proto.PartyType_purchaseInvoice.String(),
	}

	accountPaidFrom := dto.FilterOptionDto{
		Name:      t("Entity." + domain.LEDGER.Name),
		Param:     "account_paid_from_id",
		Type:      dto.FILTER_TYPE_STRING,
		Operators: dto.StringOperators,
		PartyType: proto.PartyType_ledger.String(),
	}
	accountPaidTo := dto.FilterOptionDto{
		Name:      t("Entity." + domain.LEDGER.Name),
		Param:     "account_paid_to_id",
		Type:      dto.FILTER_TYPE_STRING,
		Operators: dto.StringOperators,
		PartyType: proto.PartyType_ledger.String(),
	}

	filterOptions = append(filterOptions, status, paymentType, amount, postingDate, code, purchaseInvoice, saleInvoice,
		accountPaidFrom, accountPaidTo)
	return filterOptions
}
