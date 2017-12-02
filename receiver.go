package receiver

import (
	"errors"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

const tagKey = "request"

// Set data value from http.Request.Form
func SetData(data interface{}, r *http.Request) error {

	dataValue := reflect.ValueOf(data)
	if dataValue.Kind() != reflect.Ptr {
		return errors.New("data argument must has address eg: &result")
	}

	dataType := reflect.TypeOf(data)
	for i := 0; i < dataValue.Elem().NumField(); i++ {
		field := dataType.Elem().Field(i)
		tagValue := field.Tag.Get(tagKey)
		keyValue, validator := getRequirements(tagValue)

		//getting value from request form
		value := r.FormValue(keyValue)
		if value == "" && isRequired(validator) {
			return errors.New("required value form is empty on the key = " + keyValue)
		}

		// parse the value to the respected type field
		fieldElement := dataValue.Elem().Field(i)
		fieldAddr := fieldElement.Addr().Interface()
		if value != "" {
			switch fieldType := fieldAddr.(type) {
			case *int64:
				parseValue, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					return errors.New("Parse integer failed check your tag key " + keyValue)
				}
				*fieldType = parseValue
			case *float64:
				parseValue, err := strconv.ParseFloat(value, 64)
				if err != nil {
					return errors.New("Parse float failed check your tag key " + keyValue)
				}
				*fieldType = parseValue
			case *string:
				*fieldType = value
			case *bool:
				parseValue, err := strconv.ParseBool(value)
				if err != nil {
					return errors.New("Parse bool failed check your tag key " + keyValue)
				}
				*fieldType = parseValue
			}
		}

	}

	return nil
}

// check validator is required or not
func isRequired(validator string) bool {
	if validator == "required" {
		return true
	}

	return false
}

// getting value [0] = value [1] = required/optional
// which coming from data interface{}
// return the key value and validator in sequence order
func getRequirements(tagValue string) (string, string) {
	values := strings.Split(tagValue, ",")
	if len(values) == 2 {
		return values[0], values[1]
	}

	return "", ""

}

func getAllTags(data interface{}, r *http.Request) []string {
	var tagValues []string
	t := reflect.TypeOf(data)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tagValue := field.Tag.Get(tagKey)
		tagValues = append(tagValues, tagValue)
	}

	return tagValues
}
