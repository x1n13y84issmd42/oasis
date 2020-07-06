package params

import (
	"net/http"

	"github.com/x1n13y84issmd42/oasis/src/contract"
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
		params.Error(err)
	}

	q := req.URL.Query()

	for pt := range params.Iterate() {
		q.Add(pt.N, pt.V())
	}

	req.URL.RawQuery = q.Encode()
}
