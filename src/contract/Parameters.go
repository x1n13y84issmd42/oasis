package contract

// Parameters is an interface to a set of named values used as input parameters for an operation.
type Parameters interface {
	Load(src ParameterSource)
	Require(paramName string)
	Validate() error
	Iterate() ParameterIterator
	String() string
}

// ParameterTuple is a pair of parameter name and it's value.
type ParameterTuple struct {
	N string
	V string
}

// ParameterIterator is an iterable channel to receive tuples
// of parameter name & parameter value.
type ParameterIterator chan ParameterTuple

// ParameterSource is an interface for a parameter source.
// Parameters may come from various places, such as API specs,
// CLI arguments & another test output.
type ParameterSource interface {
	// Get(paramName string) string
	Iterate() ParameterIterator
}
