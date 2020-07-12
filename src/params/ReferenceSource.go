package params

import (
	"fmt"

	"github.com/x1n13y84issmd42/oasis/src/contract"
)

// ParameterReference is a special kind of parameter which comes from
// an operation response. When a parameter for some operation in a script
// has a value like "operationID.response.[0].user.id" this means that
// the actual value comes from JSON response of the operation "operationID"
// and it's exact location is "[0].user.id" field.
type ParameterReference struct {
	OpID     string
	Result   *contract.OperationResult
	Selector string
}

// Value returns a parameter access function which computes and returns a real value.
func (pr ParameterReference) Value() contract.ParameterAccess {
	//TODO: implement me
	return func() string { return "10" }
}

// ParameterReferenceMap is a multimap of ParameterReferences.
type ParameterReferenceMap map[string][]ParameterReference

// ReferenceSource is a parameter source backed by a response of other operations.
// When you want to use some of response data as a input for an operation test,
// yu add references here then Load() them into a parameter set.
type ReferenceSource struct {
	Refs ParameterReferenceMap
}

// NewReferenceSource creates a new ReferenceSource instance.
func NewReferenceSource() *ReferenceSource {
	return &ReferenceSource{
		Refs: make(ParameterReferenceMap),
	}
}

// AddReference adds a reference to a parameter value, located in response of op.
func (src *ReferenceSource) AddReference(pn string, opID string, result *contract.OperationResult, selector string) {
	if src.Refs[pn] == nil {
		src.Refs[pn] = []ParameterReference{}
	}

	src.Refs[pn] = append(src.Refs[pn], ParameterReference{OpID: opID, Result: result, Selector: selector})
}

// Get returns a parameter by it's name.
func (src *ReferenceSource) Get(pn string) string {
	return ""
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

func (src *ReferenceSource) Print() {
	fmt.Printf("Contents of a ref source\n")
	for p := range src.Iterate() {
		fmt.Printf("%s = %s (from %s)\n", p.N, p.V(), p.Source)
	}
}
