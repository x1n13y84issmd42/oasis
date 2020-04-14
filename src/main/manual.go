package main

import (
	"fmt"
	"os"

	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/env"
	"github.com/x1n13y84issmd42/oasis/src/log"
	"github.com/x1n13y84issmd42/oasis/src/test"
	"github.com/x1n13y84issmd42/oasis/src/utility"
)

// Manual is an entry point for manual testing mode.
func Manual(args *env.Args, logger log.ILogger) {
	params := &api.OperationParameters{}
	params.LoadFromArgs(args)

	spec, specErr := utility.Load(args.Spec, logger)

	if specErr == nil {
		logger.TestingProject(spec.GetProjectInfo())

		var specHost *api.Host

		if params.Request.HostHint == "" {
			specHost = spec.GetDefaultHost()
			logger.UsingDefaultHost()
		} else {
			specHost = spec.GetHost(params.Request.HostHint)
		}

		if specHost != nil {
			logger.UsingHost(specHost)
		} else {
			logger.HostNotFound(params.Request.HostHint)
		}

		testResult := true
		printOps := false
		for _, inOp := range args.Ops {
			specOp, specOpErr := spec.GetOperation(inOp, params)
			if specOpErr == nil {
				testResult = test.Operation(specHost, specOp, params, logger) && testResult
			} else {
				// logger.OperationNotFound(inOp)
				// fmt.Printf("ERROR: %s", specOpErr.Error())
				logger.XError(specOpErr, logger, log.Tab(0))
				printOps = true
			}
		}

		if printOps {
			// logger.PrintOperations(spec.GetOperations(params))
		}

		if !testResult {
			os.Exit(255)
		}
	} else {
		fmt.Println(specErr)
		os.Exit(255)
	}
}
