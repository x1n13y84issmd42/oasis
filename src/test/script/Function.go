package script

import (
	"regexp"

	"github.com/x1n13y84issmd42/oasis/src/strings"
)

// Dereference checks if v is a reference to another operation
// and returns a ParameterAccess function for it.
func Dereference(v string) (bool, string, string) {
	rx := regexp.MustCompile("#(?P<opRef>\\w+)\\.response\\.(?P<selector>.*)")

	if rx.Match([]byte(v)) {
		matches := strings.RxMatches(v, rx)

		//TODO: pass-through dereferencing, anyone?
		return true, matches["opRef"], matches["selector"]
	}

	return false, "", ""

}
