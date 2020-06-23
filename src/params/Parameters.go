package params

import (
	"net/http"

	"github.com/x1n13y84issmd42/oasis/src/contract"
	"github.com/x1n13y84issmd42/oasis/src/errors"
)

// ParameterFn ...
type ParameterFn func() string

// ParameterMap ...
type ParameterMap map[string]string

// Parameters is a set of named values used as input parameters for an operation.
type Parameters struct {
	Data     ParameterMap
	required []string
}

// New creates a new Parameters instance.
func New() *Parameters {
	return &Parameters{
		Data:     make(ParameterMap),
		required: []string{},
	}
}

// Load reads parameters from a source.
func (params *Parameters) Load(src contract.ParameterSource) {
	for p := range src.Iterate() {
		params.Data[p.N] = p.V
	}
}

// Require adds a parameter name to the requried parameters list.
func (params *Parameters) Require(paramName string) {
	for _, pn := range params.required {
		if pn == paramName {
			return
		}
	}

	params.required = append(params.required, paramName)
}

// Validate checks that all required parameters have their values available.
func (params *Parameters) Validate() error {
	missingParams := []string{}

	for _, rpn := range params.required {
		_, found := params.Data[rpn]
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
func (params *Parameters) Iterate() contract.ParameterIterator {
	ch := make(contract.ParameterIterator)

	go func() {
		for pn, pv := range params.Data {
			ch <- contract.ParameterTuple{N: pn, V: pv}
		}
		close(ch)
	}()

	return ch
}

// String creates an string representation of the paralmeter set.
func (params *Parameters) String() string {
	return ""
}

// Enrich applies parameters to the request.
func (params *Parameters) Enrich(req *http.Request) {
}
