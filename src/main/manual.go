package main

import (
	"os"

	"github.com/x1n13y84issmd42/oasis/src/contract"
	"github.com/x1n13y84issmd42/oasis/src/env"
	"github.com/x1n13y84issmd42/oasis/src/test"
	"github.com/x1n13y84issmd42/oasis/src/utility"
)

// Manual is an entry point for manual testing mode.
func Manual(args *env.Args, logger contract.Logger) {
	params := &contract.OperationParameters{}
	params.LoadFromArgs(args)

	spec := utility.Load(args.Spec, logger)

	logger.TestingProject(spec)
	specHost := ""

	// Resolving.
	resolver := utility.NewOperationResolver(spec, logger)
	specOps := resolver.Resolve(args.Ops)
	result := test.Success()

	if len(specOps) > 0 {
		for _, specOp := range specOps {
			// Testing.
			result = result.And(test.Operation(specHost, specOp, params, logger))
		}

	} else {
		logger.PrintOperations(spec.Operations())
	}

	if !result.Success {
		os.Exit(255)
	}
}
