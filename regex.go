package validator

import (
	"regexp"
	// "errors"
	"fmt"
	"strconv"
)

type regexpRule struct{}

func (r *regexpRule) Validate(data interface{}, field string, params map[string]string) (errorLogic, errorInput error) {

	regex := params["val"]

	allowEmpty, ok := strconv.ParseBool(params["allowEmpty"])

	fieldVal := getInterfaceValue(data, field)

	if ok == nil && allowEmpty && fieldVal == "" {
		return nil, nil
	}

	compiledRegex, err := regexp.Compile(regex)

	if err != nil {
		return fmt.Errorf("The field %s does not contain a valid regexp", field), nil
	}

	if !compiledRegex.MatchString(fieldVal.(string)) {
		return nil, fmt.Errorf("The field %s does not match regexp", field)
	}

	return nil, nil
}
