package validator_test

import (
	"bitbucket.org/simplifyourbusiness/validator"
	"testing"
)

/* Test min_length rule*/

func TestMinLengthOK(t *testing.T) {
	type data struct {
		FieldOne string `validation:"min_length:3" `
		FieldTwo int    `validation:"min_length:4" `
	}

	v := validator.New()

	err := v.Validate(
		&data{FieldOne: "foo", FieldTwo: 1456})
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
		t.Errorf("Input Errors: %s", v.Errors)
		return
	}

	inputErr := v.Errors[0]
	expectedMessage := "The field Field should have a minimum length of 4 characters"
	actualMessage := inputErr.Message()
	if actualMessage != expectedMessage {
		t.Errorf("Expecting error message to be %s, got %s", expectedMessage, actualMessage)
	}
}

/* Test max_length rule*/

func TestMaxLengthOK(t *testing.T) {
	type data struct {
		FieldOne string `validation:"max_length:3" `
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

func TestMaxLengthKO(t *testing.T) {
	type data struct {
		Field string `validation:"max_length:2" `
	}

	v := validator.New()

	err := v.Validate(
		&data{Field: "foo"})

	if err != nil {
		t.Errorf("Error during validation: %s", err.Error())
	}

	if v.Errors.Count() != 1 {
		t.Errorf("Expected exactly 1 input error, got %d", v.Errors.Count())
		t.Errorf("Input Errors: %s", v.Errors)
		return
	}

	inputErr := v.Errors[0]
	expectedMessage := "The field Field should have a maximum length of 2 characters"
	actualMessage := inputErr.Message()
	if actualMessage != expectedMessage {
		t.Errorf("Expecting error message to be %s, got %s", expectedMessage, actualMessage)
	}
}

func TestEmailOK(t *testing.T) {
	type data struct {
		Email string `validation:"email" `
	}

	v := validator.New()

	err := v.Validate(
		&data{Email: "john.williams@lso.co.uk"})

	if err != nil {
		t.Errorf("Error during validation: %s", err.Error())
	}

	if v.Errors.Count() > 0 {
		t.Errorf("Input Errors: %s", v.Errors)
	}
}

func TestEmailKO(t *testing.T) {
	type data struct {
		Email string `validation:"email" `
	}

	v := validator.New()

	err := v.Validate(
		&data{Email: "john.williams@lso"})

	if err != nil {
		t.Errorf("Error during validation: %s", err.Error())
	}

	if v.Errors.Count() != 1 {
		t.Errorf("Expecting exactly one validation error, got %d", v.Errors.Count())
		t.Errorf("Input Errors: %s", v.Errors)
		return
	}
}

func TestEmaiUnsupportedType(t *testing.T) {
	type data struct {
		Email struct{} `validation:"email" `
	}

	v := validator.New()
	err := v.Validate(&data{})
	if err != validator.ErrUnsupportedType {
		t.Errorf("Expected unsupported type error, got %+v", err)
	}
}
