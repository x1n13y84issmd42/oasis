package api

import (
	"net/http"

	"github.com/x1n13y84issmd42/oasis/src/env"
)

// OperationRequestParameters is operation request parameters.
type OperationRequestParameters struct {
	ContentTypeHint string
	Headers         http.Header
}

// OperationResponseExpectation is operation response expectations.
// Use hints to tell Oasis which response definition to pick from a spec.
type OperationResponseExpectation struct {
	StatusHint      uint
	ContentTypeHint string
	Headers         http.Header
}

// OperationSecurityParameters is operation security parameters.
// Use hints to tell Oasis which response definition to pick from a spec.
type OperationSecurityParameters struct {
	SecurityHint    string
	HTTPBasicValue  string
	HTTPDigestValue string
}

// PathParameters is a map of named path parameters.
type PathParameters map[string]string

// OperationParameters is operation parameters. Use it to fine tune Oasis
// by giving hints regarding selection of example data, overriding
// request data (path parameters, headers & bodies) & security settings.
type OperationParameters struct {
	Security OperationSecurityParameters
	Request  OperationRequestParameters
	Response OperationResponseExpectation
	Path     PathParameters
}

// LoadFromArgs fills the instance with values from the application CLI arguments object.
func (params *OperationParameters) LoadFromArgs(args *env.Args) {
	params.Request.ContentTypeHint = args.Use.CT
	params.Response.ContentTypeHint = args.Expect.CT
	params.Response.StatusHint = uint(args.Expect.Status)
	params.Security.SecurityHint = args.Use.Security
}
