package utility

import "testing"

type CLIQLBuilder struct {
	Parser *CLIQLParser
}

func NewBuilder() *CLIQLBuilder {
	return &CLIQLBuilder{
		Parser: CLIQL(),
	}
}

func (builder *CLIQLBuilder) SetFlags(flags []string) *CLIQLParser {
	for _, f := range flags {
		builder.Parser.Flag(f)
	}

	return builder.Parser
}

type TestFlagInput struct {
	Flags       []string
	Args        []string
	ExpComplete bool
	ExpProgress int
}

func TestFlag(t *testing.T) {
	inputs := []TestFlagInput{
		TestFlagInput{
			Flags:       []string{"foo", "bar", "qeq"},
			Args:        []string{"foo", "bar", "qeq"},
			ExpComplete: true,
			ExpProgress: 3,
		},

		TestFlagInput{
			Flags:       []string{"foo", "bar", "qeq"},
			Args:        []string{"foo", "bar"},
			ExpComplete: false,
			ExpProgress: 2,
		},

		TestFlagInput{
			Flags:       []string{"foo", "bar", "qeq"},
			Args:        []string{"foo", "bar", "qeq", "baz", "qyx"},
			ExpComplete: true,
			ExpProgress: 3,
		},
	}

	for _, input := range inputs {
		parser := NewBuilder().SetFlags(input.Flags).Parse(input.Args)

		if parser.Complete == input.ExpComplete &&
			parser.Progress == input.ExpProgress {
			t.Log("FINE")
		} else {
			t.Errorf("FAILED Expected Complete to be %#v but got %#v", input.ExpComplete, parser.Complete)
			t.Errorf("FAILED Expected Progress to be %d but got %d", input.ExpProgress, parser.Progress)
		}
	}
}

type ParserCtor = func() *CLIQLParser

type TestInput struct {
	Init        ParserCtor
	Args        []string
	ExpComplete bool
	ExpProgress int
}

func genericTest(t *testing.T, inputs []TestInput) {
	for _, input := range inputs {
		parser := input.Init().Parse(input.Args)

		if parser.Complete == input.ExpComplete &&
			parser.Progress == input.ExpProgress {
			t.Log("FINE")
		} else {
			t.Errorf("FAILED Expected Complete to be %#v but got %#v", input.ExpComplete, parser.Complete)
			t.Errorf("FAILED Expected Progress to be %d but got %d", input.ExpProgress, parser.Progress)
		}
	}
}

func TestAny(t *testing.T) {
	parserCtor1 := func() *CLIQLParser {
		return CLIQL().Any([]*CLIQLParser{
			CLIQL().Flag("foo").Flag("bar").Flag("qeq"),
			CLIQL().Flag("one").Flag("two").Flag("three").Flag("four"),
		})
	}

	inputs := []TestInput{
		TestInput{
			Init: parserCtor1,
			Args: []string{"foo", "bar", "qeq"},

			ExpComplete: true,
			ExpProgress: 3,
		},

		TestInput{
			Init: parserCtor1,
			Args: []string{"one", "two", "three", "four"},

			ExpComplete: true,
			ExpProgress: 4,
		},

		TestInput{
			Init: parserCtor1,
			Args: []string{"one", "two", "three", "four", "5", "6ix", "se7en"},

			ExpComplete: true,
			ExpProgress: 4,
		},

		TestInput{
			Init: parserCtor1,
			Args: []string{"one", "qyx"},

			ExpComplete: false,
			ExpProgress: 0,
		},

		TestInput{
			Init: parserCtor1,
			Args: []string{"qyx", "one"},

			ExpComplete: false,
			ExpProgress: 0,
		},
	}

	genericTest(t, inputs)
}

func TestRepeat(t *testing.T) {
	ctor1 := func() *CLIQLParser {
		return CLIQL().Repeat(CLIQL().Flag("yolo"), 0, 1)
	}

	ctor2 := func() *CLIQLParser {
		return CLIQLn("repeat").Repeat(CLIQLn("any").Any([]*CLIQLParser{
			CLIQLn("f00").Flag("f00"),
			CLIQLn("b4r").Flag("b4r"),
		}), 3, 5)
	}

	inputs := []TestInput{
		TestInput{
			Init:        ctor1,
			Args:        []string{"yolo"},
			ExpComplete: true,
			ExpProgress: 1,
		},

		TestInput{
			Init:        ctor1,
			Args:        []string{"whatever"},
			ExpComplete: true,
			ExpProgress: 0,
		},

		TestInput{
			Init:        ctor2,
			Args:        []string{"f00", "nonbar"},
			ExpComplete: false,
			ExpProgress: 0,
		},

		TestInput{
			Init:        ctor2,
			Args:        []string{"f00", "b4r"},
			ExpComplete: false,
			ExpProgress: 0,
		},

		TestInput{
			Init:        ctor2,
			Args:        []string{"f00", "b4r", "f00"},
			ExpComplete: true,
			ExpProgress: 3,
		},

		TestInput{
			Init:        ctor2,
			Args:        []string{"f00", "b4r", "f00", "f00", "b4r", "bar", "qeq"},
			ExpComplete: true,
			ExpProgress: 5,
		},
	}

	genericTest(t, inputs)
}

func TestCapture(t *testing.T) {
	var v string

	genericTest(t, []TestInput{
		TestInput{
			Init: func() *CLIQLParser {
				return CLIQL().Flag("yeet").Capture(&v)
			},
			Args:        []string{"yeet", "you"},
			ExpComplete: true,
			ExpProgress: 2,
		},
	})

	if v == "you" {
		t.Log("FINE")
	} else {
		t.Errorf("FAILURE Expected the 'you' value to get captured.")
	}
}
