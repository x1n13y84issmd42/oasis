package api

import (
	"net/http"
)

// Operation is an operation description.
type Operation struct {
	Name        string
	Description string
	ID          string
	Method      string
	Path        string
	Requests    []*Request
	Responses   Response
}

// CreateRequest creates a Request instance, already configured
// to make requests to the operation URL.
func (op *Operation) CreateRequest() *http.Request {
	req := &http.Request{}
	return req
}

func (op *Operation) addQueryParameters(req *http.Request) {
	/* useParameters := func(specParams []Parameter, container string) {
		q := req.URL.Query()
		for _, specP := range specParams {
			if specP.In != ParameterLocationQuery {
				continue
			}

			if !specP.Required {
				continue
			}

			var reqParamValue string
			hasExample := false
			if specP.Example != "" {
				reqParamValue = specP.Example
				hasExample = true
			} else if specP.Schema != nil {
				bytes, _ := test.pickExample(specP.Schema.Examples)
				if bytes != nil {
					reqParamValue = string(bytes)
					hasExample = true
				}
			}

			if hasExample {
				q.Add(specP.Name, reqParamValue)
				fmt.Printf("\tAdded a \"%s\" query parameter \"%s\".\n", specP.Name, reqParamValue)
			} else {
				test.Log.ParameterHasNoExample(&specP, container)
			}
		}
		req.URL.RawQuery = q.Encode()
	} */

	// useParameters(test.Operation.Parameters, "operation")
	// useParameters(test.Operation.Path.Parameters, "path")
}

func (op *Operation) addHeaders(req *http.Request) {
	/* useParameters := func(specParams []Parameter, container string) {
		for _, specP := range specParams {
			if specP.In != ParameterLocationHeader {
				continue
			}

			if !specP.Required {
				continue
			}

			var reqHeaderValue string
			hasExample := false
			if specP.Example != "" {
				reqHeaderValue = specP.Example
				hasExample = true
			} else if specP.Schema != nil {
				bytes, _ := test.pickExample(specP.Schema.Examples)
				if bytes != nil {
					reqHeaderValue = string(bytes)
					hasExample = true
				}
			}

			if hasExample {
				req.Header.Add(specP.Name, reqHeaderValue)
				fmt.Printf("\tAdded a \"%s\" header \"%s\".\n", specP.Name, reqHeaderValue)
			} else {
				test.Log.ParameterHasNoExample(&specP, container)
			}
		}
	} */

	// useParameters(test.Operation.Parameters, "operation")
	// useParameters(test.Operation.Path.Parameters, "path")
}
