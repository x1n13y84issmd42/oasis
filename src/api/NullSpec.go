package api

import (
	"github.com/x1n13y84issmd42/oasis/src/contract"
	"github.com/x1n13y84issmd42/oasis/src/errors"
)

// NullSpec is a placehilder spec used when we can't have a real one.
type NullSpec struct {
	errors.NullObjectPrototype
}

// NoSpec creates a new NullSpec instance.
func NoSpec(err error, log contract.Logger) NullSpec {
	return NullSpec{
		NullObjectPrototype: errors.NullObject(err, log),
	}
}

// Operations reports an error.
func (spec NullSpec) Operations() contract.OperationIterator {
	spec.Report()
	return nil
}

// GetOperation reports an error.
func (spec NullSpec) GetOperation(name string) contract.Operation {
	spec.Report()
	return nil
}

// Title reports an error.
func (spec NullSpec) Title() string {
	spec.Report()
	return ""
}

// Description reports an error.
func (spec NullSpec) Description() string {
	spec.Report()
	return ""
}

// Version reports an error.
func (spec NullSpec) Version() string {
	spec.Report()
	return ""
}
