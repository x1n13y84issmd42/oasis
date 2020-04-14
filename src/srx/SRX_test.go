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
		builder.Parser.String(f)
	}

	return builder.Parser
}

type TestFlagInput struct {
	Flags       []string
	Args        []string
	ExpComplete bool
	ExpProgress int
}

func TestString(t *testing.T) {
	inputs := []TestFlagInput{
		{
			Flags:       []string{"foo", "bar", "qeq"},
			Args:        []string{"foo", "bar", "qeq"},
			ExpComplete: true,
			ExpProgress: 3,
		},

		{
			Flags:       []string{"foo", "bar", "qeq"},
			Args:        []string{"foo", "bar"},
			ExpComplete: false,
			ExpProgress: 2,
		},

		{
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

func ExpectInt64Slice(t *testing.T, actual *[]int64, expected []int64, msg string) Expectation {
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

func ExpectFloat64(t *testing.T, actual *float64, expected float64, msg string) Expectation {
	return func() bool {
		res := *actual == expected
		if !res {
			t.Errorf("FAILED %s\n", msg)
			t.Errorf("       Expected %f\n", expected)
			t.Errorf("       Actual %f\n", *actual)
		}
		return res
	}
}

func ExpectFloat64Slice(t *testing.T, actual *[]float64, expected []float64, msg string) Expectation {
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

func ExpectBoolSlice(t *testing.T, actual *[]bool, expected []bool, msg string) Expectation {
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
		return OneOf(
			Strings("foo", "bar", "qeq"),
			Strings("one", "two", "three", "four"),
		)
	}

	inputs := []TestInput{
		{
			Init: parserCtor1,
			Args: []string{"foo", "bar", "qeq"},

			ExpComplete: true,
			ExpProgress: 3,
		},

		{
			Init: parserCtor1,
			Args: []string{"one", "two", "three", "four"},

			ExpComplete: true,
			ExpProgress: 4,
		},

		{
			Init: parserCtor1,
			Args: []string{"one", "two", "three", "four", "5", "6ix", "se7en"},

			ExpComplete: true,
			ExpProgress: 4,
		},

		{
			Init: parserCtor1,
			Args: []string{"one", "qyx"},

			ExpComplete: false,
			ExpProgress: 0,
		},

		{
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
		return Repeat(String("yolo"), 0, 1)
	}

	ctor2 := func() *SRX {
		return Repeat(OneOf(
			String("f00"),
			String("b4r"),
		), 3, 5)
	}

	inputs := []TestInput{
		{
			Init:        ctor1,
			Args:        []string{"yolo"},
			ExpComplete: true,
			ExpProgress: 1,
		},

		{
			Init:        ctor1,
			Args:        []string{"whatever"},
			ExpComplete: true,
			ExpProgress: 0,
		},

		{
			Init:        ctor2,
			Args:        []string{"f00", "nonbar"},
			ExpComplete: false,
			ExpProgress: 0,
		},

		{
			Init:        ctor2,
			Args:        []string{"f00", "b4r"},
			ExpComplete: false,
			ExpProgress: 0,
		},

		{
			Init:        ctor2,
			Args:        []string{"f00", "b4r", "f00"},
			ExpComplete: true,
			ExpProgress: 3,
		},

		{
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
		{
			Init: func() *SRX {
				return String("yeet").CaptureString(&v)
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

func TestCaptureStringSlice(t *testing.T) {
	var v []string

	genericTest(t, []TestInput{
		{
			Init: func() *SRX {
				return String("things").CaptureStringSlice(&v)
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
				}, "Expected a slice of strings to get captured."),
			},
		},

		{
			Init: func() *SRX {
				return String("things").CaptureStringSlice(&v)
			},
			Args:        []string{"things"},
			ExpComplete: false,
			ExpProgress: 1,
		},
	})
}

func TestCaptureInt64(t *testing.T) {
	var v int64

	genericTest(t, []TestInput{
		{
			Init: func() *SRX {
				return String("answer").CaptureInt64(&v)
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

func TestCaptureInt64Slice(t *testing.T) {
	var v []int64

	genericTest(t, []TestInput{
		{
			Init: func() *SRX {
				return String("answer").CaptureInt64Slice(&v)
			},
			Args:        []string{"answer", "42,43,44,45"},
			ExpComplete: true,
			ExpProgress: 2,
			Expect: []Expectation{
				ExpectInt64Slice(t, &v, []int64{42, 43, 44, 45}, "Expected a slice of int64 to get captured."),
			},
		},
	})
}

func TestCaptureInt64Fail(t *testing.T) {
	var v int64

	genericTest(t, []TestInput{
		{
			Init: func() *SRX {
				return String("answer").CaptureInt64(&v)
			},
			Args:        []string{"answer", "ab42"},
			ExpComplete: false,
			ExpProgress: 1,
			Expect: []Expectation{
				ExpectInt64(t, &v, 0, "Expected no int64 value to get captured."),
			},
		},
	})
}

func TestCaptureInt64SliceFail(t *testing.T) {
	var v []int64

	genericTest(t, []TestInput{
		{
			Init: func() *SRX {
				return String("answer").CaptureInt64Slice(&v)
			},
			Args:        []string{"answer", "42,ab,44,xxx"},
			ExpComplete: false,
			ExpProgress: 1,
			Expect: []Expectation{
				ExpectInt64Slice(t, &v, []int64{}, "Expected an empty slice of int64 to get captured."),
			},
		},
	})
}

func TestCaptureFloat64(t *testing.T) {
	var v float64

	genericTest(t, []TestInput{
		{
			Init: func() *SRX {
				return String("pie").CaptureFloat64(&v)
			},
			Args:        []string{"pie", "3.1415"},
			ExpComplete: true,
			ExpProgress: 2,
			Expect: []Expectation{
				ExpectFloat64(t, &v, 3.1415, "Expected the '3.1415' float64 value to get captured."),
			},
		},
	})
}

func TestCaptureFloat64Slice(t *testing.T) {
	var v []float64

	genericTest(t, []TestInput{
		{
			Init: func() *SRX {
				return String("answer").CaptureFloat64Slice(&v)
			},
			Args:        []string{"answer", "3.14,42.1,9000.0000001"},
			ExpComplete: true,
			ExpProgress: 2,
			Expect: []Expectation{
				ExpectFloat64Slice(t, &v, []float64{3.14, 42.1, 9000.0000001}, "Expected a slice of float64 to get captured."),
			},
		},
	})
}

func TestCaptureFloat64Fail(t *testing.T) {
	var v float64

	genericTest(t, []TestInput{
		{
			Init: func() *SRX {
				return String("answer").CaptureFloat64(&v)
			},
			Args:        []string{"answer", "3.1abc"},
			ExpComplete: false,
			ExpProgress: 1,
			Expect: []Expectation{
				ExpectFloat64(t, &v, 0, "Expected no float64 value to get captured."),
			},
		},
	})
}

func TestCaptureFloat64SliceFail(t *testing.T) {
	var v []float64

	genericTest(t, []TestInput{
		{
			Init: func() *SRX {
				return String("answer").CaptureFloat64Slice(&v)
			},
			Args:        []string{"answer", "3.one,4.4444,a.b"},
			ExpComplete: false,
			ExpProgress: 1,
			Expect: []Expectation{
				ExpectFloat64Slice(t, &v, []float64{}, "Expected an empty slice of float64 to get captured."),
			},
		},
	})
}

func TestCaptureBool(t *testing.T) {
	var v bool

	genericTest(t, []TestInput{
		{
			Init: func() *SRX {
				return String("truth").CaptureBool(&v)
			},
			Args:        []string{"truth", "true"},
			ExpComplete: true,
			ExpProgress: 2,
			Expect: []Expectation{
				ExpectBool(t, &v, true, "Expected a true bool value to get captured."),
			},
		},
	})
}

func TestCaptureBoolSlice(t *testing.T) {
	var v []bool

	genericTest(t, []TestInput{
		{
			Init: func() *SRX {
				return String("truthies").CaptureBoolSlice(&v)
			},
			Args:        []string{"truthies", "true,true,false,true,false,false"},
			ExpComplete: true,
			ExpProgress: 2,
			Expect: []Expectation{
				ExpectBoolSlice(t, &v, []bool{true, true, false, true, false, false}, "Expected a slice of bool to get captured."),
			},
		},
	})
}

func TestCaptureBoolFail(t *testing.T) {
	var v bool

	genericTest(t, []TestInput{
		{
			Init: func() *SRX {
				return String("truth").CaptureBool(&v)
			},
			Args:        []string{"truth", "I_SWEAR"},
			ExpComplete: false,
			ExpProgress: 1,
			Expect: []Expectation{
				ExpectBool(t, &v, false, "Expected no bool value to get captured."),
			},
		},
	})
}

func TestCaptureBoolSliceFail(t *testing.T) {
	var v []bool

	genericTest(t, []TestInput{
		{
			Init: func() *SRX {
				return String("truthies").CaptureBoolSlice(&v)
			},
			Args:        []string{"truthies", "true,false,I_SWEAR,ture,flase"},
			ExpComplete: false,
			ExpProgress: 1,
			Expect: []Expectation{
				ExpectBoolSlice(t, &v, []bool{}, "Expected an empty slice of bool to get captured."),
			},
		},
	})
}
