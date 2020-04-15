package api

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Spec is an interface to access specification data.
type Spec interface {
	GetProjectInfo() *ProjectInfo
	GetHost(name string) (*Host, error)
	GetDefaultHost() (*Host, error)
	GetOperations(params *OperationParameters) []*Operation
	GetOperation(name string, params *OperationParameters) (*Operation, error)
}

// ProjectInfo is a generic project information.
type ProjectInfo struct {
	Title       string
	Description string
	Version     string
}

// ExampleList is a map of maps to keep request example data in.
type ExampleList map[string]ExampleObject

// ExampleObject is an example request data for an operation.
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

// Host is an API host description.
type Host struct {
	Name string
	URL  string
}
