package test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/log"
)

// JSONResponse tests JSON response bodies.
func JSONResponse(resp *http.Response, specResp *api.Response, logger log.ILogger) bool {
	if resp.Body == nil {
		return false
	}

	respData, _ := ioutil.ReadAll(resp.Body)

	var err error = nil

	if res, err := TryJSONObjectResponse(&respData, specResp, logger); err == nil {
		return Schema(res, specResp.Schema, logger)
	}

	if res, err := TryJSONArrayResponse(&respData, specResp, logger); err == nil {
		return Schema(res, specResp.Schema, logger)
	}

	if res, err := TryJSONStringResponse(&respData, specResp, logger); err == nil {
		return Schema(res, specResp.Schema, logger)
	}

	if res, err := TryJSONNumberResponse(&respData, specResp, logger); err == nil {
		return Schema(res, specResp.Schema, logger)
	}

	if res, err := TryJSONBooleanResponse(&respData, specResp, logger); err == nil {
		return Schema(res, specResp.Schema, logger)
	}

	if err != nil {
		logger.Error(err)
	}

	return false
}

type (
	// JSONMap is a map to unmarshal JSONs into.
	JSONMap = map[string]interface{}

	// JSONArray is an array to unmarshal JSONs into.
	JSONArray = []interface{}
)

// TryJSONStringResponse tries to unmarshal respData as a string.
func TryJSONStringResponse(respData *[]byte, specResp *api.Response, logger log.ILogger) (res string, err error) {
	err = json.Unmarshal(*respData, &res)
	return
}

// TryJSONNumberResponse tries to unmarshal respData as a number.
func TryJSONNumberResponse(respData *[]byte, specResp *api.Response, logger log.ILogger) (res int64, err error) {
	err = json.Unmarshal(*respData, &res)
	return
}

// TryJSONBooleanResponse tries to unmarshal respData as a boolean.
func TryJSONBooleanResponse(respData *[]byte, specResp *api.Response, logger log.ILogger) (res bool, err error) {
	err = json.Unmarshal(*respData, &res)
	return
}

// TryJSONObjectResponse tries to unmarshal respData as an object.
func TryJSONObjectResponse(respData *[]byte, specResp *api.Response, logger log.ILogger) (res *JSONMap, err error) {
	v := make(JSONMap)
	err = json.Unmarshal(*respData, &v)

	if err != nil {
		res = &v
		return nil, err
	}

	res = &v

	return
}

// TryJSONArrayResponse tries to unmarshal respData as an array.
func TryJSONArrayResponse(respData *[]byte, specResp *api.Response, logger log.ILogger) (res *JSONArray, err error) {
	var v JSONArray
	err = json.Unmarshal(*respData, &v)

	if err != nil {
		res = &v
		return nil, err
	}

	res = &v

	return
}
