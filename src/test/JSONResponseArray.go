package test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/log"
)

// JSONResponseArray represents an array JSON response.
type JSONResponseArray struct {
	Log         log.ILogger
	APIResponse *api.Response
}

// Test attempts to unmarshal a response body as a JSON array
// and validate it's items.
func (test JSONResponseArray) Test(response *http.Response) bool {
	// Testing the basic HTTP response props first.
	HTTPOK := HTTPResponse{
		Log:         test.Log,
		APIResponse: test.APIResponse,
	}.Test(response)

	if HTTPOK {
		responseBody, _ := ioutil.ReadAll(response.Body)
		var mJSON []interface{}
		errJSON := json.Unmarshal(responseBody, &mJSON)

		if errJSON != nil {
			test.Log.ResponseExpectedArray(test.APIResponse)
			test.Log.Error(errJSON)
			return false
		}

		/* ctx := &utility.Context{
			Path: []string{"Response"},
		}

		return Schema{test.APIResponse.Schema, test.Log}.Test(mJSON, ctx) */
	}

	return false
}
