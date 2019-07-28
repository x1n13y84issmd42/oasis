package main

import (
	"fmt"

	"github.com/x1n13y84issmd42/goasis/src/api"
	"github.com/x1n13y84issmd42/goasis/src/log"
	"github.com/x1n13y84issmd42/goasis/src/test"
)

func main() {
	fmt.Println("Hello go")

	spec, specErr := api.Load("../spec/oasis.yaml")

	if specErr == nil {
		runner := test.Runner{
			Spec: spec,
			Log:  log.Simple{},
		}

		// runner.Test("Localhost", "Get boolean", "*", 0, "*")
		// runner.Test("Localhost", "Get number", "*", 0, "*")
		// runner.Test("Localhost", "Get string", "*", 0, "*")
		// runner.Test("Localhost", "List strings", "*", 0, "*")
		// runner.Test("Localhost", "List numbers", "*", 0, "*")
		runner.Test("Localhost", "List visits", "*", 0, "*")
		runner.Test("Localhost", "Single user", "*", 0, "*")
		runner.Test("Localhost", "List users", "*", 0, "*")
	} else {
		fmt.Println(specErr)
	}
}
