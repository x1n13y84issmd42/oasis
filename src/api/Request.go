package api

import (
	"fmt"
	"net/http"
	"net/url"
)

// Request describes the properties of an HTTP request.
type Request struct {
	Method  string
	Path    string
	Query   *url.Values
	Headers http.Header
}

// CreateRequest creates a Request instance, already configured
// to make requests to the operation URL.
func (specReq *Request) CreateRequest(host *Host) *http.Request {
	URL := fmt.Sprintf("%s%s", host.URL, specReq.Path)
	req, _ := http.NewRequest(specReq.Method, URL, nil)
	req.URL.RawQuery = specReq.Query.Encode()
	//TODO: add headers
	return req
}
