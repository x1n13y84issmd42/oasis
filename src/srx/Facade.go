package srx

// String is a shorthand function to create a String() expression.
func String(f string) *SRX {
	return NewSRX().String(f)
}

// Strings is a syntactic sugar for sequential String() calls.
func Strings(ss ...string) *SRX {
	p := NewSRX().String(ss[0])
	for i := 1; i < len(ss); i++ {
		p.String(ss[i])
	}
	return p
}

// CaptureString is a shorthand function to create a CaptureString() expression.
func CaptureString(v *string) *SRX {
	return NewSRX().CaptureString(v)
}

// CaptureStringSlice is a shorthand function to create a CaptureStringSlice() expression.
func CaptureStringSlice(v *[]string) *SRX {
	return NewSRX().CaptureStringSlice(v)
}

// CaptureInt64 is a shorthand function to create a CaptureInt64() expression.
func CaptureInt64(v *int64) *SRX {
	return NewSRX().CaptureInt64(v)
}

// CaptureInt64Slice is a shorthand function to create a CaptureInt64Slice() expression.
func CaptureInt64Slice(v *[]int64) *SRX {
	return NewSRX().CaptureInt64Slice(v)
}

// CaptureFloat64 is a shorthand function to create a CaptureFloat64() expression.
func CaptureFloat64(v *float64) *SRX {
	return NewSRX().CaptureFloat64(v)
}

// CaptureFloat64Slice is a shorthand function to create a CaptureFloat64Slice() expression.
func CaptureFloat64Slice(v *[]float64) *SRX {
	return NewSRX().CaptureFloat64Slice(v)
}

// CaptureBool is a shorthand function to create a CaptureBool() expression.
func CaptureBool(v *bool) *SRX {
	return NewSRX().CaptureBool(v)
}

// CaptureBoolSlice is a shorthand function to create a CaptureBoolSlice() expression.
func CaptureBoolSlice(v *[]bool) *SRX {
	return NewSRX().CaptureBoolSlice(v)
}

// OneOf is a shorthand function to create a OneOf() expression.
func OneOf(parsers ...*SRX) *SRX {
	return NewSRX().OneOf(parsers)
}

// Repeat is a shorthand function to create a Repeat() expression.
func Repeat(parser *SRX, min uint, max uint) *SRX {
	return NewSRX().Repeat(parser, min, max)
}
