package api

import "github.com/x1n13y84issmd42/oasis/src/contract"

// NullSpec ...
type NullSpec struct {
	Log   contract.Logger
	Error error
}

func (spec NullSpec) report() {
	spec.Log.Error(spec.Error)
}

// Operations ...
func (spec NullSpec) Operations() contract.OperationIterator {
	spec.report()
	return nil
}

// GetOperation ...
func (spec NullSpec) GetOperation(name string) contract.Operation {
	spec.report()
	return nil
}

// Title return project title.
func (spec NullSpec) Title() string {
	spec.report()
	return ""
}

// Description return project description.
func (spec NullSpec) Description() string {
	spec.report()
	return ""
}

// Version return project version.
func (spec NullSpec) Version() string {
	spec.report()
	return ""
}
