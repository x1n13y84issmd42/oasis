package api

// OperationDesc is a general operation description.
type OperationDesc struct {
	Name        string
	Description string
	ID          string
}

// Operation is an description of an API operation which is a subject to testing.
type Operation struct {
	OperationDesc
	Method    string
	Path      string
	Security  ISecurity
	Requests  []*Request
	Responses []*Response
}
