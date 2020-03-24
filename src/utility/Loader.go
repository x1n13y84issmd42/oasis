package utility

import (
	"fmt"

	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/api/openapi3"
)

// Load loads an API spec file.
func Load(path string, params *api.OperationParameters) (api.Spec, error) {
	fmt.Println(fmt.Sprintf("Loading %s", path))

	spec, specErr := openapi3.Load(path, params)
	if specErr != nil {
		return nil, specErr
	}

	return spec, nil
}
