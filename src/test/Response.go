package test

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/log"
)

// HTTPResponse test checks basic response properties such as status code and headers.
func HTTPResponse(resp *http.Response, specResp *api.Response, logger log.ILogger) bool {
	statusOK := specResp.StatusCode == uint64(resp.StatusCode)

	//FIXME: this is to get rid of the "; charset=utf-8" part.
	respCT := strings.Split(resp.Header.Get("Content-Type"), ";")[0]
	CTOK := (specResp.ContentType == "") || (specResp.ContentType == respCT)

	headersOK := true
	for specHeaderName, specHeaderValues := range specResp.Headers {
		headersOK = Header(specHeaderName, specHeaderValues, resp.Header.Values(specHeaderName), logger) && headersOK
	}

	if !statusOK {
		logger.ResponseHasWrongStatus(specResp, resp.StatusCode)
	}

	if !CTOK {
		logger.ResponseHasWrongContentType(specResp, resp.Header.Get("Content-Type"))
	}

	return statusOK && CTOK && headersOK
}

// Response test response bodies.
func Response(resp *http.Response, specResp *api.Response, logger log.ILogger) bool {
	httpOK := HTTPResponse(resp, specResp, logger)
	contentOK := true

	if httpOK {
		if specResp.ContentType != "" {
			//TODO: test contents with content-specific tests.
			switch specResp.ContentType {
			case "application/json":
				contentOK = JSONResponse(resp, specResp, logger)

			default:
				fmt.Printf("The Content-Type of '%s' is not supported.\n", specResp.ContentType)
			}
		}
	}

	return httpOK && contentOK
}
