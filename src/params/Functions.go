package params

// Value is the default pass-through function to provide parameters.
func Value(v string) func() string {
	return func() string {
		return v
	}
}
