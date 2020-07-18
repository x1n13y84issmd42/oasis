package errors

import (
	"fmt"
	"runtime"
	"strconv"
	gostrings "strings"

	"github.com/x1n13y84issmd42/gog/graph/collection"
	"github.com/x1n13y84issmd42/oasis/src/contract"
	"github.com/x1n13y84issmd42/oasis/src/strings"
)

// Report reports errors.
func Report(err error, source string, log contract.Logger) {
	if err != nil {
		log.Error(err)
	} else {
		log.Error(NewBase(nil, "Reporting nothing."))
	}

	panic(gostrings.Title(source) + " has panicked.\nSee the error message reported above for details.")
}

// IError is an error interface.
type IError interface {
	Error() string
	Cause() error
	Caller() string
}

// Base is a generic error used within Oasis.
// Having it's root cause is what makes it good.
type Base struct {
	TheCause  error
	Details   string
	TheCaller string
}

// NewBase creates a new Base instance.
func NewBase(cause error, details string) Base {
	caller := ""
	if _, f, l, ok := runtime.Caller(2); ok {
		caller = f + ":" + strconv.Itoa(l)
	}

	if ghi := gostrings.Index(caller, "github.com"); ghi != -1 {
		caller = caller[ghi:]
	}

	return Base{
		TheCause:  cause,
		Details:   details,
		TheCaller: caller,
	}
}

// Error returns a string representation of an error.
func (err Base) Error() string {
	return "---"
}

// Cause returns an cause error instance.
func (err Base) Cause() error {
	return err.TheCause
}

// Caller returns path t a file where the error has been created.
func (err Base) Caller() string {
	return err.TheCaller
}

// ErrOops is a generic error with a message.
type ErrOops struct {
	Base
}

func (err ErrOops) Error() string {
	return err.Details
}

// Oops creates a new ErrOops error instance.
func Oops(msg string, cause error) ErrOops {
	return ErrOops{
		Base: NewBase(cause, msg),
	}
}

// ErrOperationMalformed is returned when there is not enough data in the API spec,
// or the data is incorrect, which makes operation creation impossible.
type ErrOperationMalformed struct {
	Base
	OpID    string
	Details string
}

func (err ErrOperationMalformed) Error() string {
	return "Operation '" + err.OpID + "' has malformed or incomplete data. " + err.Base.Details
}

// OperationMalformed creates a new ErrOperationMalformed error instance.
func OperationMalformed(id string, details string, cause error) ErrOperationMalformed {
	return ErrOperationMalformed{
		Base: NewBase(cause, details),
		OpID: id,
	}
}

// ErrNoParameters happens when there are parameters missing
// which are needed to build an operation before testing.
// Usually happens in paths, queries & headers.
type ErrNoParameters struct {
	Base
	Kind          string
	MissingParams []string
}

func (err ErrNoParameters) Error() string {
	hint := ""
	if len(err.MissingParams) >= 1 {
		hint = "\nWhen in manual mode, try adding 'use path parameters " + err.MissingParams[0] + "=VALUE,...' to the command line."
		hint += "\nWhen in script mode, add a 'use:path:" + err.MissingParams[0] + ":VALUE' block to the operation in the script file."
	}

	if len(err.MissingParams) == 1 {
		return "Parameter '" + err.MissingParams[0] + "' (in " + err.Kind + ") is missing it's value." + hint
	}

	missing := strings.Map(err.MissingParams, func(s string) string {
		return "'" + s + "'"
	})

	return "Parameters " + strings.EnumJoin(missing) + "' (in " + err.Kind + ") are missing their values." + hint
}

// NoParameters creates a new ErrNoData error instance.
func NoParameters(missingParams []string, kind string, cause error) ErrNoParameters {
	return ErrNoParameters{
		Base:          NewBase(cause, ""),
		Kind:          kind,
		MissingParams: missingParams,
	}
}

// ErrNotFound may be returned from host-related methods of specs.
type ErrNotFound struct {
	Base
	What string
	Name string
}

func (err ErrNotFound) Error() string {
	return err.What + " '" + err.Name + "' not found in the spec."
}

// NotFound creates a new ErrHostNotFound error instance.
func NotFound(what string, hn string, cause error) ErrNotFound {
	return ErrNotFound{
		Base: NewBase(cause, ""),
		What: what,
		Name: hn,
	}
}

// ErrInvalidSchema may be returned from schema-related methods of specs.
type ErrInvalidSchema struct {
	Base
	SchemaName string
}

func (err ErrInvalidSchema) Error() string {
	return "Schema " + err.SchemaName + " is invalid. " + err.Base.Details
}

// InvalidSchema creates a new ErrInvalidSchema error instance.
func InvalidSchema(sn string, details string, cause error) ErrInvalidSchema {
	return ErrInvalidSchema{
		Base:       NewBase(cause, details),
		SchemaName: sn,
	}
}

// ErrInvalidResponse may be returned from response-related methods of specs.
type ErrInvalidResponse struct {
	Base
}

func (err ErrInvalidResponse) Error() string {
	return "Response spec is invalid. " + err.Base.Details
}

// InvalidResponse creates a new ErrInvalidResponse error instance.
func InvalidResponse(details string, cause error) ErrInvalidResponse {
	return ErrInvalidResponse{
		Base: NewBase(cause, details),
	}
}

// ErrSecurityNotFound may be returned from security-related methods of specs.
type ErrSecurityNotFound struct {
	Base
	Name string
}

func (err ErrSecurityNotFound) Error() string {
	return "Security definition '" + err.Name + "' not found in the spec. " + err.Details
}

// SecurityNotFound creates a new ErrSecurityNotFound error instance.
func SecurityNotFound(sn string, details string, cause error) ErrSecurityNotFound {
	return ErrSecurityNotFound{
		Base: NewBase(cause, details),
		Name: sn,
	}
}

// ErrGraphHasCycles means an execution graph is not a DAG.
type ErrGraphHasCycles struct {
	Base
	Cycles *collection.NodeStack
}

func (err ErrGraphHasCycles) Error() string {
	return err.Details
}

// GraphHasCycles creates a new ErrGraphHasCycles error instance.
func GraphHasCycles(cycle *collection.NodeStack, cause error) ErrGraphHasCycles {
	cycleString := ""

	// The comp.Cycle() function adds the starting node to the stack,
	// so if there's only two of them it means we're dealing with
	// a self-referenced node. Leaving it for clarity and nice formatting,
	// otherwise popping the clutter away.
	if len(*cycle) > 2 {
		cycle.Pop()
	}

	// Formatting the cycle nicely.
	for i, n := range *cycle {
		cycleString += "  "
		if i == 0 {
			cycleString += "┌>"
		} else if i == len(*cycle)-1 {
			cycleString += "└─"
		} else {
			cycleString += "│ "
		}

		cycleString += string(n.ID())

		if i < len(*cycle)-1 {
			cycleString += "\n"
			cycleString += "  │ ↓"
			cycleString += "\n"
		}
	}

	msg1 := "The execution graph contains a cycle.\n"
	msg2 := "\nThe sequence of operations in the script must be a DAG."

	return ErrGraphHasCycles{
		Base:   NewBase(cause, msg1+cycleString+msg2),
		Cycles: cycle,
	}
}

// ErrOutOfRange happens when an array index is out of range.
type ErrOutOfRange struct {
	Base
	i int
	a *[]interface{}
}

func (err ErrOutOfRange) Error() string {
	return fmt.Sprintf("Array index %d is out of range 0-%d.", err.i, len(*err.a)-1)
}

// OutOfRange creates a new ErrOutOfRange error instance.
func OutOfRange(i int, a *[]interface{}, cause error) ErrOutOfRange {
	return ErrOutOfRange{
		Base: NewBase(cause, ""),
		i:    i,
		a:    a,
	}
}

// ErrNoProperty happens when a property was not found in an object.
type ErrNoProperty struct {
	Base
	prop string
}

func (err ErrNoProperty) Error() string {
	return "Property '" + err.prop + "' is not found in the object."
}

// NoProperty creates a new ErrNoProperty error instance.
func NoProperty(prop string, cause error) ErrNoProperty {
	return ErrNoProperty{
		Base: NewBase(cause, ""),
		prop: prop,
	}
}

// ErrNotAn happens when something is not what it appears to be.
type ErrNotAn struct {
	Base
	whatShouldBe string
	whatIs       interface{}
}

func (err ErrNotAn) Error() string {
	return "Not an " + err.whatShouldBe + ", but a " + fmt.Sprintf("%T", err.whatIs) + "."
}

// NotAn creates a new ErrNotAn error instance.
func NotAn(whatShouldBe string, whatIs interface{}, cause error) ErrNotAn {
	return ErrNotAn{
		Base:         NewBase(cause, ""),
		whatShouldBe: whatShouldBe,
		whatIs:       whatIs,
	}
}
