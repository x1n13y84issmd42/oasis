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
	contract.EntityTrait

	Path string
}

// URL creates a new URLParameters instance.
func URL(path string, log contract.Logger) *URLParameters {
	return &URLParameters{
		ParameterProviderPrototype: contract.NewParameterProviderPrototype(),
		EntityTrait:                contract.Entity(log),
		Path:                       path,
	}
}

// Make creates a URL string value.
func (params URLParameters) Make() string {
	if err := params.Validate(); err != nil {
		params.Log.Error(err)
	}

	return ""
}
