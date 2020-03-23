package main

import (
	"fmt"
	"os"

	"github.com/x1n13y84issmd42/oasis/src/log"
	"github.com/x1n13y84issmd42/oasis/src/test"
	"github.com/x1n13y84issmd42/oasis/src/utility"
)

// Manual is an entry point for manual testing mode.
func Manual(args *Args, logger log.ILogger) {
	spec, specErr := utility.Load(args.Spec)

	if specErr == nil {
		runner := test.Runner{
			Spec: spec,
			Log:  logger,
		}

		logger.TestingProject(spec.GetProjectInfo())
		logger.PrintOperations(spec.GetOperations())

		if hostOK := runner.UseHost(args.Host); hostOK {
			testResult := true
			for _, inOp := range args.Ops {
				testResult = runner.Test(inOp, args.Use.CT, int(args.Expect.Status), args.Expect.CT) && testResult
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
}
