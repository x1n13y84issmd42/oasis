package params

import (
	"fmt"
	"sort"

	"github.com/x1n13y84issmd42/oasis/src/contract"
	"github.com/x1n13y84issmd42/oasis/src/errors"
)

// Map is a map of parameters.
type Map map[string][]contract.Parameter

// MultiSet is a set of named values used as input parameters for an operation.
// Each key can have multiple values.
type MultiSet struct {
	Name     string
	data     Map
	required []string
}

// NewMultiSet creates a new MultiSet instance.
func NewMultiSet(name string) *MultiSet {
	return &MultiSet{
		Name:     name,
		data:     make(Map),
		required: []string{},
	}
}

// Load reads parameters from a source.
func (params *MultiSet) Load(src contract.ParameterSource) {
	for p := range src.Iterate() {
		params.data[p.N] = append(params.data[p.N], p.Parameter)
	}
}

// Require adds a parameter name to the requried parameters list.
func (params *MultiSet) Require(paramName string) {
	for _, pn := range params.required {
		if pn == paramName {
			return
		}
	}

	params.required = append(params.required, paramName)
}

// Validate checks that all required parameters have their values available.
func (params *MultiSet) Validate() error {
	missingParams := []string{}

	for _, rpn := range params.required {
		_, found := params.data[rpn]
		if !found {
			missingParams = append(missingParams, rpn)
		}
	}

	if len(missingParams) > 0 {
		return errors.NoParameters(missingParams, params.Name, nil)
	}

	return nil
}

// Iterate creates an iterable channel.
func (params *MultiSet) Iterate() contract.ParameterIterator {
	ch := make(contract.ParameterIterator)

	go func() {
		keys := []string{}
		for pn := range params.data {
			keys = append(keys, pn)
		}

		sort.Strings(keys)

		for _, k := range keys {
			for _, v := range params.data[k] {
				ch <- contract.ParameterTuple{N: k, Parameter: v}
			}
		}

		close(ch)
	}()

	return ch
}

func (s *MultiSet) Print() {
	fmt.Printf("Contents of %s\n", s.Name)
	for p := range s.Iterate() {
		fmt.Printf("%s = %s (from %s)\n", p.N, p.V(), p.Source)
	}
}
