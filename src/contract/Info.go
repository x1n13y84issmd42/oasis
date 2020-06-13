package contract

// ProjectInfo is a generic project information.
type ProjectInfo interface {
	Title() string
	Description() string
	Version() string
}
