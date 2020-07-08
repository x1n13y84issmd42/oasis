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
func (err NullObjectPrototype) Report() {
	if err.Error != nil {
		err.Log.Error(err.Error)
	} else {
		err.Log.Error(errors.New("no error in a null object"))
	}

	os.Exit(1)
	// panic("A NullObject has triggered a panic.\nSee the error message reported above.\nCheck where the object comes from\nfor the source of the error.")
}
