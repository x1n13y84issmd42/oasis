package api

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type imap = map[interface{}]interface{}
type iarray = []interface{}

type mapPropsMapFn = func(k string, v imap) Property

func mapPropsMap(m imap, cb mapPropsMapFn) (props []Property) {
	for k, v := range m {
		props = append(props, cb(k.(string), v.(imap)))
	}
	return
}

func isstring(i interface{}) (s string) {
	if i != nil {
		s = i.(string)
	}

	return
}

// SpecV3 provides access to the spec data V3
type SpecV3 struct {
	data imap
}

// MediaType describes a Media Type Object.
// Used as an internal intermediary type during parsing of requests & responses.
type MediaType struct {
	Schema   *Schema
	Example  string
	Examples ExampleList
}

// GetProjectInfo returns project info.
func (spec SpecV3) GetProjectInfo() *ProjectInfo {
	id := spec.data["info"].(imap)
	return &ProjectInfo{
		Title:       id["title"].(string),
		Description: id["description"].(string),
		Version:     id["version"].(string),
	}
}

// GetOperation extracts an operation data from the spec file.
// It includes everything needed to perform the operation
// and validate it's outcome: request methods, URLs, response & request
// headers, security settings, request & response bodies, etc.
func (spec SpecV3) GetOperation(name string) *Operation {
	ymlPaths := spec.data["paths"].(imap)

	for ymlPath, ymlPathData := range ymlPaths {
		for _, method := range []string{"get", "post", "put", "delete", "patch", "options", "trace", "head"} {
			ymlPathDataM := ymlPathData.(imap)
			ymlOp := ymlPathDataM[method]

			if ymlOp != nil {
				ymlOpM := ymlOp.(imap)

				if ymlOpM["summary"] == name || ymlOpM["operationId"] == name {

					// Loading the associated security scheme
					var specSecurity *Security
					ymlOpSec := ymlOpM["security"]

					if ymlOpSec != nil {
						// `security` is an array of maps, using the first key of the first item as a security scheme name
						for osn := range ymlOpSec.(iarray)[0].(imap) {
							specSecurity = spec.GetSecurity(osn.(string))
							break
						}
					}

					// Loading the requests.
					var specRequests []Request
					if ymlOpM["requestBody"] != nil {
						if ymlOpM["requestBody"].(imap)["content"] != nil {
							ymlReqContents := ymlOpM["requestBody"].(imap)["content"].(imap)
							for ymlReqCT, ymlReqContent := range ymlReqContents {
								specMT := spec.parseMediaType(ymlReqContent.(imap))

								specRequests = append(specRequests, Request{
									ContentType: ymlReqCT.(string),
									Schema:      specMT.Schema,
									Examples:    specMT.Examples,
								})
							}

						}
					}

					// Loading the responses.
					var specResponses []Response
					ymlResps := ymlOpM["responses"]

					if ymlResps != nil {
						ymlResponsesM := ymlResps.(imap)

						// Iterating over status codes in the 'responses' map.
						for ymlStatus, ymlStatusResponse := range ymlResponsesM {
							ymlStatusContentResponses := ymlStatusResponse.(imap)["content"]
							ymlStatusHeaders := ymlStatusResponse.(imap)["headers"]
							specHeaders := HeaderBag{}

							if ymlStatusHeaders != nil {
								// Iterating over header names in the 'responses[STATUS_CODE]' map.
								for ymlHeaderName, ymlHeader := range ymlStatusHeaders.(imap) {
									ymlSHeaderName := ymlHeaderName.(string)
									specHeaders[ymlSHeaderName] = append(specHeaders[ymlSHeaderName], spec.makeHeader(ymlSHeaderName, ymlHeader.(imap)))
								}
							}

							if ymlStatusContentResponses != nil {
								// Iterate over content-type keys in the 'content' map.
								for ymlCT, ymlCTResp := range ymlStatusContentResponses.(imap) {
									specResponses = append(specResponses, spec.makeResponse(ymlCT.(string), ymlStatus.(string), ymlCTResp.(imap), specHeaders))
								}
							} else {
								// Contentless responses which have onlny status code & headers.
								specResponses = append(specResponses, spec.makeResponse("", ymlStatus.(string), nil, specHeaders))
							}
						}
					}

					var specParams []Parameter

					if ymlOpM["parameters"] != nil {
						specParams = spec.mapParameters(ymlOpM["parameters"].(iarray))
					}

					return &Operation{
						Name:       isstring(ymlOpM["summary"]),
						ID:         isstring(ymlOpM["operationId"]),
						Path:       spec.makePath(ymlPath.(string), ymlPathDataM),
						Method:     strings.ToUpper(method),
						Security:   specSecurity,
						Requests:   &specRequests,
						Responses:  &specResponses,
						Parameters: specParams,
					}
				}
			}
		}
	}

	return nil
}

func (spec SpecV3) makePath(ymlPath string, ymlPathData imap) (path Path) {
	path.Path = ymlPath
	if ymlPathData["parameters"] != nil {
		path.Parameters = spec.mapParameters(ymlPathData["parameters"].(iarray))
	}
	return
}

func (spec SpecV3) mapParameters(ymlParameters iarray) []Parameter {
	res := []Parameter{}

	for _, ymlPI := range ymlParameters {
		ymlP := ymlPI.(imap)
		specP := Parameter{
			Name:        isstring(ymlP["name"]),
			In:          spec.mapParameterLocation(isstring(ymlP["in"])),
			Description: isstring(ymlP["description"]),
			Example:     isstring(ymlP["example"]),
		}

		if ymlP["schema"] != nil {
			specP.Schema = spec.parseSchema(specP.Name, ymlP["schema"].(imap))
		}

		res = append(res, specP)
	}

	return res
}

func (spec SpecV3) makeHeader(ymlHeaderName string, ymlHeader imap) Header {
	return Header{
		Name:        ymlHeaderName,
		Schema:      spec.parseSchema(ymlHeaderName, ymlHeader["schema"].(imap)),
		Description: ymlHeader["description"].(string),
		Required:    ymlHeader["required"].(bool),
		// Example: string,
		// Value: string
	}
}

func (spec SpecV3) parseMediaType(ymlMT imap) (specMT MediaType) {
	if ymlMT != nil {
		if ymlMT["schema"] != nil {
			specMT.Schema = spec.parseSchema("unnamed", ymlMT["schema"].(imap))
		}

		if ymlMT["example"] != nil {
			specMT.Example = ymlMT["example"].(string)
		}

		if ymlMT["examples"] != nil {
			specMT.Examples = ExampleList{}
			for ymlExampleName, ymlExample := range ymlMT["examples"].(imap) {
				specMT.Examples[ymlExampleName.(string)] = ExampleObject(ymlExample.(imap)["value"].(imap))
			}
		}
	}

	return
}

func (spec SpecV3) makeResponse(ymlCT string, ymlStatus string, ymlResp imap, specHeaders HeaderBag) Response {
	specMT := spec.parseMediaType(ymlResp)

	specStatus, specStatusErr := strconv.ParseInt(ymlStatus, 10, 64)
	if specStatusErr != nil {
		specStatus = 0
	}

	return Response{
		ContentType: ymlCT,
		Example:     specMT.Example,
		Headers:     specHeaders,
		Schema:      specMT.Schema,
		StatusCode:  int(specStatus),
	}
}

func (spec SpecV3) mapDataType(ymlDT string) DataType {
	dt := map[string]DataType{
		"string":  DataTypeString,
		"object":  DataTypeObject,
		"array":   DataTypeArray,
		"boolean": DataTypeBoolean,
		"number":  DataTypeNumber,
		"integer": DataTypeInteger,
	}[ymlDT]

	if dt == "" {
		fmt.Printf("\"%s\" is an unknown data type.\n", ymlDT)
		//TODO: error?
	}

	return dt
}

func (spec SpecV3) mapSecurityType(ymlST string) SecurityType {
	dt := map[string]SecurityType{
		"http":   SecurityTypeHTTP,
		"apiKey": SecurityTypeAPIKey,
	}[ymlST]

	if dt == "" {
		fmt.Printf("\"%s\" is an unknown security type.\n", ymlST)
		//TODO: error?
	}

	return dt
}

func (spec SpecV3) mapSecurityScheme(ymlSS string) SecurityScheme {
	dt := map[string]SecurityScheme{
		"basic":  SecuritySchemeBasic,
		"digest": SecuritySchemeDigest,
		"bearer": SecuritySchemeBearer,
	}[ymlSS]

	if dt == "" {
		fmt.Printf("\"%s\" is an unknown security scheme.\n", ymlSS)
		//TODO: error?
	}

	return dt
}

func (spec SpecV3) mapParameterLocation(ymlIn string) ParameterLocation {
	pLoc := map[string]ParameterLocation{
		"path":   ParameterLocationPath,
		"query":  ParameterLocationQuery,
		"header": ParameterLocationHeader,
		"cookie": ParameterLocationCookie,
	}[ymlIn]

	if pLoc == "" {
		fmt.Printf("\"%s\" is an unknown parameter location.\n", ymlIn)
	}

	return pLoc
}

func (spec SpecV3) mapSchemaDataType(ymlRespSchema imap) DataType {
	typesMap := map[string]DataType{
		"array":  DataTypeArray,
		"object": DataTypeObject,
	}

	if ymlRespSchema["type"] != nil {
		return typesMap[ymlRespSchema["type"].(string)]
	}

	if ymlRespSchema["$ref"] != nil {
		//TODO: resolve the $ref and get actual data type info?
		return DataTypeObject
	}

	return DataTypeString
}

func (spec SpecV3) resolveShemaRef(ymlSchemaRef string) string {
	r, _ := regexp.Compile("[a-zA-Z0-9]+$")
	ymlSchemaName := r.FindString(ymlSchemaRef)

	if len(ymlSchemaName) == 0 {
		fmt.Printf("Could not find the schema referenced by %s", ymlSchemaRef)
		return ""
	}

	return ymlSchemaName
}

// GetOperations -
func (spec SpecV3) GetOperations() []Operation {
	specOps := []Operation{}

	ymlPaths := spec.data["paths"].(imap)
	for ymlPath, ymlPathData := range ymlPaths {
		for _, method := range []string{"get", "post", "put", "delete", "patch", "options", "trace", "head"} {
			ymlOp := ymlPathData.(imap)[method]
			if ymlOp != nil {
				specOps = append(specOps, Operation{
					Name:   ymlOp.(imap)["summary"].(string),
					ID:     ymlOp.(imap)["operationId"].(string),
					Path:   Path{Path: ymlPath.(string)},
					Method: method,
				})
			}
		}
	}

	return specOps
}

// GetHost finds & returns an API host by name.
func (spec SpecV3) GetHost(name string) *Host {
	ymlServers := spec.data["servers"].(iarray)
	for _, s := range ymlServers {
		ymlServerM := s.(imap)
		if ymlServerM["description"] == name {
			return &Host{
				URL:         ymlServerM["url"].(string),
				Name:        name,
				Description: name,
			}
		}
	}
	return nil
}

// GetDefaultHost returns the first host from the list as default.
func (spec SpecV3) GetDefaultHost() *Host {
	ymlServers := spec.data["servers"].(iarray)
	if len(ymlServers) > 0 {
		ymlServerM := ymlServers[0].(imap)
		return &Host{
			URL:         ymlServerM["url"].(string),
			Name:        ymlServerM["description"].(string),
			Description: ymlServerM["description"].(string),
		}
	}
	return nil
}

// GetSecurity returns a security mechanism definition.
func (spec SpecV3) GetSecurity(name string) *Security {
	ymlScheme := spec.data["components"].(imap)["securitySchemes"].(imap)[name]
	ymlSchemeM := ymlScheme.(imap)

	if ymlScheme == nil {
		fmt.Printf("Security scheme \"%s\" isn't there.\n", name)
		return nil
	}

	//	Example data.
	ymlExample := ymlSchemeM["example"]

	if ymlExample == nil || len(ymlExample.(string)) == 0 {
		fmt.Printf("Security schema \"%s\" has no example data to use during request.\n", name)
		return nil
	}

	return &Security{
		Name:           name,
		SecurityType:   spec.mapSecurityType(ymlSchemeM["type"].(string)),
		SecurityScheme: spec.mapSecurityScheme(ymlSchemeM["scheme"].(string)),
		ParamName:      isstring(ymlSchemeM["name"]),
		In:             spec.mapParameterLocation(isstring(ymlSchemeM["in"])),
		Example:        ymlExample.(string),
	}
}

// GetSchema parses and returns a data schema definition.
func (spec SpecV3) GetSchema(name string) *Schema {
	ymlSchema := spec.data["components"].(imap)["schemas"].(imap)[name]
	if ymlSchema == nil {
		fmt.Printf("Schema \"%s\" isn't there.\n", name)
		return nil
	}

	return spec.parseSchema(name, ymlSchema.(imap))
}

func (spec SpecV3) parseSchema(name string, ymlSchemaM imap) *Schema {
	if ymlSchemaM["$ref"] != nil {
		return spec.GetSchema(spec.resolveShemaRef(ymlSchemaM["$ref"].(string)))
	}

	specSchemaDT := spec.mapDataType(isstring(ymlSchemaM["type"]))

	var specProps *[]Property
	var specItems *Schema

	if specSchemaDT == DataTypeObject {
		specProps = spec.parseObjectProperties(name, ymlSchemaM)
	} else if specSchemaDT == DataTypeArray {
		specItems = spec.parseArrayItems(name, ymlSchemaM)
	}

	return &Schema{
		Name:        name,
		DataType:    specSchemaDT,
		Description: isstring(ymlSchemaM["description"]),
		Example:     isstring(ymlSchemaM["example"]),
		Properties:  specProps,
		Items:       specItems,
	}
}

func (spec SpecV3) parseObjectProperties(name string, ymlSchemaM imap) *[]Property {
	if ymlSchemaM == nil || ymlSchemaM["properties"] == nil {
		return nil
	}

	var ymlRequired iarray

	if ymlSchemaM["required"] != nil {
		ymlRequired = ymlSchemaM["required"].(iarray)
	}

	isRequired := func(propName string) bool {
		for _, ymlRequiredProp := range ymlRequired {
			if ymlRequiredProp.(string) == propName {
				return true
			}
		}

		return false
	}

	propmap := func(propName string, p imap) (specProp Property) {
		specProp.Name = propName

		if p["description"] != nil {
			specProp.Description = p["description"].(string)
		}

		specProp.Schema = spec.parseSchema(propName, p)
		specProp.Required = isRequired(propName)
		//TODO: the rest of validation rules
		return
	}

	specProps := mapPropsMap(ymlSchemaM["properties"].(imap), propmap)
	return &specProps
}

func (spec SpecV3) parseArrayItems(name string, ymlSchemaM imap) *Schema {
	if ymlSchemaM == nil || ymlSchemaM["items"] == nil {
		return nil
	}

	return spec.parseSchema(name, ymlSchemaM["items"].(imap))
}
