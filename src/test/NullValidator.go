package test

import (
	"github.com/x1n13y84issmd42/oasis/src/contract"
	"github.com/x1n13y84issmd42/oasis/src/errors"
)

// NullValidator is used whenever we can't have a real one from DataResolver.Response().
// Reports the contained error on every method call.
type NullValidator struct {
	errors.NullObjectPrototype
}

// NoValidator creates a new NullValidator instance.
func NoValidator(err error, log contract.Logger) *NullValidator {
	return &NullValidator{
		NullObjectPrototype: errors.NullObject(err, log),
	}
}

// Validate runs through the list of expectations and tests the HTTP request
// instance from the OperationResult against them.
func (v *NullValidator) Validate(r *contract.OperationResult) *contract.OperationResult {
	v.Report()
	return nil
}

// Expect adds an Expectation function to the list.
func (v *NullValidator) Expect(ex contract.Expectation) {
	v.Report()
}
