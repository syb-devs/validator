// Copyright 2014 Simplify Your Bussiness S.L. All rights reserved.

package validator

import (
	"errors"
	"fmt"
	"strconv"
	"unicode/utf8"
)

var ErrLengthParamCount = errors.New("This rule needs two mandatory params, operator and value")

// lengthRule struct holds Validate() method to satisfy the Validator interface.
type lengthRule struct{}

// Validate checks that the given data conforms to the length constraints given as parameters.
func (r *lengthRule) Validate(data interface{}, field string, params []string, namedParams map[string]string) (errorLogic, errorInput error) {

	if len(params) != 2 {
		errorLogic = ErrLengthParamCount
		return
	}
	op := params[0]

	requiredLength, errorLogic := strconv.Atoi(params[1])
	if errorLogic != nil {
		return
	}

	fieldVal := getInterfaceValue(data, field)
	length := utf8.RuneCountInString(mustStringify(fieldVal))

	var ok bool
	var opLiteral string
	switch op {
	case "=":
		ok = length == requiredLength
		opLiteral = "equal to"
	case ">":
		ok = length > requiredLength
		opLiteral = "greater than"
	case ">=":
		ok = length >= requiredLength
		opLiteral = "greater than, or equal to"
	case "<":
		ok = length < requiredLength
		opLiteral = "lower than"
	case "<=":
		ok = length < requiredLength
		opLiteral = "lower than, or equal to"
	default:
		return nil, errors.New("Invalid operator")
	}

	if !ok {
		errorInput = fmt.Errorf("The field %s should have a length %s %d. Actual length: %d", field, opLiteral, requiredLength, length)
		return
	}
	return
}
