package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/log"
)

// JSONResponse test response bodies.
func JSONResponse(resp *http.Response, specResp *api.Response, logger log.ILogger) bool {
	if resp.Body == nil {
		return false
	}

	respData, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println("JSONResponseString body: ", string(respData))
	// fmt.Println("JSONResponseString status: ", resp.StatusCode)

	var err error = nil

	if res, err := tryJSONObjectResponse(&respData, specResp, logger); err == nil {
		return Schema(res, specResp.Schema, logger)
	}

	if res, err := tryJSONArrayResponse(&respData, specResp, logger); err == nil {
		return Schema(res, specResp.Schema, logger)
	}

	if res, err := tryJSONStringResponse(&respData, specResp, logger); err == nil {
		return Schema(res, specResp.Schema, logger)
	}

	if res, err := tryJSONNumberResponse(&respData, specResp, logger); err == nil {
		return Schema(res, specResp.Schema, logger)
	}

	if res, err := tryJSONBooleanResponse(&respData, specResp, logger); err == nil {
		return Schema(res, specResp.Schema, logger)
	}

	if err != nil {
		fmt.Printf("JSONResponse Error: %s", err.Error())
	}

	return false
}

type (
	// JSONMap is a map to unmarshal JSONs into.
	JSONMap = map[string]interface{}

	// JSONArray is an array to unmarshal JSONs into.
	JSONArray = []interface{}
)

// tryJSONStringResponse tryes to unmarshal respData as a string.
func tryJSONStringResponse(respData *[]byte, specResp *api.Response, logger log.ILogger) (res string, err error) {
	err = json.Unmarshal(*respData, &res)
	return
}

// tryJSONNumberResponse tryes to unmarshal respData as a string.
func tryJSONNumberResponse(respData *[]byte, specResp *api.Response, logger log.ILogger) (res int64, err error) {
	err = json.Unmarshal(*respData, &res)
	return
}

// tryJSONBooleanResponse tryes to unmarshal respData as a string.
func tryJSONBooleanResponse(respData *[]byte, specResp *api.Response, logger log.ILogger) (res bool, err error) {
	err = json.Unmarshal(*respData, &res)
	return
}

// tryJSONObjectResponse tryes to unmarshal respData as a string.
func tryJSONObjectResponse(respData *[]byte, specResp *api.Response, logger log.ILogger) (res *JSONMap, err error) {
	v := make(JSONMap)
	err = json.Unmarshal(*respData, &v)

	if err != nil {
		res = &v
		return nil, err
	}

	res = &v

	return
}

// tryJSONArrayResponse tryes to unmarshal respData as a string.
func tryJSONArrayResponse(respData *[]byte, specResp *api.Response, logger log.ILogger) (res *JSONArray, err error) {
	var v JSONArray
	err = json.Unmarshal(*respData, &v)

	if err != nil {
		res = &v
		return nil, err
	}

	res = &v

	return
}
