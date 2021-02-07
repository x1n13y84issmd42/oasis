package params

import (
	"net/http"

	"github.com/x1n13y84issmd42/oasis/src/contract"
	"github.com/x1n13y84issmd42/oasis/src/errors"
)

// QueryParameters is the source for URL query parameters.
type QueryParameters struct {
	contract.EntityTrait
	*MultiSet
}

// Query creates a new QueryParameters instance.
func Query(log contract.Logger) *QueryParameters {
	p := &QueryParameters{
		EntityTrait: contract.Entity(log),
		MultiSet:    NewMultiSet("query"),
	}

	return p
}

// Enrich applies the parameters as query values to the request.
func (params QueryParameters) Enrich(req *http.Request, log contract.Logger) {
	if err := params.Validate(); err != nil {
		errors.Report(err, "QueryParameters", params.Log)
	}

	q := req.URL.Query()

	for p := range params.Iterate() {
		v := p.V()
		log.UsingParameterExample(p.N, "query", p.Source, v)
		q.Add(p.N, v)
	}

	req.URL.RawQuery = q.Encode()
}
