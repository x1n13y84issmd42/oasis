package contract

import (
	"net/http"
)

// Operation is an interface for accessing operation properties.
type Operation interface {
	ID() string
	Name() string
	Description() string
	Method() string
	Path() string

	Data() *OperationData
	Result() *OperationResult
	Resolve() DataResolver

	GetRequest() (*http.Request, error)
}

// OperationData is an interface to access data from various sources
// (spec path, spec op, cli input, test output) needed in order to build
// an http.Request instance.
type OperationData struct {
	URL     StringParameters
	Query   RequestEnrichmentParameters
	Headers RequestEnrichmentParameters
	Body    RequestEnrichmentParameters
}

// Load loads parameters from data2.
func (data *OperationData) Load(data2 *OperationData) {
	data.URL.Load(data2.URL)
	data.Query.Load(data2.Query)
	data.Headers.Load(data2.Headers)
	data.Body.Load(data2.Body)
}

// Reload reloads all the sources.
func (data *OperationData) Reload() {
	data.URL.Reload()
	data.Query.Reload()
	data.Headers.Reload()
	data.Body.Reload()
}
