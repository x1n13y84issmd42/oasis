package main

import (
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
		var specHostErr error

		if params.Request.HostHint == "" {
			specHost, specHostErr = spec.GetDefaultHost()
			logger.UsingDefaultHost()
		} else {
			specHost, specHostErr = spec.GetHost(params.Request.HostHint)
		}

		if specHost != nil {
			logger.UsingHost(specHost)
		} else if specHostErr != nil {
			logger.Error(specHostErr)
			os.Exit(255)
			return
		}

		result := test.Success()
		resolver := utility.NewOperationResolver(spec, logger)
		specOps := resolver.Resolve(args.Ops, params)

		if len(specOps) > 0 {
			for _, specOp := range specOps {
				result = result.And(test.Operation(specHost, specOp, params, logger))
			}

		} else {
			logger.PrintOperations(spec.GetOperations(params))
		}

		if !result.Success {
			os.Exit(255)
		}
	} else {
		logger.Error(specErr)
		os.Exit(255)
	}
}
