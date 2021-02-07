package contract

// Set is an interface to a set of named values used as input parameters for an operation.
// Operation parameters come from various sources at various stages of testing, some of them are required,
// some are optional.
type Set interface {
	Load(src ParameterSource)
	Reload()
	Require(paramName string)
	Validate() error
	Iterate() ParameterIterator
}

// StringParameters represent all the available parameters as a string value.
// At the moment it is used for building a URL.
type StringParameters interface {
	IEntityTrait
	Set
	String() string
}

// RequestEnrichmentParameters is used to enrich http.Request instances
// with parameters. Used for headers & query values.
type RequestEnrichmentParameters interface {
	Set
	RequestEnrichment
}

// ParameterAccess is a function to provide a value for an operation.
// Values can be either literal or references. Reference values are used
// in scripts and look like "#opID.response[JSON_SELECTOR]".
// They are used to get certain values from one operation response
// as an input to another operation.
type ParameterAccess func() string

// Parameter is a pair of parameter value and name of it's source.
type Parameter struct {
	V      ParameterAccess
	Source string
}

// ParameterTuple is a pair of parameter name and it's value.
type ParameterTuple struct {
	Parameter
	N string
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
