package contract

import (
	"net/http"

	"github.com/x1n13y84issmd42/oasis/src/errors"
)

// DataProvider is an interface to access data from various sources
// (spec path, spec op, cli input, test output) needed in order to build
// an http.Request instance.
type DataProvider struct {
	URL     ParameterMaker
	Query   ParameterApplicator
	Headers ParameterApplicator
}

// ParameterIterator is an iterable channel to receive tuples
// of parameter name & parameter value.
type ParameterIterator chan [2]string

// ParameterSource is an interface for parameter retrieval.
// Parameters may come from various places, such as API specs,
// CLI arguments & another test output.
type ParameterSource interface {
	Get(paramName string) string
	Iterate() ParameterIterator
}

// ParameterSourceType denotes a type of parameter provider.
// Within a parameter provider sources are stored in a map
// with the source type as key.
type ParameterSourceType string

// Parameter source types.
const (
	ParameterSourceSpecPath   = ParameterSourceType("spec_path")
	ParameterSourceSpecOp     = ParameterSourceType("spec_op")
	ParameterSourceCLI        = ParameterSourceType("cli")
	ParameterSourceTestOutput = ParameterSourceType("test")
)

// ParameterSourceMap is a map of parameter sources.
type ParameterSourceMap map[ParameterSourceType]ParameterSource

// ParameterProvider is an interface for parameter providers.
// Parameter provider cares about parameters needed for making a request,
// such as URL variables, query values & headers.
type ParameterProvider interface {
	AddSource(srcType ParameterSourceType, src ParameterSource)
	Require(paramName string)
	Validate() error
}

// ParameterApplicator is a parameter provider which applies it's data directly to a http.Request instance.
type ParameterApplicator interface {
	ParameterProvider
	Apply(req *http.Request)
}

// ParameterMaker is a parameter provider which provides it's data as a string value.
type ParameterMaker interface {
	ParameterProvider
	Make() string
}

// ParameterProviderPrototype is a prototype implementation for parameter providers.
type ParameterProviderPrototype struct {
	ParameterProvider
	Required []string
	Sources  ParameterSourceMap
}

// NewParameterProviderPrototype creates a new ParameterProviderPrototype instance.
func NewParameterProviderPrototype() *ParameterProviderPrototype {
	return &ParameterProviderPrototype{
		Required: []string{},
		Sources:  make(ParameterSourceMap),
	}
}

// AddSource adds a parameter source to the provider.
func (provider *ParameterProviderPrototype) AddSource(srcType ParameterSourceType, src ParameterSource) {
	provider.Sources[srcType] = src
}

// Require adds a parameter name to the requried parameters list.
func (provider *ParameterProviderPrototype) Require(paramName string) {
	provider.Required = append(provider.Required, paramName)
}

// Validate checks that all required parameters have their values available.
func (provider *ParameterProviderPrototype) Validate() error {
	missingParams := []string{}

	for _, rpn := range provider.Required {
		value := ""
		for _, src := range provider.Sources {
			if v := src.Get(rpn); v != "" {
				value = v
				break
			}
		}

		if value == "" {
			missingParams = append(missingParams, rpn)
		}
	}

	if len(missingParams) > 0 {
		return errors.NoParameters(missingParams, nil)
	}

	return nil
}
