package contract

import (
	"errors"
	"os"
)

// NullObjectPrototype is a prototype implementation of null entities,
// such as NullSpec and NullOperation.
type NullObjectPrototype struct {
	Log   Logger
	Error error
}

// NullObject creates a new NullONullObjectPrototype instance.
func NullObject(err error, log Logger) NullObjectPrototype {
	return NullObjectPrototype{
		Error: err,
		Log:   log,
	}
}

// Report reports the contained error.
func (spec NullObjectPrototype) Report() {
	if spec.Error != nil {
		spec.Log.Error(spec.Error)
		// panic("AAAAaaaa")
		os.Exit(128)
	} else {
		spec.Log.Error(errors.New("no error in a null object"))
	}
}
