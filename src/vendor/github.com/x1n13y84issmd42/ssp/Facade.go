package ssp

// String is a shorthand function to create a String() expression.
func String(f string) *SSP {
	return New().String(f)
}

// Strings is a syntactic sugar for sequential String() calls.
func Strings(ss ...string) *SSP {
	p := New().String(ss[0])
	for i := 1; i < len(ss); i++ {
		p.String(ss[i])
	}
	return p
}

// CaptureString is a shorthand function to create a CaptureString() expression.
func CaptureString(v *string) *SSP {
	return New().CaptureString(v)
}

// CaptureStringSlice is a shorthand function to create a CaptureStringSlice() expression.
func CaptureStringSlice(v *[]string) *SSP {
	return New().CaptureStringSlice(v)
}

// CaptureInt64 is a shorthand function to create a CaptureInt64() expression.
func CaptureInt64(v *int64) *SSP {
	return New().CaptureInt64(v)
}

// CaptureInt64Slice is a shorthand function to create a CaptureInt64Slice() expression.
func CaptureInt64Slice(v *[]int64) *SSP {
	return New().CaptureInt64Slice(v)
}

// CaptureFloat64 is a shorthand function to create a CaptureFloat64() expression.
func CaptureFloat64(v *float64) *SSP {
	return New().CaptureFloat64(v)
}

// CaptureFloat64Slice is a shorthand function to create a CaptureFloat64Slice() expression.
func CaptureFloat64Slice(v *[]float64) *SSP {
	return New().CaptureFloat64Slice(v)
}

// CaptureBool is a shorthand function to create a CaptureBool() expression.
func CaptureBool(v *bool) *SSP {
	return New().CaptureBool(v)
}

// CaptureBoolSlice is a shorthand function to create a CaptureBoolSlice() expression.
func CaptureBoolSlice(v *[]bool) *SSP {
	return New().CaptureBoolSlice(v)
}

// OneOf is a shorthand function to create a OneOf() expression.
func OneOf(parsers ...*SSP) *SSP {
	return New().OneOf(parsers)
}

// Repeat is a shorthand function to create a Repeat() expression.
func Repeat(parser *SSP, min uint, max uint) *SSP {
	return New().Repeat(parser, min, max)
}
