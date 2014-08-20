package validator_test

import (
	"bitbucket.org/simplifyourbusiness/validator"
	"testing"
)

func TestValidate(t *testing.T) {
	type data struct {
		Field string `validation:"length:>4" `
	}

	v := validator.New()
	err := v.Validate(&data{})

	if err != nil {
		t.Errorf("Holy shit! %s", err)
	}
}

func _TestStructPointer(t *testing.T) {
	type data struct {
		Field string `validation:"min_length:4" `
	}

	v := validator.New()
	err := v.Validate(data{})

	if err == nil {
		t.Errorf("Expecting error because data is not a pointer to struct")
	}

	if err != validator.ErrStructPointerExpected {
		t.Errorf("Expecting error message to be %s, got %s", validator.ErrStructPointerExpected, err)
	}
}

func _TestUnexportedField(t *testing.T) {
	type data struct {
		field string `validation:"min_length:4" `
	}

	v := validator.New()
	err := v.Validate(&data{})

	if err == nil {
		t.Errorf("Expecting error: cannot return value obtained from unexported field or method")
	}
}

func _TestEmptyValidationTag(t *testing.T) {
	type data struct {
		field string
	}

	v := validator.New()
	err := v.Validate(&data{})

	if err != nil {
		t.Errorf("Error during validation")
	}
}
