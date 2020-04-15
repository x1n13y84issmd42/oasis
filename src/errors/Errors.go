package errors

import "github.com/x1n13y84issmd42/oasis/src/strings"

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
		return "Parameter '" + err.MissingParams[0] + "' is missing it's value."
	}

	missing := strings.Map(err.MissingParams, func(s string) string {
		return "'" + s + "'"
	})

	return "Parameters " + strings.EnumJoin(missing) + " are missing their values."
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
