package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
)

var (
	ErrRuleNotFound          = errors.New("Rule not found")
	ErrStructPointerExpected = errors.New("The subject for validation must be a pointer to a struct type")
	ErrUnsupportedType       = errors.New("Unsupported type for rule")
)

var (
	defaultValidator = New()
	tagName          = "validation"
)

type Rule interface {
	Validate(data interface{}, field string, params map[string]string) error
}

type errList map[string][]error

func (e errList) String() string {
	str := ""
	for field, errors := range e {
		str = str + field + ": "
		for _, err := range errors {
			str = str + err.Error() + ", "
		}
		str = str + "\n"
	}
	return str
}

func (e errList) Len() int {
	return len(e)
}

type ruleMap map[string]Rule

type validator struct {
	registeredRules ruleMap
	data            interface{}
	errors          errList
	logicError      error
	mu              sync.RWMutex
}

func RegisterRule(name string, rule Rule) {
	defaultValidator.RegisterRule(name, rule)
}

func TagName(name string) {
	tagName = name
}

func New() *validator {
	v := &validator{
		registeredRules: make(ruleMap, 0),
		errors:          make(errList, 0),
	}
	v.RegisterRule("length", &lengthRule{})

	return v
}

func (v *validator) RegisterRule(name string, rule Rule) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.registeredRules[name] = rule
}

func (v *validator) getRule(name string) (Rule, error) {
	v.mu.RLock()
	defer v.mu.RUnlock()
	r := v.registeredRules[name]
	if r != nil {
		return r, nil
	}
	return nil, ErrRuleNotFound
}

func (v *validator) Validate(data interface{}) error {
	if isStructPointer(data) == false {
		return ErrStructPointerExpected
	}

	v.data = data
	numFields := reflect.ValueOf(v.data).Elem().NumField()

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

	tag := elem.Tag.Get(tagName)
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
			if len(tmpParam) != 2 {
				return errors.New("Invalid format for params")
			}
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
		if v.logicError != nil {
			return v.logicError
		}
	}
	return nil
}

func (v *validator) Errors() *errList {
	errors := v.errors
	if len(errors) == 0 {
		return nil
	}
	return &errors
}

func (v *validator) ErrorsByField(field string) *[]error {
	if field == "" {
		return nil
	}

	errors := v.errors[field]
	if errors == nil {
		return nil
	}
	return &errors
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
