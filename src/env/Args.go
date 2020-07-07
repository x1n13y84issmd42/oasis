package env

import (
	"os"
	"strings"

	"github.com/x1n13y84issmd42/oasis/src/contract"
	"github.com/x1n13y84issmd42/oasis/src/params"
	"github.com/x1n13y84issmd42/ssp"
)

// ParameterMap is a generic map of operation test parameters.
// & query parameters. It is must be subclassed/aliased/whatever-this-is-called-in-go
// in order to override the Iterate() and provide a source name.
type ParameterMap map[string]string

// DoIterate creates an iterable channel to read parameters.
// The name argument is used for the Parameter.Source value.
func (m ParameterMap) DoIterate(name string) contract.ParameterIterator {
	ch := make(contract.ParameterIterator)

	go func() {
		for n, v := range m {
			ch <- contract.ParameterTuple{
				N: n,
				Parameter: contract.Parameter{
					V:      params.Value(v),
					Source: name,
				},
			}
		}
		close(ch)
	}()

	return ch
}

// ParameterMapPath is a map of parameters used in paths.
type ParameterMapPath ParameterMap

// Iterate creates an iterable channel to read parameters.
func (m ParameterMapPath) Iterate() contract.ParameterIterator {
	return ParameterMap(m).DoIterate("arguments, path")
}

// ParameterMultiMap is a map of operation test parameters
// where each key can have multiple values. Used for HTTP headers
// & query parameters. It is must be subclassed/aliased/whatever-this-is-called-in-go
// in order to override the Iterate() and provide a source name.
type ParameterMultiMap map[string][]string

// DoIterate creates an iterable channel to read parameters.
// The name argument is used for the Parameter.Source value.
func (m ParameterMultiMap) DoIterate(name string) contract.ParameterIterator {
	ch := make(contract.ParameterIterator)

	go func() {
		for n, vs := range m {
			for _, v := range vs {
				ch <- contract.ParameterTuple{
					N: n,
					Parameter: contract.Parameter{
						V:      params.Value(v),
						Source: name,
					},
				}
			}
		}
		close(ch)
	}()

	return ch
}

// ParameterMultiMapQuery is a map of parameters used in query.
type ParameterMultiMapQuery ParameterMultiMap

// Iterate creates an iterable channel to read parameters.
func (m ParameterMultiMapQuery) Iterate() contract.ParameterIterator {
	return ParameterMultiMap(m).DoIterate("arguments, query")
}

// ParameterMultiMapHeaders is a map of parameters used in headers.
type ParameterMultiMapHeaders ParameterMultiMap

// Iterate creates an iterable channel to read parameters.
func (m ParameterMultiMapHeaders) Iterate() contract.ParameterIterator {
	return ParameterMultiMap(m).DoIterate("arguments, headers")
}

// ArgsUse is what goes after the "use" command line argument.
type ArgsUse struct {
	CT             string
	Security       string
	PathParameters ParameterMapPath
	Query          ParameterMultiMapQuery
	Headers        ParameterMultiMapHeaders
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
	expExecute := ssp.String("execute").CaptureString(&args.Script)
	expFrom := ssp.String("from").CaptureString(&args.Spec)
	expTest := ssp.String("test").CaptureStringSlice(&args.Ops)
	expHost := ssp.String("@").CaptureString(&args.Host)

	args.Use.PathParameters = ParameterMapPath{}

	hPathParams := func(params []string) {
		for _, pp := range params {
			pps := strings.Split(pp, "=")
			args.Use.PathParameters[pps[0]] = pps[1]
		}
	}

	args.Use.Query = ParameterMultiMapQuery{}
	args.Use.Headers = ParameterMultiMapHeaders{}

	hQueryParams := func(params []string) {
		for _, pp := range params {
			pps := strings.Split(pp, "=")
			if args.Use.Query[pps[0]] == nil {
				args.Use.Query[pps[0]] = []string{}
			}
			args.Use.Query[pps[0]] = append(args.Use.Query[pps[0]], pps[1])
		}
	}

	expUse := ssp.String("use").Repeat(ssp.OneOf(
		ssp.String("security").CaptureString(&args.Use.Security),
		ssp.Strings("path", "parameters").HandleStringSlice(hPathParams),
		ssp.String("query").HandleStringSlice(hQueryParams),
	), 0, 3)

	expExpect := ssp.String("expect").Repeat(ssp.OneOf(
		ssp.String("CT").CaptureString(&args.Expect.CT),
		ssp.String("status").CaptureInt64(&args.Expect.Status),
	), 0, 2)

	expLogLevel := ssp.Strings("at", "level").CaptureInt64(&args.LogLevel)
	expLogStyle := ssp.String("in").CaptureString(&args.LogStyle).String("style")
	expLog := ssp.String("log").Repeat(ssp.OneOf(
		expLogLevel,
		expLogStyle,
	), 1, 2)

	ssp.Repeat(ssp.OneOf(
		ssp.OneOf(
			expExecute,
			expFrom,
		),
		expTest,
		expUse,
		expExpect,
		expHost,
		expLog,
	), 1, 6).Parse(os.Args[1:])
	//    ^^^ UPDATE ME EVERY TIME YOU ADD ARGUMENTS

	// fmt.Printf("Args: %#v\n", args)
}
