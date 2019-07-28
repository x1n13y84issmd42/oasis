package api

import (
	"fmt"
	"io/ioutil"

	"github.com/go-yaml/yaml"
)

// Load loads a YAML API spec file
func Load(path string) (Spec, error) {
	fmt.Println(fmt.Sprintf("Loading %s", path))

	fileData, fileErr := ioutil.ReadFile(path)
	if fileErr != nil {
		return nil, fileErr
	}

	specData := make(map[interface{}]interface{})
	yaml.Unmarshal([]byte(fileData), &specData)

	return SpecV3{specData}, nil
}
