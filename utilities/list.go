package utils

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func RemoveValue(slice []string, value string) []string {
	var index int
	for i, v := range slice {
		if v == value {
			index = i
			break
		}
	}

	// If the value was found, remove it
	if index < len(slice) {
		slice = append(slice[:index], slice[index+1:]...)
	}

	return slice
}

func Contains[T any](slice []T, value T) bool {
	for _, v := range slice {
		if reflect.DeepEqual(v, value) {
			return true
		}
	}
	return false
}

func Map[T, U any](input []T, fn func(T) U) []U {
	result := make([]U, len(input))
	for i, v := range input {
		result[i] = fn(v)
	}
	return result
}

func MapMap[T comparable, U any, V any](input map[T]U, fn func(T, U) V) []V {
	result := []V{}
	for k, v := range input {
		result = append(result, fn(k, v))
	}
	return result
}

func ForEach[T any](input []T, fn func(T)) {
	for _, v := range input {
		fn(v)
	}
}

func AnyNil[T any](obj T, fields []string) (error, string) {
	objValue := reflect.ValueOf(obj)

	if objValue.Kind() != reflect.Struct {
		return errors.New("input is not a struct"), ""
	}

	var nilFields []string

	for _, fieldName := range fields {
		fieldValue := objValue.FieldByName(fieldName)

		if !fieldValue.IsValid() || fieldValue.IsZero() {
			if fieldName == "Tipe" {
				nilFields = append(nilFields, "Type")
			} else {
				nilFields = append(nilFields, fieldName)
			}
		}
	}

	if len(nilFields) > 0 {
		return errors.New("There are null fields"), fmt.Sprintf("Fields %v are nil", nilFields)
	}

	return nil, ""
}

func AllMatch[T any](input []T, fn func(T) bool) bool {
	for _, v := range input {
		res := fn(v)
		if !res {
			return false
		}
	}
	return true
}

func AnyMatch[T any](input []T, fn func(T) bool) bool {
	matched := false
	for _, v := range input {
		res := fn(v)
		if res {
			matched = true
		}
	}
	return matched
}
func Filter[T any](input []T, fn func(T) bool) []T {
	var filtered []T
	for _, v := range input {
		if fn(v) {
			filtered = append(filtered, v)
		}
	}
	return filtered
}

func Group[T any](data []T, groupingFunc func(p T) string) map[string][]T {
	result := make(map[string][]T)

	for _, v := range data {
		key := groupingFunc(v)
		result[key] = append(result[key], v)
	}

	return result
}

func RemoveCharactersBeforeSubstring(input, substring string) string {
	index := strings.Index(input, substring)
	if index == -1 {
		return input
	}

	return input[index:]
}
