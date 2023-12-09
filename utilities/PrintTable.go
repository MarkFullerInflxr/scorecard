package utils

import (
	"fmt"
	"reflect"
)

func PrintTable(data []interface{}, fields []string) {
	// Print header
	for _, field := range fields {
		fmt.Printf("%-15s", field)
	}
	fmt.Println()

	// Print data
	for _, item := range data {
		val := reflect.ValueOf(item)
		for _, field := range fields {
			fieldVal := val.FieldByName(field)
			// Check if the field is a pointer and dereference it
			if fieldVal.Kind() == reflect.Ptr && !fieldVal.IsNil() {
				fieldVal = fieldVal.Elem()
			}
			fmt.Printf("%-15v", fieldVal.Interface())
		}
		fmt.Println()
	}
}

func PrintTablePretty(data []interface{}, fields []string) {
	// Find maximum width for each column
	maxWidth := make(map[string]int)

	for _, field := range fields {
		maxWidth[field] = len(field) // Initialize with the header width
	}

	for _, item := range data {
		val := reflect.ValueOf(item)

		for _, field := range fields {
			fieldVal := val.FieldByName(field)

			// Check if the field is a pointer and dereference it
			if fieldVal.Kind() == reflect.Ptr && !fieldVal.IsNil() {
				fieldVal = fieldVal.Elem()
			}

			// Convert field value to string
			fieldStr := fmt.Sprintf("%v", fieldVal.Interface())

			// Update maximum width for the column
			if len(fieldStr) > maxWidth[field] {
				maxWidth[field] = len(fieldStr)
			}
		}
	}

	// Print header
	for _, field := range fields {
		fmt.Printf("%-*s", maxWidth[field]+2, field)
	}
	fmt.Println()

	// Print data
	for _, item := range data {
		val := reflect.ValueOf(item)

		for _, field := range fields {
			fieldVal := val.FieldByName(field)

			// Check if the field is a pointer and dereference it
			if fieldVal.Kind() == reflect.Ptr && !fieldVal.IsNil() {
				fieldVal = fieldVal.Elem()
			}

			// Convert field value to string
			fieldStr := fmt.Sprintf("%v", fieldVal.Interface())

			// Print field value with padding to match the maximum width
			fmt.Printf("%-*s", maxWidth[field]+2, fieldStr)
		}
		fmt.Println()
	}
}
