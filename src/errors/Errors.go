package errors

import "strings"

// Base is a generic error used within Oasis.
// It references it's root cause.
type Base struct {
	TheCause error
}

type IError interface {
	Error() string
	Cause() error
}

func (err Base) Error() string {
	return "---"
}

func (err Base) Cause() error {
	return err.TheCause
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
func NoParameters(missingParams []string, cause error) ErrNoParameters {
	return ErrNoParameters{
		Base: Base{
			TheCause: cause,
		},
		MissingParams: missingParams,
	}
}
