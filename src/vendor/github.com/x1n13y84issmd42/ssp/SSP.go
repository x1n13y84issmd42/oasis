package ssp

import (
	"strconv"
	"strings"
)

// ExpressionFn is a single parser function.
// ssp instances contain multiple of these.
type ExpressionFn func(*Cursor) bool

// SSP is regular expression on string slices.
// It allows to parse slices containing structured or patterned data
// and extract that data into variables for further use.
//
// For example, if you have a slice of strings:
//		inputSlice := []string{"er", "er", "er", "er", "yolo", "1,2,3,4"}
// you can make an expression for it and parse your slice:
//		rx := SSP.Repeat(SSP.Flag("er"), 2, 200).CaptureString(&theYoloString).CaptureStringSlice(&the1234StringSlice)
//		rx.Parse(inputSlice)
// thus making sure it fits the pattern, and extracting the "yolo" string
// and a list of "1", "2", "3" & "4" strings into variables.
type SSP struct {
	Expressions []ExpressionFn
	Progress    int
	Complete    bool
}

// New creates a new ssp instance.
func New() *SSP {
	return &SSP{
		Expressions: []ExpressionFn{},
	}
}

// Parse starts the parsing process.
func (ssp *SSP) Parse(args []string) *SSP {
	return ssp.actualParse(&Cursor{Args: args})
}

// actualParse actually parses.
func (ssp *SSP) actualParse(cursor *Cursor) *SSP {

	ssp.Complete = false
	ssp.Progress = 0

	if cursor.Left() <= 0 {
		return ssp
	}

	cpI := 0
	curParser := ssp.Expressions[cpI]

	for curParser(cursor) {
		cpI++
		ssp.Progress = cursor.I
		if cpI >= len(ssp.Expressions) {
			ssp.Complete = true
			break
		}

		if cursor.I >= cursor.Len() {
			// Consumed entire input.
		}

		curParser = ssp.Expressions[cpI]
	}

	return ssp
}

// String parses a single string value.
func (ssp *SSP) String(f string) *SSP {
	ssp.Expressions = append(ssp.Expressions, func(cursor *Cursor) bool {
		res := f == cursor.Get()
		if res {
			cursor.Shift(1)
		}
		return res
	})
	return ssp
}

//StringHandler is a function to handle string values.
type StringHandler = func(string)

//StringSliceHandler is a function to handle []string values.
type StringSliceHandler = func([]string)

//Int64Handler is a function to handle int64 values.
type Int64Handler = func(int64)

//Int64SliceHandler is a function to handle []int64 values.
type Int64SliceHandler = func([]int64)

//Float64Handler is a function to handle float64 values.
type Float64Handler = func(float64)

//Float64SliceHandler is a function to handle []float64 values.
type Float64SliceHandler = func([]float64)

//BoolHandler is a function to handle bool values.
type BoolHandler = func(bool)

//BoolSliceHandler is a function to handle bool values.
type BoolSliceHandler = func([]bool)

// HandleString invokes the provided handler function with the current item as an argument.
func (ssp *SSP) HandleString(h StringHandler) *SSP {
	ssp.Expressions = append(ssp.Expressions, func(cursor *Cursor) bool {
		h(cursor.Get())
		cursor.Shift(1)
		return true
	})
	return ssp
}

// HandleStringSlice tries to parse the current item as a comma-separated list
// of string values and call the provided string slice handler with that list as an argument.
func (ssp *SSP) HandleStringSlice(h StringSliceHandler) *SSP {
	ssp.Expressions = append(ssp.Expressions, func(cursor *Cursor) bool {
		item := cursor.Get()
		if item != "" {
			values := strings.Split(item, ",")
			if len(values) > 0 {
				h(values)
				cursor.Shift(1)
				return true
			}
		}
		return false
	})
	return ssp
}

// HandleInt64 tries to parse the current item as integer value
// and store the result in the provided string pointer.
func (ssp *SSP) HandleInt64(h Int64Handler) *SSP {
	ssp.Expressions = append(ssp.Expressions, func(cursor *Cursor) bool {
		i, ierr := strconv.ParseInt(cursor.Get(), 10, 64)
		if ierr == nil {
			h(i)
			cursor.Shift(1)
			return true
		}
		return false
	})
	return ssp
}

// HandleInt64Slice tries to parse the current item as integer value
// and store the result in the provided string pointer.
func (ssp *SSP) HandleInt64Slice(h Int64SliceHandler) *SSP {
	ssp.Expressions = append(ssp.Expressions, func(cursor *Cursor) bool {
		item := cursor.Get()
		if item != "" {
			ss := strings.Split(item, ",")
			if len(ss) > 0 {
				is := []int64{}
				for _, sv := range ss {
					iv, iverr := strconv.ParseInt(sv, 10, 64)
					if iverr == nil {
						is = append(is, iv)
					} else {
						return false
					}
				}
				if len(is) > 0 {
					h(is)
					cursor.Shift(1)
					return true
				}
			}
		}
		return false
	})
	return ssp
}

// HandleFloat64 tries to parse the current item as integer value
// and store the result in the provided string pointer.
func (ssp *SSP) HandleFloat64(h Float64Handler) *SSP {
	ssp.Expressions = append(ssp.Expressions, func(cursor *Cursor) bool {
		f, ferr := strconv.ParseFloat(cursor.Get(), 64)
		if ferr == nil {
			h(f)
			cursor.Shift(1)
			return true
		}
		return false
	})
	return ssp
}

// HandleFloat64Slice tries to parse the current item as integer value
// and store the result in the provided string pointer.
func (ssp *SSP) HandleFloat64Slice(h Float64SliceHandler) *SSP {
	ssp.Expressions = append(ssp.Expressions, func(cursor *Cursor) bool {
		item := cursor.Get()
		if item != "" {
			ss := strings.Split(item, ",")
			if len(ss) > 0 {
				fs := []float64{}
				for _, sv := range ss {
					fv, fverr := strconv.ParseFloat(sv, 64)
					if fverr == nil {
						fs = append(fs, fv)
					} else {
						return false
					}
				}
				if len(fs) > 0 {
					h(fs)
					cursor.Shift(1)
					return true
				}
			}
		}
		return false
	})
	return ssp
}

// HandleBool tries to parse the current item as integer value
// and store the result in the provided string pointer.
func (ssp *SSP) HandleBool(h BoolHandler) *SSP {
	ssp.Expressions = append(ssp.Expressions, func(cursor *Cursor) bool {
		f, ferr := strconv.ParseBool(cursor.Get())
		if ferr == nil {
			h(f)
			cursor.Shift(1)
			return true
		}
		return false
	})
	return ssp
}

// HandleBoolSlice tries to parse the current item as integer value
// and store the result in the provided string pointer.
func (ssp *SSP) HandleBoolSlice(h BoolSliceHandler) *SSP {
	ssp.Expressions = append(ssp.Expressions, func(cursor *Cursor) bool {
		item := cursor.Get()
		if item != "" {
			ss := strings.Split(item, ",")
			if len(ss) > 0 {
				bs := []bool{}
				for _, sv := range ss {
					bv, bverr := strconv.ParseBool(sv)
					if bverr == nil {
						bs = append(bs, bv)
					} else {
						return false
					}
				}
				if len(bs) > 0 {
					h(bs)
					cursor.Shift(1)
					return true
				}
			}
		}
		return false
	})
	return ssp
}

// CaptureString stores the current item in the provided string pointer.
func (ssp *SSP) CaptureString(v *string) *SSP {
	return ssp.HandleString(func(vv string) {
		*v = vv
	})
}

// CaptureStringSlice tries to parse the current item as a comma-separated list
// of string values and store the result in the provided string slice pointer.
func (ssp *SSP) CaptureStringSlice(v *[]string) *SSP {
	return ssp.HandleStringSlice(func(vv []string) {
		*v = vv
	})
}

// CaptureInt64 stores the current item in the provided int64 pointer.
func (ssp *SSP) CaptureInt64(v *int64) *SSP {
	return ssp.HandleInt64(func(vv int64) {
		*v = vv
	})
}

// CaptureInt64Slice tries to parse the current item as a comma-separated list
// of int64 values and store the result in the provided int64 slice pointer.
func (ssp *SSP) CaptureInt64Slice(v *[]int64) *SSP {
	return ssp.HandleInt64Slice(func(vv []int64) {
		*v = vv
	})
}

// CaptureFloat64 stores the current item in the provided float64 pointer.
func (ssp *SSP) CaptureFloat64(v *float64) *SSP {
	return ssp.HandleFloat64(func(vv float64) {
		*v = vv
	})
}

// CaptureFloat64Slice tries to parse the current item as a comma-separated list
// of float64 values and store the result in the provided float64 slice pointer.
func (ssp *SSP) CaptureFloat64Slice(v *[]float64) *SSP {
	return ssp.HandleFloat64Slice(func(vv []float64) {
		*v = vv
	})
}

// CaptureBool stores the current item in the provided bool pointer.
func (ssp *SSP) CaptureBool(v *bool) *SSP {
	return ssp.HandleBool(func(vv bool) {
		*v = vv
	})
}

// CaptureBoolSlice tries to parse the current item as a comma-separated list
// of bool values and store the result in the provided bool slice pointer.
func (ssp *SSP) CaptureBoolSlice(v *[]bool) *SSP {
	return ssp.HandleBoolSlice(func(vv []bool) {
		*v = vv
	})
}

// OneOf attempts to execute one parser from the supplied list of parsers.
func (ssp *SSP) OneOf(parsers []*SSP) *SSP {
	ssp.Expressions = append(ssp.Expressions, func(cursor *Cursor) bool {

		for _, parser := range parsers {
			cursor2 := cursor.Clone()
			parser.actualParse(cursor2)
			if parser.Complete {
				cursor.Shift(cursor2.I)
				return true
			}
		}

		return false
	})
	return ssp
}

// Repeat tries to execute the supplied parser at least min & at most max times.
func (ssp *SSP) Repeat(parser *SSP, min uint, max uint) *SSP {
	ssp.Expressions = append(ssp.Expressions, func(cursor *Cursor) bool {
		complete := uint(0)
		progress := 0
		cursor2 := cursor.Clone()
		for i := uint(0); i < max; i++ {
			parser.actualParse(cursor2)
			if parser.Complete {
				complete++
				progress += parser.Progress
			}
		}

		res := complete >= min && complete <= max

		if res {
			cursor.Shift(cursor2.I)
		}

		return res
	})
	return ssp
}
