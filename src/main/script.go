package main

import (
	"fmt"

	"github.com/x1n13y84issmd42/gog/graph/comp"
	"github.com/x1n13y84issmd42/oasis/src/contract"
	"github.com/x1n13y84issmd42/oasis/src/env"
	"github.com/x1n13y84issmd42/oasis/src/test/script"
)

// Script is an entry point for the scripted testing mode.
func Script(args *env.Args, log contract.Logger) {
	log.NOMESSAGE("Testing the %s script!", args.Script)

	script := script.Load(args.Script, log)
	graph := script.GetExecutionGraph()

	//TODO: check for loops.
	n0 := comp.MotherNode(graph)
	fmt.Printf("Execution starts from the node '%s'\n", n0.ID())

	//TODO: execute.
}
