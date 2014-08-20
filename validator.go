package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	ErrRuleNotFound          = errors.New("Rule not found")
	ErrStructPointerExpected = errors.New("The subject for validation must be a pointer to a struct type")
	ErrUnsupportedType       = errors.New("Unsupported type for rule")
)

var defaultValidator = New()

type Rule interface {
	Validate(data interface{}, field string, params map[string]string) error
}

type ruleMap map[string]Rule

type validator struct {
	registeredRules ruleMap
	data            interface{}
	errors          map[string][]error
	logicError      error
}

func New() *validator {
	v := &validator{
		registeredRules: make(ruleMap, 0),
		errors:          make(map[string][]error, 0),
	}
	v.RegisterRule("length", &lengthRule{})

	return v
}

func (v *validator) RegisterRule(name string, rule Rule) {
	//TODO: mutex read / write lock
	v.registeredRules[name] = rule
}

func (v *validator) getRule(name string) (Rule, error) {
	//TODO: mutex read lock
	r := v.registeredRules[name]
	if r != nil {
		return r, nil
	}
	return nil, ErrRuleNotFound
}

func RegisterRule(name string, rule Rule) {
	defaultValidator.RegisterRule(name, rule)
}

type ruleExtractor struct {
	subject   interface{}
	numFields int
	current   int
}

func newruleExtractor(data interface{}) *ruleExtractor {
	numFields := reflect.ValueOf(data).Elem().NumField()
	return &ruleExtractor{subject: data, numFields: numFields}
}

func (e *ruleExtractor) next() bool {
	e.current++
	if e.current > e.numFields {
		return false
	}
	return true
}

//func (e *ruleExtractor) extract() (string, []Rule) {
//	rules := make([]Rule, 0)
//	index := e.current - 1
//	elem := reflect.TypeOf(e.subject).Elem().Field(index)
//	fieldName := elem.Name
//	var ruleName, ruleParams string
//
//	tag := elem.Tag.Get("validation")
//	if tag == "" {
//		return fieldName, rules
//	}
//	for _, ruleStr := range strings.Split(tag, "|") {
//		ruleParts := strings.Split(ruleStr, ":")
//		ruleName = ruleParts[0]
//		if len(ruleParts) > 1 {
//			ruleParams = ruleParts[1]
//		} else {
//			ruleParams = ""
//		}
//		rule, err := getRule(ruleName, ruleParams, fieldName, e.subject)
//		if err == nil {
//			rules = append(rules, rule)
//		}
//	}
//
//	return fieldName, rules
//}

func (v *validator) Validate(data interface{}) error {
	if isStructPointer(data) == false {
		return ErrStructPointerExpected
	}

	numFields := reflect.ValueOf(data).Elem().NumField()

	for curField := 0; curField < numFields; curField++ {
		err := v.validateField(curField)
		if err != nil {
			return err
		}
	}
	return nil
}

func (v *validator) validateField(i int) error {

	elem := reflect.TypeOf(v.data).Elem().Field(i)
	fieldName := elem.Name

	tag := elem.Tag.Get("validation")
	if tag == "" {
		return nil
	}

	for _, ruleStr := range strings.Split(tag, "|") {
		var j = strings.Index(tag, ":")
		var ruleParamsStr = ruleStr[j+1:]
		var ruleParams map[string]string

		var ruleName = ruleStr[0:j]

		ruleParams = make(map[string]string, 0)

		for _, paramPart := range strings.Split(ruleParamsStr, ",") {
			var tmpParam = strings.Split(paramPart, ":")
			ruleParams[tmpParam[0]] = tmpParam[1]
		}

		var fieldCheck = func() {
			rule, err := v.getRule(ruleName)
			if err != nil {
				v.logicError = err
				return
			}
			err = rule.Validate(v.data, fieldName, ruleParams)
			if err != nil {
				v.errors[fieldName] = append(v.errors[fieldName], err)
			}
		}

		v.safeExec(fieldCheck)
	}
	return nil
}

type safeFunc func()

func (v *validator) safeExec(f safeFunc) {
	defer func() {
		if recErr := recover(); recErr != nil {
			switch errv := recErr.(type) {
			case string:
				v.logicError = errors.New(errv)
			case error:
				v.logicError = errv
			default:
				v.logicError = errors.New(fmt.Sprintf("Panic recovered with type: %+v", recErr))
			}
		}
	}()
	f()
}

func isStructPointer(data interface{}) bool {
	if reflect.TypeOf(data).Kind() != reflect.Ptr {
		return false
	}
	if reflect.ValueOf(data).Elem().Kind() != reflect.Struct {
		return false
	}
	return true
}

//func fieldPresent(data interface{}, name string) bool {
//	_, present := reflect.TypeOf(data).Elem().FieldByName(name)
//	return present
//}

func getInterfaceValue(data interface{}, name string) interface{} {
	return reflect.ValueOf(data).Elem().FieldByName(name).Interface()
}

func ruleString(ruleName, structField string, data interface{}) string {
	return fmt.Sprintf("<<Validation Rule: %s. Field: %s. Data: %s>>", ruleName, structField, fmt.Sprintf("%+v", data))
}

// toString
func toString(value interface{}) (string, bool) {
	switch v := value.(type) {
	case string, *string, int, *int, int32, *int32, int64, *int64:
		return fmt.Sprintf("%v", v), true
	default:
		return "", false
	}
}

func mustStringify(value interface{}) string {
	strVal, ok := toString(value)
	if ok == false {
		panic(ErrUnsupportedType)
	}
	return strVal
}
