package bank_repo

import (
	common "erp/api/common"
	dto "erp/api/dto"
	"erp/gen/proto"
	"strings"

	"gorm.io/gen/helper"
)

func (m *repository) Get(req *common.RequestContext, d dto.RequestEntity) (
	res dto.BankDto, err error) {
	qI := m.Q.Bank
	bankID := m.convertor.StrtoInt(d.ID)
	err = qI.WithContext(req.Ctx).Select(
		qI.ID, qI.UUID, qI.Name, qI.Status, qI.CreatedAt,
	).
		Where(
			qI.CompanyID.Eq(req.ActiveCompany.ID),
			qI.ID.Eq(bankID),
		).Scan(&res)
	return
}
func (m *repository) GetList(req *common.RequestContext, d dto.BanksRequest) (
	res []dto.BankDto, err error) {
	var (
		generateSQL strings.Builder
	)
	builder := m.Q.WithContext(req.Ctx).Bank
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
		e.id,e.uuid,e.created_at,e.name,e.status
		from banks as e `)
	whereSQL.WriteString(` e.deleted_at is null and e.company_id = ?`)
	params = append(params, req.ActiveCompany.ID)
	columnFilters := []string{
		"name",
		"created_at",
		"status",
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
