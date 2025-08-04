package helpers

import (
	"erp/internal/domain"
	"fmt"
	"strconv"
	"strings"
)

type QueryHelper interface {
	OrderAndLimitBuilder(generateSQL *strings.Builder, d map[string]string)
	ReferenceFilterBuilder(generateSQL *strings.Builder,
		whereSQL *strings.Builder, params *[]interface{},
		d map[string]string, args ...string)
	PartyFilterBuilder(generateSQL *strings.Builder,
		whereSQL *strings.Builder, params *[]interface{},
		d map[string]string, args ...string)
	FilterBuilder(whereSQL *strings.Builder, params *[]interface{},
		d map[string]string, args ...string)
	FilterMapBuilder(whereSQL *strings.Builder, params *[]interface{},
		d map[string]string, columns map[string][]string)
}

type queryHelper struct {
	convertor ConvertorHelper
}

func NewQueryHelpers(
	convertor ConvertorHelper,
) QueryHelper {
	return &queryHelper{
		convertor: convertor,
	}
}

func (h *queryHelper) FilterMapBuilder(whereSQL *strings.Builder, params *[]interface{},
	d map[string]string, columnsByAlias map[string][]string) {
	for alias, columns := range columnsByAlias {
		for _, column := range columns {
			if val, ok := d[column]; ok {
				whereSQL.WriteString(h.convertor.GetConditionFromQuery(val, fmt.Sprintf("%s.%s", alias, column), params))
			}
		}
	}
}

func (h *queryHelper) FilterBuilder(whereSQL *strings.Builder, params *[]interface{},
	d map[string]string, args ...string) {
	for _, arg := range args {
		if val, ok := d[arg]; ok {
			whereSQL.WriteString(h.convertor.GetConditionFromQuery(val, fmt.Sprintf("e.%s", arg), params))
		}
	}
}

func (h *queryHelper) ReferenceFilterBuilder(generateSQL *strings.Builder,
	whereSQL *strings.Builder, params *[]interface{},
	d map[string]string, args ...string) {
	for i, arg := range args {
		if value, ok := d[arg]; ok {
			alias := fmt.Sprintf("r%d", i)
			whereSQL.WriteString(h.convertor.GetConditionFromQuery(value, fmt.Sprintf("%s.party_id", alias), params))
			generateSQL.WriteString(fmt.Sprintf(`join party_references as %s on %s.reference_id = e.id `, alias, alias))
		}
	}
}

func (h *queryHelper) PartyFilterBuilder(generateSQL *strings.Builder,
	whereSQL *strings.Builder, params *[]interface{},
	d map[string]string, args ...string) {
	for i, arg := range args {
		if value, ok := d[arg]; ok {
			alias := fmt.Sprintf("p%d", i)
			whereSQL.WriteString(h.convertor.GetConditionFromQuery(value, fmt.Sprintf("%s.reference_id", alias), params))
			generateSQL.WriteString(fmt.Sprintf(`join party_references as %s on %s.party_id = e.id `, alias, alias))
		}
	}
}

func (h *queryHelper) OrderAndLimitBuilder(generateSQL *strings.Builder, d map[string]string) {
	if orderColumn, ok := d["column"]; ok {
		parts := strings.Split(orderColumn, ".")
		if len(parts) > 1 {
			generateSQL.WriteString(fmt.Sprintf(` order by %s`, strings.Join(strings.Fields(orderColumn), "")))
		}else {
			generateSQL.WriteString(fmt.Sprintf(` order by e.%s`, strings.Join(strings.Fields(orderColumn), "")))
		}
		
		if order, ok := d["orientation"]; ok {
			generateSQL.WriteString(fmt.Sprintf(` %s`, strings.Join(strings.Fields(order), "")))
		}

	}
	size := strconv.Itoa(domain.DEFAULT_LIMIT)
	if sizeP, ok := d["size"]; ok {
		size = sizeP
	}
	generateSQL.WriteString(fmt.Sprintf(` limit %s`, strings.Join(strings.Fields(size), "")))
	// params = append(params,r.convertor.StrtoInt(size))
}
