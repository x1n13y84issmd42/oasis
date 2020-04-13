package env

import (
	"fmt"
	"os"
	"strings"

	"github.com/x1n13y84issmd42/oasis/src/srx"
)

// ParameterMap is a generic map of test parameters.
type ParameterMap map[string]string

// ArgsUse is what goes after the "use" command line argument.
type ArgsUse struct {
	CT             string
	Security       string
	PathParameters ParameterMap
}

// ArgsExpect is what goes after the "expect" command line argument.
type ArgsExpect struct {
	CT     string
	Status int64
}

// Args is a program arguments.
type Args struct {
	Script   string
	Spec     string
	Host     string
	Ops      []string
	Use      ArgsUse
	Expect   ArgsExpect
	LogLevel int64
	LogStyle string
}

// ParseArgs parses command line arguments into the args struct.
func ParseArgs(args *Args) {
	expExecute := srx.String("execute").CaptureString(&args.Script)
	expFrom := srx.String("from").CaptureString(&args.Spec)
	expTest := srx.String("test").CaptureStringSlice(&args.Ops)
	expHost := srx.String("@").CaptureString(&args.Host)

	args.Use.PathParameters = ParameterMap{}

	hPathParams := func(params []string) {
		for _, pp := range params {
			pps := strings.Split(pp, "=")
			args.Use.PathParameters[pps[0]] = pps[1]
		}
	}

	expUse := srx.String("use").Repeat(srx.OneOf(
		srx.String("security").CaptureString(&args.Use.Security),
		srx.String("path").String("parameters").HandleStringSlice(hPathParams),
	), 0, 1)

	expExpect := srx.String("expect").Repeat(srx.OneOf(
		srx.String("CT").CaptureString(&args.Expect.CT),
		srx.String("status").CaptureInt64(&args.Expect.Status),
	), 0, 2)

	expLogLevel := srx.String("at").String("level").CaptureInt64(&args.LogLevel)
	expLogStyle := srx.String("in").CaptureString(&args.LogStyle).String("style")
	expLog := srx.String("log").Repeat(srx.OneOf(
		expLogLevel,
		expLogStyle,
	), 1, 2)

	srx.Repeat(srx.OneOf(
		expExecute,
		expFrom,
		expTest,
		expUse,
		expExpect,
		expHost,
		expLog,
	), 1, 8).Parse(os.Args[1:])
	//    ^^^ UPDATE ME EVERY TIME YOU ADD ARGUMENTS

	fmt.Printf("Args: %#v", args)
}
