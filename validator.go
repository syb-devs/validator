package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	ErrStructPointerExpected = errors.New("The subject for validation must be a pointer to a struct type")
	ErrRuleNotFound          = errors.New("Rule not found")
)

type validator struct {
	Rules  rules
	Errors inputErrors
}

func New() *validator {
	return &validator{}
}

type field struct {
	name  string
	value interface{}
}

type inputError struct {
	field   string
	message string
}

func (e inputError) Message() string {
	return e.message
}

func (e inputError) String() string {
	return e.Message()
}

type inputErrors []inputError

func (e inputErrors) String() string {
	var r string
	for _, err := range e {
		r = r + err.String()
	}
	return r
}

func (e inputErrors) Count() int {
	return len(e)
}

type Rule interface {
	Validate() (*inputError, error)
}

type rules []Rule

func (r rules) Count() int {
	return len(r)
}

type RuleConstructor func(fieldName string, params string, dataStruct interface{}) (Rule, error)

var ruleMap = make(map[string]RuleConstructor, 0)

func RegisterRule(name string, ruleFunc RuleConstructor) {
	ruleMap[name] = ruleFunc
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

func (e *ruleExtractor) extract() rules {
	rules := make(rules, 0)
	index := e.current - 1
	elem := reflect.TypeOf(e.subject).Elem().Field(index)
	fieldName := elem.Name
	var ruleName, ruleParams string

	tag := elem.Tag.Get("validation")
	if tag == "" {
		return rules
	}
	for _, ruleStr := range strings.Split(tag, "|") {
		ruleParts := strings.Split(ruleStr, ":")
		ruleName = ruleParts[0]
		if len(ruleParts) > 1 {
			ruleParams = ruleParts[1]
		} else {
			ruleParams = ""
		}
		rule, err := getRule(ruleName, ruleParams, fieldName, e.subject)
		if err == nil {
			rules = append(rules, rule)
		}
	}

	return rules
}

func (v *validator) Validate(data interface{}) error {
	if isStructPointer(data) == false {
		return ErrStructPointerExpected
	}
	v.Errors = make(inputErrors, 0)

	var rules rules
	for extractor := newruleExtractor(data); extractor.next(); {
		rules = extractor.extract()
		v.Rules = append(v.Rules, rules...)
	}

	for _, rule := range v.Rules {
		inputError, err := checkRule(rule)

		if err != nil {
			return err
		}

		if inputError != nil {
			v.Errors = append(v.Errors, *inputError)
		}
	}
	return nil
}

func checkRule(rule Rule) (ierr *inputError, err error) {
	defer func(ierr *inputError, e error) {
		if err := recover(); err != nil {
			ierr = nil
			if errString, ok := err.(string); ok {
				e = errors.New(errString)
			}
			return
		}
	}(ierr, err)

	//fmt.Printf("Checking rule: %s...\n", rule)

	ierr, err = rule.Validate()
	return
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

func getRule(name string, params string, fieldName string, data interface{}) (Rule, error) {
	ruleConstructor := ruleMap[name]
	if ruleConstructor == nil {
		return nil, ErrRuleNotFound
	}
	return ruleConstructor(fieldName, params, data)
}

func fieldPresent(data interface{}, name string) bool {
	_, present := reflect.TypeOf(data).Elem().FieldByName(name)
	return present
}

func getInterfaceValue(data interface{}, name string) interface{} {
	return reflect.ValueOf(data).Elem().FieldByName(name).Interface()
}

func ruleString(ruleName, structField string, data interface{}) string {
	return fmt.Sprintf("<<Validation Rule: %s. Field: %s. Data: %s>>", ruleName, structField, fmt.Sprintf("%+v", data))
}
