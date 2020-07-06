package params

import (
	"github.com/x1n13y84issmd42/oasis/src/contract"
)

// Set is a set of parameters used in operation testing.
// It is basically a MultiSet, the only difference being
// is that Load() overwrites each key with a [1]string
// instead of appending them.
type Set struct {
	*MultiSet
}

// NewSet creates a new MultiSet instance.
func NewSet(name string) *Set {
	return &Set{
		MultiSet: NewMultiSet(name),
	}
}

// Load reads parameters from a source.
func (params *Set) Load(src contract.ParameterSource) {
	for p := range src.Iterate() {
		params.data[p.N] = []string{p.V()}
	}
}
