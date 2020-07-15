package api

import (
	"net/http"

	"github.com/x1n13y84issmd42/oasis/src/contract"
	"github.com/x1n13y84issmd42/oasis/src/errors"
)

// NullSecurity is used whenever we can't have a real one.
// Reports the contained error on every method call.
type NullSecurity struct {
	errors.NullObjectPrototype
}

// NoSecurity creates a new NullSecurity instance.
func NoSecurity(err error, log contract.Logger) contract.Security {
	return &NullSecurity{
		NullObjectPrototype: errors.NullObject(err, log),
	}
}

//GetName reports an error.
func (sec *NullSecurity) GetName() string {
	sec.Report()
	return ""
}

//Enrich reports an error.
func (sec *NullSecurity) Enrich(req *http.Request, log contract.Logger) {
	sec.Report()
}
