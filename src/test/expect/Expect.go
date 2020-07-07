package expect

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/contract"
	"github.com/x1n13y84issmd42/oasis/src/test"
)

// Status creates an expectation as for response's status code.
func Status(status int, log contract.Logger) contract.Expectation {
	log.Expecting("status", strconv.Itoa(status))

	return func(resp *http.Response) bool {
		if resp.StatusCode == status {
			return true
		}

		log.ResponseHasWrongStatus(status, resp.StatusCode)
		return false
	}
}

// HeaderRequired creates an expectation as for response's header value.
func HeaderRequired(n string, log contract.Logger) contract.Expectation {
	log.Expecting("required header", n)

	return func(resp *http.Response) bool {
		if resp.Header.Get(n) != "" {
			return true
		}

		log.HeaderHasNoValue(n)
		return false
	}
}

// HeaderSchema creates an expectation as for response's header contents
// which must comply to the provided JSON schema.
func HeaderSchema(n string, schema *api.Schema, log contract.Logger) contract.Expectation {
	log.Expecting("header "+n+" to conform schema", schema.Name)

	return func(resp *http.Response) bool {
		return test.Schema(schema.Cast(resp.Header.Get(n)), schema, log)
	}
}

// ContentType creates an expectation as for response's content type.
func ContentType(v string, log contract.Logger) contract.Expectation {
	log.Expecting("Content-Type", v)

	return func(resp *http.Response) bool {
		// This is to get rid of the possible "; charset=utf-8" part.
		respCT := strings.Split(resp.Header.Get("Content-Type"), ";")[0]
		if respCT == v {
			return true
		}

		log.ResponseHasWrongContentType(v, resp.Header.Get("Content-Type"))
		return false
	}
}

// ContentSchema creates an expectation as for response's content body
// which must comply to the provided JSON schema.
func ContentSchema(schema *api.Schema, log contract.Logger) contract.Expectation {
	log.Expecting("content schema", schema.Name)

	return func(resp *http.Response) bool {
		respCT := strings.Split(resp.Header.Get("Content-Type"), ";")[0]

		switch respCT {
		case "application/json":
			if test.JSONResponse(resp, schema, log) {
				return true
			}

		default:
			log.NOMESSAGE("The Content-Type of '%s' is not supported.\n", respCT)
		}

		log.NOMESSAGE("Response data is wrong.")
		return false
	}
}
