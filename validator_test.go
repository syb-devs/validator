package validator_test

import (
	"bitbucket.org/simplifyourbusiness/validator"
	"testing"
)

func TestMinLength(t *testing.T) {

	// Case: valid struct
	type data struct {
		FieldOne string `validation:"min_length:3" `
	}

	inputErrs, err := validator.Validate(
		&data{FieldOne: "foo"})
	if inputErrs.Count() > 0 {
		t.Errorf("Input Errors: %s", inputErrs)
	}
	if err != nil {
		t.Errorf("Error during validation: %s", err.Error())
	}

	// Case: min lenth not satisfied
	type dataTwo struct {
		FieldTwo string `validation:"min_length:4" `
	}

	inputErrs, err = validator.Validate(
		&dataTwo{FieldTwo: "foo"})
	if inputErrs.Count() != 1 {
		t.Errorf("Expected exactly 1 input error, got %d", inputErrs.Count())
	}

	inputErr := inputErrs[0]
	expectedMessage := "The field FieldTwo should have a minimum length of 4 characters"
	actualMessage := inputErr.Message()
	if actualMessage != expectedMessage {
		t.Errorf("Expecting error message to be %s, got %s", expectedMessage, actualMessage)
	}

	if err != nil {
		t.Errorf("Error during validation: %s", err.Error())
	}

	// Case: data is not ptr to struct
	type dataThree struct {
		fieldTwo string `validation:"min_length:4" `
	}

	_, err = validator.Validate(dataThree{})

	if err == nil {
		t.Errorf("Expecting error because data is not a pointer to struct")
	}

	if err != validator.ErrStructPointerExpected {
		t.Errorf("Expecting error message to be %s, got %s", validator.ErrStructPointerExpected, err)
	}
}
