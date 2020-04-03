package api

// Header is a specification of a response HTTP header used for testing.
type Header struct {
	Name        string
	Description string
	Required    bool
	Schema      *Schema
}

// Headers is a map of headers.
type Headers map[string]*Header
