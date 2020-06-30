package test

import (
	"net/http"

	"github.com/x1n13y84issmd42/oasis/src/contract"
)

// Request represents an operation HTTP request.
// It contains native API objects, supplies data, tracks results and does logging.
type Request struct {
	contract.EntityTrait
	HTTPRequest *http.Request
	HTTPClient  *http.Client
	Result      *contract.OperationResult
}

// NewRequest creates a new Request instance.
func NewRequest(op contract.Operation, log contract.Logger) *Request {
	req := &Request{
		EntityTrait: contract.Entity(log),

		HTTPRequest: op.GetRequest(),
		HTTPClient: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},

		Result: &contract.OperationResult{
			Success: true,
		},
	}

	req.Result.HTTPRequest = req.HTTPRequest

	return req
}

// Enrich extends the internal request with additional data.
func (req *Request) Enrich(en contract.RequestEnrichment) {
	en.Enrich(req.HTTPRequest, req.Log)
}

// Execute executes the request.
func (req *Request) Execute() *contract.OperationResult {
	req.Log.Requesting(req.HTTPRequest.Method, req.HTTPRequest.URL.String())
	response, err := req.HTTPClient.Do(req.HTTPRequest)

	req.Result.HTTPResponse = response

	if err != nil {
		req.Log.Error(err)
		req.Result.Success = false
	}

	return req.Result
}
