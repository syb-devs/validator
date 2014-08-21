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

func TestNotStruct(t *testing.T) {
	v := validator.New()
	err := v.Validate("string")

	if err != validator.ErrStructExpected {
		t.Errorf("Expected: %s, got: %s", validator.ErrStructExpected, err)
	}
}

func TestEmbeddedStruct(t *testing.T) {
	type embed struct{}
	type data struct {
		field embed
	}

	v := validator.New()
	err := v.Validate(data{field: embed{}})
	if err != nil {
		t.Error(err.Error())
	}
}

type jar struct{}

var isStructTests = []struct {
	data     interface{} // input
	expected bool        // expected result
}{
	{"gopher", false},
	{1845, false},
	{jar{}, true},
	{&jar{}, true},
}

func TestIsStruct(t *testing.T) {
	for _, test := range isStructTests {
		actual := validator.IsStruct(test.data)
		if test.expected != actual {
			t.Errorf("IsStruct(%v): expected %v, actual %v", test.data, test.expected, actual)
		}
	}
}
