package openapi3

import (
	"encoding/json"
	"strconv"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/contract"
	"github.com/x1n13y84issmd42/oasis/src/errors"
	"github.com/x1n13y84issmd42/oasis/src/params"
	"github.com/x1n13y84issmd42/oasis/src/test"
	"github.com/x1n13y84issmd42/oasis/src/test/expect"
)

// DataResolver provides spec data based on user input.
// OAS spec file may contain multiple definitions for everything,
// so this resolver helps to choose the right one among htem
// based on hints.
type DataResolver struct {
	contract.EntityTrait
	Spec          *openapi3.Swagger
	SpecResponses *openapi3.Responses
}

// NewDataResolver creates a new DataResolver instance.
func NewDataResolver(log contract.Logger, spec *openapi3.Swagger, resps *openapi3.Responses) *DataResolver {
	return &DataResolver{
		EntityTrait:   contract.Entity(log),
		Spec:          spec,
		SpecResponses: resps,
	}
}

// Host returns a ParameterSource which contains a host name
// under the params.KeyHost key to be used in the URL parameter set.
func (r *DataResolver) Host(hostHint string) contract.ParameterSource {
	var host *string
	if hostHint == "" {
		if len(r.Spec.Servers) > 0 {
			host = &(r.Spec.Servers)[0].URL
		}
	} else {
		for _, oasServer := range r.Spec.Servers {
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
// Since there may be multiple responses in a OAS spec file, it selects
// on of them based on the arguments.
// If no status is supplied then 200 is used by default.
// If no CT is supplied then "application/json" is used by default.
func (r *DataResolver) Response(status int64, CT string) contract.Validator {
	v := test.NewValidator(r.Log)

	// Responses a grouped under status codes, so selecting the status code first.
	// When no particular status is expected, trying to use the 200 as default.
	specStatus, specResp, err := func() (int64, *openapi3.Response, error) {
		if status == 0 {
			//TODO: log using default status
			status = 200
		}

		sstatus := strconv.Itoa(int(status))

		r := (*r.SpecResponses)[sstatus]

		if r != nil && r.Value != nil {
			return status, r.Value, nil
		}

		return 0, nil, errors.NotFound("spec response", sstatus, nil)
	}()

	if err != nil {
		return test.NoValidator(err, r.Log)
	}

	// Under status code keys there are Content-Typed responses.
	// Selecting the needed one (or application/json as default).
	ct, _, err := func() (string, *openapi3.MediaType, error) {
		if CT == "" {
			CT = "application/json"
			//TODO: log using default CT
		}

		mt := specResp.Content[CT]
		if mt != nil {
			return CT, mt, nil
		}

		return "", nil, errors.NotFound("spec response", CT, nil)
	}()
	if err != nil {
		return test.NoValidator(err, r.Log)
	}

	v.Expect(expect.Status(int(specStatus)))
	v.Expect(expect.ContentType(ct))

	err = r.Headers(specResp, v)
	if err != nil {
		return test.NoValidator(err, r.Log)
	}

	return v
}

// Headers populates the provided validator with expectations for HTTP headers.
func (r *DataResolver) Headers(specResp *openapi3.Response, v contract.Validator) error {
	for headerName, specHeaderRef := range specResp.Headers {
		if specHeaderRef.Value != nil {
			specHeader := specHeaderRef.Value
			if specHeader.Schema != nil && specHeader.Schema.Value != nil {

				specSchema, specSchemaErr := r.MakeSchema(headerName, specHeader.Schema.Value)
				if specSchemaErr != nil {
					return specSchemaErr
				}

				v.Expect(expect.HeaderSchema(headerName, specSchema))
			}

			if specHeader.Required {
				v.Expect(expect.HeaderRequired(headerName))
			}
		}
	}

	return nil
}

// Content populates the provided validator with expectations for HTTP response body.
func (r *DataResolver) Content(mt *openapi3.MediaType, CT string, v contract.Validator) error {
	if mt.Schema != nil && mt.Schema.Value != nil {

		specSchema, specSchemaErr := r.MakeSchema("Response", mt.Schema.Value)
		if specSchemaErr != nil {
			return errors.InvalidResponse("Failed to create a '"+CT+"' response body schema.", specSchemaErr)
		}

		v.Expect(expect.ContentSchema(specSchema))
	}

	return nil
}

// MakeSchema creates an api.Schema instance from available operation spec data.
// The schema is later used to test response contents (headers & bodies) against it.
func (r *DataResolver) MakeSchema(
	oasSchemaName string,
	oasSchema *openapi3.Schema,
) (*api.Schema, error) {
	jsonSchema, jsonSchemaErr := json.Marshal(oasSchema)
	if jsonSchemaErr == nil {
		sch := make(api.JSONSchema)
		jsonSchemaErr = json.Unmarshal(jsonSchema, &sch)
		if jsonSchemaErr == nil {
			// Adding the components object to the JSON schema object because of $refs
			jsonComps, jsonCompsErr := r.Spec.Components.MarshalJSON()
			if jsonCompsErr == nil {
				comps := make(map[string]interface{})
				jsonCompsErr = json.Unmarshal(jsonComps, &comps)
				if jsonCompsErr == nil {
					sch["components"] = comps
				} else {
					return nil, errors.InvalidSchema(oasSchemaName, "Failed to unmarshal Components.", jsonCompsErr)
				}
			} else {
				return nil, errors.InvalidSchema(oasSchemaName, "Failed to marshal Components.", jsonCompsErr)
			}

			return &api.Schema{
				JSONSchema: sch,
				Name:       oasSchemaName,
			}, nil
		}

		return nil, errors.InvalidSchema(oasSchemaName, "Failed to unmarshal the schema.", jsonSchemaErr)
	}

	return nil, errors.InvalidSchema(oasSchemaName, "Failed to marshal the schema.", jsonSchemaErr)
}
