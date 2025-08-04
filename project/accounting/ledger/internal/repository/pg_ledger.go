package ledger_repo

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/gen/proto"
	"erp/internal/app/service/helpers"
	"erp/pkg/db"
	"fmt"
	"strings"

	"gorm.io/gen/helper"
	"gorm.io/gorm"
)

type LedgerRepository interface {
	CreateLedger(req *common.RequestContext, i dto.LedgerData) (dto.LedgerDto, error)
	GetLedgersAccounts(req *common.RequestContext, d dto.LedgersRequest) (
		[]dto.LedgerDto, error)
	GetLedgerAccountsTree(req *common.RequestContext) (
		[]dto.TreeEntryDto, error)
	GetLedgerDetail(req *common.RequestContext, i *dto.RequestEntity) (
		dto.ResultEntity[dto.LedgerDetailDto], error)
	GetGeneralLedgerReport(req *common.RequestContext, i *dto.RequestGeneralLedger) (
		[]dto.GeneralLedgerEntryDto, error)
	EditLedger(req *common.RequestContext, d dto.LedgerData) (err error)
}

type ledgerRepository struct {
	convertor helpers.ConvertorHelper
	query     helpers.QueryHelper
	Q         *query.Query
	DB        *gorm.DB
	dbHelper  db.DbHelper
}

func NewLedgerRepository(
	conn db.Connection,
	helpers *helpers.Helpers,
) LedgerRepository {
	return &ledgerRepository{
		convertor: helpers.Convertor,
		DB:        conn.GetDB(),
		Q:         conn.GetQ(),
		dbHelper:  conn.GetDbHelper(),
		query:     helpers.Query,
	}
}

func (r *ledgerRepository) GetLedgerAccountsTree(req *common.RequestContext) (
	res []dto.TreeEntryDto, err error) {
	query := `WITH RECURSIVE data_cte AS (
		SELECT 
			id, 
			uuid,
			ledger_parent as parent_id,
			is_group,
			name
		FROM ledgers
		WHERE ledger_parent IS NULL and company_id = ?
		UNION ALL 
		SELECT 
			l.id, 
			l.uuid,
			l.ledger_parent as parent_id,
			l.is_group,
			l.name
		FROM ledgers l
		INNER JOIN data_cte d
			ON l.ledger_parent = d.id  
	)
	SELECT * FROM data_cte;`
	err = r.DB.WithContext(req.Ctx).Raw(query, req.ActiveCompany.ID).Scan(&res).Error
	return
}

func (r *ledgerRepository) GetGeneralLedgerReport(req *common.RequestContext, i *dto.RequestGeneralLedger) (
	res []dto.GeneralLedgerEntryDto, err error) {
	var (
		optConds   string
		partyConds string
	)
	if i.VoucherNo != "" {
		optConds += fmt.Sprintf("and tx.voucher_code = '%s'", i.VoucherNo)
	}

	joinCustomer := fmt.Sprintf(`LEFT JOIN 
			customers AS cust 
			ON (cust.id = p.id AND p.party_type_code = '%s')`, proto.PartyType_customer.String())
	joinSupplier := fmt.Sprintf(`LEFT JOIN 
			suppliers AS suppl 
			ON (suppl.id = p.id AND p.party_type_code = '%s')`, proto.PartyType_supplier.String())
	if i.PartyType != "" || i.Party != "" {
		partyConds = fmt.Sprintf(`LEFT JOIN 
			parties AS p ON p.id = %d `, r.convertor.StrtoInt(i.Party))
	} else {
		partyConds = fmt.Sprintf(`LEFT JOIN parties AS p ON p.id = tx.party_id`)
	}
	partyConds = fmt.Sprintf(`%s %s %s`, partyConds, joinCustomer, joinSupplier)

	query := fmt.Sprintf(`
		SELECT 
			tx.created_at AS posting_date,
			acc1.name AS account,
			tx.debit,
			tx.credit,
			tx.balance,
			acc2.name AS against_account,
			tx.voucher_type,
			tx.voucher_subtype,
			tx.voucher_code as voucher_no,
			COALESCE(p.party_type_code, '') AS party_type,
			COALESCE(suppl.name, cust.name) AS party_name,
			tx.currency
		FROM 
			transaction_ledgers AS tx
		JOIN 
			ledgers AS acc1 ON acc1.id = tx.ledger and acc1.company_id = $3
		JOIN 
			ledgers AS acc2 ON acc2.id = tx.ledger_against
		%s
		where tx.created_at::date between $1 and $2 %s
		order by tx.created_at`, partyConds, optConds)

	// fmt.Println("QUERY",query)

	err = r.DB.WithContext(req.Ctx).Raw(query, i.FromDate, i.ToDate, req.ActiveCompany.ID).Scan(&res).Error

	return
}

func (r *ledgerRepository) GetLedgerDetail(req *common.RequestContext, i *dto.RequestEntity) (
	res dto.ResultEntity[dto.LedgerDetailDto], err error) {
	id := r.convertor.StrtoInt(i.ID)
	ledgerQ := r.Q.Ledger
	ledgerParentQ := r.Q.Ledger.As("parent")
	err = r.Q.Ledger.WithContext(req.Ctx).Select(
		ledgerQ.ID, ledgerQ.UUID, ledgerQ.LedgerNo, ledgerQ.Name, ledgerQ.IsGroup,
		ledgerQ.CreatedAt, ledgerQ.AccountType, ledgerQ.Status, ledgerQ.AccountRootType,
		ledgerQ.ReportType, ledgerQ.CashFlowSection, ledgerQ.IsOffsetAccount,
		ledgerParentQ.Name.As("parent"), ledgerParentQ.ID.As("parent_id"), ledgerParentQ.UUID.As("parent_uuid"),
	).LeftJoin(ledgerParentQ, ledgerParentQ.ID.EqCol(ledgerQ.LedgerParent)).Where(
		ledgerQ.ID.Eq(id),
		ledgerQ.CompanyID.Eq(req.ActiveCompany.ID),
	).Scan(&res.Entity)

	return res, err
}

func (r *ledgerRepository) GetLedgersAccounts(req *common.RequestContext, d dto.LedgersRequest) (
	res []dto.LedgerDto, err error) {
	var (
		generateSQL strings.Builder
	)
	builder := r.Q.WithContext(req.Ctx).Ledger
	queryData := r.convertor.GenerateQueryMap(d)
	// fmt.Println("LEDGER REQUEST", d)
	params := r.ledgersQuery(req, queryData, &generateSQL)
	// fmt.Println("LEDGER QUERY", generateSQL.String(),params)
	err = builder.UnderlyingDB().Raw(generateSQL.String(), params...).Scan(&res).Error
	return
}

func (r *ledgerRepository) ledgersQuery(req *common.RequestContext, d map[string]string,
	generateSQL *strings.Builder) (params []interface{}) {
	var (
		whereSQL strings.Builder
	)
	generateSQL.WriteString(`SELECT 
		e.id,e.uuid,e.ledger_no,e.created_at,e.name,e.is_group,e.updated_at,
		e.account_type,e.status,
		la.currency,la.can_credit,la.can_debit
		from ledgers as e 
		left join ledger_accounts as la on la.ledger_id = e.id
		`)
	whereSQL.WriteString(` e.deleted_at is null and e.company_id = ?`)
	params = append(params, req.ActiveCompany.ID)
	columnFilters := map[string][]string{
		"e": {
			"name",
			"created_at",
			"status",
			"is_group",
			"account_type",
		},
	}
	r.query.FilterMapBuilder(&whereSQL, &params, d, columnFilters)
	helper.JoinWhereBuilder(generateSQL, whereSQL)
	r.query.OrderAndLimitBuilder(generateSQL, d)
	return
}

func (r *ledgerRepository) EditLedger(req *common.RequestContext, d dto.LedgerData) (
	err error) {
	tx := r.Q.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	data, err := r.convertor.DataMap(d.Fields)
	if err != nil {
		return
	}
	fmt.Println("LEDGER",d.ID)
	err = tx.Ledger.UnderlyingDB().WithContext(req.Ctx).Model(
		&model.Ledger{ID: d.ID},
	).Updates(data).Error
	if err != nil {
		return
	}

	err = tx.Ledger.InsertActivity(d.ID, req.Profile.ID, proto.ActivityType_EDIT.String(), nil)
	if err != nil {
		return
	}

	err = tx.Commit()
	return
}

func (r *ledgerRepository) CreateLedger(req *common.RequestContext, i dto.LedgerData) (
	res dto.LedgerDto, err error) {
	tx := r.Q.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	ledger, err := r.createLedger(req, tx, i)
	if err != nil {
		return res, err
	}

	//Only create ledger accounts for ledgers that are not marked as group accounts.
	if !ledger.IsGroup {
		err = r.createLedgerAccount(req, tx, ledger, i)
		if err != nil {
			return
		}
	}
	err = tx.Commit()
	if err != nil {
		return
	}
	res = dto.LedgerDtoFromModel(&ledger)
	return
}

func (r *ledgerRepository) createLedger(req *common.RequestContext, tx *query.QueryTx, i dto.LedgerData) (model.Ledger, error) {
	var (
		ledger model.Ledger
		err    error
	)
	partyId, err := tx.Ledger.InsertParty(proto.PartyType_ledger.String())
	if err != nil {
		return ledger, err
	}
	ledger.ID = partyId
	ledger.CompanyID = req.ActiveCompany.ID
	ledger.Status = proto.State_ENABLED.String()
	fields := i.Fields
	if err = r.convertor.CopyStructData(fields, &ledger); err != nil {
		return ledger, err
	}
	// ledger.Description = i.Description
	err = tx.Ledger.WithContext(req.Ctx).Save(&ledger)
	return ledger, err
}

func (r *ledgerRepository) createLedgerAccount(req *common.RequestContext, tx *query.QueryTx,
	ledger model.Ledger, i dto.LedgerData) error {
	ledgerAccount := model.LedgerAccount{}
	ledgerAccount.LedgerID = ledger.ID
	ledgerAccount.Currency = req.CompanyDefaults.Currency

	err := tx.LedgerAccount.WithContext(req.Ctx).Save(&ledgerAccount)
	return err
}
