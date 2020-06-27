package openapi3

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/contract"
	"github.com/x1n13y84issmd42/oasis/src/errors"
	"github.com/x1n13y84issmd42/oasis/src/params"
)

// DataResolver provides spec data based on user input.
type DataResolver struct {
	contract.EntityTrait
	SpecHosts *openapi3.Servers
}

// NewDataResolver creates a new DataResolver instance.
func NewDataResolver(log contract.Logger, hosts *openapi3.Servers) *DataResolver {
	return &DataResolver{
		EntityTrait: contract.Entity(log),
		SpecHosts:   hosts,
	}
}

// Host returns a ParameterSource which contains a host name
// under the params.KeyHost key.
func (r *DataResolver) Host(hostHint string) contract.ParameterSource {
	var host *string
	if hostHint == "" {
		if len(*r.SpecHosts) > 0 {
			host = &(*r.SpecHosts)[0].URL
		}
	} else {
		for _, oasServer := range *r.SpecHosts {
			if oasServer.Description == hostHint {
				host = &oasServer.URL
			}
		}
	}

	if host != nil {
		src := params.NewMemorySource()
		src.Add(params.KeyHost, *host)
		return src
	}

	return params.NoSource(errors.NotFound("Host", hostHint, nil), r.Log)
}

// Security ...
func (r *DataResolver) Security(name string) contract.Security {
	return nil
}

// Response returns a Validator instance to test response correctness.
func (r *DataResolver) Response(status int, CT string) contract.Validator {
	return &api.Validator{}
}
