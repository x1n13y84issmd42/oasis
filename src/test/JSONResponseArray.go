package test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/x1n13y84issmd42/goasis/src/api"
	"github.com/x1n13y84issmd42/goasis/src/log"
	"github.com/x1n13y84issmd42/goasis/src/utility"
)

// JSONResponseArray --
type JSONResponseArray struct {
	Log         log.ILogger
	APIResponse *api.Response
}

// Test tests.
func (test JSONResponseArray) Test(response *http.Response) bool {
	//TODO: HTTPReponse{}.Test(response)
	responseBody, _ := ioutil.ReadAll(response.Body)

	var mJSON []interface{}
	errJSON := json.Unmarshal(responseBody, &mJSON)

	if errJSON != nil {
		test.Log.ResponseExpectedArray(test.APIResponse)
		test.Log.Error(errJSON)
		return false
	}

	ctx := &utility.Context{
		Path: []string{"Response"},
	}

	return Schema{test.APIResponse.Schema, test.Log}.Test(mJSON, ctx)
}
