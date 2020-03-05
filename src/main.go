package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/x1n13y84issmd42/goasis/src/api"
	"github.com/x1n13y84issmd42/goasis/src/log"
	"github.com/x1n13y84issmd42/goasis/src/test"
)

func main() {
	//	Command line args
	var fSpec string
	var fScript string
	var fHost string
	var fOp string

	flag.StringVar(&fSpec, "spec", "", "A path to an OAS file.")
	flag.StringVar(&fScript, "script", "", "A path to a test script file.")
	flag.StringVar(&fHost, "host", "", "A hostname from the spec file.")
	flag.StringVar(&fOp, "op", "", "An operation name from the spec file.")
	flag.Parse()

	if fSpec == "" && fScript == "" {
		flag.Usage()
		os.Exit(0)
	}

	if fScript != "" {
		//	Executing a test script
		//	TODO
	} else if fSpec != "" {
		//	Running a single test
		spec, specErr := api.Load(fSpec)

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

			// runner.Test("Localhost", "List visits", "*", 0, "*")
			// runner.Test("Localhost", "Single user", "*", 0, "*")
			// runner.Test("Localhost", "List users", "*", 0, "*")
			// runner.Test("Localhost", "List visits failing", "*", 0, "*")
			// runner.Test("Localhost", "Single user failing", "*", 0, "*")
			// runner.Test("Localhost", "List users failing", "*", 0, "*")
			runner.Test(fHost, fOp, "*", 0, "*")
		} else {
			fmt.Println(specErr)
		}
	}

}
