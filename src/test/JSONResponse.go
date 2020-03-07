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

// JSONResponseString --
type JSONResponseString struct {
	Log         log.ILogger
	APIResponse *api.Response
}

// Test tests.
func (test JSONResponseString) Test(response *http.Response) bool {
	var err error
	responseBody, _ := ioutil.ReadAll(response.Body)

	fmt.Println("JSONResponseString body: ", string(responseBody))
	fmt.Println("JSONResponseString status: ", response.StatusCode)

	var v string
	err = json.Unmarshal(responseBody, &v)

	if err == nil {
		return (SchemaString{test.APIResponse.Schema, test.Log}).Test(v, &utility.Context{Path: []string{"Response"}})
	}

	test.Log.Error(err)

	return false
}

// JSONResponseNumber --
type JSONResponseNumber struct {
	Log         log.ILogger
	APIResponse *api.Response
}

// Test tests.
func (test JSONResponseNumber) Test(response *http.Response) bool {
	var err error
	responseBody, _ := ioutil.ReadAll(response.Body)

	fmt.Println("JSONResponseNumber body: ", string(responseBody))
	fmt.Println("JSONResponseNumber status: ", response.StatusCode)

	var v int64
	err = json.Unmarshal(responseBody, &v)

	if err == nil {
		return (SchemaNumber{test.APIResponse.Schema, test.Log}).Test(v, &utility.Context{Path: []string{"Response"}})
	}

	test.Log.Error(err)

	return false
}

// JSONResponseBoolean --
type JSONResponseBoolean struct {
	Log         log.ILogger
	APIResponse *api.Response
}

// Test tests.
func (test JSONResponseBoolean) Test(response *http.Response) bool {
	var err error
	responseBody, _ := ioutil.ReadAll(response.Body)

	fmt.Println("JSONResponseBoolean body: ", string(responseBody))
	fmt.Println("JSONResponseBoolean status: ", response.StatusCode)

	var v bool
	err = json.Unmarshal(responseBody, &v)

	if err == nil {
		return (SchemaBoolean{test.APIResponse.Schema, test.Log}).Test(v, &utility.Context{Path: []string{"Response"}})
	}

	test.Log.Error(err)

	return false
}
