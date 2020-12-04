package contract

// Graph is an interface for unweighted graphs.
type Graph interface {
	NodeAccess
	DFS(n NodeID, traverse TraversalOrder) NChannel
	BFS(n NodeID) NChannel
	RDFS(n NodeID, traverse TraversalOrder) NChannel
	RBFS(n NodeID) NChannel

	Len() uint

	AddNode(n Node)
	AddEdge(v1 NodeID, v2 NodeID)
	//TODO: RemoveEdge(v1 Node, v2 Node)
}

// WGraph is an interface for weighted graphs.
type WGraph interface {
	NodeAccess
	AddNode(n Node)
	AddEdge(v1 NodeID, w float64, v2 NodeID)
	Len() uint
	OutEdges(n NodeID) []WEdge
}
