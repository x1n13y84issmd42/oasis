package openapi3

import (
	"encoding/json"
	"sort"
	"strconv"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/api/security"
	"github.com/x1n13y84issmd42/oasis/src/contract"
	"github.com/x1n13y84issmd42/oasis/src/errors"
	"github.com/x1n13y84issmd42/oasis/src/params"
	"github.com/x1n13y84issmd42/oasis/src/test"
	"github.com/x1n13y84issmd42/oasis/src/test/expect"

	secAPIKey "github.com/x1n13y84issmd42/oasis/src/api/security/APIKey"
	secHTTP "github.com/x1n13y84issmd42/oasis/src/api/security/HTTP"
)

// DataResolver provides spec data based on user input.
// OAS spec file may contain multiple definitions for everything,
// so this resolver helps to choose the right one among htem
// based on hints.
type DataResolver struct {
	contract.EntityTrait
	Spec          *openapi3.Swagger
	Op            *Operation
	SpecResponses *openapi3.Responses
}

// ResolverExpectedHeader is header to expect.
type ResolverExpectedHeader struct {
	Name     string
	Schema   *api.Schema
	Required bool
}

// NewDataResolver creates a new DataResolver instance.
func NewDataResolver(log contract.Logger, spec *openapi3.Swagger, op *Operation, resps *openapi3.Responses) *DataResolver {
	return &DataResolver{
		EntityTrait:   contract.Entity(log),
		Spec:          spec,
		SpecResponses: resps,
		Op:            op,
	}
}

// Host returns a ParameterSource which contains a host name
// under the params.KeyHost key to be used in the URL parameter set.
func (resolver *DataResolver) Host(hostHint string) contract.ParameterSource {
	var host *string
	if hostHint == "" {
		if len(resolver.Spec.Servers) > 0 {
			host = &(resolver.Spec.Servers)[0].URL
		}
	} else {
		for _, oasServer := range resolver.Spec.Servers {
			if oasServer.Description == hostHint {
				host = &oasServer.URL
			}
		}
	}

	if host != nil {
		src := params.NewMemorySource("resolver")
		src.Add(params.KeyHost, *host)
		return src
	}

	return params.NoSource(errors.NotFound("Host", hostHint, nil), resolver.Log)
}

// SecurityName figures security scheme name from the operation or from global settings.
func (resolver *DataResolver) SecurityName(name string) string {
	beepbop := func(secReqs *openapi3.SecurityRequirements) string {
		if secReqs != nil {
			for _, opSecRef := range *secReqs {
				for specSecName := range opSecRef {
					if name != "" {
						if specSecName == name {
							return specSecName
						}
					}

					if name == "" {
						return specSecName
					}
				}
			}

		}

		return ""
	}

	result := beepbop(resolver.Op.SpecOp.Security)

	if result == "" {
		return beepbop(&resolver.Spec.Security)
	}

	return result
}

// SecurityCredentials returns a username, password and token when available.
func (resolver *DataResolver) SecurityCredentials(scheme *openapi3.SecurityScheme) (string, string, string, error) {
	username := ""
	password := ""
	token := ""

	extract := func(n string) (string, error) {
		if scheme.Extensions[n] != nil {
			ev := scheme.Extensions[n]
			if jre, ok := ev.(json.RawMessage); ok {
				v := ""
				err := json.Unmarshal(jre, &v)
				if err != nil {
					return "", errors.Oops("Cannot unmarshal the '"+n+"' field.", err)
				}

				return v, nil
			}
		}

		return "", nil
	}

	username, err := extract("x-username")
	if err != nil {
		return "", "", "", err
	}

	password, err = extract("x-password")
	if err != nil {
		return "", "", "", err
	}

	token, err = extract("x-token")
	if err != nil {
		return "", "", "", err
	}

	return username, password, token, nil
}

// Security returns a security object to use in request.
func (resolver *DataResolver) Security(name string) contract.Security {
	secName := resolver.SecurityName(name)

	if secName == "" {
		return security.Insecurity(resolver.Log)
		// return api.NoSecurity(errors.NotFound("Security", name, nil), resolver.Log)
	}

	// Retrieving security scheme details from the spec components.
	if secSchemeRef, ok := resolver.Spec.Components.SecuritySchemes[secName]; ok {
		if secScheme := secSchemeRef.Value; secScheme != nil {
			username, password, token, err := resolver.SecurityCredentials(secScheme)

			if err != nil {
				return api.NoSecurity(err, resolver.Log)
			}

			switch secScheme.Type {
			case "apiKey":
				return secAPIKey.New(secName, secScheme.In, secScheme.Name, token, resolver.Log)

			case "http":
				return secHTTP.New(secName, secScheme.Scheme, token, username, password, resolver.Log)
			}
		}
	}

	return security.Insecurity(resolver.Log)
}

// Response returns a Validator instance to test response correctness.
// Since there may be multiple responses in a OAS spec file, it selects
// on of them based on the arguments.
// If no status is supplied then 200 is used by default.
// If no CT is supplied then "application/json" is used by default.
func (resolver *DataResolver) Response(status int64, CT string) contract.Validator {
	v := test.NewValidator(resolver.Log)

	specStatus, specCT, specMT, specResp, err := resolver.MetaData(status, CT)

	if err != nil {
		return test.NoValidator(err, resolver.Log)
	}

	v.Expect(expect.Status(specStatus, resolver.Log))
	v.Expect(expect.ContentType(specCT, resolver.Log))

	err = resolver.Headers(specResp, v)
	if err != nil {
		return test.NoValidator(err, resolver.Log)
	}

	if specMT != nil {
		err = resolver.Content(specMT, specCT, v)
		if err != nil {
			return test.NoValidator(err, resolver.Log)
		}
	}

	return v
}

// MetaData populates the provided validator with expectations for HTTP status & content type.
func (resolver *DataResolver) MetaData(status int64, CT string) (
	int,
	string,
	*openapi3.MediaType,
	*openapi3.Response,
	error,
) {
	// Responses are grouped under status codes, so selecting the status code first.
	// When no particular status is expected, trying to use the 200 as default.
	specStatus, specResp, err := func() (int64, *openapi3.Response, error) {
		if status == 0 {
			//TODO: log using default status
			//TODO: use 200 as default only if it presents in the spec
			//TODO: otherwise pick something
			status = 200
		}

		sstatus := strconv.Itoa(int(status))

		r := (*resolver.SpecResponses)[sstatus]

		if r != nil && r.Value != nil {
			return status, r.Value, nil
		}

		return 0, nil, errors.NotFound("spec response", sstatus, nil)
	}()

	if err != nil {
		return 0, "", nil, nil, err
	}

	// Under status code keys there are Content-Typed responses.
	// Selecting the needed one (or application/json as default).
	ct, mt, err := func() (string, *openapi3.MediaType, error) {
		if CT == "" {
			CT = "application/json"
			//TODO: log using default CT
		}

		if len(specResp.Content) > 0 {
			mt := specResp.Content[CT]
			if mt != nil {
				return CT, mt, nil
			}
			return "", nil, errors.NotFound("spec response", CT, nil)
		}

		return CT, nil, nil
	}()

	if err != nil {
		return 0, "", nil, nil, err
	}

	return int(specStatus), ct, mt, specResp, nil
}

// Headers populates the provided validator with expectations for HTTP headers.
func (resolver *DataResolver) Headers(specResp *openapi3.Response, v contract.Validator) error {
	headers, err := resolver.CollectHeaders(specResp)

	if err != nil {
		return err
	}

	for _, eh := range headers {
		v.Expect(expect.HeaderSchema(eh.Name, eh.Schema, resolver.Log))

		if eh.Required {
			v.Expect(expect.HeaderRequired(eh.Name, resolver.Log))
		}
	}

	return nil
}

// CollectHeaders collects the headers to expect in response.
func (resolver *DataResolver) CollectHeaders(specResp *openapi3.Response) (
	[]ResolverExpectedHeader,
	error,
) {
	res := []ResolverExpectedHeader{}

	keys := []string{}
	for pn := range specResp.Headers {
		keys = append(keys, pn)
	}

	sort.Strings(keys)

	for _, headerName := range keys {
		specHeaderRef := specResp.Headers[headerName]
		if specHeaderRef.Value != nil {
			eh := ResolverExpectedHeader{
				Name: headerName,
			}

			specHeader := specHeaderRef.Value
			if specHeader.Schema != nil && specHeader.Schema.Value != nil {

				specSchema, specSchemaErr := resolver.MakeSchema(headerName, specHeader.Schema.Value)
				if specSchemaErr != nil {
					return []ResolverExpectedHeader{}, specSchemaErr
				}

				eh.Schema = specSchema
			}

			eh.Required = specHeader.Required

			res = append(res, eh)
		}
	}

	return res, nil
}

// Content populates the provided validator with expectations for HTTP response body structure.
func (resolver *DataResolver) Content(mt *openapi3.MediaType, CT string, v contract.Validator) error {
	if mt.Schema != nil && mt.Schema.Value != nil {

		schemaName := "Response"
		if mt.Schema.Value.Title != "" {
			schemaName = mt.Schema.Value.Title
		}

		specSchema, specSchemaErr := resolver.MakeSchema(schemaName, mt.Schema.Value)
		if specSchemaErr != nil {
			return errors.InvalidResponse("Failed to create a '"+CT+"' response body schema.", specSchemaErr)
		}

		v.Expect(expect.ContentSchema(specSchema, resolver.Log))
	}

	return nil
}

// MakeSchema creates an api.Schema instance from available operation spec data.
// The schema is later used to test response contents (headers & bodies) against it.
func (resolver *DataResolver) MakeSchema(
	oasSchemaName string,
	oasSchema *openapi3.Schema,
) (*api.Schema, error) {
	jsonSchema, jsonSchemaErr := json.Marshal(oasSchema)
	if jsonSchemaErr == nil {
		sch := make(api.JSONSchema)
		jsonSchemaErr = json.Unmarshal(jsonSchema, &sch)
		if jsonSchemaErr == nil {
			// Adding the components object to the JSON schema object because of $refs
			jsonComps, jsonCompsErr := resolver.Spec.Components.MarshalJSON()
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
