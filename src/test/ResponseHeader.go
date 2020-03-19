package test

import (
	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/log"
	"github.com/x1n13y84issmd42/oasis/src/utility"
)

// ResponseHeader tests HTTP response headers.
// If the header API spec says the header is required, it must present in the response with any value.
// If the header API spec has a schema for the header, the response header value must conform to it.
type ResponseHeader struct {
	Log        log.ILogger
	APIHeaders []api.Header
}

// NewResponseHeader ceates a new ResponseHeader test instance.
func NewResponseHeader(apiHdrs []api.Header, log log.ILogger) ResponseHeader {
	return ResponseHeader{
		Log:        log,
		APIHeaders: apiHdrs,
	}
}

// Test validates header values against their shemas
// and applies addition validation as per OAS spec.
func (test ResponseHeader) Test(respHeaderValues []string) bool {
	OK := true
	requiredOK := true

	for _, apiHeader := range test.APIHeaders {
		// Testing the 'required'.
		// It is considered valid if respHeaderValues is not empty.
		if apiHeader.Required {
			requiredOK = (len(respHeaderValues) > 0)

			if !requiredOK {
				test.Log.HeaderHasNoValue(&apiHeader)
			}

			OK = OK && requiredOK
		}

		// Testing the schema.
		if requiredOK && apiHeader.Schema != nil {
			schemaOK := true
			for _, respHeader := range respHeaderValues {
				schemaOK = Schema{apiHeader.Schema, test.Log}.Test(respHeader, utility.NewContext("Header")) && schemaOK
			}

			if !schemaOK {
				test.Log.SchemaFail(apiHeader.Schema)
			}

			OK = OK && schemaOK
		}
	}

	return OK
}
