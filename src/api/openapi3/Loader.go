package openapi3

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/x1n13y84issmd42/oasis/src/log"
)

// Load reads the spec file at path, parses it and returns parsed spec data.
func Load(path string, logger log.ILogger) (*Spec, error) {
	oas, oasErr := openapi3.NewSwaggerLoader().LoadSwaggerFromFile(path)
	if oasErr == nil {
		return &Spec{
			OAS: oas,
			Log: logger,
		}, nil
	}

	return nil, oasErr
}
