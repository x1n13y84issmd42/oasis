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
	spec := utility.Load(args.Spec, logger)

	logger.TestingProject(spec)

	// Resolving.
	specOps := utility.NewOperationResolver(spec, logger).Resolve(args.Ops)
	result := test.Success()

	if len(specOps) > 0 {
		for _, specOp := range specOps {
			// Stuffing it with data.
			specOp.Data().URL.Load(args.Use.PathParameters)
			specOp.Data().URL.Load(specOp.Host(args.Host))

			//TODO:
			// specOp.Data().URL.Load(HostResolver(args.Use.Host))
			// HostResolver implements the ParameterSource contract.
			// HostResolver's logic is a spec's domain, so it comes either from spec
			// or from op.

			specOp.Data().Query.Load(args.Use.Query)
			specOp.Data().Headers.Load(args.Use.Headers)

			enrichment := []contract.RequestEnrichment{
				specOp.Data().Query,
				specOp.Data().Headers,
			}

			// Testing.
			result = result.And(test.Operation(specOp, &enrichment, logger))
		}

	} else {
		logger.PrintOperations(spec.Operations())
	}

	if !result.Success {
		os.Exit(255)
	}
}
