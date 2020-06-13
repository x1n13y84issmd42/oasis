package test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/log"
)

func TestHTTPResponse_False_Status(T *testing.T) {
	logger := log.NewFestive(0)

	httpResp := http.Response{
		StatusCode: 500,
	}

	specResp := api.Response{
		StatusCode:  200,
		ContentType: "application/json",
	}

	actual := HTTPResponse(&httpResp, &specResp, logger)
	assert.False(T, actual)
}

func TestHTTPResponse_False_CT(T *testing.T) {
	logger := log.NewFestive(0)

	httpResp := http.Response{
		StatusCode: 200,
		Header: http.Header{
			"Content-Type": []string{
				"application/xml",
			},
		},
	}

	specResp := api.Response{
		StatusCode:  200,
		ContentType: "application/json",
	}

	actual := HTTPResponse(&httpResp, &specResp, logger)
	assert.False(T, actual)
}

func TestHTTPResponse_False_Headers(T *testing.T) {
	logger := log.NewFestive(0)

	httpResp := http.Response{
		StatusCode: 200,
		Header: http.Header{
			"Content-Type": []string{
				"application/json",
			},
			"X-Index": []string{
				"3.1415",
			},
		},
	}

	specResp := api.Response{
		StatusCode:  200,
		ContentType: "application/json",
		Headers: api.Headers{
			"X-Index": &api.Header{
				Name:        "X-Index",
				Description: "whatever",
				Required:    true,
				Schema: &api.Schema{
					JSONSchema: api.JSONSchema{
						"type": "integer",
					},
				},
			},
		},
	}

	actual := HTTPResponse(&httpResp, &specResp, logger)
	assert.False(T, actual)
}

func TestHTTPResponse(T *testing.T) {
	logger := log.NewFestive(0)

	httpResp := http.Response{
		StatusCode: 200,
		Header: http.Header{
			"Content-Type": []string{
				"application/json",
			},
			"X-Index": []string{
				"3",
			},
		},
	}

	specResp := api.Response{
		StatusCode:  200,
		ContentType: "application/json",
		Headers: api.Headers{
			"X-Index": &api.Header{
				Name:        "X-Index",
				Description: "whatever",
				Required:    true,
				Schema: &api.Schema{
					JSONSchema: api.JSONSchema{
						"type": "integer",
					},
				},
			},
		},
	}

	actual := HTTPResponse(&httpResp, &specResp, logger)
	assert.True(T, actual)
}

func TestResponse_NoBody(T *testing.T) {
	logger := log.NewFestive(0)

	httpResp := http.Response{
		StatusCode: 200,
		Header: http.Header{
			"Content-Type": []string{
				"application/json",
			},
			"X-Index": []string{
				"3",
			},
		},
	}

	specResp := api.Response{
		StatusCode:  200,
		ContentType: "application/json",
		Schema: &api.Schema{
			JSONSchema: api.JSONSchema{
				"type": "integer",
			},
		},
		Headers: api.Headers{
			"X-Index": &api.Header{
				Name:        "X-Index",
				Description: "whatever",
				Required:    true,
				Schema: &api.Schema{
					JSONSchema: api.JSONSchema{
						"type": "integer",
					},
				},
			},
		},
	}

	actual := Response(&httpResp, &specResp, logger)
	assert.False(T, actual)
}

func TestResponse_SchemaError(T *testing.T) {
	logger := log.NewFestive(0)

	httpResp := http.Response{
		StatusCode: 200,
		Header: http.Header{
			"Content-Type": []string{
				"application/json",
			},
			"X-Index": []string{
				"3",
			},
		},
		Body: ioutil.NopCloser(bytes.NewReader([]byte("3.1415"))),
	}

	specResp := api.Response{
		StatusCode:  200,
		ContentType: "application/json",
		Schema: &api.Schema{
			JSONSchema: api.JSONSchema{
				"type": "integer",
			},
		},
		Headers: api.Headers{
			"X-Index": &api.Header{
				Name:        "X-Index",
				Description: "whatever",
				Required:    true,
				Schema: &api.Schema{
					JSONSchema: api.JSONSchema{
						"type": "integer",
					},
				},
			},
		},
	}

	actual := Response(&httpResp, &specResp, logger)
	assert.False(T, actual)
}

func TestResponse(T *testing.T) {
	logger := log.NewFestive(0)

	httpResp := http.Response{
		StatusCode: 200,
		Header: http.Header{
			"Content-Type": []string{
				"application/json",
			},
			"X-Index": []string{
				"3",
			},
		},
		Body: ioutil.NopCloser(bytes.NewReader([]byte("42"))),
	}

	specResp := api.Response{
		StatusCode:  200,
		ContentType: "application/json",
		Schema: &api.Schema{
			JSONSchema: api.JSONSchema{
				"type": "integer",
			},
		},
		Headers: api.Headers{
			"X-Index": &api.Header{
				Name:        "X-Index",
				Description: "whatever",
				Required:    true,
				Schema: &api.Schema{
					JSONSchema: api.JSONSchema{
						"type": "integer",
					},
				},
			},
		},
	}

	actual := Response(&httpResp, &specResp, logger)
	assert.True(T, actual)
}
