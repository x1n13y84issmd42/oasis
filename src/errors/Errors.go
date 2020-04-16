package errors

import (
	"github.com/x1n13y84issmd42/oasis/src/strings"
)

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

// ErrOperationMalformed is returned when there is not enough data in the API spec,
// or the data is incorrect, which makes operation creation impossible.
type ErrOperationMalformed struct {
	Base
	OpID    string
	Details string
}

func (err ErrOperationMalformed) Error() string {
	return "Operation '" + err.OpID + "' has malformed or incomplete data. " + err.Details
}

// OperationMalformed creates a new ErrOperationMalformed error instance.
func OperationMalformed(id string, details string, cause error) ErrOperationMalformed {
	return ErrOperationMalformed{
		Base: Base{
			TheCause: cause,
		},
		OpID:    id,
		Details: details,
	}
}

// ErrOperationNotFound is returned from Spec.GetOperation() in case
// there is no operation in the spec with requested name/id.
type ErrOperationNotFound struct {
	Base
	OpID string
}

func (err ErrOperationNotFound) Error() string {
	return "Operation '" + err.OpID + "' not found."
}

// OperationNotFound creates a new ErrOperationNotFound error instance.
func OperationNotFound(id string, cause error) ErrOperationNotFound {
	return ErrOperationNotFound{
		Base: Base{
			TheCause: cause,
		},
		OpID: id,
	}
}

// ErrNoData happens when there are parameters missing
// which are needed to build an operation before testing.
// Usually happens in paths, queries & headers.
type ErrNoData struct {
	Base
	MissingParams []string
}

func (err ErrNoData) Error() string {
	if len(err.MissingParams) == 1 {
		return "Parameter '" + err.MissingParams[0] + "' is missing it's value."
	}

	missing := strings.Map(err.MissingParams, func(s string) string {
		return "'" + s + "'"
	})

	return "Parameters " + strings.EnumJoin(missing) + " are missing their values."
}

// NoData creates a new ErrNoData error instance.
func NoData(missingParams []string, cause error) ErrNoData {
	return ErrNoData{
		Base: Base{
			TheCause: cause,
		},
		MissingParams: missingParams,
	}
}

// ErrHostNotFound may be returned from host-related methods of specs.
type ErrHostNotFound struct {
	Base
	HostName string
}

func (err ErrHostNotFound) Error() string {
	return "Host '" + err.HostName + "' not found in the spec."
}

// HostNotFound creates a new ErrHostNotFound error instance.
func HostNotFound(hn string, cause error) ErrHostNotFound {
	return ErrHostNotFound{
		Base:     Base{TheCause: cause},
		HostName: hn,
	}
}
