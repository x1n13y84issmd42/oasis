package params

import (
	"github.com/x1n13y84issmd42/oasis/src/contract"
)

// ReferenceMap is a multimap of ParameterReferences.
type ReferenceMap map[string][]Reference

// ReferenceSource is a parameter source backed by a response of other operations.
// When you want to use some of response data as a input for an operation test,
// you add references here then Load() them into a target parameter set.
type ReferenceSource struct {
	contract.EntityTrait
	Refs ReferenceMap
}

// NewReferenceSource creates a new ReferenceSource instance.
func NewReferenceSource(log contract.Logger) *ReferenceSource {
	return &ReferenceSource{
		EntityTrait: contract.Entity(log),
		Refs:        make(ReferenceMap),
	}
}

// AddReference adds a reference to a parameter value, located in response of op.
func (src *ReferenceSource) AddReference(pn string, opID string, result *contract.OperationResult, selector string) {
	if src.Refs[pn] == nil {
		src.Refs[pn] = []Reference{}
	}

	src.Refs[pn] = append(src.Refs[pn], Reference{
		OpID:     opID,
		Result:   result,
		Selector: selector,
		Log:      src.Log,
	})
}

// Iterate creates an iterable channel.
func (src *ReferenceSource) Iterate() contract.ParameterIterator {
	ch := make(contract.ParameterIterator)

	go func() {
		for pn, pvs := range src.Refs {
			for _, pv := range pvs {
				ch <- contract.ParameterTuple{
					N: pn,
					Parameter: contract.Parameter{
						V:      pv.Value(),
						Source: pv.OpID,
					},
				}
			}
		}
		close(ch)
	}()

	return ch
}
