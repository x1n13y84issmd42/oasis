package test

import (
	"encoding/json"
	"fmt"
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
	// fmt.Println("JSONResponseArray is working")
	responseBody, _ := ioutil.ReadAll(response.Body)

	// fmt.Println("Response body: ", string(responseBody))
	// fmt.Println("Response status: ", response.StatusCode)

	var mJSON []interface{}
	errJSON := json.Unmarshal(responseBody, &mJSON)

	if errJSON != nil {
		fmt.Println("Something happened with JSON: ", errJSON)
		return false
	}

	ctx := &utility.Context{
		Path: []string{"Response"},
	}

	return Schema{tResp.APIResponse.Schema, tResp.Log}.Test(mJSON, ctx)
}
