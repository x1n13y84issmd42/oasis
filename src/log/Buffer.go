package log

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/x1n13y84issmd42/oasis/src/contract"
	"github.com/xeipuuv/gojsonschema"
)

// Buffer accumulates log output to flush it all at once.
// Useful for concurrent script execution.
type Buffer struct {
	Logger contract.Logger
	Buffer bytes.Buffer
}

// NewBuffer creates a new Buffer instance.
func NewBuffer(logger contract.Logger) contract.Logger {
	return &Buffer{
		Logger: logger,
	}
}

// Capture captures the terminal output.
func (buffer *Buffer) Capture(f func()) {
	reader, writer, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	stdout := os.Stdout

	restore := func() {
		os.Stdout = stdout
		// fmt.Println("Restored the os.Stdout")
	}

	os.Stdout = writer

	out := make(chan string)
	wg := new(sync.WaitGroup)
	wg.Add(1)

	go func() {
		var buf bytes.Buffer
		wg.Done()
		io.Copy(&buf, reader)
		out <- buf.String()
	}()

	wg.Wait()
	// fmt.Println("Calling f()")
	f()
	// fmt.Println("Returned from f()")
	writer.Close()

	output := <-out
	restore()
	// fmt.Printf("Buffered output: '%s'\n", output)
	buffer.Buffer.WriteString(output)
}

// Flush flushes the accumulated output to stdout.
func (buffer *Buffer) Flush() {
	// fmt.Println("Flushing")
	// fmt.Printf("The buffer is '%s'\n", buffer.Buffer.String())
	fmt.Print(buffer.Buffer.String())
}

// NOMESSAGE ...
func (buffer *Buffer) NOMESSAGE(msg string, args ...interface{}) {
	buffer.Capture(func() {
		buffer.Logger.NOMESSAGE(msg, args...)
	})
}

// Parameters ...
func (buffer *Buffer) Parameters(name string, params contract.ParameterSource) {
	buffer.Capture(func() {
		buffer.Logger.Parameters(name, params)
	})
}

// Print ...
func (buffer *Buffer) Print(l int64, msg string, args ...interface{}) {
	buffer.Capture(func() {
		buffer.Logger.Print(l, msg, args)
	})
}

// Usage ...
func (buffer *Buffer) Usage() {
	buffer.Capture(func() {
		buffer.Logger.Usage()
	})
}

// Error ...
func (buffer *Buffer) Error(err error) {
	buffer.Capture(func() {
		buffer.Logger.Error(err)
	})
}

// LoadingSpec ...
func (buffer *Buffer) LoadingSpec(path string) {
	buffer.Capture(func() {
		buffer.Logger.LoadingSpec(path)
	})
}

// LoadingScript ...
func (buffer *Buffer) LoadingScript(path string) {
	buffer.Capture(func() {
		buffer.Logger.LoadingScript(path)
	})
}

// PrintOperations ...
func (buffer *Buffer) PrintOperations(ops contract.OperationIterator) {
	buffer.Capture(func() {
		buffer.Logger.PrintOperations(ops)
	})
}

// TestingProject ...
func (buffer *Buffer) TestingProject(p contract.ProjectInfo) {
	buffer.Capture(func() {
		buffer.Logger.TestingProject(p)
	})
}

// TestingOperation ...
func (buffer *Buffer) TestingOperation(res contract.Operation) {
	buffer.Capture(func() {
		buffer.Logger.TestingOperation(res)
	})
}

// UsingSecurity ...
func (buffer *Buffer) UsingSecurity(sec contract.Security) {
	buffer.Capture(func() {
		buffer.Logger.UsingSecurity(sec)
	})
}

// SecurityHasNoData ...
func (buffer *Buffer) SecurityHasNoData(sec contract.Security) {
	buffer.Capture(func() {
		buffer.Logger.SecurityHasNoData(sec)
	})
}

// Requesting ...
func (buffer *Buffer) Requesting(method string, url string) {
	buffer.Capture(func() {
		buffer.Logger.Requesting(method, url)
	})
}

// UsingParameterExample ...
func (buffer *Buffer) UsingParameterExample(paramName string, in string, container string, value string) {
	buffer.Capture(func() {
		buffer.Logger.UsingParameterExample(paramName, in, container, value)
	})
}

// Expecting ...
func (buffer *Buffer) Expecting(what string, v string) {
	buffer.Capture(func() {
		buffer.Logger.Expecting(what, v)
	})
}

// ExpectingProperty ...
func (buffer *Buffer) ExpectingProperty(what string, v string) {
	buffer.Capture(func() {
		buffer.Logger.ExpectingProperty(what, v)
	})
}

// HeaderHasNoValue ...
func (buffer *Buffer) HeaderHasNoValue(hdr string) {
	buffer.Capture(func() {
		buffer.Logger.HeaderHasNoValue(hdr)
	})
}

// ResponseHasWrongStatus ...
func (buffer *Buffer) ResponseHasWrongStatus(expectedStatus int, actualStatus int) {
	buffer.Capture(func() {
		buffer.Logger.ResponseHasWrongStatus(expectedStatus, actualStatus)
	})
}

// ResponseHasWrongContentType ...
func (buffer *Buffer) ResponseHasWrongContentType(expectedCT string, actualCT string) {
	buffer.Capture(func() {
		buffer.Logger.ResponseHasWrongContentType(expectedCT, actualCT)
	})
}

// ResponseHasWrongPropertyValue ...
func (buffer *Buffer) ResponseHasWrongPropertyValue(propName string, expected string, actual string) {
	buffer.Capture(func() {
		buffer.Logger.ResponseHasWrongPropertyValue(propName, expected, actual)
	})
}

// OperationOK ...
func (buffer *Buffer) OperationOK() {
	buffer.Capture(func() {
		buffer.Logger.OperationOK()
	})
}

// OperationFail ...
func (buffer *Buffer) OperationFail() {
	buffer.Capture(func() {
		buffer.Logger.OperationFail()
	})
}

// SchemaOK ...
func (buffer *Buffer) SchemaOK(schemaName string) {
	buffer.Capture(func() {
		buffer.Logger.SchemaOK(schemaName)
	})
}

// SchemaFail ...
func (buffer *Buffer) SchemaFail(schemaName string, errors []gojsonschema.ResultError) {
	buffer.Capture(func() {
		buffer.Logger.SchemaFail(schemaName, errors)
	})
}

// ScriptExecutionStart ...
func (buffer *Buffer) ScriptExecutionStart(node string) {
	buffer.Capture(func() {
		buffer.Logger.ScriptExecutionStart(node)
	})
}

// XError ...
func (buffer *Buffer) XError(err error, style contract.LogStyle, tab contract.TabFn) {
	buffer.Capture(func() {
		buffer.Logger.XError(err, style, tab)
	})
}
