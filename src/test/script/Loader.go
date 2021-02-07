package script

import (
	"io/ioutil"

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
	}

	yaml.Unmarshal([]byte(fileData), script)

	specs := make(map[string]contract.OperationAccess)

	for k, v := range script.SpecPaths {
		spec := utility.Load(v, script.Log)
		specs[k] = spec
	}

	script.OperationCache = api.NewOperationCache(specs)

	//TODO: some validation is required
	// Like unique op IDs.

	return script
}
