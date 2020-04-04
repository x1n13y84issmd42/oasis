package utility

import (
	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/api/openapi3"
	"github.com/x1n13y84issmd42/oasis/src/log"
)

// Load loads an API spec file.
func Load(path string, logger log.ILogger) (api.Spec, error) {
	logger.LoadingSpec(path)

	spec, specErr := openapi3.Load(path, logger)
	if specErr != nil {
		return nil, specErr
	}

	return spec, nil
}
