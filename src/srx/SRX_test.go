package srx

import "testing"

type SRXBuilder struct {
	Parser *SRX
}

func NewBuilder() *SRXBuilder {
	return &SRXBuilder{
		Parser: NewSRX(),
	}
}

func (builder *SRXBuilder) SetFlags(flags []string) *SRX {
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

		if parser.Complete != input.ExpComplete {
			t.Errorf("FAILED Expected Complete to be %#v but got %#v", input.ExpComplete, parser.Complete)
		}

		if parser.Progress != input.ExpProgress {
			t.Errorf("FAILED Expected Progress to be %d but got %d", input.ExpProgress, parser.Progress)
		}
	}
}

type ParserCtor = func() *SRX
type Expectation = func() bool

type TestInput struct {
	Init        ParserCtor
	Args        []string
	ExpComplete bool
	ExpProgress int
	Expect      []Expectation
}

func ExpectString(t *testing.T, actual *string, expected string, msg string) Expectation {
	return func() bool {
		res := *actual == expected
		if !res {
			t.Errorf("FAILED %s\n", msg)
			t.Errorf("       Expected %s\n", expected)
			t.Errorf("       Actual %s\n", *actual)
		}
		return res
	}
}

func ExpectBool(t *testing.T, actual *bool, expected bool, msg string) Expectation {
	return func() bool {
		res := *actual == expected
		if !res {
			t.Errorf("FAILED %s\n", msg)
			t.Errorf("       Expected %#v\n", expected)
			t.Errorf("       Actual %#v\n", *actual)
		}
		return res
	}
}

func ExpectInt64(t *testing.T, actual *int64, expected int64, msg string) Expectation {
	return func() bool {
		res := *actual == expected
		if !res {
			t.Errorf("FAILED %s\n", msg)
			t.Errorf("       Expected %d\n", expected)
			t.Errorf("       Actual %d\n", *actual)
		}
		return res
	}
}

func ExpectStringSlice(t *testing.T, actual *[]string, expected []string, msg string) Expectation {
	return func() bool {
		res := len(*actual) == len(expected)

		if res {
			for aI, a := range *actual {
				res = (a == expected[aI])
			}
		}

		if !res {
			t.Errorf("FAILED %s\n", msg)
			t.Errorf("       Expected %#v\n", expected)
			t.Errorf("       Actual %#v\n", *actual)
		}
		return res
	}
}

func genericTest(t *testing.T, inputs []TestInput) {
	for _, input := range inputs {
		parser := input.Init().Parse(input.Args)

		result := true

		if len(input.Expect) > 0 {
			for _, expectation := range input.Expect {
				result = result && expectation()
			}
		}

		if parser.Complete != input.ExpComplete {
			t.Errorf("FAILED Expected Complete to be %#v but got %#v", input.ExpComplete, parser.Complete)
		}

		if parser.Progress != input.ExpProgress {
			t.Errorf("FAILED Expected Progress to be %d but got %d", input.ExpProgress, parser.Progress)
		}
	}
}

func TestOneOf(t *testing.T) {
	parserCtor1 := func() *SRX {
		return OneOf([]*SRX{
			Flag("foo").Flag("bar").Flag("qeq"),
			Flag("one").Flag("two").Flag("three").Flag("four"),
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
	ctor1 := func() *SRX {
		return Repeat(Flag("yolo"), 0, 1)
	}

	ctor2 := func() *SRX {
		return Repeat(OneOf([]*SRX{
			Flag("f00"),
			Flag("b4r"),
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

func TestCaptureString(t *testing.T) {
	var v string

	genericTest(t, []TestInput{
		TestInput{
			Init: func() *SRX {
				return Flag("yeet").CaptureString(&v)
			},
			Args:        []string{"yeet", "you"},
			ExpComplete: true,
			ExpProgress: 2,
			Expect: []Expectation{
				ExpectString(t, &v, "you", "Expected the 'you' string value to get captured."),
			},
		},
	})
}

func TestCaptureInt64(t *testing.T) {
	var v int64

	genericTest(t, []TestInput{
		TestInput{
			Init: func() *SRX {
				return Flag("answer").CaptureInt64(&v)
			},
			Args:        []string{"answer", "42"},
			ExpComplete: true,
			ExpProgress: 2,
			Expect: []Expectation{
				ExpectInt64(t, &v, 42, "Expected the '42' int64 value to get captured."),
			},
		},
	})
}

func TestCaptureStringSlice(t *testing.T) {
	var v []string

	genericTest(t, []TestInput{
		TestInput{
			Init: func() *SRX {
				return Flag("things").CaptureStringSlice(&v)
			},
			Args:        []string{"things", "item,object,gizmo,article"},
			ExpComplete: true,
			ExpProgress: 2,
			Expect: []Expectation{
				ExpectStringSlice(t, &v, []string{
					"item",
					"object",
					"gizmo",
					"article",
				}, "Expected the '42' int64 value to get captured."),
			},
		},

		TestInput{
			Init: func() *SRX {
				return Flag("things").CaptureStringSlice(&v)
			},
			Args:        []string{"things"},
			ExpComplete: false,
			ExpProgress: 1,
		},
	})
}
