package utility

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
		// fmt.Printf("aP (%s): Nothing left (%#v)\n", cliql.L, args)
		return cliql
	}

	cpI := 0
	curParser := cliql.ExpressionParsers[cpI]

	// fmt.Printf("aP (%s): Parsing the args %#v\n", cliql.L, args.Tail())

	for curParser(args) {
		cpI++
		cliql.Progress = args.I
		if cpI >= len(cliql.ExpressionParsers) {
			// fmt.Printf("aP (%s): Complete, args.I %d\n", cliql.L, args.I)
			cliql.Complete = true
			break
		}

		if args.I >= args.Len() {
			// fmt.Printf("aP (%s): Out of input\n", cliql.L)
		}

		curParser = cliql.ExpressionParsers[cpI]
	}

	// fmt.Printf("aP (%s): Done parsing, Complete is %#v, Progress is %d\n", cliql.L, cliql.Complete, cliql.Progress)

	return cliql
}

// Flag --
func (cliql *CLIQLParser) Flag(f string) *CLIQLParser {
	cliql.ExpressionParsers = append(cliql.ExpressionParsers, func(args *CLIQLArgs) bool {
		// fmt.Printf("\tFlag('%s') argument '%s'\n", f, args.Get())
		res := f == args.Get()
		if res {
			args.Add(1)
		}
		return res
	})
	return cliql
}

// Capture --
func (cliql *CLIQLParser) Capture(v *string) *CLIQLParser {
	cliql.ExpressionParsers = append(cliql.ExpressionParsers, func(args *CLIQLArgs) bool {
		// fmt.Printf("\tCapturing the argument '%s'\n", args.Get())
		*v = args.Get()
		args.Add(1)
		return true
	})
	return cliql
}

// Any --
func (cliql *CLIQLParser) Any(parsers []*CLIQLParser) *CLIQLParser {
	cliql.ExpressionParsers = append(cliql.ExpressionParsers, func(args *CLIQLArgs) bool {
		// fmt.Printf("  Trying one of few parsers\n")

		for _, parser := range parsers {
			args2 := args.Clone()
			// fmt.Printf("  Any parser #%d:\n", pI)
			parser.actualParse(args2)
			// fmt.Printf("  Any parser #%d Progress is %d\n", pI, parser.Progress)
			if parser.Complete {
				// fmt.Printf("  Any parser #%d is complete\n", pI)
				// args.Add(parser.Progress)
				args.Add(args2.I)
				// args.I = args2.I
				// fmt.Printf("  Any parser #%d args.I %d\n", pI, args.I)
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
		// fmt.Printf("    Repeating from %d to %d, args.I %d\n", min, max, args.I)
		args2 := args.Clone()
		for i := uint(0); i < max; i++ {
			// fmt.Printf("    Repeat #%d in, args2.I %d\n", i, args2.I)
			parser.actualParse(args2)
			if parser.Complete {
				// fmt.Printf("    Repeat #%d Complete, Progress %d\n", i, parser.Progress)
				complete++
				progress += parser.Progress
				// args2.Add(parser.Progress)
				// fmt.Printf("    Repeat #%d out, args2.I %d, progress %d\n", i, args2.I, progress)
			}
		}

		// fmt.Printf("    Succesful repeats: %d\n", complete)
		res := complete >= min && complete <= max

		if res {
			// args.Add(progress)
			args.Add(args2.I)
			// args.I = args2.I
		}

		return res
	})
	return cliql
}
