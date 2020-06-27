package contract

import (
	"net/http"
)

// Operation ...
type Operation interface {
	ID() string
	Name() string
	Description() string
	Method() string
	Path() string

	Data() *OperationData
	Resolve() DataResolver

	GetRequest() *http.Request
}

// OperationData is an interface to access data from various sources
// (spec path, spec op, cli input, test output) needed in order to build
// an http.Request instance.
type OperationData struct {
	URL     StringParameters
	Query   RequestEnrichmentParameters
	Headers RequestEnrichmentParameters
}
