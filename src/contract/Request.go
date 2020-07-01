package contract

// Request represents an operation HTTP request.
type Request interface {
	Enrich(en RequestEnrichment)
	Execute() *OperationResult
}
