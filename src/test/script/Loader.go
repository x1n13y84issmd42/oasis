package script

import (
	"io/ioutil"

	"github.com/go-yaml/yaml"
	"github.com/x1n13y84issmd42/oasis/src/contract"
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

	//TODO: some validation is required
	// Like unique op IDs.

	return script
}
