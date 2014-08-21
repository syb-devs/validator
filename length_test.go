// Copyright 2014 Simplify Your Bussiness S.L. All rights reserved.

package validator_test

import (
	"bitbucket.org/simplifyourbusiness/validator"
	"strings"
	"testing"
)

func TestMinLengthOK(t *testing.T) {
	type data struct {
		Field string `validation:"length:op:>=,val:4" `
	}

	v := validator.New()
	err := v.Validate(data{Field: "fool"})

	if err != nil {
		t.Errorf(err.Error())
	}

	if v.Errors() != nil {
		t.Errorf("Unexpected validation error: %s", v.Errors())
	}
}

func TestMinLengthKO(t *testing.T) {
	type data struct {
		Field string `validation:"length:op:>,val:4" `
	}

	v := validator.New()
	err := v.Validate(data{})

	if err != nil {
		t.Errorf(err.Error())
	}

	if v.Errors() == nil || v.Errors().Len() != 1 {
		t.Errorf("Expecting exactly one validation error")
	}
}

func TestBadOperator(t *testing.T) {
	type data struct {
		Field string `validation:"length:op:ar,val:4" `
	}
	v := validator.New()
	err := v.Validate(data{})

	if err != nil {
		t.Errorf(err.Error())
	}

	if v.Errors() == nil || !strings.Contains(v.Errors().String(), "Invalid operator") {
		t.Errorf("Expected 'Invalid operator' error. Got : %s", v.Errors())
	}
}
