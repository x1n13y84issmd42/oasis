package openapi3

import (
	"github.com/getkin/kin-openapi/openapi3"
)

// Load reads the spec file at path, parses it and returns parsed spec data.
func Load(path string) (*Spec, error) {
	oas, oasErr := openapi3.NewSwaggerLoader().LoadSwaggerFromFile(path)
	if oasErr == nil {
		return &Spec{
			OAS: oas,
		}, nil
	}

	return nil, oasErr
}
