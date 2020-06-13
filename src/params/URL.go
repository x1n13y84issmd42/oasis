package params

import "github.com/x1n13y84issmd42/oasis/src/contract"

// URLParameters is the default source for URL path parameters.
// It treats it's sources hierarchically, i.e. values are taken
// in this order:
// 		test output
// 		CLI input
// 		spec op
// 		spec path
type URLParameters struct {
	*contract.ParameterProviderPrototype
}

// URL creates a new URLParameters instance.
func URL() *URLParameters {
	return &URLParameters{
		ParameterProviderPrototype: contract.NewParameterProviderPrototype(),
	}
}

// Make creates a URL string value.
func (params URLParameters) Make() string {
	if err := params.Validate(); err != nil {

	}

	return ""
}
