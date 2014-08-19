package validator

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"unicode/utf8"
)

func init() {
	RegisterRule("min_length", newMinLengthRule)
	RegisterRule("max_length", newMaxLengthRule)
	RegisterRule("email", newEmailRule)

	emailRegexp = regexp.MustCompile(`^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$`)
}

var (
	ErrUnsupportedType = errors.New("Unsupported type for rule")
)

var (
	emailRegexp *regexp.Regexp
)

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
	fInterface := getInterfaceValue(r.data, r.field)
	str, ok := toString(fInterface)
	if ok == false {
		return nil, ErrUnsupportedType
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
	fInterface := getInterfaceValue(r.data, r.field)
	str, ok := toString(fInterface)
	if ok == false {
		return nil, ErrUnsupportedType
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

type emailRule struct {
	field string
	data  interface{}
}

func newEmailRule(fieldName string, ruleValue string, data interface{}) (Rule, error) {
	return &emailRule{field: fieldName, data: data}, nil
}

func (r *emailRule) Validate() (*inputError, error) {
	strVal := mustStringify(getInterfaceValue(r.data, r.field))

	if emailRegexp.MatchString(strVal) {
		return nil, nil
	}
	return &inputError{field: r.field, message: "Not a valid email"}, nil
}

func (r *emailRule) String() string {
	return ruleString("max length", r.field, r.data)
}
