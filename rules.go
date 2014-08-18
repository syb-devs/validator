package validator

import (
	"errors"
	"fmt"
	"strconv"
	"unicode/utf8"
)

func init() {
	RegisterRule("min_length", newMinLengthRule)
}

type minLengthRule struct {
	field  string
	length int
	data   interface{}
}

func newMinLengthRule(fieldName string, params string, data interface{}) (Rule, error) {
	//Todo: strict string to int parsing
	lengthParsed, err := strconv.Atoi(params)
	if err != nil {
		return nil, errors.New("rule minLength must be an integer")
	}
	return &minLengthRule{field: fieldName, length: lengthParsed, data: data}, nil
}

func (r *minLengthRule) Validate() (*inputError, error) {
	if !fieldPresent(r.data, r.field) {
		return nil, fmt.Errorf("field %s not present and tried to evaluate", r.field)
	}

	fInterface := getInterfaceValue(r.data, r.field)

	var length int

	switch v := fInterface.(type) {
	case string:
		length = utf8.RuneCountInString(v)
	case int:
		length = utf8.RuneCountInString(strconv.Itoa(v))
	default:
		return nil, errors.New("Unsupported type for min_length rule")
	}

	if length < r.length {
		message := fmt.Sprintf("The field %s should have a minimum length of %d characters", r.field, r.length)
		return &inputError{field: r.field, message: message}, nil
	}

	return nil, nil
}
