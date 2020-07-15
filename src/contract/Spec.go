package contract

// OperationIterator is an iterable channel to receive operations from.
type OperationIterator chan Operation

// ProjectInfo is a generic project information.
type ProjectInfo interface {
	Title() string
	Description() string
	Version() string
}

// OperationAccess is an interafce for operation access.
type OperationAccess interface {
	Operations() OperationIterator
	GetOperation(id string) Operation
}

// Spec is an interface to access specification data.
type Spec interface {
	ProjectInfo
	OperationAccess
}
