package params

import (
	"net/http"

	"github.com/x1n13y84issmd42/oasis/src/contract"
)

// HeadersParameters is the source for request header parameters.
type HeadersParameters struct {
	contract.EntityTrait
	*MultiSet
}

// Headers creates a new HeadersParameters instance.
func Headers(log contract.Logger) *HeadersParameters {
	p := &HeadersParameters{
		EntityTrait: contract.Entity(log),
		MultiSet:    NewMultiSet("headers"),
	}

	return p
}

// Enrich applies the parameters as header values to the request.
func (params HeadersParameters) Enrich(req *http.Request, log contract.Logger) {
	if err := params.Validate(); err != nil {
		params.Error(err)
	}

	for pt := range params.Iterate() {
		req.Header.Add(pt.N, pt.V())
	}
}
