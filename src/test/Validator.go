package test

import "github.com/x1n13y84issmd42/oasis/src/contract"

// Validator tests an operation response for compliance with the expectations.
type Validator struct {
	expectations []contract.Expectation
	log          contract.Logger
}

// NewValidator creates a new Validator instance.
func NewValidator(log contract.Logger) *Validator {
	return &Validator{
		expectations: []contract.Expectation{},
		log:          log,
	}
}

// Validate runs through the list of expectations and tests the HTTP request
// instance from the OperationResult against them.
func (v *Validator) Validate(r *contract.OperationResult) *contract.OperationResult {
	for _, ex := range v.expectations {
		r.Success = ex(r.HTTPResponse) && r.Success
	}
	return r
}

// Expect adds an Expectation function to the list.
func (v *Validator) Expect(ex contract.Expectation) {
	v.expectations = append(v.expectations, ex)
}
