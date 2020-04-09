package api

// Operation is an description of an API operation which is a subject to testing.
type Operation struct {
	Name        string
	Description string
	ID          string
	Method      string
	Path        string
	Security    ISecurity
	Requests    []*Request
	Responses   []*Response
}
