package test

import "github.com/x1n13y84issmd42/oasis/src/contract"

// NullRequest ...
type NullRequest struct {
	contract.NullObjectPrototype
}

// NoRequest creates a new NullRequest instance.
func NoRequest(err error, log contract.Logger) contract.Request {
	return &NullRequest{
		NullObjectPrototype: contract.NullObject(err, log),
	}
}

// Enrich reports an error.
func (req *NullRequest) Enrich(en contract.RequestEnrichment) {
	req.Report()
}

// Execute reports an error.
func (req *NullRequest) Execute() *contract.OperationResult {
	req.Report()
	return nil
}
