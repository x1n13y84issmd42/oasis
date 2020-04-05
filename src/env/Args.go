package env

import (
	"os"

	"github.com/x1n13y84issmd42/oasis/src/srx"
)

// ArgsUse is what goes after the "use" command line argument.
type ArgsUse struct {
	CT       string
	Security string
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
	expExecute := srx.Flag("execute").CaptureString(&args.Script)
	expFrom := srx.Flag("from").CaptureString(&args.Spec)
	expTest := srx.Flag("test").CaptureStringSlice(&args.Ops)
	expHost := srx.Flag("@").CaptureString(&args.Host)

	expUse := srx.Flag("use").Repeat(srx.OneOf([]*srx.SRX{
		srx.Flag("security").CaptureString(&args.Use.Security),
	}), 0, 1)

	expExpect := srx.Flag("expect").Repeat(srx.OneOf([]*srx.SRX{
		srx.Flag("CT").CaptureString(&args.Expect.CT),
		srx.Flag("status").CaptureInt64(&args.Expect.Status),
	}), 0, 2)

	expLogLevel := srx.Flag("loglevel").CaptureInt64(&args.LogLevel)

	expLogStyle := srx.Flag("logstyle").CaptureString(&args.LogStyle)

	srx.Repeat(srx.OneOf([]*srx.SRX{
		expExecute,
		expFrom,
		expTest,
		expUse,
		expExpect,
		expHost,
		expLogLevel,
		expLogStyle,
	}), 1, 4).Parse(os.Args[1:])
}
