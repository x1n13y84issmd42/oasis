package contract

// TraversalOrder is a traversal order function for regular graphs.
type TraversalOrder func(func(), func())

// BinaryTraversalOrder is a traversal order function for binary trees.
type BinaryTraversalOrder func(func(), func(), func())

// TraversalDirection determines traversal direction by providing sets of nodes
// to "consider adjacent" to a particular node N.
// Implementations include forward & backward (reversed or upwards) traversal orders.
type TraversalDirection func(NodeAccess, NodeID) Nodes
