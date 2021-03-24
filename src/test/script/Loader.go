package script

import (
	"io/ioutil"
	"path/filepath"

	"github.com/go-yaml/yaml"
	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/contract"
	"github.com/x1n13y84issmd42/oasis/src/utility"
)

// Load loads a script file.
func Load(path string, log contract.Logger) contract.Script {
	fileData, fileErr := ioutil.ReadFile(path)
	if fileErr != nil {
		return NoScript(fileErr, log)
	}

	script := &Script{
		EntityTrait: contract.Entity(log),
		Sec:         make(map[string]*contract.SecurityAccess),
	}

	yaml.Unmarshal([]byte(fileData), script)

	specs := make(map[string]contract.OperationAccess)

	for k, v := range script.SpecPaths {
		specPath, _ := filepath.Abs(filepath.Join("script", v))
		spec := utility.Load(specPath, script.Log)
		specs[k] = spec
	}

	script.OperationCache = api.NewOperationCache(specs)

	//TODO: some validation is required
	// Like unique node IDs.

	return script
}
