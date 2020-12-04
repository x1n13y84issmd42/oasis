package contract

//NodeID is a value, unique within a graph, which determines a node's identity.
type NodeID string

// Node is an intefrace for basic graph nodes.
type Node interface {
	ID() NodeID
}
