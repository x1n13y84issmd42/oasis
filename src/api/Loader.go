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

	// fmt.Printf("--- The specData map:\n%v\n\n", specData)

	/* b, err := json.Marshal(specData["components"])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b)) */

	return SpecV3{specData}, nil
}
