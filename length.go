// Copyright 2014 Simplify Your Bussiness S.L. All rights reserved.

package validator

import (
	"errors"
	"fmt"
	"strconv"
	"unicode/utf8"
)

// lengthRule struct holds Validate() method to satisfy the Validator interface.
type lengthRule struct{}

// Validate checks that the given data conforms to the length constraints given as parameters.
func (r *lengthRule) Validate(data interface{}, field string, params map[string]string) (errorLogic, errorInput error) {
	op := params["op"]

	requiredLength, errorLogic := strconv.Atoi(params["val"])
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
