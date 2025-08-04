package dto

type FilterOptionType string

const (
	FILTER_TYPE_STRING FilterOptionType = "string"
	FILTER_TYPE_DATE   FilterOptionType = "date"
	FILTER_TYPE_NUMBER   FilterOptionType = "number"
)

type (
	FilterOptionDto struct {
		Param     string           `json:"param"`
		Name      string           `json:"name"`
		Type      FilterOptionType `json:"type"`
		Operators []string         `json:"operators"`
		Options   []string         `json:"options"`
		PartyType string           `json:"party_type"`
		// Fields []FilterFieldDto   `json:"fields"`
	}
	FilterFieldDto struct {
		Operator string   `json:"operator"`
		Options  []string `json:"options"`
	}
)

type FilterOperators struct {
	DateOperators   []string
	StringOperators []string
}

var (
	DateOperators   = []string{"=", "!=", ">", "<", ">=", "<=", "between", "in"}
	StringOperators = []string{"=", "!="}
	// StringOptionsOperators = []string{"=", "!=", "in"}
	NumberOperators = []string{"=", "!=", ">", "<", ">=", "<=", "between"}
)
