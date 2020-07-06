package params

import (
	"sort"

	"github.com/x1n13y84issmd42/oasis/src/contract"
	"github.com/x1n13y84issmd42/oasis/src/errors"
)

// MultiSet is a set of named values used as input parameters for an operation.
// Each key can have multiple values.
type MultiSet struct {
	data     map[string][]string
	required []string
}

// NewMultiSet creates a new MultiSet instance.
func NewMultiSet() *MultiSet {
	return &MultiSet{
		data:     make(map[string][]string),
		required: []string{},
	}
}

// Load reads parameters from a source.
func (params *MultiSet) Load(src contract.ParameterSource) {
	for p := range src.Iterate() {
		params.data[p.N] = append(params.data[p.N], p.V)
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
		return errors.NoParameters(missingParams, nil)
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
				ch <- contract.ParameterTuple{N: k, V: v}
			}
		}

		close(ch)
	}()

	return ch
}
