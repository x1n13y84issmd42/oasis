package srx

import (
	"strconv"
	"strings"
)

// ExpressionFn is a single parser function.
// SRX instances contain multiple of these.
type ExpressionFn func(*Cursor) bool

// SRX is regular expression on slices.
// It allows to parse slices containing structured or patterned data and extract that.
// For example, if you have a slice of strings:
//		inputSlice := []string{"er", "er", "er", "er", "yolo", "1,2,3,4"}
// you can make an expression for it:
//		rx := srx.Repeat(srx.Flag("er"), 2, 200).CaptureString(&theYoloString).CaptureStringSlice(&the1234StringSlice)
//		rx.Parse(inputSlice)
// then you can parse your slice and, first, make sure it fits the pattern, and second extract the "yolo" string
// and a slice of "1", "2", "3" & "4" strings.
type SRX struct {
	Expressions []ExpressionFn
	Progress    int
	Complete    bool
}

// NewSRX creates a new SRX instance.
func NewSRX() *SRX {
	return &SRX{
		Expressions: []ExpressionFn{},
	}
}

// Parse starts the parsing process.
func (srx *SRX) Parse(args []string) *SRX {
	return srx.actualParse(&Cursor{Args: args})
}

func (srx *SRX) actualParse(cursor *Cursor) *SRX {

	srx.Complete = false
	srx.Progress = 0

	if cursor.Left() <= 0 {
		return srx
	}

	cpI := 0
	curParser := srx.Expressions[cpI]

	for curParser(cursor) {
		cpI++
		srx.Progress = cursor.I
		if cpI >= len(srx.Expressions) {
			srx.Complete = true
			break
		}

		if cursor.I >= cursor.Len() {
			// Consumed entire input.
		}

		curParser = srx.Expressions[cpI]
	}

	return srx
}

// Flag parses a single string value.
func (srx *SRX) Flag(f string) *SRX {
	srx.Expressions = append(srx.Expressions, func(cursor *Cursor) bool {
		res := f == cursor.Get()
		if res {
			cursor.Shift(1)
		}
		return res
	})
	return srx
}

// CaptureString stores the current item in the provided string pointer.
func (srx *SRX) CaptureString(v *string) *SRX {
	srx.Expressions = append(srx.Expressions, func(cursor *Cursor) bool {
		*v = cursor.Get()
		cursor.Shift(1)
		return true
	})
	return srx
}

// CaptureInt64 tries to parse the current item as integer value
// and store the result in the provided string pointer.
func (srx *SRX) CaptureInt64(v *int64) *SRX {
	srx.Expressions = append(srx.Expressions, func(cursor *Cursor) bool {
		i, ierr := strconv.ParseInt(cursor.Get(), 10, 64)
		if ierr == nil {
			*v = i
			cursor.Shift(1)
		}
		return true
	})
	return srx
}

// CaptureStringSlice tries to parse the current item as a list of string values
// and store the result in the provided string slice pointer.
func (srx *SRX) CaptureStringSlice(v *[]string) *SRX {
	srx.Expressions = append(srx.Expressions, func(cursor *Cursor) bool {
		item := cursor.Get()
		if item != "" {
			values := strings.Split(item, ",")
			if len(values) > 0 {
				*v = values
				cursor.Shift(1)
				return true
			}
		}

		return false
	})
	return srx
}

// OneOf attempts to execute one parser from the supplied list of parsers.
func (srx *SRX) OneOf(parsers []*SRX) *SRX {
	srx.Expressions = append(srx.Expressions, func(cursor *Cursor) bool {

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
	return srx
}

// Repeat tries to execute the supplied parser at least min & at most max times.
func (srx *SRX) Repeat(parser *SRX, min uint, max uint) *SRX {
	srx.Expressions = append(srx.Expressions, func(cursor *Cursor) bool {
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
	return srx
}
