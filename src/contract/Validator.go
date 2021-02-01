package contract

// Validator implements validation/testing logic.
// It is specific to spec formats and comes from there (via Operation).
type Validator interface {
	IEntityTrait
	Validate(*OperationResult) *OperationResult
	Expect(Expectation)
}
