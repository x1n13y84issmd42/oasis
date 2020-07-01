package expect

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/contract"
	"github.com/x1n13y84issmd42/oasis/src/test"
)

// Status creates an expectation as for response's status code.
func Status(status int) contract.Expectation {
	fmt.Printf("\tExpecting status %d\n", status)
	return func(resp *http.Response, log contract.Logger) bool {
		if resp.StatusCode == status {
			return true
		}

		log.ResponseHasWrongStatus(status, resp.StatusCode)
		return false
	}
}

// HeaderRequired creates an expectation as for response's header value.
func HeaderRequired(n string) contract.Expectation {
	fmt.Println("\tExpecting required header " + n)
	return func(resp *http.Response, log contract.Logger) bool {
		if resp.Header.Get(n) != "" {
			return true
		}

		log.HeaderHasNoValue(n)
		return false
	}
}

// HeaderSchema creates an expectation as for response's header contents
// which must comply to the provided JSON schema.
func HeaderSchema(n string, schema *api.Schema) contract.Expectation {
	fmt.Println("\tExpecting header " + n + " to conform schema " + schema.Name)
	return func(resp *http.Response, log contract.Logger) bool {
		return test.Schema(schema.Cast(resp.Header.Get(n)), schema, log)
	}
}

// ContentType creates an expectation as for response's content type.
func ContentType(v string) contract.Expectation {
	fmt.Println("\tExpecting Content-Type " + v)
	return func(resp *http.Response, log contract.Logger) bool {
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
func ContentSchema(schema *api.Schema) contract.Expectation {
	fmt.Println("\tExpecting content schema " + schema.Name)
	return func(resp *http.Response, log contract.Logger) bool {
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
