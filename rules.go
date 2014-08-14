package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
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
	lengthParsed, err := strconv.Atoi(params)
	if err != nil {
		return nil, errors.New("rule minLength must be an integer")
	}
	return &minLengthRule{field: fieldName, length: lengthParsed, data: data}, nil
}

func (r *minLengthRule) Validate() (*inputError, error) {
	//TODO: validate
	field, present := reflect.TypeOf(r.data).Elem().FieldByName(r.field)
	if !present {
		return nil, fmt.Errorf("field %s not present and tried to evaluate", r.field)
	}

	fmt.Printf("%+v", field)
	// fieldValueParsed, ok := string(r.data[r.field])
	return nil, nil
}
