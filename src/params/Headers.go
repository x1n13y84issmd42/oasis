package params

import (
	"net/http"

	"github.com/x1n13y84issmd42/oasis/src/contract"
)

// HeadersParameters is the default source for URL Headers parameters.
// It treats it's sources hierarchically, i.e. values are taken
// in this order:
// 		test output
// 		CLI input
// 		spec op
// 		spec path
type HeadersParameters struct {
	*contract.ParameterProviderPrototype
}

// Headers creates a new HeadersParameters instance.
func Headers() *HeadersParameters {
	return &HeadersParameters{
		ParameterProviderPrototype: contract.NewParameterProviderPrototype(),
	}
}

// Enrich applies the parameters to the req's Headers.
func (params HeadersParameters) Enrich(req *http.Request) {
	if err := params.Validate(); err != nil {

	}
}
