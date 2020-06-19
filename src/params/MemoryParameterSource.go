package params

import "github.com/x1n13y84issmd42/oasis/src/contract"

// MemoryParameterSource is a parameter source which uses a native map as a storage.
type MemoryParameterSource struct {
	Data map[string]string
}

// NewMemoryParameterSource creates a new MemoryParameterSource instance.
func NewMemoryParameterSource() *MemoryParameterSource {
	return &MemoryParameterSource{
		Data: map[string]string{},
	}
}

// Get returns a parameter by it's name.
func (ds *MemoryParameterSource) Get(pn string) string {
	return ds.Data[pn]
}

// Iterate returns an iterable channel to receive parameter tuples.
func (ds *MemoryParameterSource) Iterate() contract.ParameterIterator {
	ch := make(contract.ParameterIterator)

	go func() {
		for pn, pv := range ds.Data {
			ch <- contract.ParameterTuple{pn, pv}
		}

		close(ch)
	}()

	return ch
}
