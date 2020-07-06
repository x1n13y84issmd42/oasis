package params

import (
	"github.com/x1n13y84issmd42/oasis/src/contract"
)

// ParameterReference is a special kind of parameter which comes from
// an operation response. When a parameter for some operation in a script
// has a value like "operationID.response.[0].user.id" this means that
// the actual value comes from JSON response of the operation "operationID"
// and it's exact location is "[0].user.id" field.
type ParameterReference struct {
	Operation contract.Operation
	Selector  string
}

// ParameterReferenceMap is a multimap of ParameterReferences.
type ParameterReferenceMap map[string][]ParameterReference

// ReferenceSource is a parameter source backed by a response of other operations.
// When you want to use some of response data as a input for an operation test,
// yu add references here then Load() them into a parameter set.
type ReferenceSource struct {
	Name string
	Refs ParameterReferenceMap
}

// NewReferenceSource creates a new ReferenceSource instance.
func NewReferenceSource(name string) *ReferenceSource {
	return &ReferenceSource{
		Name: name,
		Refs: make(ParameterReferenceMap),
	}
}

// AddReference returns a parameter by it's name.
func (src *ReferenceSource) AddReference(pn string, op contract.Operation, selector string) {
	if src.Refs[pn] == nil {
		src.Refs[pn] = []ParameterReference{}
	}

	src.Refs[pn] = append(src.Refs[pn], ParameterReference{Operation: op, Selector: selector})
}

// Get returns a parameter by it's name.
func (src *ReferenceSource) Get(pn string) string {
	return ""
}

// Iterate creates an iterable channel.
func (src *ReferenceSource) Iterate() contract.ParameterIterator {
	ch := make(contract.ParameterIterator)

	go func() {

		close(ch)
	}()

	return ch
}
