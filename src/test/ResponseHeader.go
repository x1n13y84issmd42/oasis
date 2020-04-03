package test

import (
	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/log"
)

// ResponseHeader tests HTTP response headers.
// If the header API spec says the header is required, it must present in the response with any value.
// If the header API spec has a schema for the header, the response header value must conform to it.
func ResponseHeader(specHeaderName string, specHeader *api.Header, respHeaderValues []string, logger log.ILogger) bool {
	OK := true
	requiredOK := true

	// Testing the 'required'.
	// It is considered valid if respHeaderValues is not empty.
	if specHeader.Required {
		requiredOK = (len(respHeaderValues) > 0)

		if !requiredOK {
			logger.HeaderHasNoValue(specHeader)
		}

		OK = OK && requiredOK
	}

	// Testing the schema.
	if requiredOK && specHeader.Schema != nil {
		schemaOK := true
		for _, respHeader := range respHeaderValues {
			schemaOK = ValidateJSONSchema(respHeader, specHeader.Schema) && schemaOK
		}

		if !schemaOK {
			logger.SchemaFail(specHeader.Schema)
		}

		OK = OK && schemaOK
	}

	return OK
}
