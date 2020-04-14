package errors

import "strings"

// Base is a generic error used within Oasis.
// It references it's root cause.
type Base struct {
	TheCause IError
}

type IError interface {
	Error() string
	Cause() IError
}

func (err Base) Error() string {
	return "---"
}

func (err Base) Cause() IError {
	return err.TheCause
}

type WError struct {
	Base
	or error
}

func (err WError) Error() string {
	return err.or.Error()
}

func Wrap(err error) IError {
	return WError{
		or: err,
	}
}

// ErrNoParameters ...
type ErrNoParameters struct {
	Base
	MissingParams []string
}

func (err ErrNoParameters) Error() string {
	if len(err.MissingParams) == 1 {
		return "Parameter " + err.MissingParams[0] + " is missing it's value."
	}

	return "Parameters " + strings.Join(err.MissingParams, ", ") + " are missing their values."
}

// NoParameters ...
func NoParameters(missingParams []string, cause IError) ErrNoParameters {
	return ErrNoParameters{
		Base: Base{
			TheCause: cause,
		},
		MissingParams: missingParams,
	}
}
