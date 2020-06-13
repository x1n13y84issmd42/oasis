package api

import "github.com/x1n13y84issmd42/oasis/src/contract"

// NullSpec ...
type NullSpec struct {
	contract.NullObjectPrototype

	Log   contract.Logger
	Error error
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
