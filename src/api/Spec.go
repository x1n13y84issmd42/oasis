package api

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Spec is an interface to access specification data.
type Spec interface {
	GetProjectInfo() *ProjectInfo
	GetOperation(name string) *Operation
	GetOperations() []Operation
	GetHost(name string) *Host
	GetDefaultHost() *Host
	GetSecurity(name string) *Security
	GetSchema(name string) *Schema
}

// ProjectInfo is a generic project information.
type ProjectInfo struct {
	Title       string
	Description string
	Version     string
}

// ParameterLocation describes a parameter location in requests and security settings.
type ParameterLocation string

// SecurityType is a type of security mechanism used.
type SecurityType string

// SecurityScheme is specific exclusively to the http security type
// and describes the HTTP authentication mechanism used.
type SecurityScheme string

// DataType is a schema data type, FFS.
type DataType string

// DataFormat is an additional rule, specific for each type (string can be emails, as well as dates).
type DataFormat string

// ExampleList is a map of maps to keep request example data in.
type ExampleList map[string]ExampleObject

// ExampleObject --
type ExampleObject map[interface{}]interface{}

// MarshalJSON encodes an example map from the OAS spec as a JSON string.
func (ex ExampleObject) MarshalJSON() ([]byte, error) {
	props := []string{}

	for propKey, propVal := range ex {
		jp, err := json.Marshal(propVal)
		if err != nil {
			return nil, err
		}

		props = append(props, fmt.Sprintf("\"%s\":%s", propKey, jp))
	}

	return []byte(fmt.Sprintf("{%s}", strings.Join(props, ","))), nil
}

// Types of security mechanisms
const (
	SecurityTypeHTTP   = SecurityType("http")
	SecurityTypeAPIKey = SecurityType("apiKey")
	// SecurityTypeOAuth2 = SecurityType("oauth2")
	// SecurityTypeOpenID = SecurityType("openIdConnect")
)

// Subtypes of security mechanisms
const (
	SecuritySchemeBasic  = SecurityScheme("basic")
	SecuritySchemeDigest = SecurityScheme("digest")
	SecuritySchemeBearer = SecurityScheme("bearer")
	// SecuritySchemeHoba        = SecurityScheme("hoba")
	// SecuritySchemeMutual      = SecurityScheme("mutual")
	// SecuritySchemeNegotiate   = SecurityScheme("negotiate")
	// SecuritySchemeOauth       = SecurityScheme("oauth")
	// SecuritySchemeScramSHA1   = SecurityScheme("scram-sha-1")
	// SecuritySchemeScramSHA256 = SecurityScheme("scram-sha-256")
	// SecuritySchemeVapid       = SecurityScheme("vapid")
)

// Parameter locations
const (
	ParameterLocationPath   = ParameterLocation("path")
	ParameterLocationQuery  = ParameterLocation("query")
	ParameterLocationHeader = ParameterLocation("header")
	ParameterLocationCookie = ParameterLocation("cookie")
)

// Data types
const (
	DataTypeString  = DataType("string")
	DataTypeObject  = DataType("object")
	DataTypeArray   = DataType("array")
	DataTypeBoolean = DataType("boolean")
	DataTypeNumber  = DataType("number")
	DataTypeInteger = DataType("integer")
)

// Data formats
const (
	DataFormatDate     = DataFormat("date")
	DataFormatDateTime = DataFormat("date-time")
)

// Operation is an operation description.
type Operation struct {
	Name       string
	ID         string
	Path       Path
	Method     string
	Headers    *HeaderBag
	Security   *Security
	Requests   *[]Request
	Responses  *[]Response
	Parameters []Parameter
}

// Security is a description of a security mechanism used on some Operation.
type Security struct {
	Name           string
	SecurityType   SecurityType
	SecurityScheme SecurityScheme
	ParamName      string
	In             ParameterLocation
	Example        string
}

// Parameter is a description of a URL parameter.
type Parameter struct {
	Name        string
	In          ParameterLocation
	Description string
	Example     string
	Schema      *Schema
}

// Path represents an operation path from OAS spec.
type Path struct {
	Path       string
	Parameters []Parameter
}

// Host is an API host description.
type Host struct {
	Name        string
	URL         string
	Description string
}

// Schema is a description of structured data used in requests, responses & security.
// The Properties property must be set if DataType is DataTypeObject.
// The Items property must be set if DataType is DataTypeArray.
type Schema struct {
	Name        string
	Description string
	DataType    DataType
	Properties  *[]Property
	Items       *Schema
	Example     string
	Examples    ExampleList
}

// Property describes a property of an object being described by a Schema.
type Property struct {
	Name        string
	Description string
	Schema      *Schema
	Required    bool
}

// ParameterStyle defines parameter serialization styles from OAS.
type ParameterStyle string

// OAS parameter serialization styles.
const (
	ParameterStyleMatrix         = ParameterStyle("matrix")
	ParameterStyleLabel          = ParameterStyle("label")
	ParameterStyleForm           = ParameterStyle("form")
	ParameterStyleSimple         = ParameterStyle("simple")
	ParameterStyleSpaceDelimited = ParameterStyle("spaceDelimited")
	ParameterStylePipeDelimited  = ParameterStyle("pipeDelimited")
	ParameterStyleDeepObject     = ParameterStyle("deepObject")
)

// Header is a description of an HTTP header.
type Header struct {
	Name        string
	Schema      *Schema
	Description string
	Required    bool
	AllowEmpty  bool
	Explode     bool
	Style       ParameterStyle
	Example     string
	Value       string
}

// HeaderBag is a set of headers to used in requests and responses.
type HeaderBag = map[string][]Header

// Request describes a requestBody item from the OAS spec.
type Request struct {
	ContentType string
	Headers     HeaderBag
	Schema      *Schema
	Examples    ExampleList
}

// Response describes a generic HTTP response.
type Response struct {
	StatusCode  int
	ContentType string
	Headers     HeaderBag
	Schema      *Schema
	Example     string
}
