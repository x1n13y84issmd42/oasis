package params

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/x1n13y84issmd42/oasis/src/contract"
	"github.com/x1n13y84issmd42/oasis/src/errors"
)

// BodyParameters is the source for request header parameters.
type BodyParameters struct {
	contract.EntityTrait
	*MultiSet
}

// Body creates a new BodyParameters instance.
func Body(log contract.Logger) *BodyParameters {
	p := &BodyParameters{
		EntityTrait: contract.Entity(log),
		MultiSet:    NewMultiSet("body"),
	}

	return p
}

// Enrich applies the parameters as header values to the request.
func (params BodyParameters) Enrich(req *http.Request, log contract.Logger) {
	if err := params.Validate(); err != nil {
		errors.Report(err, "BodyParameters", params.Log)
	}
	// fmt.Printf();

	data := make(map[string]interface{})
	fd := url.Values{}

	for p := range params.Iterate() {
		v := p.V()
		params.Log.UsingParameterExample(p.N, "body", p.Source, v)
		data[p.N] = v
		fd.Set(p.N, v)
	}

	// TODO: this actually depends on request content-type
	req.Body = ioutil.NopCloser(strings.NewReader(fd.Encode()))
}
