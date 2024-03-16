package usecase

import (
	"errors"
	"strings"
)

// ValidationErrors is an implementation of error that contains a map of validation errors
type ValidationErrors struct {
	// fieldErrors contains a map of field name, to the errors on that field
	fieldErrors map[string][]error
}

func (e ValidationErrors) Error() string {
	var errMsg []string
	for _, errs := range e.fieldErrors {
		for _, err := range errs {
			errMsg = append(errMsg, err.Error())
		}
	}
	return strings.Join(errMsg, " | ")
}

func (e ValidationErrors) GetErrors() map[string][]error {
	return e.fieldErrors
}

func NewValidationError(errors map[string][]error) error {
	return ValidationErrors{fieldErrors: errors}
}

var (
	UserInvalidLogin  = errors.New("invalid phone number or password")
	UserInvalidToken  = errors.New("invalid / expired token, please login again")
	UserNotFoundError = errors.New("user not found")
	UserConflictError = errors.New("user record conflict, phone number must be unique")
)
