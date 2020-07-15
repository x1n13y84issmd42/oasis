package api

import (
	"github.com/x1n13y84issmd42/oasis/src/contract"
	"github.com/x1n13y84issmd42/oasis/src/errors"
)

// NullSpec ...
type NullSpec struct {
	errors.NullObjectPrototype
}

// NoSpec creates a new NullSpec instance.
func NoSpec(err error, log contract.Logger) NullSpec {
	return NullSpec{
		NullObjectPrototype: errors.NullObject(err, log),
	}
}

// Operations ...
func (spec NullSpec) Operations() contract.OperationIterator {
	spec.Report()
	return nil
}

// GetOperation ...
func (spec NullSpec) GetOperation(name string) contract.Operation {
	spec.Report()
	return nil
}

// Title return project title.
func (spec NullSpec) Title() string {
	spec.Report()
	return ""
}

// Description return project description.
func (spec NullSpec) Description() string {
	spec.Report()
	return ""
}

// Version return project version.
func (spec NullSpec) Version() string {
	spec.Report()
	return ""
}
