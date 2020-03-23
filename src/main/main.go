package main

import (
	"fmt"
	"os"

	"github.com/x1n13y84issmd42/oasis/src/log"
	"github.com/x1n13y84issmd42/oasis/src/test"
	"github.com/x1n13y84issmd42/oasis/src/utility"
)

func main() {
	args := &Args{}
	ParseArgs(args)

	logger := log.Simple{}

	if args.Script != "" {
		//	Executing a test script
		Script(args)
	} else if args.Spec != "" {
		//	Running a single test
		spec, specErr := utility.Load(args.Spec)

		if specErr == nil {
			runner := test.Runner{
				Spec: spec,
				Log:  logger,
			}

			runner.Log.TestingProject(spec.GetProjectInfo())
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
	} else {
		logger.Usage()
	}
}
