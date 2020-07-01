package params

import (
	"github.com/x1n13y84issmd42/oasis/src/contract"
)

// NullSource is a placeholder source used when we can't have a real one.
type NullSource struct {
	contract.NullObjectPrototype
}

// NoSource creates a new NullSource instance.
func NoSource(err error, log contract.Logger) *NullSource {
	return &NullSource{
		NullObjectPrototype: contract.NullObject(err, log),
	}
}

// Get reports an error.
func (ds *NullSource) Get(pn string) string {
	ds.Report()
	return ""
}

// Iterate reports an error.
func (ds *NullSource) Iterate() contract.ParameterIterator {
	ds.Report()
	return nil
}
