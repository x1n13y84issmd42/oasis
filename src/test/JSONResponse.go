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

// JSONResponseString represents a JSON string response.
type JSONResponseString struct {
	Log         log.ILogger
	APIResponse *api.Response
}

// Test attempts to unmarshal response bodies as JSON strings then validate their contents.
func (test JSONResponseString) Test(response *http.Response) bool {
	// Testing the basic HTTP response props first.
	HTTPOK := HTTPResponse{
		Log:         test.Log,
		APIResponse: test.APIResponse,
	}.Test(response)

	if HTTPOK {
		var err error
		responseBody, _ := ioutil.ReadAll(response.Body)

		fmt.Println("JSONResponseString body: ", string(responseBody))
		fmt.Println("JSONResponseString status: ", response.StatusCode)

		var v string
		err = json.Unmarshal(responseBody, &v)

		if err == nil {
			return (SchemaString{test.APIResponse.Schema, test.Log}).Test(v, &utility.Context{Path: []string{"Response"}})
		}
		//TODO: log.FailedToParseJSONResponse(response, "string")

		test.Log.Error(err)
	}

	return false
}

// JSONResponseNumber represents a JSON number response.
type JSONResponseNumber struct {
	Log         log.ILogger
	APIResponse *api.Response
}

// Test attempts to unmarshal response bodies as JSON numbers then validate their contents.
func (test JSONResponseNumber) Test(response *http.Response) bool {
	// Testing the basic HTTP response props first.
	HTTPOK := HTTPResponse{
		Log:         test.Log,
		APIResponse: test.APIResponse,
	}.Test(response)

	if HTTPOK {
		var err error
		responseBody, _ := ioutil.ReadAll(response.Body)

		fmt.Println("JSONResponseNumber body: ", string(responseBody))
		fmt.Println("JSONResponseNumber status: ", response.StatusCode)

		var v int64
		err = json.Unmarshal(responseBody, &v)

		if err == nil {
			return (SchemaNumber{test.APIResponse.Schema, test.Log}).Test(v, &utility.Context{Path: []string{"Response"}})
		}
		//TODO: log.FailedToParseJSONResponse(response, "number")

		test.Log.Error(err)
	}

	return false
}

// JSONResponseBoolean represents a JSON boolean response.
type JSONResponseBoolean struct {
	Log         log.ILogger
	APIResponse *api.Response
}

// Test attempts to unmarshal response bodies as JSON booleans then validate their contents.
func (test JSONResponseBoolean) Test(response *http.Response) bool {
	// Testing the basic HTTP response props first.
	HTTPOK := HTTPResponse{
		Log:         test.Log,
		APIResponse: test.APIResponse,
	}.Test(response)

	if HTTPOK {
		var err error
		responseBody, _ := ioutil.ReadAll(response.Body)

		fmt.Println("JSONResponseBoolean body: ", string(responseBody))
		fmt.Println("JSONResponseBoolean status: ", response.StatusCode)

		var v bool
		err = json.Unmarshal(responseBody, &v)

		if err == nil {
			return (SchemaBoolean{test.APIResponse.Schema, test.Log}).Test(v, &utility.Context{Path: []string{"Response"}})
		}
		//TODO: log.FailedToParseJSONResponse(response, "boolean")

		test.Log.Error(err)
	}

	return false
}
