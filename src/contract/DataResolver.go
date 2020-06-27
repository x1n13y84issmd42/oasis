package contract

// DataResolver is an interface to resolution of internal spec data
// based on user input such as host name or ID, security name, etc.
type DataResolver interface {

	// Host attempts to come up with a valid host name based on the input
	// host hint, which is a spec-specific host identifier.
	// F.e. it's a name for OAS3, but may be index or a literal host for other
	// spec standards.
	Host(hostHint string) ParameterSource
	Security(secName string) Security
	Response(status int, CT string) Validator
}
