package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/x1n13y84issmd42/goasis/src/api"
	"github.com/x1n13y84issmd42/goasis/src/log"
	"github.com/x1n13y84issmd42/goasis/src/test"
	"github.com/x1n13y84issmd42/goasis/src/utility"
)

func main() {
	//	Command line args
	var fSpec string
	var fScript string
	var fHost string
	var fOp string
	var fReqCT string = "*"
	var fRespCT string = "*"
	var fRespStatus int64 = 0

	/* flag.StringVar(&fSpec, "spec", "", "A `path` to an OAS file.")
	flag.StringVar(&fScript, "script", "", "A `path` to a test script file.")
	flag.StringVar(&fHost, "host", "", "A `hostname` from the spec file.")
	flag.StringVar(&fOp, "op", "", "An `operation` name from the spec file.")
	flag.StringVar(&fReqCT, "reqCT", "*", "A request `Content-Type` to use for test.")
	flag.StringVar(&fRespCT, "respCT", "*", "A response `Content-Type` to test against.")
	flag.Int64Var(&fRespStatus, "respStatus", 0, "An expected response status `code`. (by default any works)")
	flag.Parse()

	if fSpec == "" && fScript == "" {
		flag.Usage()
		os.Exit(255)
	} */

	var expCT string
	var expStatus string
	var useSecurity string

	expFrom := utility.CLIQL().Flag("from").Capture(&fSpec)
	expTest := utility.CLIQL().Flag("test").Capture(&fOp)
	expHost := utility.CLIQL().Flag("@").Capture(&fHost)

	expUse := utility.CLIQL().Flag("use").Repeat(utility.CLIQL().Any([]*utility.CLIQLParser{
		utility.CLIQL().Flag("security").Capture(&useSecurity),
	}), 0, 1)

	expExpect := utility.CLIQL().Flag("expect").Repeat(utility.CLIQL().Any([]*utility.CLIQLParser{
		utility.CLIQL().Flag("CT").Capture(&expCT),
		utility.CLIQL().Flag("status").Capture(&expStatus),
	}), 0, 2)

	utility.CLIQL().Repeat(utility.CLIQL().Any([]*utility.CLIQLParser{
		expFrom,
		expTest,
		expUse,
		expExpect,
		expHost,
	}), 1, 4).Parse(os.Args[1:])

	fRespStatus, _ = strconv.ParseInt(expStatus, 10, 64)
	fRespCT = expCT

	if fScript != "" {
		//	Executing a test script
		//	TODO
	} else if fSpec != "" {
		//	Running a single test
		spec, specErr := api.Load(fSpec)

		if specErr == nil {
			runner := test.Runner{
				Spec: spec,
				Log:  log.Simple{},
			}

			testResult := runner.Test(fHost, fOp, fReqCT, int(fRespStatus), fRespCT)
			if !testResult {
				os.Exit(255)
			}
		} else {
			fmt.Println(specErr)
			os.Exit(255)
		}
	} else {
		fmt.Println("Please specify at least a spec file & an operation to test.")
		fmt.Println("Example:")
		fmt.Println("goasis from path/to/oas_spec.yaml test operation_id")
	}
}
