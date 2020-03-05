package api

import (
	"fmt"
	"regexp"
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

func istring(i interface{}) (s string) {
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

// GetOperation -
func (spec SpecV3) GetOperation(name string) *Operation {
	paths := spec.data["paths"].(imap)

	for path, pd := range paths {
		pathData := pd.(imap)

		for method, op := range pathData {
			operation := op.(imap)

			if operation["summary"] == name {

				// Loading the associated security scheme
				var security *Security
				// `security` is an array of maps, using the first key of the first item as a security scheme name
				ymlOpSec := operation["security"]
				if ymlOpSec != nil {
					ymlOpSec0 := ymlOpSec.([]interface{})[0].(imap)
					for osn := range ymlOpSec0 {
						spec.GetSecurity(osn.(string))
						break
					}
				}

				// Loading the responses.
				var responses []Response
				ymlResps := operation["responses"]
				if ymlResps != nil {
					ymlMResps := ymlResps.(imap)

					for ymlStatus, ymlStatusResponses := range ymlMResps {
						ymlMStatusCTResponses := ymlStatusResponses.(imap)["content"].(imap)

						for ymlCT, ymlCTResp := range ymlMStatusCTResponses {
							responses = append(responses, spec.makeResponse(ymlCT.(string), ymlStatus.(int), ymlCTResp.(imap)))
						}
					}
				}

				return &Operation{
					Name:      name,
					Path:      spec.assemblePath(operation, path.(string)),
					Method:    strings.ToUpper(method.(string)),
					Security:  security,
					Responses: &responses,
				}
			}
		}
	}

	return nil
}

func (spec SpecV3) assemblePath(ymlOp imap, p string) (path string) {
	path = p
	if ymlOp["parameters"] != nil {
		ymlParams := ymlOp["parameters"].(iarray)
		for _, ymlIP := range ymlParams {
			ymlP := ymlIP.(imap)
			if ymlP["example"] != nil {
				RX, _ := regexp.Compile("\\{" + ymlP["name"].(string) + "\\}")
				path = string(RX.ReplaceAll([]byte(path), []byte(ymlP["example"].(string))))
			} else {
				fmt.Printf("Tha path '%s' parameter %s has no example value to use.", p, ymlP["name"].(string))
			}
		}
	}

	return
}

func (spec SpecV3) makeResponse(ymlCT string, ymlStatus int, ymlResp imap) Response {
	ymlRespSchema := ymlResp["schema"].(imap)
	ymlRespExample := ""

	if ymlResp["example"] != nil {
		ymlRespExample = ymlResp["example"].(string)
	}

	// respDataType := spec.mapSchemaDataType(ymlRespSchema)
	// respModelSchemaName := spec.resolveShemaRef(ymlRespSchema, respDataType)
	respModelSchema := spec.parseSchema("unnamed", ymlRespSchema)

	return Response{
		ContentType: ymlCT,
		// DataType:    respDataType, //TODO: move to schema
		Example:    ymlRespExample,
		Headers:    nil, //TODO
		Schema:     respModelSchema,
		StatusCode: ymlStatus,
	}
}

func (spec SpecV3) mapDataType(ymlDT string) DataType {
	return map[string]DataType{
		"string":  DataTypeString,
		"object":  DataTypeObject,
		"array":   DataTypeArray,
		"boolean": DataTypeBoolean,
		"number":  DataTypeNumber,
	}[ymlDT]
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
	ops := []Operation{}

	paths := spec.data["paths"].(imap)
	for path, pd := range paths {
		pathData := pd.(imap)

		for method, op := range pathData {
			operation := op.(imap)
			ops = append(ops, Operation{
				Name:   operation["summary"].(string),
				Path:   path.(string),
				Method: method.(string),
			})
		}
	}

	return ops
}

// GetHost -
func (spec SpecV3) GetHost(name string) *Host {
	servers := spec.data["servers"].([]interface{})
	for _, s := range servers {
		server := s.(imap)
		if server["description"] == name {
			return &Host{
				URL:         server["url"].(string),
				Name:        name,
				Description: name,
			}
		}
	}
	return nil
}

// GetDefaultHost -
func (spec SpecV3) GetDefaultHost() *Host {
	servers := spec.data["servers"].([]interface{})
	for _, s := range servers {
		server := s.(imap)
		return &Host{
			URL:         server["url"].(string),
			Name:        server["description"].(string),
			Description: server["description"].(string),
		}
	}
	return nil
}

// GetSecurity -
func (spec SpecV3) GetSecurity(name string) *Security {
	ymlScheme := spec.data["components"].(imap)["securitySchemes"].(imap)[name]

	if ymlScheme == nil {
		fmt.Printf("Security scheme \"%s\" isn't there.\n", name)
		return nil
	}

	ymlExample := ymlScheme.(imap)["example"]

	if ymlExample == nil || len(ymlExample.(string)) == 0 {
		fmt.Printf("Security schema \"%s\" has no example data to use during request.\n", name)
		return nil
	}

	pLoc := map[string]ParameterLocation{
		"query":  ParameterLocationQuery,
		"header": ParameterLocationHeader,
		"cookie": ParameterLocationCookie,
	}[ymlScheme.(imap)["in"].(string)]

	if len(pLoc) == 0 {
		fmt.Printf("\"%s\" is an unknown parameter location.", ymlScheme.(imap)["in"])
		return nil
	}

	return &Security{
		Name:         name,
		SecurityType: ymlScheme.(imap)["type"].(string),
		ParamName:    ymlScheme.(imap)["name"].(string),
		In:           pLoc,
		Example:      ymlExample.(string),
	}
}

// GetSchema -
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

	schemaDT := spec.mapDataType(istring(ymlSchemaM["type"]))

	var props *[]Property
	var items *Schema

	if schemaDT == DataTypeObject {
		props = spec.parseObjectProperties(name, ymlSchemaM)
	} else if schemaDT == DataTypeArray {
		items = spec.parseArrayItems(name, ymlSchemaM)
	}

	return &Schema{
		Name:        name,
		DataType:    schemaDT,
		Description: istring(ymlSchemaM["description"]),
		Example:     istring(ymlSchemaM["example"]),
		Properties:  props,
		Items:       items,
	}
}

func (spec SpecV3) parseObjectProperties(name string, ymlSchemaM imap) *[]Property {
	if ymlSchemaM == nil || ymlSchemaM["properties"] == nil {
		return nil
	}

	var ymlRequired []interface{}

	if ymlSchemaM["required"] != nil {
		ymlRequired = ymlSchemaM["required"].([]interface{})
	}

	isRequired := func(propName string) bool {
		for _, ymlRequiredProp := range ymlRequired {
			if ymlRequiredProp.(string) == propName {
				return true
			}
		}

		return false
	}

	propmap := func(propName string, p imap) (prop Property) {
		prop.Name = propName

		if p["description"] != nil {
			prop.Description = p["description"].(string)
		}

		prop.Schema = spec.parseSchema(propName, p)
		prop.Required = isRequired(propName)
		return
	}

	props := mapPropsMap(ymlSchemaM["properties"].(imap), propmap)
	return &props
}

func (spec SpecV3) parseArrayItems(name string, ymlSchemaM imap) *Schema {
	if ymlSchemaM == nil || ymlSchemaM["items"] == nil {
		return nil
	}

	return spec.parseSchema(name, ymlSchemaM["items"].(imap))
}
