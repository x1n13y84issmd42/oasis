package utility

import (
	"strconv"
)

// CLIQLArgs is a pair of argument list and an index of a current argument.
type CLIQLArgs struct {
	Args []string
	I    int
}

func (args *CLIQLArgs) Clone() *CLIQLArgs {
	return &CLIQLArgs{
		Args: args.Tail(),
	}
}

func (args *CLIQLArgs) Tail() []string {
	if args.I >= len(args.Args) {
		return []string{}
	}

	return args.Args[args.I:]
}

func (args *CLIQLArgs) Get() string {
	if args.I >= len(args.Args) {
		return ""
	}

	return args.Args[args.I]
}

func (args *CLIQLArgs) Len() int {
	return len(args.Args)
}

func (args *CLIQLArgs) Left() int {
	return len(args.Args) - args.I
}

func (args *CLIQLArgs) Add(v int) {
	args.I += v
}

type ExpressionParser func(*CLIQLArgs) bool

// CLIQLParser is a parser for a CLI Query language.
type CLIQLParser struct {
	ExpressionParsers []ExpressionParser
	Progress          int
	Complete          bool
	L                 string
}

// CLIQL creates a new CLIQLParser instance.
func CLIQL() *CLIQLParser {
	return &CLIQLParser{
		ExpressionParsers: []ExpressionParser{},
	}
}

// CLIQLn creates a new CLIQLParser instance.
func CLIQLn(n string) *CLIQLParser {
	return &CLIQLParser{
		ExpressionParsers: []ExpressionParser{},
		L:                 n,
	}
}

// Parse --
func (cliql *CLIQLParser) Parse(args []string) *CLIQLParser {
	return cliql.actualParse(&CLIQLArgs{Args: args})
}

func (cliql *CLIQLParser) actualParse(args *CLIQLArgs) *CLIQLParser {

	cliql.Complete = false
	cliql.Progress = 0

	if args.Left() <= 0 {
		return cliql
	}

	cpI := 0
	curParser := cliql.ExpressionParsers[cpI]

	for curParser(args) {
		cpI++
		cliql.Progress = args.I
		if cpI >= len(cliql.ExpressionParsers) {
			cliql.Complete = true
			break
		}

		if args.I >= args.Len() {
			// Consumed entire input.
		}

		curParser = cliql.ExpressionParsers[cpI]
	}

	return cliql
}

// Flag --
func (cliql *CLIQLParser) Flag(f string) *CLIQLParser {
	cliql.ExpressionParsers = append(cliql.ExpressionParsers, func(args *CLIQLArgs) bool {
		res := f == args.Get()
		if res {
			args.Add(1)
		}
		return res
	})
	return cliql
}

// Capture stores the current item in the provided string pointer.
func (cliql *CLIQLParser) Capture(v *string) *CLIQLParser {
	cliql.ExpressionParsers = append(cliql.ExpressionParsers, func(args *CLIQLArgs) bool {
		*v = args.Get()
		args.Add(1)
		return true
	})
	return cliql
}

// CaptureInt64 tries to parse the current item as integer value
// and store the result in the provided string pointer.
func (cliql *CLIQLParser) CaptureInt64(v *int64) *CLIQLParser {
	cliql.ExpressionParsers = append(cliql.ExpressionParsers, func(args *CLIQLArgs) bool {
		i, ierr := strconv.ParseInt(args.Get(), 10, 64)
		if ierr == nil {
			*v = i
			args.Add(1)
		}
		return true
	})
	return cliql
}

// Any --
func (cliql *CLIQLParser) Any(parsers []*CLIQLParser) *CLIQLParser {
	cliql.ExpressionParsers = append(cliql.ExpressionParsers, func(args *CLIQLArgs) bool {

		for _, parser := range parsers {
			args2 := args.Clone()
			parser.actualParse(args2)
			if parser.Complete {
				args.Add(args2.I)
				return true
			}
		}

		return false
	})
	return cliql
}

// Repeat --
func (cliql *CLIQLParser) Repeat(parser *CLIQLParser, min uint, max uint) *CLIQLParser {
	cliql.ExpressionParsers = append(cliql.ExpressionParsers, func(args *CLIQLArgs) bool {
		complete := uint(0)
		progress := 0
		args2 := args.Clone()
		for i := uint(0); i < max; i++ {
			parser.actualParse(args2)
			if parser.Complete {
				complete++
				progress += parser.Progress
			}
		}

		res := complete >= min && complete <= max

		if res {
			args.Add(args2.I)
		}

		return res
	})
	return cliql
}
