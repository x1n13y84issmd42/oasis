package params

import (
	"net/http"

	"github.com/x1n13y84issmd42/oasis/src/contract"
	"github.com/x1n13y84issmd42/oasis/src/errors"
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
		errors.Report(err, "HeadersParameters", params.Log)
	}

	// (*params.MultiSet)["Content-type"] = "application/json"

	for p := range params.Iterate() {
		v := p.V()
		params.Log.UsingParameterExample(p.N, "header", p.Source, v)
		req.Header.Add(p.N, v)
	}

	// req.Header.Add("Content-type", "application/json")
}
