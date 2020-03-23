package srx

// Flag is a shorthand function to create a Flag() expression.
func Flag(f string) *SRX {
	return NewSRX().Flag(f)
}

// CaptureString is a shorthand function to create a CaptureString() expression.
func CaptureString(v *string) *SRX {
	return NewSRX().CaptureString(v)
}

// CaptureInt64 is a shorthand function to create a CaptureInt64() expression.
func CaptureInt64(v *int64) *SRX {
	return NewSRX().CaptureInt64(v)
}

// CaptureStringSlice is a shorthand function to create a CaptureStringSlice() expression.
func CaptureStringSlice(v *[]string) *SRX {
	return NewSRX().CaptureStringSlice(v)
}

// OneOf is a shorthand function to create a OneOf() expression.
func OneOf(parsers []*SRX) *SRX {
	return NewSRX().OneOf(parsers)
}

// Repeat is a shorthand function to create a Repeat() expression.
func Repeat(parser *SRX, min uint, max uint) *SRX {
	return NewSRX().Repeat(parser, min, max)
}
