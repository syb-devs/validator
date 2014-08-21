// Copyright 2014 Simplify Your Bussiness S.L. All rights reserved.

package validator

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

var (
	ErrEmptyRegexp      = errors.New("No valid regexp was found in the first parameter")
	ErrRegexpParamCount = errors.New("This rule needs one mandatory param, the regular expression")
)

type regexpRule struct{}

// Validate checks that the field value matches the regexp passed in the val parameter
func (r *regexpRule) Validate(data interface{}, field string, params []string, namedParams map[string]string) (errorLogic, errorInput error) {

	if len(params) != 1 {
		errorLogic = ErrRegexpParamCount
		return
	}

	regex := params[0]
	if regex == "" {
		errorLogic = ErrEmptyRegexp
		return
	}

	allowEmpty, errorLogic := strconv.ParseBool(namedParams["allowEmpty"])
	if errorLogic != nil {
		return
	}

	fieldVal := getInterfaceValue(data, field)

	if allowEmpty && fieldVal == "" {
		return
	}

	compiledRegex, err := regexp.Compile(regex)

	if err != nil {
		errorLogic = fmt.Errorf("The field %s does not contain a valid regexp", field)
		return
	}

	if !compiledRegex.MatchString(fieldVal.(string)) {
		errorInput = fmt.Errorf("The field %s does not match regexp", field)
		return
	}

	return
}
