package contract

import (
	"net/http"
	"net/url"

	"github.com/x1n13y84issmd42/oasis/src/env"
)

// PathParameters is a map of named path parameters.
type PathParameters map[string]string

// QueryValues is an alias type for url.Values to add the ForEach method to it.
type QueryValues url.Values

// HTTPHeaders is an alias type for http.Header to add the ForEach method to it.
type HTTPHeaders http.Header

// OperationRequestParameters is operation request parameters.
type OperationRequestParameters struct {
	HostHint        string
	ContentTypeHint string
	Headers         HTTPHeaders
}

// OperationResponseExpectation is operation response expectations.
// Use hints to tell Oasis which response definition to pick from a spec.
type OperationResponseExpectation struct {
	StatusHint      uint64
	ContentTypeHint string
	Headers         HTTPHeaders
}

// OperationSecurityParameters is operation security parameters.
// Use hints to tell Oasis which security definition to pick from a spec.
type OperationSecurityParameters struct {
	SecurityHint  string
	HTTPAuthValue string
}

// OperationParameters is operation parameters. Use it to fine tune Oasis
// by giving hints regarding selection of example data, overriding
// request data (path parameters, headers & bodies) & security settings.
// It is used both during spec retrieval (to fix operations' paths) and
// in testing.
type OperationParameters struct {
	Security OperationSecurityParameters
	Request  OperationRequestParameters
	Response OperationResponseExpectation
	Path     PathParameters
	Query    QueryValues
}

// LoadFromArgs fills the instance with values from the application CLI arguments object.
func (params *OperationParameters) LoadFromArgs(args *env.Args) {
	params.Request.ContentTypeHint = args.Use.CT
	params.Request.HostHint = args.Host
	params.Response.ContentTypeHint = args.Expect.CT
	params.Response.StatusHint = uint64(args.Expect.Status)
	params.Security.SecurityHint = args.Use.Security
	params.Path = PathParameters(args.Use.PathParameters)
	params.Query = QueryValues(args.Use.Query)

	// fmt.Printf("Operation params: %#v\n", params)
}

// ParameterHandler is a handler function type to use in ForEach methods of input
// parameter types (PathParameters, QueryValues, Headers).
type ParameterHandler func(n string, v string)

// ForEach iterates over path parameters and invokes the provided handler function
// with name & value of each of them.
func (pp *PathParameters) ForEach(handler ParameterHandler) {
	for ppn, ppv := range *pp {
		handler(ppn, ppv)
	}
}

// ForEach iterates over query values and invokes the provided handler function
// with name & value of each of them.
func (values *QueryValues) ForEach(handler ParameterHandler) {
	for qpn, qpvs := range *values {
		for _, qpv := range qpvs {
			handler(qpn, qpv)
		}
	}
}

// ForEach iterates over header values and invokes the provided handler function
// with name & value of each of them.
func (headers *HTTPHeaders) ForEach(handler ParameterHandler) {
	for hn, hvs := range *headers {
		for _, hv := range hvs {
			handler(hn, hv)
		}
	}
}
