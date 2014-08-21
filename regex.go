package validator

import (
	"fmt"
	"regexp"
)

type regexpRule struct{}

func (r *regexpRule) Validate(data interface{}, field string, params map[string]string) (errorLogic, errorInput error) {

	regex := params["val"]

	compiledRegex, err := regexp.Compile(regex)

	if err != nil {
		return fmt.Errorf("The field %s does not contain a valid regexp", field), nil
	}

	fieldVal := getInterfaceValue(data, field)
	if !compiledRegex.MatchString(fieldVal.(string)) {
		return nil, fmt.Errorf("The field %s does not match regexp", field)
	}

	return nil, nil
}
