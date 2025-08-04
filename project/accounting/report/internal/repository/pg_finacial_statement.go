package acct_report_repo

import (
	"context"
	"erp/api/common"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/pkg/db"
	"fmt"
	"strings"

	"gorm.io/gen/helper"
	"gorm.io/gorm"
)

type FinancialStatementRepo interface {
	ProfitAndLossStatement(req *common.RequestContext, d *dto.RequestFinancialStatement) (
		res []dto.ProfitAndLossEntryDto, err error)
	CashFlowStatement(req *common.RequestContext, d *dto.RequestFinancialStatement) (
		res []dto.CashFlowEntryDto, err error)
	BalanceSheetStatement(req *common.RequestContext, d *dto.RequestFinancialStatement) (
		res []dto.BalanceSheetEntryDto, err error)
}

type financialStatementRepo struct {
	DB        *gorm.DB
	convertor helpers.ConvertorHelper
}

func NewAcctFinancialStatementRepo(
	db db.Connection,
	helpers *helpers.Helpers,
) FinancialStatementRepo {
	return &financialStatementRepo{
		DB:        db.GetDB(),
		convertor: helpers.Convertor,
	}
}

func (r *financialStatementRepo) BalanceSheetStatement(req *common.RequestContext, d *dto.RequestFinancialStatement) (
	res []dto.BalanceSheetEntryDto, err error) {
	var (
		generateSQL strings.Builder
		whereSQL    strings.Builder
		params      []interface{}
	)
	generateSQL.WriteString(`
	select l.account_type,l.name as account_name,l.account_root_type,
	SUM(tx.debit) as debit,SUM(tx.credit) as credit
	from transaction_ledgers as tx
	join ledgers as l on l.id = tx.ledger 
	and l.report_type = 'BALANCE_SHEET'
	`)
	params = append(params, d.FromDate, d.ToDate, req.ActiveCompany.ID)
	whereSQL.WriteString(`tx.deleted_at is null and tx.posting_date between ? and ? and l.company_id = ?`)
	if d.Currency != "" {
		params = append(params, d.Currency)
		whereSQL.WriteString(` and tx.currency = ?`)
	}
	if d.ProjectID != "" {
		params = append(params, d.ProjectID)
		whereSQL.WriteString(` and tx.project_id = ?`)
	}
	if d.CostCenterID != "" {
		params = append(params, d.CostCenterID)
		whereSQL.WriteString(` and tx.cost_center_id = ?`)
	}
	helper.JoinWhereBuilder(&generateSQL, whereSQL)
	generateSQL.WriteString(` group by l.account_type,l.account_root_type,l.name`)
	err = r.DB.Raw(generateSQL.String(), params...).Scan(&res).Error
	return
}

func (r *financialStatementRepo) CashFlowStatement(req *common.RequestContext, d *dto.RequestFinancialStatement) (
	res []dto.CashFlowEntryDto, err error) {
	var (
		generateSQL strings.Builder
		whereSQL    strings.Builder
		params      []interface{}
	)
	generateSQL.WriteString(`
	select l.account_type,l.name as account_name,l.cash_flow_section,
	(SUM(tx.credit)-SUM(tx.debit)) as amount from transaction_ledgers as tx
	join ledgers as l on l.id = tx.ledger 
	and l.cash_flow_section = any(array['OPERATING','INVESTING','FINANCING'])
	`)
	params = append(params, d.FromDate, d.ToDate, req.ActiveCompany.ID)
	whereSQL.WriteString(`tx.deleted_at is null and tx.posting_date between ? and ? and l.company_id = ? `)
	if d.Currency != "" {
		params = append(params, d.Currency)
		whereSQL.WriteString(` and  tx.currency = ?`)
	}
	fmt.Println("PROJECT", d.ProjectID)
	if d.ProjectID != "" {
		params = append(params, d.ProjectID)
		whereSQL.WriteString(` and tx.project_id = ?`)
	}
	if d.CostCenterID != "" {
		params = append(params, d.CostCenterID)
		whereSQL.WriteString(` and tx.cost_center_id = ?`)
	}
	helper.JoinWhereBuilder(&generateSQL, whereSQL)
	generateSQL.WriteString(` group by l.account_type,l.cash_flow_section,l.name`)
	err = r.DB.Raw(generateSQL.String(), params...).Scan(&res).Error
	if err != nil {
		return
	}
	profitCasgEntry, err := r.getProfit(req.Ctx, whereSQL, params)
	if err != nil {
		return
	}
	res = append([]dto.CashFlowEntryDto{profitCasgEntry}, res...)
	return
}

func (r *financialStatementRepo) getProfit(ctx context.Context, whereSQL strings.Builder, params []interface{}) (
	res dto.CashFlowEntryDto, err error) {
	var (
		generateSQL strings.Builder
	)
	generateSQL.WriteString(`
	select ('Net Income') as account_type,
	'Income' as account_name,
	'OPERATING' as cash_flow_section,
	coalesce(sum(tx.credit)-sum(tx.debit),0) as amount
		from transaction_ledgers as tx
		join ledgers as l on l.id = tx.ledger and l.report_type = 'PROFIT_AND_LOSS'
	`)
	helper.JoinWhereBuilder(&generateSQL, whereSQL)
	err = r.DB.WithContext(ctx).Raw(generateSQL.String(), params...).Scan(&res).Error
	return
}

func (r *financialStatementRepo) ProfitAndLossStatement(req *common.RequestContext, d *dto.RequestFinancialStatement) (
	res []dto.ProfitAndLossEntryDto, err error) {
	var (
		generateSQL strings.Builder
		filteredSQL strings.Builder
		whereSQL    strings.Builder
		params      []interface{}
	)
	filteredSQL.WriteString(`
	WITH filtered_transactions AS (
	SELECT 
        tx.credit,
        tx.debit,
        l.account_type,
        l.name,
        date_trunc(?, tx.posting_date)::date AS posting_date  
    FROM transaction_ledgers AS tx
	JOIN ledgers AS l ON l.company_id = ? AND tx.ledger = l.id
	`)
	params = append(params, d.TimeUnit, req.ActiveCompany.ID)
	whereSQL.WriteString(`tx.deleted_at is null and l.is_group = false AND report_type = 'PROFIT_AND_LOSS'
	and tx.posting_date between ? and ?`)
	params = append(params, d.FromDate, d.ToDate)
	if d.Currency != "" {
		whereSQL.WriteString(` and currency = ?`)
		params = append(params, d.Currency)
	}
	if d.ProjectID != "" {
		whereSQL.WriteString(` and tx.project_id = ?`)
		params = append(params, r.convertor.StrtoInt(d.ProjectID))
	}
	if d.CostCenterID != "" {
		whereSQL.WriteString(` and tx.cost_center_id = ?`)
		params = append(params, r.convertor.StrtoInt(d.CostCenterID))
	}
	helper.JoinWhereBuilder(&filteredSQL, whereSQL)
	filteredSQL.WriteString(")")
	generateSQL.WriteString(filteredSQL.String())
	generateSQL.WriteString(`
	SELECT account_type,name,posting_date,SUM(credit - debit) AS balance
	FROM 
		filtered_transactions
		GROUP BY account_type,name,posting_date
	`)
	err = r.DB.Raw(generateSQL.String(), params...).Scan(&res).Error
	fmt.Println(res, generateSQL.String())
	return
}
