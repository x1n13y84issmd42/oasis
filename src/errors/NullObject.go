package errors

import (
	"github.com/x1n13y84issmd42/oasis/src/contract"
)

// NullObjectPrototype is a prototype implementation of null entities,
// such as NullSpec and NullOperation.
type NullObjectPrototype struct {
	Log   contract.Logger
	Error error
}

// NullObject creates a new NullONullObjectPrototype instance.
func NullObject(err error, log contract.Logger) NullObjectPrototype {
	return NullObjectPrototype{
		Error: err,
		Log:   log,
	}
}

// Report reports the contained error.
func (null NullObjectPrototype) Report() {
	Report(null.Error, "NullObject", null.Log)
}
