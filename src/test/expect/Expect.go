package expect

import (
	"encoding/json"
	"strconv"
	"strings"

	gcontract "github.com/x1n13y84issmd42/gog/graph/contract"
	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/contract"
	"github.com/x1n13y84issmd42/oasis/src/params"
	"github.com/x1n13y84issmd42/oasis/src/test"
)

// Status creates an expectation as for response's status code.
func Status(status int, log contract.Logger) contract.Expectation {
	log.Expecting("status", strconv.Itoa(status))

	return func(result *contract.OperationResult) bool {
		if result.HTTPResponse == nil {
			return false
		}

		if result.HTTPResponse.StatusCode == status {
			return true
		}

		log.ResponseHasWrongStatus(status, result.HTTPResponse.StatusCode)
		return false
	}
}

// HeaderRequired creates an expectation as for response's header value.
func HeaderRequired(n string, log contract.Logger) contract.Expectation {
	log.Expecting("required header", n)

	return func(result *contract.OperationResult) bool {
		if result.HTTPResponse.Header.Get(n) != "" {
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

	return func(result *contract.OperationResult) bool {
		return test.Schema(schema.Cast(result.HTTPResponse.Header.Get(n)), schema, log)
	}
}

// ContentType creates an expectation as for response's content type.
func ContentType(v string, log contract.Logger) contract.Expectation {
	log.Expecting("Content-Type", v)

	return func(result *contract.OperationResult) bool {
		if result.HTTPResponse == nil {
			return false
		}

		// This is to get rid of the possible "; charset=utf-8" part.
		respCT := strings.Split(result.HTTPResponse.Header.Get("Content-Type"), ";")[0]
		if respCT == v {
			return true
		}

		log.ResponseHasWrongContentType(v, result.HTTPResponse.Header.Get("Content-Type"))
		return false
	}
}

// ContentSchema creates an expectation as for response's body
// structure which must comply to the provided JSON schema.
func ContentSchema(schema *api.Schema, log contract.Logger) contract.Expectation {
	log.Expecting("content schema", schema.Name)

	return func(result *contract.OperationResult) bool {
		if result.HTTPResponse == nil {
			return false
		}

		respCT := strings.Split(result.HTTPResponse.Header.Get("Content-Type"), ";")[0]

		switch respCT {
		case "application/json":
			if test.JSONResponse(result, schema, log) {
				return true
			}

		default:
			log.NOMESSAGE("The Content-Type of '%s' is not supported.\n", respCT)
		}

		return false
	}
}

// JSONBody creates an expectation as for response's
// body properties values.
func JSONBody(props contract.Set, graph gcontract.Graph, log contract.Logger) contract.Expectation {

	return func(result *contract.OperationResult) bool {
		if result.HTTPResponse == nil {
			return false
		}

		respCT := strings.Split(result.HTTPResponse.Header.Get("Content-Type"), ";")[0]

		switch respCT {
		case "application/json":
			data := make(map[string]interface{})
			err := json.Unmarshal(result.ResponseBytes, &data)
			if err != nil {
				log.Error(err)
				return false
			}

			result := true

			for ebp := range props.Iterate() {
				expected := ebp.V()
				actual := params.Cast(data[ebp.N])

				if expected == actual {
					///
				} else {
					log.ResponseHasWrongPropertyValue(ebp.N, expected, actual)
					result = false
				}
			}

			return result

		default:
			log.NOMESSAGE("The Content-Type of '%s' is not supported.\n", respCT)
		}

		return false
	}
}
