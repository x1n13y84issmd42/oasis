package log

// IStyle is an interface to log styling.
type IStyle interface {
	Default(...interface{}) string
	URL(...interface{}) string
	Method(...interface{}) string
	Op(...interface{}) string
	OK(...interface{}) string
	Failure(...interface{}) string
	Success(...interface{}) string
	Error(...interface{}) string
	ID(...interface{}) string
	ValueExpected(...interface{}) string
	ValueActual(...interface{}) string
}
