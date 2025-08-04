package helpers

import (
	"encoding/json"
	"erp/internal/domain"
	"erp/pkg/api"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/lib/pq"
)

type ConvertorHelper interface {
	StrtoInt(s string) int64
	StrToBool(s string) bool

	StrToArrayInt64(s string) []int64
	ToPaginationParams(page, size string) (limit, offset int)
	ArrayToString(a any) string
	DecodeString(a string) string
	GenerateQueryMap(params interface{}) map[string]string
	GetConditionFromQuery(q string, column string, params *[]interface{}) (res string)
	DataMap(data interface{}) (map[string]interface{}, error)
	CopyStructData(src interface{}, dst interface{}) error 
	UpdateStruct(original, new interface{}) error
}

type convertorHelper struct {
	parseQuery api.ParseQuery
}

func NewConvertorHelper() ConvertorHelper {
	parseQuery := api.NewParseQuery()
	return &convertorHelper{
		parseQuery: parseQuery,
	}
}

func (h *convertorHelper)UpdateStruct(original, new interface{}) error {
	vOriginal := reflect.ValueOf(original)
	vNew := reflect.ValueOf(new)

	// Ensure the original is a pointer to a struct
	if vOriginal.Kind() != reflect.Ptr || vOriginal.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("original must be a pointer to a struct")
	}
	if vNew.Kind() != reflect.Struct {
		return fmt.Errorf("new must be a struct")
	}

	vOriginal = vOriginal.Elem() // Dereference the original pointer

	// Iterate over each field of the struct
	for i := 0; i < vOriginal.NumField(); i++ {
		fieldOriginal := vOriginal.Field(i)
		fieldNew := vNew.Field(i)

		// Check if the field in the new struct is not empty (zero value)
		if !fieldNew.IsZero() {
			// Set the field in the original struct to the value from the new struct
			if fieldOriginal.CanSet() {
				fieldOriginal.Set(fieldNew)
			}
		}
	}

	return nil
}

func (h *convertorHelper)CopyStructData(src interface{}, dst interface{}) error {
	// Get the reflect values of both structs
	srcVal := reflect.ValueOf(src)
	dstVal := reflect.ValueOf(dst)

	// Ensure both are pointers to structs
	if srcVal.Kind() != reflect.Struct || dstVal.Kind() != reflect.Ptr || dstVal.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("both arguments must be structs, with the destination being a pointer to a struct")
	}

	// Dereference the destination pointer to access the struct value
	dstVal = dstVal.Elem()

	// Iterate over the fields of the source struct
	for i := 0; i < srcVal.NumField(); i++ {
		// Get the field name and value from the source struct
		fieldName := srcVal.Type().Field(i).Name
		fieldValue := srcVal.Field(i)

		// Check if the destination struct has a matching field
		dstField := dstVal.FieldByName(fieldName)
		if dstField.IsValid() && dstField.CanSet() {
			// Copy the value from the source to the destination
			dstField.Set(fieldValue)
		}
	}

	return nil
}

func (h *convertorHelper) DataMap(data interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	val := reflect.ValueOf(data)
	if val.Kind() != reflect.Struct {
		return result, domain.INVALID_TYPE
	}
	for i := 0; i < val.NumField(); i++ {
		field:= val.Type().Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "-" || jsonTag == "" {
			// Skip fields that are explicitly ignored
			continue
		}
		fieldValue := val.Field(i).Interface()
		result[jsonTag] = fieldValue
	}
	return result,nil
}

func (h *convertorHelper) GenerateQueryMap(params interface{}) map[string]string {
	// Create a map to store the query parameters
	queryMap := make(map[string]string)

	// Use reflection to iterate over the struct fields
	v := reflect.ValueOf(params)
	t := reflect.TypeOf(params)

	// Ensure the input is a struct (not a pointer)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	// Iterate over each field in the struct
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// If the field is an embedded struct, we need to iterate over its fields too
		if field.Kind() == reflect.Struct {
			// Recursively call GenerateQueryMap for the embedded struct
			embeddedQueryMap := h.GenerateQueryMap(field.Interface())
			for k, v := range embeddedQueryMap {
				queryMap[k] = v
			}
		} else {
			// Get the tag for the field (the query parameter name)
			queryTag := fieldType.Tag.Get("query")

			// Check if the field has the "query" tag and if it's not empty
			if queryTag != "" && field.String() != "" {
				// Add the field to the query map
				queryMap[queryTag] = field.String()
			}
		}
	}

	return queryMap
}

func (h *convertorHelper) DecodeString(a string) string {
	decoded, err := url.QueryUnescape(a)
	if err != nil {
		fmt.Println("Error decoding:", err)
		return ""
	}
	return decoded
}

func (h *convertorHelper) ArrayToString(a any) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", ",", -1), "[]")
	//return strings.Trim(strings.Join(strings.Split(fmt.Sprint(a), " "), delim), "[]")
	//return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(a)), delim), "[]")
}

func (h *convertorHelper) StrToArrayInt64(s string) []int64 {
	var result []int64

	// Unmarshal the input JSON string into the result slice
	err := json.Unmarshal([]byte(s), &result)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
	}
	return result
}

func (h *convertorHelper) GetConditionFromQuery(q string, column string, params *[]interface{}) (res string) {
	fmt.Println("Q", q)
	operator, value := h.parseQuery.GetValue(q)
	fmt.Println("Q d", operator, value)
	// if operator == "" || len(value) == 0 {
	// 	return
	// }
	column = strings.Join(strings.Fields(column), "")
	switch operator {
	case "=":
		*params = append(*params, value[0])
		return fmt.Sprintf(` and %s = ?`, column)
	case "!=":
		*params = append(*params, value[0])
		return fmt.Sprintf(` and %s != ?`, column)
	case ">":
		*params = append(*params, value[0])
		return fmt.Sprintf(` and %s > ?`, column)
	case ">=":
		*params = append(*params, value[0])
		return fmt.Sprintf(` and %s >= ?`, column)
	case "<":
		*params = append(*params, value[0])
		return fmt.Sprintf(` and %s < ?`, column)
	case "<=":
		*params = append(*params, value[0])
		return fmt.Sprintf(` and %s <= ?`, column)
	case "like":
		v, ok := value[0].(string)
		if ok {
			q := "%"+v+"%"
			*params = append(*params, strings.ToLower(q))
			return fmt.Sprintf(` and lower(%s) like ?`, column)
		}
		return ""
	case "between":
		fmt.Println("LEN", len(value))
		if len(value) < 2 {
			return
		}
		fmt.Println("LEN VAL", len(value))
		*params = append(*params, value[0], value[1])
		return fmt.Sprintf(` and %s between ? and ?`, column)
	case "in":
		fmt.Println("LEN", len(value))
		if len(value) < 2 {
			return
		}
		fmt.Println("LEN VAL", len(value))
		*params = append(*params, value[0], value[1])
		return fmt.Sprintf(` and %s in (?)`, pq.Array(value))
	default:
		*params = append(*params, q)
		return fmt.Sprintf(` and %s = ?`, column)
	}
}

func (h *convertorHelper) StrtoInt(s string) int64 {
	res, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return int64(res)
}

func (h *convertorHelper) StrToBool(s string) bool {

	res, _ := strconv.ParseBool(s)
	return res
}

func (h *convertorHelper) ToPaginationParams(page, size string) (limit, offset int) {
	p := h.StrtoInt(page)
	s := h.StrtoInt(size)
	limit = int(s)
	offset = int(p) * int(s)
	return
}

// func (h *convertorHelper) OrderExpression()

// func (h *ConvertorHelper) IntToStr(s int) string {
// 	return strconv.Itoa(s)
// }
