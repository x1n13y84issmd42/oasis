package log

import (
	"fmt"
)

// BufferedStdOut is a buffered output for log.
// It accumulates log output to flush it all at once.
// Useful for concurrent script execution.
type BufferedStdOut struct {
	data string
}

// NewBufferedStdOut creates a new Buffer instance.
func NewBufferedStdOut() *BufferedStdOut {
	return &BufferedStdOut{}
}

// Flush flushes the accumulated output to stdout.
func (buffer *BufferedStdOut) Flush() {
	// fmt.Print("Flushing\n")
	fmt.Printf(buffer.data)
	buffer.data = ""
}

// Print collects the messages into an internal buffer.
func (buffer *BufferedStdOut) Print(msg string, args ...interface{}) {
	data := fmt.Sprintf(msg, args...)
	buffer.data += data
}

// StdOut is a standard terminal output for Log.
type StdOut struct {
	///
}

// NewStdOut creates a new StdOut instance.
func NewStdOut() StdOut {
	return StdOut{}
}

// Print collects the messages into an internal StdOut.
func (out StdOut) Print(msg string, args ...interface{}) {
	fmt.Printf(msg, args...)
}

// Flush does nothing for regular output.
func (out StdOut) Flush() {
}
