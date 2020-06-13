package contract

import "errors"

// NullObjectPrototype is a prototype implementation of null entities,
// such as NullSpec and NullOperation.
type NullObjectPrototype struct {
	Log   Logger
	Error error
}

// Report reports the contained error.
func (spec NullObjectPrototype) Report() {
	if spec.Error != nil {
		spec.Log.Error(spec.Error)
	} else {
		spec.Log.Error(errors.New("no error in a null object"))
	}
}
