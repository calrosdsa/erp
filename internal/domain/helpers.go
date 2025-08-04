package domain

import (
	"fmt"
	"reflect"
)

func AssignIfNotEmpty(field interface{}, value interface{}) {
	// Check if the value is not empty based on its type
	if !isEmpty(value) {
		// Ensure the field is a pointer
		fieldVal := reflect.ValueOf(field)
		if fieldVal.Kind() != reflect.Ptr {
			fmt.Printf("Error: 'field' must be a pointer, got %s\n", fieldVal.Kind())
			return
		}
		fieldVal.Set(reflect.ValueOf(value))
		// Ensure the field is a pointer to a valid type
		fieldElem := fieldVal.Elem()
		fmt.Println("FIELD ELEM",fieldElem)
		if !fieldElem.CanSet() {
			fmt.Println("Error: Cannot set the value to field (possibly because it's unaddressable).")
			return
		}

		// Check if the value can be assigned to the field
		if canAssign(fieldElem, value) {
			// Set the value to the field
			fieldElem.Set(reflect.ValueOf(value))
			fmt.Println("Assigned value:", value)
		} else {
			fmt.Println("Error: Incompatible types between field and value.")
		}
	}
}

// isEmpty returns true if the value is considered "empty" for its type.
func isEmpty(value interface{}) bool {
	v := reflect.ValueOf(value)

	// Check for zero value (empty strings, 0 for numbers, nil for pointers, empty slices/maps, etc.)
	switch v.Kind() {
	case reflect.String:
		// Empty string is considered "empty"
		return v.Len() == 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// 0 for integer types is considered "empty"
		return v.Int() == 0
	case reflect.Float32, reflect.Float64:
		// 0.0 for float types is considered "empty"
		return v.Float() == 0
	case reflect.Bool:
		// false for booleans is considered "empty"
		return !v.Bool()
	case reflect.Ptr:
		// nil pointers are considered "empty"
		return v.IsNil()
	case reflect.Slice, reflect.Array:
		// Empty slices and arrays are considered "empty"
		return v.Len() == 0
	case reflect.Map:
		// Empty maps are considered "empty"
		return v.IsNil() || v.Len() == 0
	case reflect.Interface:
		// Empty interfaces are considered "empty" if they are nil
		return v.IsNil()
	default:
		// For other types, we check if the value is the "zero" value (e.g., 0 for numbers, nil for pointers)
		return v.IsZero()
	}
}

// canAssign checks if a value can be assigned to a field (checks if types are compatible).
func canAssign(field reflect.Value, value interface{}) bool {
	// Get the type of the value being assigned
	valueVal := reflect.ValueOf(value)

	// Handle the case where both field and value are pointers
	if field.Kind() == reflect.Ptr && valueVal.Kind() == reflect.Ptr {
		// Check if the types of the pointers are compatible (pointer-to-pointer)
		return field.Type().Elem().AssignableTo(valueVal.Type().Elem())
	}

	// Handle case where field is a pointer but value is not (pointer-to-value)
	if field.Kind() == reflect.Ptr {
		// Check if the field type can hold the value type
		return field.Type().Elem().AssignableTo(valueVal.Type())
	}

	// Handle case where both field and value are not pointers
	return field.Type().AssignableTo(valueVal.Type())
}