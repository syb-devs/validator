package validator

import (
	"errors"
	"fmt"
	"strconv"
	"unicode/utf8"
)

func init() {
	RegisterRule("min_length", newMinLengthRule)
	RegisterRule("max_length", newMaxLengthRule)
}

type minLengthRule struct {
	field  string
	length int
	data   interface{}
}

/* MIN LENGHT */
func newMinLengthRule(fieldName string, ruleValue string, data interface{}) (Rule, error) {
	//Todo: strict string to int parsing
	lengthParsed, err := strconv.Atoi(ruleValue)
	if err != nil {
		return nil, err
	}
	return &minLengthRule{field: fieldName, length: lengthParsed, data: data}, nil
}

func (r *minLengthRule) Validate() (*inputError, error) {
	if !fieldPresent(r.data, r.field) {
		return nil, fmt.Errorf("field %s not present and tried to evaluate", r.field)
	}

	fInterface := getInterfaceValue(r.data, r.field)
	str, ok := toString(fInterface)
	if ok == false {
		return nil, errors.New("Unsupported type for min_length rule")
	}

	length := utf8.RuneCountInString(str)

	if length < r.length {
		message := fmt.Sprintf("The field %s should have a minimum length of %d characters", r.field, r.length)
		return &inputError{field: r.field, message: message}, nil
	}

	return nil, nil
}

func (r *minLengthRule) String() string {
	return ruleString("min length", r.field, r.data)
}

/* MAX LENGHT */
type maxLengthRule struct {
	field  string
	length int
	data   interface{}
}

func newMaxLengthRule(fieldName string, ruleValue string, data interface{}) (Rule, error) {
	lengthParsed, err := strconv.Atoi(ruleValue)
	if err != nil {
		return nil, errors.New("rule maxLength must be an integer")
	}
	return &maxLengthRule{field: fieldName, length: lengthParsed, data: data}, nil
}

func (r *maxLengthRule) Validate() (*inputError, error) {
	if !fieldPresent(r.data, r.field) {
		return nil, fmt.Errorf("field %s not present and tried to evaluate", r.field)
	}

	fInterface := getInterfaceValue(r.data, r.field)
	str, ok := toString(fInterface)
	if ok == false {
		return nil, errors.New("Unsupported type for min_length rule")
	}

	length := utf8.RuneCountInString(str)

	if length > r.length {
		message := fmt.Sprintf("The field %s should have a maximum length of %d characters", r.field, r.length)
		return &inputError{field: r.field, message: message}, nil
	}

	return nil, nil
}

func (r *maxLengthRule) String() string {
	return ruleString("max length", r.field, r.data)
}
