package validator

func init() {
	RegisterRule("min_length", newMinLengthRule)
}

type minLengthRule struct {
	field  string
	length int
}

func newMinLengthRule(fieldName string, params string, obj interface{}) (Rule, error) {
	return &minLengthRule{}, nil
}

func (r *minLengthRule) Validate() (*inputError, error) {
	//TODO: validate
	return nil, nil
}
