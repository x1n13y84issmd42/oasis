package openapi3

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
	"github.com/x1n13y84issmd42/oasis/src/api"
	APIKey "github.com/x1n13y84issmd42/oasis/src/api/security/APIKey"
	HTTP "github.com/x1n13y84issmd42/oasis/src/api/security/HTTP"
	"github.com/x1n13y84issmd42/oasis/src/errors"
	"github.com/x1n13y84issmd42/oasis/src/log"
)

func TestCreatePath(T *testing.T) {
	path := "/foo/{param1}/bar/{param2}/{param3}"
	expectedPath := "/foo/p1_from_op/bar/p2_from_path/p3_from_override"

	OAS := &openapi3.Swagger{}
	oasOperation := openapi3.Operation{
		Parameters: openapi3.Parameters{
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					In:       "path",
					Name:     "param1",
					Required: true,
					Example:  "p1_from_op",
				},
			},

			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					In:      "query",
					Name:    "param2",
					Example: "p2_from_op",
				},
			},

			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					In:       "path",
					Name:     "param3",
					Required: true,
					Example:  "p3_from_op",
				},
			},
		},
	}
	oasPathItem := &openapi3.PathItem{
		Get: &oasOperation,
		Parameters: openapi3.Parameters{
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					In:       "path",
					Name:     "param1",
					Required: true,
					Example:  "p1_from_path",
				},
			},

			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					In:       "path",
					Name:     "param2",
					Required: true,
					Example:  "p2_from_path",
				},
			},
		},
	}
	OAS.Paths = openapi3.Paths{}
	OAS.Paths[path] = oasPathItem

	params := api.OperationParameters{
		Path: api.PathParameters{
			"param1": "",
			"param3": "p3_from_override",
		},
	}

	spec := Spec{
		Log: log.NewFestive(0),
		OAS: OAS,
	}

	actualPath, err := spec.CreatePath(path, oasPathItem, &oasOperation, &params)

	assert.Equal(T, expectedPath, actualPath)
	assert.Nil(T, err)
}

func TestCreatePath_ErrNoParameters(T *testing.T) {
	path := "/foo/{param1}/bar/{param2}/{param3}/qeq/{param4}"
	expectedPath := "/foo/p1_from_op/bar/{param2}/{param3}/qeq/{param4}"
	expectedErr := errors.NoParameters([]string{"param2", "param3", "param4"}, nil)

	OAS := &openapi3.Swagger{}
	oasOperation := openapi3.Operation{
		Parameters: openapi3.Parameters{
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					In:       "path",
					Name:     "param1",
					Required: true,
					Example:  "p1_from_op",
				},
			},

			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					In:      "query",
					Name:    "param2",
					Example: "p2_from_op",
				},
			},
		},
	}
	oasPathItem := &openapi3.PathItem{
		Get: &oasOperation,
	}
	OAS.Paths = openapi3.Paths{}
	OAS.Paths[path] = oasPathItem

	params := api.OperationParameters{}

	spec := Spec{
		Log: log.NewFestive(0),
		OAS: OAS,
	}

	actualPath, err := spec.CreatePath(path, oasPathItem, &oasOperation, &params)

	assert.Equal(T, expectedPath, actualPath)
	assert.Equal(T, err.Error(), expectedErr.Error())
}

func TestCreateQuery(T *testing.T) {
	path := "/foo/{param1}/bar/{param2}/{param3}"
	expectedQuery := &url.Values{
		"q1": []string{
			"q1_from_override",
			"q1_from_op",
			"q1_from_path",
		},

		"q2": []string{
			"q2_from_op",
		},

		"q3": []string{
			"q3_from_override",
			"q3_from_path",
		},
	}

	OAS := &openapi3.Swagger{}
	oasOperation := openapi3.Operation{
		Parameters: openapi3.Parameters{
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					In:       "query",
					Required: true,
					Name:     "q1",
					Example:  "q1_from_op",
				},
			},

			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					In:       "query",
					Required: true,
					Name:     "q2",
					Example:  "q2_from_op",
				},
			},

			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					In:      "query",
					Name:    "q3",
					Example: "q3_from_op",
				},
			},
		},
	}
	oasPathItem := &openapi3.PathItem{
		Get: &oasOperation,
		Parameters: openapi3.Parameters{
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					In:       "query",
					Required: true,
					Name:     "q1",
					Example:  "q1_from_path",
				},
			},

			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					In:       "query",
					Required: true,
					Name:     "q3",
					Example:  "q3_from_path",
				},
			},
		},
	}
	OAS.Paths = openapi3.Paths{}
	OAS.Paths[path] = oasPathItem

	params := api.OperationParameters{
		Query: api.QueryValues{
			"q1": []string{
				"q1_from_override",
			},
			"q2": nil,
			"q3": []string{
				"q3_from_override",
			},
		},
	}

	spec := Spec{
		Log: log.NewFestive(0),
		OAS: OAS,
	}

	actualQuery, err := spec.CreateQuery(oasPathItem, &oasOperation, &params)

	assert.Equal(T, expectedQuery, actualQuery)
	assert.Nil(T, err)
}

func TestCreateQuery_ErrNoParameters(T *testing.T) {
	path := "/foo/{param1}/bar/{param2}/{param3}"

	expectedQuery := &url.Values{
		"q1": []string{
			"q1_from_op",
		},
	}

	expectedErr := errors.NoParameters([]string{"q2", "q3", "q4"}, nil)

	OAS := &openapi3.Swagger{}
	oasOperation := openapi3.Operation{
		Parameters: openapi3.Parameters{
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					In:       "query",
					Required: true,
					Name:     "q1",
					Example:  "q1_from_op",
				},
			},

			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					In:       "query",
					Required: true,
					Name:     "q2",
				},
			},

			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					In:      "query",
					Name:    "q3",
					Example: "q3_from_op",
				},
			},
		},
	}
	oasPathItem := &openapi3.PathItem{
		Get: &oasOperation,
		Parameters: openapi3.Parameters{
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					In:       "query",
					Required: true,
					Name:     "q1",
				},
			},

			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					In:       "query",
					Required: true,
					Name:     "q3",
				},
			},

			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					In:       "query",
					Required: true,
					Name:     "q4",
				},
			},
		},
	}
	OAS.Paths = openapi3.Paths{}
	OAS.Paths[path] = oasPathItem

	params := api.OperationParameters{}

	spec := Spec{
		Log: log.NewFestive(0),
		OAS: OAS,
	}

	actualQuery, actualErr := spec.CreateQuery(oasPathItem, &oasOperation, &params)

	assert.Equal(T, expectedQuery, actualQuery)
	assert.Equal(T, expectedErr.Error(), actualErr.Error())
}

func TestGetHost(T *testing.T) {
	expectedHost := &api.Host{
		URL:  "stage.localhost",
		Name: "Staging localhost",
	}

	OAS := &openapi3.Swagger{}
	OAS.Servers = append(OAS.Servers, &openapi3.Server{
		URL:         "localhost",
		Description: "Localhost",
	}, &openapi3.Server{
		URL:         "stage.localhost",
		Description: "Staging localhost",
	}, &openapi3.Server{
		URL:         "dev.localhost",
		Description: "Development localhost",
	})

	spec := Spec{
		Log: log.NewFestive(0),
		OAS: OAS,
	}

	actualHost, actualErr := spec.GetHost("Staging localhost")

	assert.Equal(T, expectedHost, actualHost)
	assert.Nil(T, actualErr)
}

func TestGetHost_HostNotFound(T *testing.T) {
	expectedErr := errors.ErrNotFound{
		What: "Host",
		Name: "Production localhost",
	}

	OAS := &openapi3.Swagger{}
	OAS.Servers = append(OAS.Servers, &openapi3.Server{
		URL:         "localhost",
		Description: "Localhost",
	}, &openapi3.Server{
		URL:         "stage.localhost",
		Description: "Staging localhost",
	}, &openapi3.Server{
		URL:         "dev.localhost",
		Description: "Development localhost",
	})

	spec := Spec{
		Log: log.NewFestive(0),
		OAS: OAS,
	}

	actualHost, actualErr := spec.GetHost("Production localhost")

	assert.Nil(T, actualHost)
	assert.Equal(T, expectedErr.Error(), actualErr.Error())
}

func TestGetDefaultHost(T *testing.T) {
	expectedHost := &api.Host{
		URL:  "localhost",
		Name: "Localhost",
	}

	OAS := &openapi3.Swagger{}
	OAS.Servers = append(OAS.Servers, &openapi3.Server{
		URL:         "localhost",
		Description: "Localhost",
	}, &openapi3.Server{
		URL:         "stage.localhost",
		Description: "Staging localhost",
	}, &openapi3.Server{
		URL:         "dev.localhost",
		Description: "Development localhost",
	})

	spec := Spec{
		Log: log.NewFestive(0),
		OAS: OAS,
	}

	actualHost, actualErr := spec.GetDefaultHost()

	assert.Equal(T, expectedHost, actualHost)
	assert.Nil(T, actualErr)
}

func TestGetDefaultHost_NotFound(T *testing.T) {
	expectedErr := errors.ErrNotFound{
		What: "Host",
		Name: "Default",
	}

	OAS := &openapi3.Swagger{}
	spec := Spec{
		Log: log.NewFestive(0),
		OAS: OAS,
	}

	actualHost, actualErr := spec.GetDefaultHost()

	assert.Nil(T, actualHost)
	assert.Equal(T, expectedErr.Error(), actualErr.Error())
}

func TestMakeJSONSchema_ErrSchemaMarshal(T *testing.T) {
	schemaName := "schema_one"
	expectedErr := errors.InvalidSchema(schemaName, "Failed to marshal the schema.", nil)

	schema := openapi3.Schema{
		Type:    "object",
		Example: make(chan int),
	}

	OAS := &openapi3.Swagger{
		Components: openapi3.Components{},
	}

	spec := Spec{
		Log: log.NewFestive(0),
		OAS: OAS,
	}

	actual, actualErr := spec.MakeJSONSchema(schemaName, &schema)

	assert.Nil(T, actual)
	assert.Equal(T, expectedErr.Error(), actualErr.Error())
}

func TestMakeJSONSchema_ErrComponentsMarshal(T *testing.T) {
	schemaName := "schema_one"
	expectedErr := errors.InvalidSchema(schemaName, "Failed to marshal Components.", nil)

	schema := openapi3.Schema{
		Type: "object",
	}

	OAS := &openapi3.Swagger{
		Components: openapi3.Components{
			Schemas: map[string]*openapi3.SchemaRef{
				"Errorneous component": {
					Value: &openapi3.Schema{
						Example: make(chan int),
					},
				},
			},
		},
	}

	spec := Spec{
		Log: log.NewFestive(0),
		OAS: OAS,
	}

	actual, actualErr := spec.MakeJSONSchema(schemaName, &schema)

	assert.Nil(T, actual)
	assert.Equal(T, expectedErr.Error(), actualErr.Error())
}

func TestMakeJSONSchema(T *testing.T) {
	schemaName := "schema_one"
	// expectedSchema := api.JSONSchema{}
	expectedSchema := api.JSONSchema{
		"components": map[string]interface{}{
			"schemas": map[string]interface{}{
				"User": map[string]interface{}{
					"example": "yolo",
				},
			},
		},
		"properties": map[string]interface{}{
			"x": map[string]interface{}{
				"type": "integer",
			},
			"y": map[string]interface{}{
				"$ref": "#/components/schemas/User",
			},
		},
		"type": "object",
	}

	schema := openapi3.Schema{
		Type: "object",
		Properties: map[string]*openapi3.SchemaRef{
			"x": {
				Value: &openapi3.Schema{
					Type: "integer",
				},
			},
			"y": {
				Ref: "#/components/schemas/User",
			},
		},
	}

	OAS := &openapi3.Swagger{
		Components: openapi3.Components{
			Schemas: map[string]*openapi3.SchemaRef{
				"User": {
					Value: &openapi3.Schema{
						Example: "yolo",
					},
				},
			},
		},
	}

	spec := Spec{
		Log: log.NewFestive(0),
		OAS: OAS,
	}

	actual, actualErr := spec.MakeJSONSchema(schemaName, &schema)

	assert.Equal(T, expectedSchema, actual)
	assert.Nil(T, actualErr)
}

func TestMakeSchema(T *testing.T) {
	schemaName := "schema_one"
	// expectedSchema := api.JSONSchema{}
	expectedSchema := &api.Schema{
		Name: schemaName,
		JSONSchema: api.JSONSchema{
			"components": map[string]interface{}{
				"schemas": map[string]interface{}{
					"User": map[string]interface{}{
						"example": "yolo",
					},
				},
			},
			"properties": map[string]interface{}{
				"x": map[string]interface{}{
					"type": "integer",
				},
				"y": map[string]interface{}{
					"$ref": "#/components/schemas/User",
				},
			},
			"type": "object",
		},
	}

	schema := openapi3.Schema{
		Type: "object",
		Properties: map[string]*openapi3.SchemaRef{
			"x": {
				Value: &openapi3.Schema{
					Type: "integer",
				},
			},
			"y": {
				Ref: "#/components/schemas/User",
			},
		},
	}

	OAS := &openapi3.Swagger{
		Components: openapi3.Components{
			Schemas: map[string]*openapi3.SchemaRef{
				"User": {
					Value: &openapi3.Schema{
						Example: "yolo",
					},
				},
			},
		},
	}

	spec := Spec{
		Log: log.NewFestive(0),
		OAS: OAS,
	}

	actual, actualErr := spec.MakeSchema(schemaName, &schema)

	assert.Equal(T, expectedSchema, actual)
	assert.Nil(T, actualErr)
}

func TestMakeSchema_ErrPassThrough(T *testing.T) {
	schemaName := "schema_one"
	expectedErr := errors.InvalidSchema(schemaName, "Failed to marshal the schema.", nil)

	schema := openapi3.Schema{
		Type:    "object",
		Example: make(chan int),
	}

	OAS := &openapi3.Swagger{
		Components: openapi3.Components{},
	}

	spec := Spec{
		Log: log.NewFestive(0),
		OAS: OAS,
	}

	actual, actualErr := spec.MakeSchema(schemaName, &schema)

	assert.Nil(T, actual)
	assert.Equal(T, expectedErr.Error(), actualErr.Error())
}

func TestMakeSchema_ErrNotFound(T *testing.T) {
	schemaName := "schema_one"
	expectedErr := errors.NotFound("Schema", schemaName, nil)

	OAS := &openapi3.Swagger{
		Components: openapi3.Components{},
	}

	spec := Spec{
		Log: log.NewFestive(0),
		OAS: OAS,
	}

	actual, actualErr := spec.MakeSchema(schemaName, nil)

	assert.Nil(T, actual)
	assert.Equal(T, expectedErr.Error(), actualErr.Error())
}

func TestMakeHeader_Err(T *testing.T) {
	headerName := "x-tested"
	expectedErr := errors.InvalidSchema(headerName, "Failed to marshal the schema.", nil)

	oasHeader := &openapi3.Header{
		Description: "i am a header",
		Required:    true,
		Schema: &openapi3.SchemaRef{
			Value: &openapi3.Schema{
				Type:    "object",
				Example: make(chan int),
			},
		},
	}

	OAS := &openapi3.Swagger{
		Components: openapi3.Components{},
	}

	spec := Spec{
		Log: log.NewFestive(0),
		OAS: OAS,
	}

	actual, actualErr := spec.MakeHeader(headerName, oasHeader)

	assert.Nil(T, actual)
	assert.Equal(T, expectedErr.Error(), actualErr.Error())
}

func TestMakeHeader(T *testing.T) {
	headerName := "x-tested"
	expected := &api.Header{
		Name:        headerName,
		Description: "i am a header",
		Required:    true,
		Schema: &api.Schema{
			Name: headerName,
			JSONSchema: api.JSONSchema{
				"type":       "integer",
				"example":    float64(42),
				"components": map[string]interface{}{},
			},
		},
	}

	oasHeader := &openapi3.Header{
		Description: "i am a header",
		Required:    true,
		Schema: &openapi3.SchemaRef{
			Value: &openapi3.Schema{
				Type:    "integer",
				Example: 42,
			},
		},
	}

	OAS := &openapi3.Swagger{
		Components: openapi3.Components{},
	}

	spec := Spec{
		Log: log.NewFestive(0),
		OAS: OAS,
	}

	actual, actualErr := spec.MakeHeader("x-tested", oasHeader)

	assert.Equal(T, expected, actual)
	assert.Nil(T, actualErr)
}

func TestMakeResponses_NoBodies(T *testing.T) {
	expected := []*api.Response{
		{
			Description: "A successful test response",
			StatusCode:  200,
			Headers:     api.Headers{},
		},
	}

	OAS := &openapi3.Swagger{
		Components: openapi3.Components{},
	}

	oasResp := &openapi3.Response{
		Description: "A successful test response",
	}

	spec := Spec{
		Log: log.NewFestive(0),
		OAS: OAS,
	}

	actual, actualErr := spec.MakeResponses("200", oasResp)
	assert.Equal(T, expected, actual)
	assert.Nil(T, actualErr)
}

func TestMakeResponses_Bodies(T *testing.T) {
	expected := []*api.Response{
		{
			Description: "test responses",
			ContentType: "application/json",
			StatusCode:  200,
			Schema: &api.Schema{
				Name: "Response",
				JSONSchema: api.JSONSchema{
					"components": map[string]interface{}{},
					"type":       "integer",
				},
			},
			Headers: api.Headers{},
		},
		{
			Description: "test responses",
			ContentType: "application/xml",
			StatusCode:  200,
			Schema: &api.Schema{
				Name: "Response",
				JSONSchema: api.JSONSchema{
					"components": map[string]interface{}{},
					"type":       "integer",
				},
			},
			Headers: api.Headers{},
		},
	}

	OAS := &openapi3.Swagger{
		Components: openapi3.Components{},
	}

	oasResp := &openapi3.Response{
		Description: "test responses",
		Content: openapi3.Content{
			"application/json": &openapi3.MediaType{
				Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: "integer",
					},
				},
			},

			"application/xml": &openapi3.MediaType{
				Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: "integer",
					},
				},
			},
		},
	}

	spec := Spec{
		Log: log.NewFestive(0),
		OAS: OAS,
	}

	actual, actualErr := spec.MakeResponses("200", oasResp)

	assert.Equal(T, expected, actual)
	assert.Nil(T, actualErr)
}

func TestMakeResponses_Err_Headers(T *testing.T) {
	headerName := "x-bugged"
	expectedErr := errors.InvalidResponse("Failed to create a response header '"+headerName+"' schema.", nil)

	OAS := &openapi3.Swagger{
		Components: openapi3.Components{},
	}

	oasResp := &openapi3.Response{
		Headers: map[string]*openapi3.HeaderRef{
			headerName: {
				Value: &openapi3.Header{
					Description: "i am a header",
					Required:    true,
					Schema: &openapi3.SchemaRef{
						Value: &openapi3.Schema{
							Type:    "object",
							Example: make(chan int),
						},
					},
				},
			},
		},
	}

	spec := Spec{
		Log: log.NewFestive(0),
		OAS: OAS,
	}

	actual, actualErr := spec.MakeResponses("200", oasResp)
	assert.Nil(T, actual)
	assert.Equal(T, expectedErr.Error(), actualErr.Error())
}

func TestMakeResponses_Err_Content(T *testing.T) {
	respCT := "application/json"
	expectedErr := errors.InvalidResponse("Failed to create a '"+respCT+"' response body schema.", nil)

	OAS := &openapi3.Swagger{
		Components: openapi3.Components{},
	}

	oasResp := &openapi3.Response{
		Content: openapi3.Content{
			respCT: &openapi3.MediaType{
				Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type:    "object",
						Example: make(chan int),
					},
				},
			},
		},
	}

	spec := Spec{
		Log: log.NewFestive(0),
		OAS: OAS,
	}

	actual, actualErr := spec.MakeResponses("200", oasResp)
	assert.Nil(T, actual)
	assert.Equal(T, expectedErr.Error(), actualErr.Error())
}

func TestMakeSecurity_Err_NotFound_unnamed(T *testing.T) {
	expectedErr := errors.SecurityNotFound("[unnamed]", "No security name has been supplied.", nil)

	OAS := &openapi3.Swagger{
		Components: openapi3.Components{},
	}

	oasSecReqs := &openapi3.SecurityRequirements{}

	params := &api.OperationParameters{}

	spec := Spec{
		Log: log.NewFestive(0),
		OAS: OAS,
	}

	actual, actualErr := spec.MakeSecurity(oasSecReqs, params)
	assert.Nil(T, actual)
	assert.Equal(T, expectedErr.Error(), actualErr.Error())
}

func TestMakeSecurity_Err_NotFound(T *testing.T) {
	secName := "nonexistent_sec"
	expectedErr := errors.SecurityNotFound(secName, "", nil)

	OAS := &openapi3.Swagger{
		Components: openapi3.Components{
			SecuritySchemes: map[string]*openapi3.SecuritySchemeRef{
				"sec_1": {
					Value: &openapi3.SecurityScheme{},
				},
			},
		},
	}

	oasSecReqs := &openapi3.SecurityRequirements{}

	params := &api.OperationParameters{
		Security: api.OperationSecurityParameters{
			SecurityHint: secName,
		},
	}

	spec := Spec{
		Log: log.NewFestive(0),
		OAS: OAS,
	}

	actual, actualErr := spec.MakeSecurity(oasSecReqs, params)
	assert.Nil(T, actual)
	assert.Equal(T, expectedErr.Error(), actualErr.Error())
}

func TestMakeSecurity_Err_NotFound_unknown(T *testing.T) {
	secName := "sec_1"
	secType := "weirdsec"
	expectedErr := errors.SecurityNotFound(secName, "Security type '"+secType+"' is unknown.", nil)

	OAS := &openapi3.Swagger{
		Components: openapi3.Components{
			SecuritySchemes: map[string]*openapi3.SecuritySchemeRef{
				secName: {
					Value: &openapi3.SecurityScheme{
						Type: secType,
					},
				},
			},
		},
	}

	oasSecReqs := &openapi3.SecurityRequirements{}

	params := &api.OperationParameters{
		Security: api.OperationSecurityParameters{
			SecurityHint: secName,
		},
	}

	spec := Spec{
		Log: log.NewFestive(0),
		OAS: OAS,
	}

	actual, actualErr := spec.MakeSecurity(oasSecReqs, params)
	assert.Nil(T, actual)
	assert.Equal(T, expectedErr.Error(), actualErr.Error())
}

func TestMakeSecurity_Err_JSON(T *testing.T) {
	secName := "sec_1"
	expectedErr := errors.Oops("The 'x-example' does not contain any JSON data.", nil)

	OAS := &openapi3.Swagger{
		Components: openapi3.Components{
			SecuritySchemes: map[string]*openapi3.SecuritySchemeRef{
				secName: {
					Value: &openapi3.SecurityScheme{
						Type: "http",
						ExtensionProps: openapi3.ExtensionProps{
							Extensions: map[string]interface{}{
								"x-example": make(chan int),
							},
						},
					},
				},
			},
		},
	}

	oasSecReqs := &openapi3.SecurityRequirements{}

	params := &api.OperationParameters{
		Security: api.OperationSecurityParameters{
			SecurityHint: secName,
		},
	}

	spec := Spec{
		Log: log.NewFestive(0),
		OAS: OAS,
	}

	actual, actualErr := spec.MakeSecurity(oasSecReqs, params)
	assert.Nil(T, actual)
	assert.Equal(T, expectedErr.Error(), actualErr.Error())
}

func TestMakeSecurity_Err_Unmarshal(T *testing.T) {
	secName := "sec_1"
	expectedErr := errors.Oops("Cannot unmarshal the 'x-example' field.", nil)

	OAS := &openapi3.Swagger{
		Components: openapi3.Components{
			SecuritySchemes: map[string]*openapi3.SecuritySchemeRef{
				secName: {
					Value: &openapi3.SecurityScheme{
						Type: "http",
						ExtensionProps: openapi3.ExtensionProps{
							Extensions: map[string]interface{}{
								"x-example": json.RawMessage("\\\\ a totally invalid JSON //"),
							},
						},
					},
				},
			},
		},
	}

	oasSecReqs := &openapi3.SecurityRequirements{}

	params := &api.OperationParameters{
		Security: api.OperationSecurityParameters{
			SecurityHint: secName,
		},
	}

	spec := Spec{
		Log: log.NewFestive(0),
		OAS: OAS,
	}

	actual, actualErr := spec.MakeSecurity(oasSecReqs, params)
	assert.Nil(T, actual)
	assert.Equal(T, expectedErr.Error(), actualErr.Error())
}

func TestMakeSecurity_Params(T *testing.T) {
	logger := log.NewFestive(0)
	secName := "sec_1"
	expected := HTTP.Basic{
		Security: HTTP.Security{
			Name:  secName,
			Token: "42",
			Log:   logger,
		},
	}

	OAS := &openapi3.Swagger{
		Components: openapi3.Components{
			SecuritySchemes: map[string]*openapi3.SecuritySchemeRef{
				"sec_99": {
					Value: &openapi3.SecurityScheme{
						Type: "apikey",
						ExtensionProps: openapi3.ExtensionProps{
							Extensions: map[string]interface{}{
								"x-example": json.RawMessage("\"9966\""),
							},
						},
					},
				},
				secName: {
					Value: &openapi3.SecurityScheme{
						Type:   "http",
						Scheme: "basic",
						ExtensionProps: openapi3.ExtensionProps{
							Extensions: map[string]interface{}{
								"x-example": json.RawMessage("\"42\""),
							},
						},
					},
				},
			},
		},
	}

	oasSecReqs := &openapi3.SecurityRequirements{
		{
			"sec_99": nil,
		},
	}

	params := &api.OperationParameters{
		Security: api.OperationSecurityParameters{
			SecurityHint: secName,
		},
	}

	spec := Spec{
		Log: logger,
		OAS: OAS,
	}

	actual, actualErr := spec.MakeSecurity(oasSecReqs, params)
	assert.EqualValues(T, expected, actual)
	assert.Nil(T, actualErr)
}

func TestMakeSecurity_Requirements(T *testing.T) {
	logger := log.NewFestive(0)
	secName := "sec_99"
	paramName := "creds"
	expected := APIKey.Query{
		Security: APIKey.Security{
			Name:      secName,
			Log:       logger,
			ParamName: paramName,
			Value:     "9966",
		},
	}

	OAS := &openapi3.Swagger{
		Components: openapi3.Components{
			SecuritySchemes: map[string]*openapi3.SecuritySchemeRef{
				secName: {
					Value: &openapi3.SecurityScheme{
						Type: "apiKey",
						In:   "query",
						Name: paramName,
						ExtensionProps: openapi3.ExtensionProps{
							Extensions: map[string]interface{}{
								"x-example": json.RawMessage("\"9966\""),
							},
						},
					},
				},
				"sec_1": {
					Value: &openapi3.SecurityScheme{
						Type:   "http",
						Scheme: "basic",
						ExtensionProps: openapi3.ExtensionProps{
							Extensions: map[string]interface{}{
								"x-example": json.RawMessage("\"42\""),
							},
						},
					},
				},
			},
		},
	}

	oasSecReqs := &openapi3.SecurityRequirements{
		{
			secName: nil,
		},
	}

	params := &api.OperationParameters{}

	spec := Spec{
		Log: logger,
		OAS: OAS,
	}

	actual, actualErr := spec.MakeSecurity(oasSecReqs, params)
	assert.EqualValues(T, expected, actual)
	assert.Nil(T, actualErr)
}

func TestMakeSecurity_Params_Value(T *testing.T) {
	logger := log.NewFestive(0)
	secName := "sec_99"
	paramName := "creds"
	expected := APIKey.Query{
		Security: APIKey.Security{
			Name:      secName,
			Log:       logger,
			ParamName: paramName,
			Value:     "YOLO",
		},
	}

	OAS := &openapi3.Swagger{
		Components: openapi3.Components{
			SecuritySchemes: map[string]*openapi3.SecuritySchemeRef{
				secName: {
					Value: &openapi3.SecurityScheme{
						Type: "apiKey",
						In:   "query",
						Name: paramName,
						ExtensionProps: openapi3.ExtensionProps{
							Extensions: map[string]interface{}{
								"x-example": json.RawMessage("\"9966\""),
							},
						},
					},
				},
				"sec_1": {
					Value: &openapi3.SecurityScheme{
						Type:   "http",
						Scheme: "basic",
						ExtensionProps: openapi3.ExtensionProps{
							Extensions: map[string]interface{}{
								"x-example": json.RawMessage("\"42\""),
							},
						},
					},
				},
			},
		},
	}

	oasSecReqs := &openapi3.SecurityRequirements{
		{
			secName: nil,
		},
	}

	params := &api.OperationParameters{
		Security: api.OperationSecurityParameters{
			HTTPAuthValue: "YOLO",
		},
	}

	spec := Spec{
		Log: logger,
		OAS: OAS,
	}

	actual, actualErr := spec.MakeSecurity(oasSecReqs, params)
	assert.EqualValues(T, expected, actual)
	assert.Nil(T, actualErr)
}

func TestCreateRequest(T *testing.T) {
	logger := log.NewFestive(0)

	method := "GET"
	path := "/foo/bar"
	CT := "application/json"
	query := url.Values{
		"qp1": []string{
			"v1",
			"v_two",
		},
	}

	params := api.OperationParameters{
		Request: api.OperationRequestParameters{
			Headers: api.HTTPHeaders{
				"x-ghosts": []string{
					"casper",
				},
			},
		},
	}

	oasOp := openapi3.Operation{
		Parameters: openapi3.Parameters{
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					In:       "header",
					Required: true,
					Name:     "x-ghosts",
					Example:  "stay puft",
				},
			},

			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					In:      "header",
					Name:    "x-ghosts",
					Example: "slimer",
				},
			},
		},
	}

	oasPathItem := openapi3.PathItem{
		Parameters: openapi3.Parameters{
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					In:       "header",
					Required: true,
					Name:     "x-ghosts",
					Example:  "ghost rider",
				},
			},
		},
	}

	oasMT := openapi3.MediaType{}

	expected := &api.Request{
		Method: method,
		Path:   path,
		Query:  &query,
		Headers: http.Header{
			"X-Ghosts": []string{
				"casper",
				"stay puft",
				"ghost rider",
			},
		},
		//TODO: body
	}

	spec := Spec{
		Log: logger,
		OAS: &openapi3.Swagger{
			Components: openapi3.Components{},
		},
	}

	actual := spec.MakeRequest(method, path, &query, &oasOp, &oasPathItem, CT, &oasMT, &params)
	assert.EqualValues(T, expected, actual)
}
