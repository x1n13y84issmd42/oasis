package errors

import (
	"github.com/x1n13y84issmd42/oasis/src/strings"
)

// Base is a generic error used within Oasis.
// Having it's root cause is what makes it good.
type Base struct {
	TheCause error
	Details  string
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

// ErrOops is a generic error with a message.
type ErrOops struct {
	Base
}

func (err ErrOops) Error() string {
	return err.Details
}

// Oops creates a new ErrOops error instance.
func Oops(msg string, cause error) ErrOops {
	return ErrOops{
		Base: Base{
			TheCause: cause,
			Details:  msg,
		},
	}
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
			Details:  details,
		},
		OpID: id,
	}
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

// NoParameters creates a new ErrNoData error instance.
func NoParameters(missingParams []string, cause error) ErrNoParameters {
	return ErrNoParameters{
		Base: Base{
			TheCause: cause,
		},
		MissingParams: missingParams,
	}
}

// ErrNotFound may be returned from host-related methods of specs.
type ErrNotFound struct {
	Base
	What string
	Name string
}

func (err ErrNotFound) Error() string {
	return err.What + " '" + err.Name + "' not found in the spec."
}

// NotFound creates a new ErrHostNotFound error instance.
func NotFound(what string, hn string, cause error) ErrNotFound {
	return ErrNotFound{
		Base: Base{TheCause: cause},
		What: what,
		Name: hn,
	}
}

// ErrInvalidSchema may be returned from schema-related methods of specs.
type ErrInvalidSchema struct {
	Base
	SchemaName string
}

func (err ErrInvalidSchema) Error() string {
	return "Schema " + err.SchemaName + " is invalid. " + err.Details
}

// InvalidSchema creates a new ErrInvalidSchema error instance.
func InvalidSchema(sn string, details string, cause error) ErrInvalidSchema {
	return ErrInvalidSchema{
		Base: Base{
			TheCause: cause,
			Details:  details,
		},
		SchemaName: sn,
	}
}

// ErrInvalidResponse may be returned from response-related methods of specs.
type ErrInvalidResponse struct {
	Base
}

func (err ErrInvalidResponse) Error() string {
	return "Response spec is invalid. " + err.Details
}

// InvalidResponse creates a new ErrInvalidResponse error instance.
func InvalidResponse(details string, cause error) ErrInvalidResponse {
	return ErrInvalidResponse{
		Base: Base{
			TheCause: cause,
			Details:  details,
		},
	}
}

// ErrSecurityNotFound may be returned from security-related methods of specs.
type ErrSecurityNotFound struct {
	Base
	Name string
}

func (err ErrSecurityNotFound) Error() string {
	return "Security definition '" + err.Name + "' not found in the spec."
}

// SecurityNotFound creates a new ErrSecurityNotFound error instance.
func SecurityNotFound(sn string, details string, cause error) ErrSecurityNotFound {
	return ErrSecurityNotFound{
		Base: Base{
			TheCause: cause,
			Details:  details,
		},
		Name: sn,
	}
}
