package log

import (
	"fmt"
)

// Plain - a plain & simple test execution logger style.
type Plain struct{}

// NewPlain creates a new Plain logger style instance.
func NewPlain(level int64) ILogger {
	return Log{
		Level: level,
		Style: Plain{},
	}
}

// Default marks up the default output.
func (log Plain) Default(args ...interface{}) string {
	return fmt.Sprint(args...)
}

// URL marks up URLs.
func (log Plain) URL(args ...interface{}) string {
	return fmt.Sprint(args...)
}

// Method marks up HTTP methods.
func (log Plain) Method(args ...interface{}) string {
	return fmt.Sprint(args...)
}

// Op marks up operation names.
func (log Plain) Op(args ...interface{}) string {
	return fmt.Sprint(args...)
}

// OK marks up the successful operation testing result.
func (log Plain) OK(args ...interface{}) string {
	return fmt.Sprint(args...)
}

// Failure marks up the unsuccesfful operation testing result.
func (log Plain) Failure(args ...interface{}) string {
	return fmt.Sprint(args...)
}

// Success marks up the successful testing of operation properties.
func (log Plain) Success(args ...interface{}) string {
	return fmt.Sprint(args...)
}

// Success marks up the unsuccessful testing of operation properties.
func (log Plain) Error(args ...interface{}) string {
	return fmt.Sprint(args...)
}

// ID marks up indentifiers.
func (log Plain) ID(args ...interface{}) string {
	return fmt.Sprint(args...)
}

// ValueExpected marks up expected values.
func (log Plain) ValueExpected(args ...interface{}) string {
	return fmt.Sprint(args...)
}

// ValueActual marks up actual received values.
func (log Plain) ValueActual(args ...interface{}) string {
	return fmt.Sprint(args...)
}
