package main

import (
	"github.com/x1n13y84issmd42/oasis/src/contract"
	"github.com/x1n13y84issmd42/oasis/src/env"
	"github.com/x1n13y84issmd42/oasis/src/test/script"
)

// Script is an entry point for the scripted testing mode.
func Script(args *env.Args, log contract.Logger) {
	log.LoadingScript(args.Script)

	s := script.Load(args.Script, log)
	graph := s.GetExecutionGraph()

	script.NewExecutor(log).Execute(graph.(*script.ExecutionGraph))
}
