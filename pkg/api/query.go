package api

import (
	"encoding/json"
)

type ParseQuery interface {
	GetValue(v string) (operator string,value []interface{})
}

type parseQuery struct{}

func NewParseQuery() ParseQuery {
	return &parseQuery{}
}


func (p *parseQuery) GetValue(v string) (operator string, value []interface{}) {
	var values []interface{}

	// Unmarshal the JSON string into the slice
	err := json.Unmarshal([]byte(v), &values)
	if err != nil {
		// If there's an error, return empty operator and value
		return "", nil
	}

	// If there are fewer than 2 values, return empty operator and value
	if len(values) < 2 {
		return "", nil
	}

	// Extract the operator (first value)
	operator, ok := values[0].(string)
	if !ok {
		// If the first value is not a string, return empty operator and value
		return "", nil
	}

	// If the operator is "between", treat the second value as a slice
	if operator == "between" || operator == "in" {
		// Ensure the second value is a slice of interface{}
		value, ok = values[1].([]interface{})
		if !ok {
			// If it's not a slice, return empty operator and value
			return "", nil
		}
	} else {
		// For other operators, treat the second value as a single interface{}
		value = []interface{}{values[1]}
	}

	return operator, value
}
