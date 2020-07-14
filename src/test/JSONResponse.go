package test

import (
	"encoding/json"

	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/contract"
)

// JSONResponse tests JSON response bodies.
func JSONResponse(result *contract.OperationResult, schema *api.Schema, log contract.Logger) bool {
	if len(result.ResponseBytes) == 0 {
		return false
	}

	var err error = nil

	if res, err := TryJSONObjectResponse(&result.ResponseBytes, log); err == nil {
		return Schema(res, schema, log)
	}

	if res, err := TryJSONArrayResponse(&result.ResponseBytes, log); err == nil {
		return Schema(res, schema, log)
	}

	if res, err := TryJSONStringResponse(&result.ResponseBytes, log); err == nil {
		return Schema(res, schema, log)
	}

	if res, err := TryJSONNumberResponse(&result.ResponseBytes, log); err == nil {
		return Schema(res, schema, log)
	}

	if res, err := TryJSONBooleanResponse(&result.ResponseBytes, log); err == nil {
		return Schema(res, schema, log)
	}

	if err != nil {
		log.Error(err)
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
func TryJSONStringResponse(respData *[]byte, log contract.Logger) (res string, err error) {
	err = json.Unmarshal(*respData, &res)
	return
}

// TryJSONNumberResponse tries to unmarshal respData as a number.
func TryJSONNumberResponse(respData *[]byte, log contract.Logger) (res int64, err error) {
	err = json.Unmarshal(*respData, &res)
	return
}

// TryJSONBooleanResponse tries to unmarshal respData as a boolean.
func TryJSONBooleanResponse(respData *[]byte, log contract.Logger) (res bool, err error) {
	err = json.Unmarshal(*respData, &res)
	return
}

// TryJSONObjectResponse tries to unmarshal respData as an object.
func TryJSONObjectResponse(respData *[]byte, log contract.Logger) (res *JSONMap, err error) {
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
func TryJSONArrayResponse(respData *[]byte, log contract.Logger) (res *JSONArray, err error) {
	var v JSONArray
	err = json.Unmarshal(*respData, &v)

	if err != nil {
		res = &v
		log.Error(err)
		return nil, err
	}

	res = &v

	return
}
