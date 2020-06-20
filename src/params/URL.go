package params

import (
	"net/http"
	"net/url"
	"regexp"

	"github.com/x1n13y84issmd42/oasis/src/contract"
)

// URLParameters is the default source for URL path parameters.
// It treats it's sources hierarchically, i.e. values are taken
// in this order:
// 		test output
// 		CLI input
// 		spec op
// 		spec path
//
// URLParameters have an implicit requirement for the @HOSTNAME parameter
// which is an API host name.
type URLParameters struct {
	*contract.ParameterProviderPrototype
	contract.EntityTrait

	Path string
}

// URL creates a new URLParameters instance.
func URL(path string, log contract.Logger) *URLParameters {
	p := &URLParameters{
		ParameterProviderPrototype: contract.NewParameterProviderPrototype(),
		EntityTrait:                contract.Entity(log),
		Path:                       path,
	}

	p.Require("@HOSTNAME")

	return p
}

// Enrich creates a URL string value from path template
// and sets it to the internal http.Request instance.
func (params URLParameters) Enrich(req *http.Request) {
	if err := params.Validate(); err != nil {
		params.Error(err)
	}

	tpl := "@HOSTNAME/" + params.Path

	for pt := range params.Iterate() {
		RX := regexp.MustCompile("\\{" + pt[0] + "\\}")

		if RX.Match([]byte(tpl)) {
			if pt[1] != "" {
				tpl = string(RX.ReplaceAll([]byte(tpl), []byte(pt[1])))
				// spec.Log.UsingParameterExample(pt[0], "path", container)
			} else {
				// spec.Log.ParameterHasNoExample(pt[0], "path", container)
			}
		}
	}

	u, err := url.Parse(tpl)
	if err != nil {
		params.Error(err)
	} else {
		req.URL = u
		req.Host = u.Host
	}
}
