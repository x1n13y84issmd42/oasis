package contract

// Set is an interface to a set of named values used as input parameters for an operation.
// Operation parameters come from various sources at various stages of testing, some of them are required,
// some are optional.
type Set interface {
	Load(src ParameterSource)
	Require(paramName string)
	Validate() error
	Iterate() ParameterIterator
}

// StringParameters represent all the available parameters as a string value.
// At the moment it is used for building a URL.
type StringParameters interface {
	Set
	String() string
}

// RequestEnrichmentParameters is used to enrich http.Request instances
// with parameters. Ued for headers & query values.
type RequestEnrichmentParameters interface {
	Set
	RequestEnrichment
}

// ParameterTuple is a pair of parameter name and it's value.
type ParameterTuple struct {
	N string
	V string
}

// ParameterIterator is an iterable channel to receive tuples
// of parameter name & value.
type ParameterIterator chan ParameterTuple

// ParameterSource is an interface for a parameter source.
// Parameters may come from various places, such as API specs,
// CLI arguments & another test output.
type ParameterSource interface {
	Iterate() ParameterIterator
}
