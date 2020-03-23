package main

import (
	"github.com/x1n13y84issmd42/oasis/src/log"
)

func main() {
	args := &Args{}
	ParseArgs(args)

	logger := &log.Simple{}

	if args.Script != "" {
		//	Executing a test script.
		Script(args, logger)
	} else if args.Spec != "" {
		//	Running a single test.
		Manual(args, logger)
	} else {
		logger.Usage()
	}
}
