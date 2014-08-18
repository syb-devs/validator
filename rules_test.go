package validator_test

import (
	"bitbucket.org/simplifyourbusiness/validator"
	"testing"
)

func TestMinLengthOK(t *testing.T) {
	type data struct {
		FieldOne string `validation:"min_length:3" `
	}

	v := validator.New()

	err := v.Validate(
		&data{FieldOne: "foo"})
	if err != nil {
		t.Errorf("Error during validation: %s", err.Error())
	}

	if v.Errors.Count() > 0 {
		t.Errorf("Input Errors: %s", v.Errors)
	}
}

func TestMinLengthKO(t *testing.T) {
	type data struct {
		Field string `validation:"min_length:4" `
	}

	v := validator.New()

	err := v.Validate(
		&data{Field: "foo"})

	if err != nil {
		t.Errorf("Error during validation: %s", err.Error())
	}

	if v.Errors.Count() != 1 {
		t.Errorf("Expected exactly 1 input error, got %d", v.Errors.Count())
	}

	inputErr := v.Errors[0]
	expectedMessage := "The field Field should have a minimum length of 4 characters"
	actualMessage := inputErr.Message()
	if actualMessage != expectedMessage {
		t.Errorf("Expecting error message to be %s, got %s", expectedMessage, actualMessage)
	}
}
