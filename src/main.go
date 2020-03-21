package main

import (
	"fmt"
	"os"

	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/log"
	"github.com/x1n13y84issmd42/oasis/src/srx"
	"github.com/x1n13y84issmd42/oasis/src/test"
)

func main() {
	//	Command line args
	var inSpec string
	var inScript string
	var inHost string
	// var inOp string
	var inOps []string
	var useCT string = "*"
	var useSecurity string
	var expCT string = "*"
	var expStatus int64 = 0

	expExecute := srx.Flag("execute").CaptureString(&inScript)
	expFrom := srx.Flag("from").CaptureString(&inSpec)
	expTest := srx.Flag("test").CaptureStringSlice(&inOps)
	expHost := srx.Flag("@").CaptureString(&inHost)

	expUse := srx.Flag("use").Repeat(srx.OneOf([]*srx.SRX{
		srx.Flag("security").CaptureString(&useSecurity),
	}), 0, 1)

	expExpect := srx.Flag("expect").Repeat(srx.OneOf([]*srx.SRX{
		srx.Flag("CT").CaptureString(&expCT),
		srx.Flag("status").CaptureInt64(&expStatus),
	}), 0, 2)

	srx.Repeat(srx.OneOf([]*srx.SRX{
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

			runner.Log.TestingProject(spec.GetProjectInfo())

			if hostOK := runner.UseHost(inHost); hostOK {
				testResult := true
				for _, inOp := range inOps {
					testResult = runner.Test(inOp, useCT, int(expStatus), expCT) && testResult
				}
				if !testResult {
					os.Exit(255)
				}
			} else {
				os.Exit(255)
			}
		} else {
			fmt.Println(specErr)
			os.Exit(255)
		}
	} else {
		fmt.Println("Please specify at least a spec file & an operation to test.")
		fmt.Println("Example:")
		fmt.Println("oasis from path/to/oas_spec.yaml test operation_id")
	}
}
