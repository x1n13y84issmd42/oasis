package log

import (
	"github.com/gookit/color"
	"github.com/x1n13y84issmd42/oasis/src/contract"
)

// Festive - a colorized test execution log style.
type Festive struct{}

// NewFestive creates a new Festive logger style instance.
func NewFestive(level int64) contract.Logger {
	return &Log{
		Level: level,
		Style: Festive{},
	}
}

// Default marks up the default output.
func (log Festive) Default(args ...interface{}) string {
	return color.New(38, 5, 218).Sprint(args...) + "\x1b[K"
}

// URL marks up URLs.
func (log Festive) URL(args ...interface{}) string {
	return color.New(color.FgCyan, color.OpUnderscore).Sprint(args...) + "\x1b[K"
}

// Method marks up HTTP methods.
func (log Festive) Method(args ...interface{}) string {
	return color.New(color.FgCyan).Sprint(args...) + "\x1b[K"
}

// Op marks up operation names.
func (log Festive) Op(args ...interface{}) string {
	return color.New(color.FgYellow).Sprint(args...) + "\x1b[K"
}

// OK marks up the successful operation testing result.
func (log Festive) OK(args ...interface{}) string {
	return color.New(color.FgLightWhite, color.BgGreen).Sprint(args...) + "\x1b[K"
}

// Failure marks up the unsuccesfful operation testing result.
func (log Festive) Failure(args ...interface{}) string {
	return color.New(color.FgLightWhite, color.BgRed).Sprint(args...) + "\x1b[K"
}

// Success marks up the successful testing of operation properties.
func (log Festive) Success(args ...interface{}) string {
	return color.New(color.FgGreen).Sprint(args...) + "\x1b[K"
}

// Success marks up the unsuccessful testing of operation properties.
func (log Festive) Error(args ...interface{}) string {
	return color.New(color.FgRed).Sprint(args...) + "\x1b[K"
}

// ID marks up indentifiers.
func (log Festive) ID(args ...interface{}) string {
	return color.New(color.FgLightWhite, color.OpUnderscore).Sprint(args...) + "\x1b[K"
}

// ValueExpected marks up expected values.
func (log Festive) ValueExpected(args ...interface{}) string {
	return color.New(48, 5, 2, 38, 5, 0).Sprint(args...) + "\x1b[K"
}

// ValueActual marks up actual received values.
func (log Festive) ValueActual(args ...interface{}) string {
	return color.New(color.FgLightWhite, color.BgRed).Sprint(args...) + "\x1b[K"
}

// Value marks up values.
func (log Festive) Value(args ...interface{}) string {
	return color.New(48, 5, 240, 38, 5, 255).Sprint(args...) + "\x1b[K"
}
