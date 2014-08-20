package validator

import (
	"errors"
	"fmt"
	"strconv"
	"unicode/utf8"
)

type lengthRule struct {
}

func (r *lengthRule) Validate(data interface{}, field string, params map[string]string) error {
	op := params["op"]

	requiredLength, err := strconv.Atoi(params["val"])
	if err != nil {
		return err
	}

	fieldVal := getInterfaceValue(data, field)
	length := utf8.RuneCountInString(mustStringify(fieldVal))

	var ok bool
	var errOp string
	switch op {
	case "=":
		ok = length == requiredLength
		errOp = "equal to"
	case ">":
		ok = length > requiredLength
		errOp = "greater than"
	case ">=":
		ok = length >= requiredLength
		errOp = "greater than, or equal to"
	case "<":
		ok = length < requiredLength
		errOp = "lower than"
	case "<=":
		ok = length < requiredLength
		errOp = "lower than, or equal to"
	default:
		return errors.New("Invalid operator")
	}

	if ok {
		return nil
	}
	return fmt.Errorf("The field %s should have a length %d %d", field, errOp, requiredLength)
}
