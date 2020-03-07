package api

// Spec is an interface to access specification data
type Spec interface {
	GetProjectInfo() *ProjectInfo
	GetOperation(name string) *Operation
	GetOperations() []Operation
	GetHost(name string) *Host
	GetDefaultHost() *Host
	GetSecurity(name string) *Security
	GetSchema(name string) *Schema
}

// Operation is an operation description from a spec
type Operation struct {
	Name      string
	Path      string
	Method    string
	Headers   *HeaderBag
	Security  *Security
	Requests  *[]Request
	Responses *[]Response
}

// ProjectInfo is a generic project information
type ProjectInfo struct {
	Title       string
	Description string
	Version     string
}

// ParameterLocation describes a parameter location in requests and security settings
type ParameterLocation = string

// SecurityType is a type of security mechanism used
type SecurityType = string

// DataType is data type, FFS
type DataType = string

// Types of security mechanisms
const (
	SecurityTypeBasic  = "basic"
	SecurityTypeAPIKey = "apiKey"
)

// Parameter locations
const (
	ParameterLocationQuery  = "query"
	ParameterLocationHeader = "header"
	ParameterLocationCookie = "cookie"
)

// Data types
const (
	DataTypeString  = "string"
	DataTypeObject  = "object"
	DataTypeArray   = "array"
	DataTypeBoolean = "boolean"
	DataTypeNumber  = "number"
	//TODO: the rest of primitive types
)

// Security is a description of a security mechanism used on some Operation
type Security struct {
	Name         string
	SecurityType SecurityType
	ParamName    string
	In           ParameterLocation
	Example      string
}

// Host is an API host desciption
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
}

// Property describes a property of an object being described by a Schema
type Property struct {
	Name        string
	Description string
	Schema      *Schema
	Required    bool
}

// Header is a description of an HTTP header
type Header struct {
	Name        string
	Schema      *Schema
	Description string
	Required    bool
	Example     string
	Value       string
}

// HeaderBag is a set of headers to used in requests and responses
type HeaderBag = map[string][]Header

// Request describes a generic HTTP request
type Request struct {
	ContentType string
	Headers     HeaderBag
	Schema      *Schema
}

// Response describes a generic HTTP response
type Response struct {
	StatusCode  int
	ContentType string
	Headers     HeaderBag
	Schema      *Schema
	Example     string
}
