package params

import (
	"sort"

	"github.com/x1n13y84issmd42/oasis/src/contract"
)

// MemorySource is a parameter source which uses a native map as a source storage.
type MemorySource struct {
	Data map[string]string
}

// NewMemorySource creates a new MemoryParameterSource instance.
func NewMemorySource() *MemorySource {
	return &MemorySource{
		Data: map[string]string{},
	}
}

// Add stores a kv pair in the source.
func (ds *MemorySource) Add(n string, v string) {
	ds.Data[n] = v
}

// Get returns a parameter by it's name.
func (ds *MemorySource) Get(pn string) string {
	return ds.Data[pn]
}

// Iterate returns an iterable channel to receive parameter tuples.
func (ds *MemorySource) Iterate() contract.ParameterIterator {
	ch := make(contract.ParameterIterator)

	go func() {
		keys := []string{}
		for pn := range ds.Data {
			keys = append(keys, pn)
		}

		sort.Strings(keys)

		for _, pn := range keys {
			ch <- contract.ParameterTuple{N: pn, V: ds.Data[pn]}
		}

		close(ch)
	}()

	return ch
}
