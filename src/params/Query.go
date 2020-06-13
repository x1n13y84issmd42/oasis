package params

import (
	"net/http"

	"github.com/x1n13y84issmd42/oasis/src/contract"
)

// QueryParameters is the default source for URL query parameters.
// It treats it's sources hierarchically, i.e. values are taken
// in this order:
// 		test output
// 		CLI input
// 		spec op
// 		spec path
type QueryParameters struct {
	*contract.ParameterProviderPrototype
}

// Query creates a new QueryParameters instance.
func Query() *QueryParameters {
	return &QueryParameters{
		ParameterProviderPrototype: contract.NewParameterProviderPrototype(),
	}
}

// Apply applies the parameters to the req's Query.
func (params QueryParameters) Apply(req *http.Request) {
	if err := params.Validate(); err != nil {

	}
}
