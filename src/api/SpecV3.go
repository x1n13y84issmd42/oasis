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

// GetProjectInfo -
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
		ymlPathDataM := ymlPathData.(imap)

		for _, method := range []string{"get", "post", "put", "delete", "patch", "options", "trace", "head"} {
			ymlOp := ymlPathDataM[method]

			if ymlOp != nil {
				ymlOpM := ymlOp.(imap)

				if ymlOpM["summary"] == name || ymlOpM["operationId"] == name {

					// Loading the associated security scheme
					var specSecurity *Security
					// `security` is an array of maps, using the first key of the first item as a security scheme name
					ymlOpSec := ymlOpM["security"]
					if ymlOpSec != nil {
						ymlOpSec0 := ymlOpSec.(iarray)[0].(imap)
						for osn := range ymlOpSec0 {
							specSecurity = spec.GetSecurity(osn.(string))
							break
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

					return &Operation{
						Name:      ymlOpM["summary"].(string),
						Path:      spec.assemblePath(ymlOpM, ymlPath.(string)),
						Method:    strings.ToUpper(method),
						Security:  specSecurity,
						Responses: &specResponses,
					}
				}
			}
		}
	}

	return nil
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

func (spec SpecV3) assemblePath(ymlOp imap, p string) (path string) {
	path = p
	//TODO: the path itself can also have parameters with examples, use those too.
	if ymlOp["parameters"] != nil {
		ymlParams := ymlOp["parameters"].(iarray)
		for _, ymlIP := range ymlParams {
			ymlP := ymlIP.(imap)
			if ymlP["example"] != nil {
				RX, _ := regexp.Compile("\\{" + ymlP["name"].(string) + "\\}")
				path = string(RX.ReplaceAll([]byte(path), []byte(ymlP["example"].(string))))
			} else {
				fmt.Printf("The path '%s' parameter %s has no example value to use.", p, ymlP["name"].(string))
			}
		}
	}

	return
}

func (spec SpecV3) makeResponse(ymlCT string, ymlStatus string, ymlResp imap, specHeaders HeaderBag) Response {
	ymlRespExample := ""
	var specRespSchema *Schema = nil

	if ymlResp != nil {
		ymlRespSchema := ymlResp["schema"].(imap)

		if ymlResp["example"] != nil {
			ymlRespExample = ymlResp["example"].(string)
		}

		specRespSchema = spec.parseSchema("unnamed", ymlRespSchema)
	}

	specStatus, specStatusErr := strconv.ParseInt(ymlStatus, 10, 64)
	if specStatusErr != nil {
		specStatus = 0
	}

	return Response{
		ContentType: ymlCT,
		Example:     ymlRespExample,
		Headers:     specHeaders,
		Schema:      specRespSchema,
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
		dt = ymlDT
	}

	return dt
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
		for ymlMethod, ymlOp := range ymlPathData.(imap) {
			specOps = append(specOps, Operation{
				Name:   ymlOp.(imap)["summary"].(string),
				Path:   ymlPath.(string),
				Method: ymlMethod.(string),
			})
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

	var specParamLoc ParameterLocation

	//	Parameter location.
	if ymlSchemeM["in"] != nil {
		pLoc := map[string]ParameterLocation{
			"query":  ParameterLocationQuery,
			"header": ParameterLocationHeader,
			"cookie": ParameterLocationCookie,
		}[ymlSchemeM["in"].(string)]

		if len(pLoc) == 0 {
			fmt.Printf("\"%s\" is an unknown parameter location.", ymlSchemeM["in"])
			return nil
		}
	}

	return &Security{
		Name:           name,
		SecurityType:   ymlSchemeM["type"].(SecurityType),
		SecurityScheme: ymlSchemeM["scheme"].(SecurityScheme),
		ParamName:      isstring(ymlSchemeM["name"]),
		In:             specParamLoc,
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
