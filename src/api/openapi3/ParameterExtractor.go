package openapi3

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/x1n13y84issmd42/oasis/src/errors"
	"github.com/x1n13y84issmd42/oasis/src/strings"
)

// ParameterExtractor is a utility type concerned with extraction of parameters from Open API specs.
type ParameterExtractor struct {
	Present map[string]bool
	Missing strings.SIMap
}

// NewParameterExtractor creates a new ParameterExtractor instance.
func NewParameterExtractor() *ParameterExtractor {
	return &ParameterExtractor{
		Present: make(map[string]bool),
		Missing: make(strings.SIMap),
	}
}

// ExtractorFunc is a function to handle parameters.
type ExtractorFunc = func(p *openapi3.Parameter, example string)

// Extract extracts parameters and passes them to the provided handler function.
func (pe *ParameterExtractor) Extract(params *openapi3.Parameters, in string, handler ExtractorFunc) {
	for _, specP := range *params {
		if specP == nil || specP.Value == nil || specP.Value.In != in || !specP.Value.Required {
			continue
		}

		example := ""
		if specP.Value.Example != nil {
			example = specP.Value.Example.(string)
			pe.Present[specP.Value.Name] = true
		} else if !pe.Present[specP.Value.Name] {
			pe.Missing[specP.Value.Name] = true
			continue
		}

		handler(specP.Value, example)
	}
}

// Error returns an ErrNoParameters error instance if there were missing parameters.
func (pe *ParameterExtractor) Error() error {
	if len(pe.Missing) > 0 {
		return errors.NoParameters(pe.Missing.Keys(), nil)
	}

	return nil
}
