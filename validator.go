package validator

import (
	"errors"
	// "fmt"
	"reflect"
	"strings"
)

type field struct {
	name  string
	value interface{}
}

type inputError struct {
	field   field
	message string
}

type inputErrors []inputError

type Rule interface {
	Validate() (*inputError, error)
}

type rules []Rule

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

func newruleExtractor(obj interface{}) *ruleExtractor {
	numFields := reflect.ValueOf(obj).Elem().NumField()
	return &ruleExtractor{subject: obj, numFields: numFields}
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
	for _, ruleStr := range strings.Split(tag, "|") {
		ruleSplit := strings.Split(ruleStr, ":")
		ruleName = ruleSplit[0]
		if len(ruleSplit) > 1 {
			ruleParams = ruleSplit[1]
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

func Validate(obj interface{}) (inputErrors, error) {
	//TODO: validate obj is a pointer to struct, error otherwise
	errors := make(inputErrors, 0)
	var rules rules

	for extractor := newruleExtractor(obj); extractor.next(); {
		rules = extractor.extract()
		for _, rule := range rules {
			inputError, err := rule.Validate()

			if err != nil {
				return errors, err
			}

			if inputError != nil {
				errors = append(errors, *inputError)
			}
		}
	}
	return errors, nil
}

func getRule(name string, params string, fieldName string, obj interface{}) (Rule, error) {
	ruleConstructor := ruleMap[name]
	if ruleConstructor == nil {
		//TODO: specify which rule (or use custom error struct)
		return nil, errors.New("Rule not found")
	}
	return ruleConstructor(fieldName, params, obj)
}
