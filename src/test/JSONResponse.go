package test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/contract"
)

// JSONResponse tests JSON response bodies.
func JSONResponse(resp *http.Response, schema *api.Schema, logger contract.Logger) bool {
	if resp.Body == nil {
		return false
	}

	respData, _ := ioutil.ReadAll(resp.Body)

	var err error = nil

	if res, err := TryJSONObjectResponse(&respData, logger); err == nil {
		return Schema(res, schema, logger)
	}

	if res, err := TryJSONArrayResponse(&respData, logger); err == nil {
		return Schema(res, schema, logger)
	}

	if res, err := TryJSONStringResponse(&respData, logger); err == nil {
		return Schema(res, schema, logger)
	}

	if res, err := TryJSONNumberResponse(&respData, logger); err == nil {
		return Schema(res, schema, logger)
	}

	if res, err := TryJSONBooleanResponse(&respData, logger); err == nil {
		return Schema(res, schema, logger)
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
func TryJSONStringResponse(respData *[]byte, logger contract.Logger) (res string, err error) {
	err = json.Unmarshal(*respData, &res)
	return
}

// TryJSONNumberResponse tries to unmarshal respData as a number.
func TryJSONNumberResponse(respData *[]byte, logger contract.Logger) (res int64, err error) {
	err = json.Unmarshal(*respData, &res)
	return
}

// TryJSONBooleanResponse tries to unmarshal respData as a boolean.
func TryJSONBooleanResponse(respData *[]byte, logger contract.Logger) (res bool, err error) {
	err = json.Unmarshal(*respData, &res)
	return
}

// TryJSONObjectResponse tries to unmarshal respData as an object.
func TryJSONObjectResponse(respData *[]byte, logger contract.Logger) (res *JSONMap, err error) {
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
func TryJSONArrayResponse(respData *[]byte, logger contract.Logger) (res *JSONArray, err error) {
	var v JSONArray
	err = json.Unmarshal(*respData, &v)

	if err != nil {
		res = &v
		return nil, err
	}

	res = &v

	return
}
