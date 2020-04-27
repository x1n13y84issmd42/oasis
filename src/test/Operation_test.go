package test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/x1n13y84issmd42/oasis/src/api"
)

//TODO: mock logger & test the Operation() function.

func TestSelectRequest(T *testing.T) {
	specMethod := "GET"
	specPath := "/foo/bar"

	specReq1 := &api.Request{
		Method: specMethod,
		Path:   specPath,
		Query: &url.Values{
			"qp1": []string{
				"one",
				"two",
				"three",
			},
		},
	}

	specReq2 := &api.Request{
		Method: specMethod,
		Path:   specPath,
		Query: &url.Values{
			"mode": []string{
				"yolo",
			},
		},
		Headers: http.Header{
			"Content-Type": []string{
				"application/xml",
			},
		},
	}

	specOp := api.Operation{
		Requests: []*api.Request{
			specReq1,
			specReq2,
		},
	}

	cases := []struct {
		ExReq  *api.Request
		Params *api.OperationParameters
	}{
		{
			ExReq:  specReq1,
			Params: &api.OperationParameters{},
		},
		{
			ExReq: specReq2,
			Params: &api.OperationParameters{
				Request: api.OperationRequestParameters{
					ContentTypeHint: "application/xml",
				},
			},
		},
		{
			ExReq: nil,
			Params: &api.OperationParameters{
				Request: api.OperationRequestParameters{
					ContentTypeHint: "text/csv",
				},
			},
		},
	}

	for _, c := range cases {
		actualReq := SelectRequest(&specOp, c.Params)
		assert.Equal(T, c.ExReq, actualReq)
	}
}

func TestSelectResponse(T *testing.T) {
	specResp1 := &api.Response{
		ContentType: "application/json",
		Description: "Successful JSON response",
		StatusCode:  200,
	}

	specResp2 := &api.Response{
		ContentType: "application/xml",
		Description: "Successful XML response",
		StatusCode:  200,
		Headers: api.Headers{
			"Content-Type": &api.Header{
				Name:     "Content-Type",
				Required: true,
			},
		},
	}

	specResp3 := &api.Response{
		ContentType: "application/json",
		Description: "Unsuccessful JSON response",
		StatusCode:  400,
		Headers: api.Headers{
			"Content-Type": &api.Header{
				Name:     "Content-Type",
				Required: true,
			},
		},
	}

	specOp := api.Operation{
		Responses: []*api.Response{
			specResp1,
			specResp2,
			specResp3,
		},
	}

	cases := []struct {
		ExResp *api.Response
		Params *api.OperationParameters
	}{
		{
			ExResp: specResp1,
			Params: &api.OperationParameters{},
		},
		{
			ExResp: specResp2,
			Params: &api.OperationParameters{
				Response: api.OperationResponseExpectation{
					ContentTypeHint: "application/xml",
				},
			},
		},
		{
			ExResp: specResp1,
			Params: &api.OperationParameters{
				Response: api.OperationResponseExpectation{
					ContentTypeHint: "application/json",
				},
			},
		},
		{
			ExResp: specResp3,
			Params: &api.OperationParameters{
				Response: api.OperationResponseExpectation{
					ContentTypeHint: "application/json",
					StatusHint:      400,
				},
			},
		},
		{
			ExResp: nil,
			Params: &api.OperationParameters{
				Response: api.OperationResponseExpectation{
					ContentTypeHint: "text/csv",
				},
			},
		},
	}

	for _, c := range cases {
		actualResp := SelectResponse(&specOp, c.Params)
		assert.Equal(T, c.ExResp, actualResp)
	}
}
