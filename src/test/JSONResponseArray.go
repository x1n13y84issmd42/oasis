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
func (tResp JSONResponseArray) Test(response *http.Response) bool {
	//TODO: HTTPReponse{}.Test(response)
	responseBody, _ := ioutil.ReadAll(response.Body)

	var mJSON []interface{}
	errJSON := json.Unmarshal(responseBody, &mJSON)

	if errJSON != nil {
		tResp.Log.ResponseExpectedArray(tResp.APIResponse)
		tResp.Log.Error(errJSON)
		return false
	}

	ctx := &utility.Context{
		Path: []string{"Response"},
	}

	return Schema{tResp.APIResponse.Schema, tResp.Log}.Test(mJSON, ctx)
}
