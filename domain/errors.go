package domain

import "errors"

var (
	// ErrInternalServerError will throw if any the Internal Server Error happen
	ErrInternalServerError = errors.New("Internal Server Error")
	// ErrNotFound will throw if the requested items is not exists
	ErrNotFound = errors.New("Your requested Items is Not Found")
	// ErrBalanceRequired will throw if the balance suuficient
	ErrBalanceRequired = errors.New("Your balance is not sufficient")
	// ErrConflict will throw if the current action already exists
	ErrConflict = errors.New("Your item already exists")
	// ErrBadParamInput will throw if the given request-body or params is not valid
	ErrBadParamInput = errors.New("Given param is not valid")
	// ErrNotAccordingSpecifications will throws if current item not according specifications
	ErrNotAccordingSpecifications = errors.New("Item not according to specifications")
)
