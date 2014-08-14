package validator_test

import (
	"bitbucket.org/simplifyourbusiness/validator"
	"testing"
)

type User struct {
	Email string `validation:"min_length:1" `
}

func TestMinLength(t *testing.T) {
	validator.Validate(&User{Email: "Manolo"})
}