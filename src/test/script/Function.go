package script

import (
	"regexp"
)

// Dereference checks if v is a reference to another operation
// and returns a ParameterAccess function for it.
func Dereference(v string) (bool, string, string) {
	rx := regexp.MustCompile("#(?P<opRef>\\w+)\\.response\\.(?P<selector>.*)")

	if rx.Match([]byte(v)) {
		match := rx.FindStringSubmatch(v)
		matches := make(map[string]string)
		for i, name := range rx.SubexpNames() {
			if i != 0 && name != "" {
				matches[name] = match[i]
			}
		}

		// fmt.Printf("Parameter references the '%s' data in the '%s' operation \n", matches["selector"], matches["opRef"])

		//TODO: pass-through dereferencing, anyone?
		return true, matches["opRef"], matches["selector"]
	}

	return false, "", ""

}
