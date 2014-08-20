package validator

import ()

type requiredRule struct{}

func (r *requiredRule) Validate(data interface{}, field string, params map[string]string) error {

	//fieldVal := getInterfaceValue(data, field)

	return nil
}
