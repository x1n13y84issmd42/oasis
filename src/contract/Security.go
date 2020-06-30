package contract

// Security is an interface for security mechanisms.
type Security interface {
	RequestEnrichment
	GetName() string
}
