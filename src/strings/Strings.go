package strings

import gostrings "strings"

// StringMapFn is a function to map strings to strings.
type StringMapFn = func(string) string

// Map creates a new string slice by applying mapfn to strings.
func Map(strings []string, mapfn StringMapFn) []string {
	res := []string{}
	for _, s := range strings {
		res = append(res, mapfn(s))
	}

	return res
}

// EnumJoin joins a string slice into a human-readable enumeration,
// so a [1, 2, 3] slice becomes a "1, 2 and 3" string.
func EnumJoin(strings []string) string {
	s1 := strings[:len(strings)-1]
	s2 := strings[len(strings)-1]
	return gostrings.Join(s1, ", ") + " and " + s2
}
