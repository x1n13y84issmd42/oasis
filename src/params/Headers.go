package params

import (
	"net/http"

	"github.com/x1n13y84issmd42/oasis/src/contract"
)

// HeadersParameters is the source for URL path parameters.
// HeadersParameters have an implicit requirement for the @HOSTNAME parameter
// which is an API host name.
type HeadersParameters struct {
	contract.EntityTrait
	*Parameters
}

// Headers creates a new HeadersParameters instance.
func Headers(log contract.Logger) *HeadersParameters {
	p := &HeadersParameters{
		EntityTrait: contract.Entity(log),
		Parameters:  New(),
	}

	return p
}

// Enrich applies the parameters as Headers values to the request..
func (params HeadersParameters) Enrich(req *http.Request) {
	if err := params.Validate(); err != nil {
		params.Error(err)
	}

	for pt := range params.Iterate() {
		req.Header.Set(pt.N, pt.V)
	}
}
