package params

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
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
		errors.Report(err, "BodyParameters", log)
	}
	// fmt.Printf();

	data := make(map[string]interface{})
	fd := url.Values{}

	for p := range params.Iterate() {
		v := p.V()
		log.UsingParameterExample(p.N, "body", p.Source, v)
		// TODO: type & format from the spec should be used here
		vi, err := strconv.Atoi(v)
		if err == nil && vi != -1 {
			data[p.N] = vi
		} else {
			data[p.N] = v
		}
		fd.Set(p.N, v)

		// log.NOMESSAGE("BodyParameters.Enrich: Added body params %s=%s", p.N, v)
	}

	// log.NOMESSAGE("BodyParameters.Enrich: Path: %s", req.RequestURI)
	// log.NOMESSAGE("BodyParameters.Enrich: CT: %#v", req.Header["Content-Type"])

	if len(req.Header["Content-Type"]) == 0 {
		return
	}

	// TODO: move this to some kind of body encoder
	switch req.Header["Content-Type"][0] {
	case "application/x-www-form-urlencoded":
		// log.NOMESSAGE("BodyParameters.Enrich: Body: application/x-www-form-urlencoded")
		req.Body = ioutil.NopCloser(strings.NewReader(fd.Encode()))
		break

	case "application/json":
		// log.NOMESSAGE("BodyParameters.Enrich: Body: application/json")
		json, jsonErr := json.Marshal(data)
		if jsonErr != nil {
			errors.Report(jsonErr, "BodyParameters", log)
		}
		// log.NOMESSAGE("BodyParameters.Enrich: Body:\n%s", string(json))

		req.Body = ioutil.NopCloser(strings.NewReader(string(json)))

		if len(req.Header["Content-Length"]) == 0 {
			req.Header["Content-Length"] = []string{}
		}

		// CL := strconv.Itoa(len(string(json)))
		CL := "127"

		req.Header.Add("Content-Length", CL)
		log.UsingParameterExample("Content-Length", "header", "computed", CL)

		break
	}
}
