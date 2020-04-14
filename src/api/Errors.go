package api

import (
	"github.com/x1n13y84issmd42/oasis/src/errors"
)

// ErrOperationMalformed is returned when there is not enough data in the API spec,
// or the data is incorrect, which makes operation creation impossible.
type ErrOperationMalformed struct {
	errors.Base
	OpID    string
	Details string
}

func (err ErrOperationMalformed) Error() string {
	return "Operation '" + err.OpID + "' has malformed or incomplete data. " + err.Details
}

// OperationMalformed creates a new ErrOperationMalformed error instance.
func OperationMalformed(id string, details string, cause errors.IError) ErrOperationMalformed {
	return ErrOperationMalformed{
		Base: errors.Base{
			TheCause: cause,
		},
		OpID:    id,
		Details: details,
	}
}

// ErrOperationNotFound is returned from Spec.GetOperation() in case
// there is no operation in the spec with requested name/id.
type ErrOperationNotFound struct {
	errors.Base
	OpID string
}

func (err ErrOperationNotFound) Error() string {
	return "Operation '" + err.OpID + "' not found."
}

// OperationNotFound creates a new ErrOperationNotFound error instance.
func OperationNotFound(id string, cause errors.IError) ErrOperationNotFound {
	return ErrOperationNotFound{
		Base: errors.Base{
			TheCause: cause,
		},
		OpID: id,
	}
}
