package acct_report_repo

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/proto"
	"erp/internal/app/service/helpers"
	"erp/pkg/db"
	"strings"

	"gorm.io/gen/helper"
	"gorm.io/gorm"
)

type AcctReportRepository interface {
	GetGeneralLedgerReport(req *common.RequestContext, i *dto.RequestGeneralLedger) (
		dto.GeneralLedgerData, error)
	GetAccountPayable(req *common.RequestContext, i *dto.RequestAccountPayable) (
		[]dto.AccountPayableEntryDto, error)
	GetAccountPayableSumary(req *common.RequestContext, i *dto.RequestAccountPayable) (
		[]dto.SumaryEntryDto, error)
	GetAccountReceivableSumary(req *common.RequestContext, i *dto.RequestAccountReceivable) (
		res []dto.SumaryEntryDto, err error)
	GetAccountReceivable(req *common.RequestContext, i *dto.RequestAccountReceivable) (
		res []dto.AccountReceivableEntryDto, err error)

	GetAccountBalance(req *common.RequestContext, d *dto.RequestAccountBalance) (res dto.GeneralLedgerOpening, err error)
}

type acctReportRepository struct {
	DB        *gorm.DB
	convertor helpers.ConvertorHelper
}

func NewAcctReportRepository(
	conn db.Connection,
	helpers *helpers.Helpers,
) AcctReportRepository {
	return &acctReportRepository{
		DB:        conn.GetDB(),
		convertor: helpers.Convertor,
	}
}

func (r *acctReportRepository) GetAccountReceivableSumary(req *common.RequestContext, i *dto.RequestAccountReceivable) (
	res []dto.SumaryEntryDto, err error) {
	var (
		params      []interface{}
		generateSQL strings.Builder
		whereSQL    strings.Builder
	)

	generateSQL.WriteString(`
		SELECT 
		tx.party_id, 
		'customer' AS party_type,
		cust.name AS party_name,
		cust.uuid AS party_uuid,
		SUM(tx.debit)  AS total_invoiced_amount,
		SUM(tx.credit)  AS total_paid_amount,
		tx.currency
		FROM transaction_ledgers AS tx
		JOIN customers AS cust 
		ON cust.id = tx.party_id 
		JOIN ledgers AS l 
			ON l.id = tx.ledger
		`)
	whereSQL.WriteString(`tx.deleted_at is null and l.company_id = ?`)
	params = append(params, req.ActiveCompany.ID)

	if i.Party != "" {
		// parties := r.convertor.StrToArrayInt64(i.Party)
		// partyCond = fmt.Sprintf("and tx.party_id = ANY (array[%s])", r.convertor.ArrayToString(parties))
		whereSQL.WriteString(` and tx.party_id = ?`)
		params = append(params, r.convertor.StrtoInt(i.Party))
	}

	if i.ProjectID != "" {
		params = append(params, i.ProjectID)
		whereSQL.WriteString(` and tx.project_id = ?`)
	}

	if i.CostCenterID != "" {
		params = append(params, i.CostCenterID)
		whereSQL.WriteString(` and tx.cost_center_id = ?`)
	}

	helper.JoinWhereBuilder(&generateSQL, whereSQL)
	generateSQL.WriteString(` GROUP BY 
			tx.party_id, cust.name, cust.uuid,tx.currency`)

	err = r.DB.WithContext(req.Ctx).Raw(generateSQL.String(), params...).Scan(&res).Error
	return
}

func (r *acctReportRepository) GetAccountReceivable(req *common.RequestContext, i *dto.RequestAccountReceivable) (
	res []dto.AccountReceivableEntryDto, err error) {
	var (
		generateSQL strings.Builder
		whereSQL    strings.Builder
		params      []interface{}
	)

	generateSQL.WriteString(`
	SELECT 
		tx.created_at AS posting_date,
		? AS party_type,
		cust.name AS party_name, 
		cust.uuid AS party_uuid,
		l.name AS receivable_account,
		l.uuid AS receivable_account_uuid,
		tx.voucher_type,
		tx.voucher_code AS voucher_no,
		tx.debit invoiced_amount,
		tx.credit paid_amount,
		tx.currency
		FROM transaction_ledgers AS tx
		JOIN customers AS cust ON cust.id = tx.party_id
		JOIN ledgers AS l ON l.id = tx.ledger
	`)
	params = append(params, proto.PartyType_customer.String())
	whereSQL.WriteString(`tx.deleted_at is null and tx.created_at::date between ? and ? and l.company_id = ?`)
	params = append(params, i.FromDate, i.ToDate, req.ActiveCompany.ID)
	if i.Party != "" {
		// parties := r.convertor.StrToArrayInt64(i.Party)
		// partyCond = fmt.Sprintf("and tx.party_id = ANY (array[%s])", r.convertor.ArrayToString(parties))
		whereSQL.WriteString(` and tx.party_id = ?`)
		params = append(params, r.convertor.StrtoInt(i.Party))
	}

	if i.ProjectID != "" {
		params = append(params, i.ProjectID)
		whereSQL.WriteString(` and tx.project_id = ?`)
	}

	if i.CostCenterID != "" {
		params = append(params, i.CostCenterID)
		whereSQL.WriteString(` and tx.cost_center_id = ?`)
	}

	helper.JoinWhereBuilder(&generateSQL, whereSQL)
	generateSQL.WriteString(` order by tx.id`)
	// ;`, proto.PartyType_customer.String(), proto.PartyType_saleInvoice.String(),
	// proto.PartyType_saleInvoice.String(), partyCond,
	// )
	err = r.DB.WithContext(req.Ctx).Raw(generateSQL.String(), params...).Scan(&res).Error
	return
}

func (r *acctReportRepository) GetAccountPayableSumary(req *common.RequestContext, i *dto.RequestAccountPayable) (
	res []dto.SumaryEntryDto, err error) {
	var (
		params      []interface{}
		generateSQL strings.Builder
		whereSQL    strings.Builder
	)

	generateSQL.WriteString(`
	SELECT 
    tx.party_id, 
    'supplier' AS party_type,
    suppl.name AS party_name,
    suppl.uuid AS party_uuid,
    SUM(tx.credit) AS total_invoiced_amount,
    SUM(tx.debit)  AS total_paid_amount,
    tx.currency
	FROM transaction_ledgers AS tx
	JOIN suppliers AS suppl 
		ON suppl.id = tx.party_id 
	JOIN ledgers AS l 
		ON l.id = tx.ledger
	`)
	whereSQL.WriteString(`tx.deleted_at is null and l.company_id = ?`)
	params = append(params, req.ActiveCompany.ID)
	if i.Party != "" {
		// parties := r.convertor.StrToArrayInt64(i.Party)
		// partyCond = fmt.Sprintf("and tx.party_id = ANY (array[%s])", r.convertor.ArrayToString(parties))
		whereSQL.WriteString(` and tx.party_id = ?`)
		params = append(params, r.convertor.StrtoInt(i.Party))
	}
	if i.ProjectID != "" {
		params = append(params, i.ProjectID)
		whereSQL.WriteString(` and tx.project_id = ?`)
	}
	if i.CostCenterID != "" {
		params = append(params, i.CostCenterID)
		whereSQL.WriteString(` and tx.cost_center_id = ?`)
	}

	helper.JoinWhereBuilder(&generateSQL, whereSQL)
	generateSQL.WriteString(` GROUP BY 
		tx.party_id, suppl.name, suppl.uuid,tx.currency`)

	err = r.DB.WithContext(req.Ctx).Raw(generateSQL.String(), params...).Scan(&res).Error
	return
}

func (r *acctReportRepository) GetAccountPayable(req *common.RequestContext, i *dto.RequestAccountPayable) (
	res []dto.AccountPayableEntryDto, err error) {
	var (
		generateSQL strings.Builder
		whereSQL    strings.Builder
		params      []interface{}
	)

	generateSQL.WriteString(`
		SELECT 
		tx.created_at AS posting_date,
		? AS party_type,
		suppl.name AS party_name, 
		suppl.uuid AS party_uuid,
		l.name AS receivable_account,
		l.uuid AS receivable_account_uuid,
		tx.voucher_type,
		tx.voucher_code AS voucher_no,
		tx.credit invoiced_amount,
		tx.debit paid_amount,
		tx.currency
		FROM transaction_ledgers AS tx
		JOIN suppliers AS suppl ON suppl.id = tx.party_id
		JOIN ledgers AS l ON l.id = tx.ledger
		`)
	params = append(params, proto.PartyType_customer.String())
	whereSQL.WriteString(`tx.deleted_at is null and tx.created_at::date between ? and ? and l.company_id = ?`)
	params = append(params, i.FromDate, i.ToDate, req.ActiveCompany.ID)
	if i.Party != "" {
		// parties := r.convertor.StrToArrayInt64(i.Party)
		// partyCond = fmt.Sprintf("and tx.party_id = ANY (array[%s])", r.convertor.ArrayToString(parties))
		whereSQL.WriteString(` and tx.party_id = ?`)
		params = append(params, r.convertor.StrtoInt(i.Party))
	}
	if i.ProjectID != "" {
		params = append(params, i.ProjectID)
		whereSQL.WriteString(` and tx.project_id = ?`)
	}
	if i.CostCenterID != "" {
		params = append(params, i.CostCenterID)
		whereSQL.WriteString(` and tx.cost_center_id = ?`)
	}
	helper.JoinWhereBuilder(&generateSQL, whereSQL)
	generateSQL.WriteString(` order by tx.id`)
	// ;`, proto.PartyType_customer.String(), proto.PartyType_saleInvoice.String(),
	// proto.PartyType_saleInvoice.String(), partyCond,
	// )
	err = r.DB.WithContext(req.Ctx).Raw(generateSQL.String(), params...).Scan(&res).Error
	return
}

func (r *acctReportRepository) GetGeneralLedgerReport(req *common.RequestContext, i *dto.RequestGeneralLedger) (
	res dto.GeneralLedgerData, err error) {
	accountID := r.convertor.StrtoInt(i.Account)
	partyID := r.convertor.StrtoInt(i.Party)
	var params []interface{}
	var cteSQL strings.Builder
	var whereSQL strings.Builder
	if accountID != 0 {
		cteSQL.WriteString(`
			SELECT 
			SUM(tx.debit) AS debit,
			SUM(tx.credit) AS credit,
			(SUM(tx.debit) - SUM(tx.credit) )AS opening_balance
			FROM 
			transaction_ledgers AS tx
			`)
		params = append(params, accountID, i.FromDate)
		whereSQL.WriteString("tx.ledger = ? AND tx.created_at < ? ")
		if partyID != 0 {
			whereSQL.WriteString("tx.party_id = ?")
			params = append(params, partyID)
		}
		helper.JoinWhereBuilder(&cteSQL, whereSQL)
		err = r.DB.Raw(cteSQL.String(), params...).Scan(&res.Opening).Error
		if err != nil {
			return
		}
		params = []interface{}{}
	}
	whereSQL.Reset()
	var generateSQL strings.Builder
	generateSQL.WriteString(` SELECT 
			tx.created_at AS posting_date,
			acc1.name AS account,
			tx.debit,
			tx.credit,
			COALESCE(acc2.name,'') AS against_account,
			tx.voucher_type,
			tx.voucher_subtype,
			tx.voucher_code as voucher_no,
			COALESCE(p.party_type_code, '') AS party_type,
			COALESCE(suppl.name, cust.name) AS party_name,
			tx.currency
		FROM 
			transaction_ledgers AS tx
		JOIN 
			ledgers AS acc1 ON acc1.id = tx.ledger and acc1.company_id = ?
		LEFT JOIN 
			ledgers AS acc2 ON acc2.id = tx.ledger_against
		`)
	params = append(params, req.ActiveCompany.ID)
	if i.Party != "" || i.PartyType != "" {
		generateSQL.WriteString(` LEFT JOIN 
			parties AS p ON p.id = ?`)
		params = append(params, partyID)
	} else {
		generateSQL.WriteString(` LEFT JOIN parties AS p ON p.id = tx.party_id `)
	}
	generateSQL.WriteString(` LEFT JOIN 
			customers AS cust 
			ON (cust.id = p.id AND p.party_type_code = 'customer')
		LEFT JOIN 
			suppliers AS suppl 
			ON (suppl.id = p.id AND p.party_type_code = 'supplier')`)
	whereSQL.WriteString(` tx.deleted_at is null and tx.created_at::date between ? and ?`)
	params = append(params, i.FromDate, i.ToDate)
	if i.VoucherNo != "" {
		whereSQL.WriteString(` and tx.voucher_code = ?`)
		params = append(params, i.VoucherNo)
	}
	if accountID != 0 {
		whereSQL.WriteString(` and tx.ledger = ?`)
		params = append(params, accountID)
	}
	if i.Project != "" {
		whereSQL.WriteString(` and tx.project_id = ?`)
		params = append(params, r.convertor.StrtoInt(i.Project))
	}
	if i.CostCenter != "" {
		whereSQL.WriteString(` and tx.cost_center_id = ?`)
		params = append(params, r.convertor.StrtoInt(i.CostCenter))
	}
	helper.JoinWhereBuilder(&generateSQL, whereSQL)
	generateSQL.WriteString(" order by tx.id asc")
	err = r.DB.Raw(generateSQL.String(), params...).Scan(&res.Entries).Error
	return
}

func (r *acctReportRepository) GetAccountBalance(req *common.RequestContext, d *dto.RequestAccountBalance) (res dto.GeneralLedgerOpening, err error) {
	var cteSQL strings.Builder
	var whereSQL strings.Builder
	id := r.convertor.StrtoInt(d.ID)
	var params []interface{}
	if id != 0 {
		cteSQL.WriteString(`
			SELECT 
			SUM(tx.debit) AS debit,
			SUM(tx.credit) AS credit,
			(SUM(tx.debit) - SUM(tx.credit) )AS opening_balance
			FROM 
			transaction_ledgers AS tx
			`)
		params = append(params, id)
		whereSQL.WriteString("tx.ledger = ?")

		helper.JoinWhereBuilder(&cteSQL, whereSQL)
		err = r.DB.Raw(cteSQL.String(), params...).Scan(&res).Error
		if err != nil {
			return
		}
	}
	return
}
