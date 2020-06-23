package env

import (
	"os"
	"strings"

	"github.com/x1n13y84issmd42/oasis/src/contract"
	"github.com/x1n13y84issmd42/oasis/src/srx"
)

// ParameterMap is a generic map of operation test parameters.
type ParameterMap map[string]string

// Iterate creates an iterable channel to read parameters.
func (m ParameterMap) Iterate() contract.ParameterIterator {
	ch := make(contract.ParameterIterator)

	go func() {
		for n, v := range m {
			ch <- contract.ParameterTuple{N: n, V: v}
		}
		close(ch)
	}()

	return ch
}

// ParameterMultiMap is a map of operation test parameters where each key can have multiple values.
// Used for HTTP headers & query parameters.
type ParameterMultiMap map[string][]string

// Iterate creates an iterable channel to read parameters.
func (m ParameterMultiMap) Iterate() contract.ParameterIterator {
	ch := make(contract.ParameterIterator)

	go func() {
		for n, vs := range m {
			for _, v := range vs {
				ch <- contract.ParameterTuple{N: n, V: v}
			}
		}
		close(ch)
	}()

	return ch
}

// ArgsUse is what goes after the "use" command line argument.
type ArgsUse struct {
	CT             string
	Security       string
	PathParameters ParameterMap
	Query          ParameterMultiMap
	Headers        ParameterMultiMap
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

	args.Use.Query = ParameterMultiMap{}
	args.Use.Headers = ParameterMultiMap{}

	hQueryParams := func(params []string) {
		for _, pp := range params {
			pps := strings.Split(pp, "=")
			if args.Use.Query[pps[0]] == nil {
				args.Use.Query[pps[0]] = []string{}
			}
			args.Use.Query[pps[0]] = append(args.Use.Query[pps[0]], pps[1])
		}
	}

	expUse := srx.String("use").Repeat(srx.OneOf(
		srx.String("security").CaptureString(&args.Use.Security),
		srx.Strings("path", "parameters").HandleStringSlice(hPathParams),
		srx.String("query").HandleStringSlice(hQueryParams),
	), 0, 3)

	expExpect := srx.String("expect").Repeat(srx.OneOf(
		srx.String("CT").CaptureString(&args.Expect.CT),
		srx.String("status").CaptureInt64(&args.Expect.Status),
	), 0, 2)

	expLogLevel := srx.Strings("at", "level").CaptureInt64(&args.LogLevel)
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

	// fmt.Printf("Args: %#v\n", args)
}
