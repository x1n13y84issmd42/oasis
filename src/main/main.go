package main

import (
	"github.com/x1n13y84issmd42/oasis/src/env"
	"github.com/x1n13y84issmd42/oasis/src/log"
)

func main() {
	args := &env.Args{}
	env.ParseArgs(args)

	logger := &log.Simple{}

	if args.Script != "" {
		Script(args, logger)
	} else if args.Spec != "" {
		Manual(args, logger)
	} else {
		logger.Usage()
	}
}