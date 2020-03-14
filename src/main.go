package main

import (
	"fmt"
	"os"

	"github.com/x1n13y84issmd42/goasis/src/api"
	"github.com/x1n13y84issmd42/goasis/src/log"
	"github.com/x1n13y84issmd42/goasis/src/test"
	"github.com/x1n13y84issmd42/goasis/src/utility"
)

func main() {
	//	Command line args
	var inSpec string
	var inScript string
	var inHost string
	var inOp string
	var useCT string = "*"
	var useSecurity string
	var expCT string = "*"
	var expStatus int64 = 0

	expExecute := utility.CLIQL().Flag("execute").Capture(&inScript)
	expFrom := utility.CLIQL().Flag("from").Capture(&inSpec)
	expTest := utility.CLIQL().Flag("test").Capture(&inOp)
	expHost := utility.CLIQL().Flag("@").Capture(&inHost)

	expUse := utility.CLIQL().Flag("use").Repeat(utility.CLIQL().Any([]*utility.CLIQLParser{
		utility.CLIQL().Flag("security").Capture(&useSecurity),
	}), 0, 1)

	expExpect := utility.CLIQL().Flag("expect").Repeat(utility.CLIQL().Any([]*utility.CLIQLParser{
		utility.CLIQL().Flag("CT").Capture(&expCT),
		utility.CLIQL().Flag("status").CaptureInt64(&expStatus),
	}), 0, 2)

	utility.CLIQL().Repeat(utility.CLIQL().Any([]*utility.CLIQLParser{
		expExecute,
		expFrom,
		expTest,
		expUse,
		expExpect,
		expHost,
	}), 1, 4).Parse(os.Args[1:])

	if inScript != "" {
		//	Executing a test script
		//	TODO
	} else if inSpec != "" {
		//	Running a single test
		spec, specErr := api.Load(inSpec)

		if specErr == nil {
			runner := test.Runner{
				Spec: spec,
				Log:  log.Simple{},
			}

			testResult := runner.Test(inHost, inOp, useCT, int(expStatus), expCT)
			if !testResult {
				os.Exit(255)
			}
		} else {
			fmt.Println(specErr)
			os.Exit(255)
		}
	} else {
		fmt.Println("Please specify at least a spec file & an operation to test.")
		fmt.Println("Example:")
		fmt.Println("goasis from path/to/oas_spec.yaml test operation_id")
	}
}
