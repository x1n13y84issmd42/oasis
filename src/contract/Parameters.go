package contract

import (
	"net/http"
	"sort"

	"github.com/x1n13y84issmd42/oasis/src/errors"
)

// DataProvider is an interface to access data from various sources
// (spec path, spec op, cli input, test output) needed in order to build
// an http.Request instance.
type DataProvider struct {
	URL     ParameterMaker
	Query   ParameterEnrichment
	Headers ParameterEnrichment
}

// ParameterTuple is a pair of parameter name and it's value.
type ParameterTuple [2]string

// ParameterIterator is an iterable channel to receive tuples
// of parameter name & parameter value.
type ParameterIterator chan ParameterTuple

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
	Get(pn string) string
	Iterate() ParameterIterator
}

// ParameterEnrichment is a parameter provider which applies it's data directly to a http.Request instance.
type ParameterEnrichment interface {
	ParameterProvider
	Enrich(req *http.Request)
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

// NewPOCTParameterProvider creates a new POCTParameterProvider instance.
func NewPOCTParameterProvider() *POCTParameterProvider {
	return &POCTParameterProvider{
		ParameterProviderPrototype: NewParameterProviderPrototype(),
	}
}

// AddSource adds a parameter source to the provider. Each type of source may appear only once.
func (provider *ParameterProviderPrototype) AddSource(srcType ParameterSourceType, src ParameterSource) {
	provider.Sources[srcType] = src
}

// Require adds a parameter name to the requried parameters list.
func (provider *ParameterProviderPrototype) Require(paramName string) {
	for _, pn := range provider.Required {
		if pn == paramName {
			return
		}
	}

	provider.Required = append(provider.Required, paramName)
}

// Validate checks that all required parameters have their values available.
func (provider *ParameterProviderPrototype) Validate() error {
	missingParams := []string{}

	for _, rpn := range provider.Required {
		found := false
		for _, src := range provider.Sources {
			if v := src.Get(rpn); v != "" {
				found = true
				break
			}
		}

		if !found {
			missingParams = append(missingParams, rpn)
		}
	}

	if len(missingParams) > 0 {
		return errors.NoParameters(missingParams, nil)
	}

	return nil
}

// POCTParameterProvider returns parameters from the sources,
// and takes them in the specified hierarchical order:
//		Test output
//		CLI input
//		spec Operation
//		spec Path
// hence the POCT. This means that if 'foo' is requested,
// the provider first tries to get it from the 'Test output' source,
// then from 'CLI' one, and so on.
type POCTParameterProvider struct {
	*ParameterProviderPrototype
}

// SourceFn is a handler function used in the EachSource method.
// The loop continues unless it returns false.
type SourceFn func(ParameterSource) bool

// Get retrieves a parameter value by it's name.
func (provider *POCTParameterProvider) Get(pn string) (res string) {
	provider.EachSource(func(src ParameterSource) bool {
		if pv := src.Get(pn); pv != "" {
			res = pv
			return false
		}
		return true
	})

	return
}

// Iterate returns an iterable channel to iterate over parameters
// in lexicographic order.
func (provider *POCTParameterProvider) Iterate() ParameterIterator {
	ch := make(ParameterIterator)
	go func() {
		keys := []string{}
		values := provider.Values()
		for pn := range values {
			keys = append(keys, pn)
		}

		sort.Strings(keys)

		for _, pn := range keys {
			ch <- ParameterTuple{pn, values[pn]}
		}

		close(ch)
	}()
	return ch
}

// Values returns a map of top layer values from the sources.
func (provider *POCTParameterProvider) Values() (res map[string]string) {
	res = map[string]string{}

	provider.EachSource(func(src ParameterSource) bool {
		for pt := range src.Iterate() {
			if _, exists := res[pt[0]]; !exists {
				res[pt[0]] = pt[1]
			}
		}
		return true
	})
	return
}

// EachSource iterates over the available sources in the predefined order,
// which is
//		spec Path
//		spec Operation
//		CLI input
//		Test output
// and passes them to the provided handler function.
func (provider *POCTParameterProvider) EachSource(fn SourceFn) {
	for _, src := range []ParameterSource{
		provider.Sources[ParameterSourceTestOutput],
		provider.Sources[ParameterSourceCLI],
		provider.Sources[ParameterSourceSpecOp],
		provider.Sources[ParameterSourceSpecPath],
	} {
		if src != nil {
			if doContinue := fn(src); !doContinue {
				break
			}
		}
	}
}
