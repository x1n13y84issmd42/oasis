package params

import (
	"regexp"

	"github.com/x1n13y84issmd42/oasis/src/contract"
)

// URLParameters is the source for URL path parameters.
// URLParameters have an implicit requirement for the @HOSTNAME parameter
// which is an API host name.
type URLParameters struct {
	contract.EntityTrait
	*Set

	Path string
}

// URL creates a new URLParameters instance.
func URL(path string, log contract.Logger) *URLParameters {
	p := &URLParameters{
		EntityTrait: contract.Entity(log),
		Set:         NewSet(),
		Path:        path,
	}

	p.Require(KeyHost)

	return p
}

// Make creates a URL string value from path template.
func (params URLParameters) String() string {
	if err := params.Validate(); err != nil {
		params.Error(err)
	}

	tpl := "{" + KeyHost + "}" + params.Path

	for pt := range params.Iterate() {
		rx := regexp.MustCompile("\\{" + pt.N + "\\}")

		if rx.Match([]byte(tpl)) {
			if pt.V != "" {
				tpl = string(rx.ReplaceAll([]byte(tpl), []byte(pt.V)))
			}
		}
	}

	return tpl
}
