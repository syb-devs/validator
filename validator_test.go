package validator_test

import (
	"bitbucket.org/simplifyourbusiness/validator"
	"testing"
)

func TestValidate(t *testing.T) {
	type data struct {
		Field string `validation:"length:op:>,val:4" `
	}

	v := validator.New()
	err := v.Validate(data{})

	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestUnexportedField(t *testing.T) {
	type data struct {
		field string `validation:"length:op:>,val:4" `
	}

	v := validator.New()
	err := v.Validate(data{})

	if err == nil {
		t.Errorf("Expecting error: cannot return value obtained from unexported field or method")
	}
}

func TestEmptyValidationTag(t *testing.T) {
	type data struct {
		field string
	}

	v := validator.New()
	err := v.Validate(data{})

	if err != nil {
		t.Errorf("Error during validation")
	}
}
