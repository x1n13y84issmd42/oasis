package strings

import (
	"regexp"
	"sort"
	gostrings "strings"
)

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
	if len(strings) > 2 {
		sort.Strings(strings)
		s1 := strings[:len(strings)-1]
		s2 := strings[len(strings)-1]
		return gostrings.Join(s1, ", ") + " and " + s2
	} else if len(strings) > 1 {
		return strings[0]
	} else {
		return ""
	}
}

// SIMap is a string-interface map.
type SIMap map[string]interface{}

// Keys returns a slice containing keys of m.
func (m SIMap) Keys() []string {
	res := []string{}

	for k := range m {
		res = append(res, k)
	}
	return res
}

// RxMatches returns a map of named capture groups.
func RxMatches(s string, rx *regexp.Regexp) map[string]string {
	matches := make(map[string]string)

	if rx.Match([]byte(s)) {
		match := rx.FindStringSubmatch(s)
		for i, name := range rx.SubexpNames() {
			if i != 0 && name != "" {
				matches[name] = match[i]
			}
		}
	}

	return matches
}

// Split a string "s" to a slice of substrings, separated by delimiter "d"
func Split(s string, d string) []string {
	return gostrings.Split(s, d)
}

// Cut returns a part of string "s" between "k" and "n" positions
func Cut(s string, k int, n int) string {
	return s[k:n]
}
