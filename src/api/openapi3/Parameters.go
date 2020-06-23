package openapi3

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/x1n13y84issmd42/oasis/src/contract"
)

// SpecParameterSource provides access to spec data.
type SpecParameterSource struct {
	Params *openapi3.Parameters
	In     string
}

// Get retrieves the requested parameters from the spec parameter list.
func (ds *SpecParameterSource) Get(n string) string {
	for _, specP := range *ds.Params {
		if specP == nil || specP.Value == nil || specP.Value.In != ds.In {
			continue
		}

		if specP.Value.Example != nil {
			return specP.Value.Example.(string)
		}
	}

	return ""
}

// Iterate returns an iterable channel to read parameter values.
func (ds *SpecParameterSource) Iterate() contract.ParameterIterator {
	ch := make(contract.ParameterIterator)

	go func() {
		close(ch)
	}()

	return ch
}

// PathParameterSource creates a parameter source concerned with extracting the "path" parameters from a spec.
func PathParameterSource(p *openapi3.Parameters) *SpecParameterSource {
	return &SpecParameterSource{
		Params: p,
		In:     "path",
	}
}

// QueryParameterSource creates a parameter source concerned with extracting the "query" parameters from a spec.
func QueryParameterSource(p *openapi3.Parameters) *SpecParameterSource {
	return &SpecParameterSource{
		Params: p,
		In:     "query",
	}
}

// HeadersParameterSource creates a parameter source concerned with extracting the "header" parameters from a spec.
func HeadersParameterSource(p *openapi3.Parameters) *SpecParameterSource {
	return &SpecParameterSource{
		Params: p,
		In:     "header",
	}
}
