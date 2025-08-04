package bank_account_repo

import (
	common "erp/api/common"
	dto "erp/api/dto"
	"erp/gen/proto"
	"strings"

	"gorm.io/gen/helper"
)

func (m *repository) Get(req *common.RequestContext, d dto.RequestEntity) (
	res dto.BankAccountDto, err error) {
	id := m.convertor.StrtoInt(d.ID)
	qI := m.Q.BankAccount
	bankQ := m.Q.Bank
	ledgerQ := m.Q.Ledger

	partyQ := m.Q.Party
	//parties assoiated with payments
	supplierQ := m.Q.Supplier
	customerQ := m.Q.Customer
	err = qI.WithContext(req.Ctx).Select(
		qI.ID, qI.UUID, qI.AccountName, qI.Status, qI.CreatedAt,
		qI.BankID, bankQ.Name.As("bank"), bankQ.UUID.As("bank_uuid"),
		qI.PartyID, partyQ.PartyTypeCode.As("party_type"),
		supplierQ.Name.As("party"), supplierQ.UUID.As("party_uuid"),
		customerQ.Name.As("party"), customerQ.UUID.As("party_uuid"),
		qI.BankAccountType, qI.BankAccountNumber, qI.Iban, qI.BranchCode,
		qI.IsCompanyAccount,
		qI.CompanyAccountID, ledgerQ.Name.As("company_account_"), ledgerQ.UUID.As("company_account_uuid"),
	).
		Join(bankQ, bankQ.ID.EqCol(qI.BankID)).
		LeftJoin(partyQ, partyQ.ID.EqCol(qI.PartyID)).
		LeftJoin(ledgerQ, ledgerQ.ID.EqCol(qI.CompanyAccountID)).
		LeftJoin(supplierQ, partyQ.PartyTypeCode.Eq(proto.PartyType_supplier.String()), supplierQ.ID.EqCol(partyQ.ID)).
		LeftJoin(customerQ, partyQ.PartyTypeCode.Eq(proto.PartyType_customer.String()), customerQ.ID.EqCol(partyQ.ID)).
		Where(
			qI.CompanyID.Eq(req.ActiveCompany.ID),
			qI.ID.Eq(id),
		).Scan(&res)
	return
}
func (m *repository) GetList(req *common.RequestContext, d dto.BankAccountsRequest) (
	res []dto.BankAccountDto, err error) {
	var (
		generateSQL strings.Builder
	)
	builder := m.Q.WithContext(req.Ctx).BankAccount
	queryData := m.convertor.GenerateQueryMap(d)
	params := m.bankListQuery(req, queryData, &generateSQL)
	err = builder.UnderlyingDB().Raw(generateSQL.String(), params...).Scan(&res).Error
	return
}

func (m *repository) bankListQuery(req *common.RequestContext, d map[string]string,
	generateSQL *strings.Builder) (params []interface{}) {
	var (
		whereSQL strings.Builder
	)
	generateSQL.WriteString(`SELECT 
		e.id,e.uuid,e.created_at,e.account_name,e.status,e.bank_account_type,
		e.bank_account_number,
		company_account_id,lg.name as company_account,lg.uuid as company_account_uuid,
		lga.currency as company_account_currency
		from bank_accounts as e 
		left join ledgers as lg on lg.id = e.company_account_id 
		left join ledger_accounts as lga on lga.ledger_id = e.company_account_id  `)
	whereSQL.WriteString(` e.deleted_at is null and e.company_id = ?`)
	params = append(params, req.ActiveCompany.ID)
	columnFilters := []string{
		"name",
		"created_at",
		"status",
		"party_id",
		"party_id",
		"is_company_account",
	}
	m.query.FilterBuilder(&whereSQL, &params, d, columnFilters...)
	helper.JoinWhereBuilder(generateSQL, whereSQL)
	m.query.OrderAndLimitBuilder(generateSQL, d)
	return
}

func (r *repository) GetFilterOptions(lng string) []dto.FilterOptionDto {
	filterOptions := []dto.FilterOptionDto{}
	status := dto.FilterOptionDto{
		Param:     "status",
		Name:      "Estado",
		Type:      dto.FILTER_TYPE_STRING,
		Operators: dto.StringOperators,
		Options:   []string{proto.State_ENABLED.String(), proto.State_DISABLED.String()},
	}
	createdAt := dto.FilterOptionDto{
		Name:      "Fecha de Creacion",
		Param:     "created_at",
		Type:      dto.FILTER_TYPE_DATE,
		Operators: dto.DateOperators,
	}

	name := dto.FilterOptionDto{
		Name:      "Nombre",
		Param:     "name",
		Type:      dto.FILTER_TYPE_STRING,
		Operators: dto.StringOperators,
	}

	filterOptions = append(filterOptions, status, createdAt, name)
	return filterOptions
}
