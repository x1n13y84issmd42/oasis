package utility

import (
	"fmt"
	"reflect"
	"strings"
)

// Context represents the current path when validating nested data structures.
type Context struct {
	Path []string
	V    interface{}
}

func NewContext(n string) *Context {
	return &Context{
		Path: []string{n},
	}
}

// PushIndex creates a new Context with an added array index value to the path.
func (ctx Context) PushIndex(i int, v interface{}) *Context {
	return &Context{append(ctx.Path, fmt.Sprintf("[%d]", i)), v}
}

// PushProperty creates a new Context with an added object property value to the path.
func (ctx Context) PushProperty(s string, v interface{}) *Context {
	return &Context{append(ctx.Path, fmt.Sprintf(".%s", s)), v}
}

func (ctx Context) String() string {
	return strings.Join(ctx.Path, "")
}

// CurrentValueType return the underlying type of V.
func (ctx Context) CurrentValueType() string {
	if ctx.V != nil {
		t := reflect.TypeOf(ctx.V)
		return t.Name()
	}

	return "null"
}
