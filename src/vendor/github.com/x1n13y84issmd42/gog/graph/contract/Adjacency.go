package contract

// NodeAccess is an interface to node access within graph.
type NodeAccess interface {
	Nodes() Nodes
	Node(nID NodeID) Node
	AdjacentNodes(nID NodeID) Nodes
	UpstreamNodes(nID NodeID) Nodes
}

// IAdjacency is an interface to access node's adjacency info.
type IAdjacency interface {
	NodeAccess
	AddNode(n Node)
	AddEdge(v1 NodeID, v2 NodeID)
	Len() uint
}

// WEdge is a weighted edge.
type WEdge struct {
	A Node
	B Node
	W float64
}

// IWAdjacency the same as IAdjacency, but for weighted nodes.
type IWAdjacency interface {
	NodeAccess
	AddNode(n Node)
	AddEdge(v1 NodeID, w float64, v2 NodeID)
	OutEdges(v NodeID) []WEdge
	Len() uint
}
