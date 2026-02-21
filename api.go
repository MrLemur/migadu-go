package migadu

import (
	"fmt"
	"reflect"
	"strings"
)

// Transform converts a struct to a map[string]interface{} for writeable API operations,
// filtering fields based on their api tags and the specified operation.
func Transform(input interface{}, operation string) (interface{}, error) {
	if input == nil {
		return nil, fmt.Errorf("input must not be nil")
	}
	if operation != "create" && operation != "update" {
		return nil, fmt.Errorf("operation must be either 'create' or 'update'")
	}
	t := reflect.TypeOf(input)
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input must be a struct")
	}
	output := make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		apiTag := field.Tag.Get("api")
		if apiTag == "read-only" {
			continue
		}
		if apiTag == "create-only" && operation == "update" {
			continue
		}

		value := reflect.ValueOf(input).Field(i).Interface()

		// Convert []string to comma-separated string
		if field.Type.Kind() == reflect.Slice && field.Type.Elem().Kind() == reflect.String {
			if slice, ok := value.([]string); ok {
				value = strings.Join(slice, ",")
			}
		}

		jsonTag := field.Tag.Get("json")
		if jsonTag == "" {
			jsonTag = strings.ToLower(field.Name)
		} else {
			jsonTag, _, _ = strings.Cut(jsonTag, ",")
		}

		output[jsonTag] = value
	}
	return output, nil
}
