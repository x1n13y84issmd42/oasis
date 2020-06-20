package contract

import "net/http"

// RequestEnrichment is a utility type that extends a request instance with additional data.
type RequestEnrichment interface {
	Enrich(request *http.Request, log Logger)
}
