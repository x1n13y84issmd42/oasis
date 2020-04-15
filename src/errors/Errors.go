package errors

import "github.com/x1n13y84issmd42/oasis/src/strings"

// Base is a generic error used within Oasis.
// Having it's root cause is what makes it good.
type Base struct {
	TheCause error
}

// IError is an error interface.
type IError interface {
	Error() string
	Cause() error
}

// Error return a string representation of an error.
func (err Base) Error() string {
	return "---"
}

// Cause return an cause error instance.
func (err Base) Cause() error {
	return err.TheCause
}

// ErrNoParameters happens when there are parameters missing
// which are needed to build an operation before testing.
// Usually happens in paths, queries & headers.
type ErrNoParameters struct {
	Base
	MissingParams []string
}

func (err ErrNoParameters) Error() string {
	if len(err.MissingParams) == 1 {
		return "Parameter '" + err.MissingParams[0] + "' is missing it's value."
	}

	missing := strings.Map(err.MissingParams, func(s string) string {
		return "'" + s + "'"
	})

	return "Parameters " + strings.EnumJoin(missing) + " are missing their values."
}

// NoParameters creates a new ErrNoParameters error instance.
func NoParameters(missingParams []string, cause error) ErrNoParameters {
	return ErrNoParameters{
		Base: Base{
			TheCause: cause,
		},
		MissingParams: missingParams,
	}
}
