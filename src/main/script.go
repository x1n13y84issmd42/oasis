package main

import (
	"fmt"

	"github.com/x1n13y84issmd42/oasis/src/contract"
	"github.com/x1n13y84issmd42/oasis/src/env"
	"github.com/x1n13y84issmd42/oasis/src/test/script"
)

// Script is an entry point for the scripted testing mode.
func Script(args *env.Args, log contract.Logger) {
	log.NOMESSAGE("Testing the %s script!", args.Script)

	script := script.Load(args.Script, log)
	graph := script.GetGraph()

	fmt.Printf("The script is:\n%#v\n", script)
	fmt.Printf("The op graph is:\n%#v\n", graph)
}
