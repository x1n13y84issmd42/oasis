package contract

// Security is an interface for security mechanisms.
type Security interface {
	RequestEnrichment

	GetName() string

	// GetValue() string
	// GetToken() string
	// GetUsername() string
	// GetPassword() string

	SetValue(v ParameterAccess)
	SetToken(v ParameterAccess)
	SetUsername(v ParameterAccess)
	SetPassword(v ParameterAccess)
}
