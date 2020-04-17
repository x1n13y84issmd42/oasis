package openapi3

import (
	"encoding/json"
	goerrors "errors"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/x1n13y84issmd42/oasis/src/api"
	secAPIKey "github.com/x1n13y84issmd42/oasis/src/api/security/APIKey"
	secHTTP "github.com/x1n13y84issmd42/oasis/src/api/security/HTTP"
	"github.com/x1n13y84issmd42/oasis/src/errors"
	"github.com/x1n13y84issmd42/oasis/src/log"
	"github.com/x1n13y84issmd42/oasis/src/strings"
)

func isstring(i interface{}) (s string) {
	if i != nil {
		if iv, ok := i.(string); ok {
			s = iv
		}

		if iv, ok := i.(int64); ok {
			s = string(iv)
		}

		if iv, ok := i.(uint64); ok {
			s = string(iv)
		}

		//TODO: fix this
		// Long float values (as 9223372036854663000) work incorrectly here.
		if iv, ok := i.(float64); ok {
			s = strconv.FormatFloat(iv, 'f', 0, 64)
		}
	}

	return
}

// Spec is an OAS3-backed API test spec.
type Spec struct {
	Log log.ILogger
	OAS *openapi3.Swagger
}

// GetOperations returns a list of all available test operations from the spec.
func (spec *Spec) GetOperations(params *api.OperationParameters) []*api.Operation {
	ops := []*api.Operation{}

	addOp := func(oasOp *openapi3.Operation, method string, oasPath string, oasPathItem *openapi3.PathItem) {
		if oasOp != nil {
			specOp, specOpErr := spec.makeOperation(method, oasOp, oasPath, oasPathItem, params)
			if specOpErr == nil {
				ops = append(ops, specOp)
			}
		}
	}

	for oasPath, oasPathItem := range spec.OAS.Paths {
		addOp(oasPathItem.Get, "GET", oasPath, oasPathItem)
		addOp(oasPathItem.Post, "POST", oasPath, oasPathItem)
		addOp(oasPathItem.Put, "PUT", oasPath, oasPathItem)
		addOp(oasPathItem.Delete, "DELETE", oasPath, oasPathItem)
		addOp(oasPathItem.Patch, "PATCH", oasPath, oasPathItem)
		addOp(oasPathItem.Head, "HEAD", oasPath, oasPathItem)
		addOp(oasPathItem.Options, "OPTIONS", oasPath, oasPathItem)
		addOp(oasPathItem.Trace, "TRACE", oasPath, oasPathItem)
		addOp(oasPathItem.Connect, "CONNECT", oasPath, oasPathItem)
	}

	return ops
}

// GetOperation returns a list of all available test operations from the spec.
func (spec *Spec) GetOperation(name string, params *api.OperationParameters) (*api.Operation, error) {
	filterOp := func(oasOp *openapi3.Operation) bool {
		return (oasOp != nil && (oasOp.Summary == name || oasOp.OperationID == name))
	}

	for oasPath, oasPathItem := range spec.OAS.Paths {
		if filterOp(oasPathItem.Get) {
			return spec.makeOperation("GET", oasPathItem.Get, oasPath, oasPathItem, params)
		}
		if filterOp(oasPathItem.Post) {
			return spec.makeOperation("POST", oasPathItem.Post, oasPath, oasPathItem, params)
		}
		if filterOp(oasPathItem.Put) {
			return spec.makeOperation("PUT", oasPathItem.Put, oasPath, oasPathItem, params)
		}
		if filterOp(oasPathItem.Delete) {
			return spec.makeOperation("DELETE", oasPathItem.Delete, oasPath, oasPathItem, params)
		}
		if filterOp(oasPathItem.Patch) {
			return spec.makeOperation("PATCH", oasPathItem.Patch, oasPath, oasPathItem, params)
		}
		if filterOp(oasPathItem.Head) {
			return spec.makeOperation("HEAD", oasPathItem.Head, oasPath, oasPathItem, params)
		}
		if filterOp(oasPathItem.Options) {
			return spec.makeOperation("OPTIONS", oasPathItem.Options, oasPath, oasPathItem, params)
		}
		if filterOp(oasPathItem.Connect) {
			return spec.makeOperation("CONNECT", oasPathItem.Connect, oasPath, oasPathItem, params)
		}
		if filterOp(oasPathItem.Trace) {
			return spec.makeOperation("TRACE", oasPathItem.Trace, oasPath, oasPathItem, params)
		}
	}

	return nil, errors.NotFound("Operation", name, nil)
}

func (spec *Spec) makeOperation(
	method string,
	oasOp *openapi3.Operation,
	oasPath string,
	oasPathItem *openapi3.PathItem,
	params *api.OperationParameters,
) (*api.Operation, error) {
	specPath, err := spec.CreatePath(oasPath, oasPathItem, oasOp, params)
	if err != nil {
		return nil, errors.OperationMalformed(oasOp.OperationID, "Could not create operation path.", err)
	}

	specQuery, err := spec.CreateQuery(oasPathItem, oasOp, params)
	if err != nil {
		return nil, errors.OperationMalformed(oasOp.OperationID, "Could not create operation query.", err)
	}

	specSecurity, err := spec.MakeSecurity(oasOp.Security, params)
	if err != nil {
		return nil, errors.OperationMalformed(oasOp.OperationID, "Could not create operation security.", err)
	}

	specRequests := []*api.Request{}
	specResponses := []*api.Response{}

	// In case there is a reponse body with certain Content-Type is expected,
	// collect them all.
	if oasOp.RequestBody != nil && oasOp.RequestBody.Value != nil {
		for oasCT, oasMT := range oasOp.RequestBody.Value.Content {
			specRequests = append(specRequests, spec.MakeRequest(method, specPath, specQuery, oasOp, oasPathItem, oasCT, oasMT, params))
		}
	} else {
		// Otherwise making a single contentless response.
		specRequests = append(specRequests, spec.MakeRequest(method, specPath, specQuery, oasOp, oasPathItem, "", nil, params))
	}

	for oasRespStatus, oasResp := range oasOp.Responses {
		if oasResp.Value != nil {
			specResps, specRespsErr := spec.MakeResponses(oasRespStatus, oasResp.Value)
			if specRespsErr == nil {
				specResponses = append(specResponses, specResps...)
			} else {
				return nil, errors.OperationMalformed(oasOp.OperationID, "Failed to create response bodies.", specRespsErr)
			}
		}
	}

	return &api.Operation{
		OperationDesc: api.OperationDesc{
			ID:          oasOp.OperationID,
			Name:        oasOp.Summary,
			Description: oasOp.Description,
		},
		Method:    method,
		Path:      specPath,
		Security:  specSecurity,
		Requests:  specRequests,
		Responses: specResponses,
	}, nil
}

// MakeSecurity creates an api.Request instance from available operation spec data.
func (spec *Spec) MakeSecurity(
	oasSecReqs *openapi3.SecurityRequirements,
	params *api.OperationParameters,
) (api.ISecurity, error) {

	var oasSec *openapi3.SecurityScheme
	var oasSecErr error
	oasSecurityToken := ""
	oasSecurityName := "[unnamed]"

	getSecurity := func(n string) (*openapi3.SecurityScheme, error) {
		oasSecR := spec.OAS.Components.SecuritySchemes[n]

		if oasSecR == nil || oasSecR.Value == nil {
			return nil, errors.SecurityNotFound(n, "", nil)
		}

		oasSecurityName = n
		return oasSecR.Value, nil
	}

	if params.Security.SecurityHint != "" {
		oasSec, oasSecErr = getSecurity(params.Security.SecurityHint)
	} else if oasSecReqs != nil {
		for _, oasSecReq := range *oasSecReqs {
			for oasSecName := range oasSecReq {
				oasSec, oasSecErr = getSecurity(oasSecName)
				break
			}
		}
	}

	if oasSecErr != nil {
		return nil, oasSecErr
	}

	if params.Security.HTTPAuthValue != "" {
		oasSecurityToken = params.Security.HTTPAuthValue
	} else if oasSec != nil && oasSec.Extensions["x-example"] != nil {
		jre, ok := oasSec.Extensions["x-example"].(json.RawMessage)
		if ok {
			tokenErr := json.Unmarshal(jre, &oasSecurityToken)
			if tokenErr != nil {
				return nil, errors.NoParameters([]string{"x-example"}, nil)
			}
		} else {
			return nil, errors.NoParameters([]string{"x-example"}, nil)
		}
	}

	if oasSec != nil {
		switch oasSec.Type {
		case "apiKey":
			return secAPIKey.New(oasSecurityName, oasSec.In, oasSec.Name, oasSecurityToken, spec.Log), nil

		case "http":
			return secHTTP.New(oasSecurityName, oasSec.Scheme, oasSecurityToken, spec.Log), nil

		default:
			return nil, errors.SecurityNotFound(oasSecurityName, "The last chance error.", nil)
		}
	}

	//TODO: this looks wrong, rethink this entire flow.
	return nil, errors.SecurityNotFound(oasSecurityName, "The last chance error.", nil)
}

// MakeRequest creates an api.Request instance from available operation spec data.
func (spec *Spec) MakeRequest(
	method string,
	specPath string,
	specQuery *url.Values,
	oasOp *openapi3.Operation,
	oasPathItem *openapi3.PathItem,
	oasCT string,
	oasMT *openapi3.MediaType,
	params *api.OperationParameters,
) *api.Request {

	specReq := &api.Request{}
	specReq.Method = method
	specReq.Path = specPath
	specReq.Query = specQuery

	specHeaders := http.Header{}

	for _, oasOpParam := range oasOp.Parameters {
		if oasOpParam.Value == nil || oasOpParam.Value.In != "header" {
			continue
		}

		specHeaders.Add(oasOpParam.Value.Name, isstring(oasOpParam.Value.Example))
	}

	//TODO: same params from oasPathItem & params

	specReq.Headers = specHeaders

	return specReq
}

// MakeResponses creates a list of api.Response instances from available operation spec data.
func (spec *Spec) MakeResponses(
	oasRespStatus string,
	oasResp *openapi3.Response,
) ([]*api.Response, error) {

	specResponses := []*api.Response{}
	specDescription := oasResp.Description
	specStatusCode, _ := strconv.ParseUint(oasRespStatus, 10, 64)
	specHeaders := api.Headers{}

	// Headers.
	for oasHeaderName, oasHeader := range oasResp.Headers {
		if oasHeader.Value != nil {
			var specHeaderErr error
			specHeaders[oasHeaderName], specHeaderErr = spec.MakeHeader(oasHeaderName, oasHeader.Value)
			if specHeaderErr != nil {
				return nil, errors.InvalidResponse("Failed to create a response header '"+oasHeaderName+"'schema.", specHeaderErr)
			}
		}
	}

	// Bodies.
	if len(oasResp.Content) > 0 {
		for oasRespContentType, oasRespContent := range oasResp.Content {
			var specSchema *api.Schema
			var specSchemaErr error
			if oasRespContent.Schema != nil && oasRespContent.Schema.Value != nil {
				specSchema, specSchemaErr = spec.MakeSchema("Response", oasRespContent.Schema.Value)
				if specSchemaErr != nil {
					return nil, errors.InvalidResponse("Failed to create a response body schema.", specSchemaErr)
				}
			}
			specResponses = append(specResponses, &api.Response{
				Description: specDescription,
				StatusCode:  specStatusCode,
				Headers:     specHeaders,
				ContentType: oasRespContentType,
				Schema:      specSchema,
			})
		}
	} else {
		specResponses = append(specResponses, &api.Response{
			Description: specDescription,
			StatusCode:  specStatusCode,
			Headers:     specHeaders,
		})
	}

	return specResponses, nil
}

// MakeHeader creates an api.Rrequest instance from available operation spec data.
func (spec *Spec) MakeHeader(
	oasHeaderName string,
	oasHeader *openapi3.Header,
) (*api.Header, error) {
	var specSchema *api.Schema
	var specSchemaErr error

	if oasHeader.Schema != nil && oasHeader.Schema.Value != nil {
		specSchema, specSchemaErr = spec.MakeSchema(oasHeaderName, oasHeader.Schema.Value)
		if specSchemaErr != nil {
			return nil, specSchemaErr
		}
	}

	return &api.Header{
		Name:        oasHeaderName,
		Description: oasHeader.Description,
		Required:    oasHeader.Required,
		Schema:      specSchema,
	}, nil
}

// MakeSchema creates an api.Schema instance from available operation spec data.
func (spec *Spec) MakeSchema(
	name string,
	oasSchema *openapi3.Schema,
) (*api.Schema, error) {

	if oasSchema != nil {
		jsonSchema, jsonSchemaErr := spec.MakeJSONSchema(name, oasSchema)
		if jsonSchemaErr == nil {
			return &api.Schema{
				Name:       name,
				JSONSchema: jsonSchema,
			}, nil
		}

		return nil, jsonSchemaErr
	}

	return nil, errors.NotFound("Schema", name, nil)
}

// MakeJSONSchema creates an api.Schema instance from available operation spec data.
func (spec *Spec) MakeJSONSchema(
	oasSchemaName string,
	oasSchema *openapi3.Schema,
) (api.JSONSchema, error) {
	jsonSchema, jsonSchemaErr := json.Marshal(oasSchema)
	if jsonSchemaErr == nil {
		sch := make(api.JSONSchema)
		jsonSchemaErr = json.Unmarshal(jsonSchema, &sch)
		if jsonSchemaErr == nil {
			// Adding the components object to the JSON schema object because of $refs
			jsonComps, jsonCompsErr := spec.OAS.Components.MarshalJSON()
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

			return sch, nil
		}

		return nil, errors.InvalidSchema(oasSchemaName, "Failed to unmarshal the schema.", jsonSchemaErr)
	}

	return nil, errors.InvalidSchema(oasSchemaName, "Failed to marshal the schema.", jsonSchemaErr)
}

// CreatePath creates an operation path with parameters replaced by actual values from `example`.
// Path parameters can be found in several places, so there is a priority list:
// 		params.Path
// 		operation-level parameters
// 		path-level parameters
func (spec *Spec) CreatePath(
	oasPath string,
	oasPathItem *openapi3.PathItem,
	oasOp *openapi3.Operation,
	params *api.OperationParameters,
) (string, error) {
	path := oasPath

	fixPath := func(ppn string, ppv string, container string) {
		RX, _ := regexp.Compile("\\{" + ppn + "\\}")

		if RX.Match([]byte(path)) {
			if ppv != "" {
				path = string(RX.ReplaceAll([]byte(path), []byte(ppv)))
				spec.Log.UsingParameterExample(ppn, "path", container)
			} else {
				spec.Log.ParameterHasNoExample(ppn, "path", container)
			}
		}
	}

	useParameters := func(specParams openapi3.Parameters, container string) {
		for _, specP := range specParams {
			if specP == nil || specP.Value == nil || specP.Value.In != "path" {
				continue
			}

			fixPath(specP.Value.Name, isstring(specP.Value.Example), container)
		}
	}

	for ppn, ppv := range params.Path {
		fixPath(ppn, ppv, "override")
	}

	useParameters(oasOp.Parameters, "spec operation")
	useParameters(oasPathItem.Parameters, "spec path")

	// Checking for leftover params.
	RX, _ := regexp.Compile("\\{[\\w\\d-_]+\\}")
	lops := RX.FindAllString(path, -1)
	if len(lops) > 0 {
		return path, errors.NoParameters(strings.Map(lops, func(lop string) string {
			return lop[1 : len(lop)-1]
		}), goerrors.New("XYNTA"))
	}

	return path, nil
}

// CreateQuery creates a query string for an operation request using values from `example`.
// Query parameters can be found in several places, so there is a priority list:
// 		params.Query
// 		operation-level parameters
// 		path-level parameters
func (spec *Spec) CreateQuery(
	oasPathItem *openapi3.PathItem,
	oasOp *openapi3.Operation,
	params *api.OperationParameters,
) (*url.Values, error) {
	qry := make(url.Values)

	add := func(qpn string, qpv string, container string) {
		// if qry.Get(qpn) == "" {
		if qpv != "" {
			qry.Add(qpn, qpv)
			spec.Log.UsingParameterExample(qpn, "query", container)
		} else {
			spec.Log.ParameterHasNoExample(qpn, "query", container)
		}
		// }
	}

	var err error
	// These are for error reporting.
	present := make(map[string]bool)
	missing := make(strings.SIMap)

	useParameters := func(specParams openapi3.Parameters, container string) {
		for _, specP := range specParams {
			if specP == nil || specP.Value == nil || specP.Value.In != "query" || !specP.Value.Required {
				continue
			}

			example := ""
			if specP.Value.Example != nil {
				example = specP.Value.Example.(string)
				present[specP.Value.Name] = true
			} else if !present[specP.Value.Name] {
				missing[specP.Value.Name] = true
				continue
			}

			add(specP.Value.Name, example, container)
		}

	}

	for qpn, qpvs := range params.Query {
		for _, qpv := range qpvs {
			add(qpn, qpv, "override")
		}
	}

	useParameters(oasOp.Parameters, "spec operation")
	useParameters(oasPathItem.Parameters, "spec path")

	if len(missing) > 0 {
		err = errors.NoParameters(missing.Keys(), nil)
	}

	return &qry, err
}

// GetProjectInfo returns project info from the spec.info object.
func (spec *Spec) GetProjectInfo() *api.ProjectInfo {
	return &api.ProjectInfo{
		Title:       spec.OAS.Info.Title,
		Description: spec.OAS.Info.Description,
		Version:     spec.OAS.Info.Version,
	}
}

// GetHost returns an API host by requested description
// from the spec.servers list.
func (spec *Spec) GetHost(name string) (*api.Host, error) {
	for _, oasServer := range spec.OAS.Servers {
		if oasServer.Description == name {
			return &api.Host{
				Name: oasServer.Description,
				URL:  oasServer.URL,
			}, nil
		}
	}

	return nil, errors.NotFound("Host", name, nil)
}

// GetDefaultHost returns the fisr host from the spec.servers list as default.
func (spec *Spec) GetDefaultHost() (*api.Host, error) {
	if len(spec.OAS.Servers) > 0 {
		return &api.Host{
			Name: spec.OAS.Servers[0].Description,
			URL:  spec.OAS.Servers[0].URL,
		}, nil
	}

	return nil, errors.NotFound("Host", "Default", nil)
}
