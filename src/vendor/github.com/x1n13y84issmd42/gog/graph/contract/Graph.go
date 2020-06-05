package contract

// Graph is an interface for unweighted graphs.
type Graph interface {
	NodeAccess
	DFS(n NodeID, traverse TraversalOrder) NChannel
	BFS(n NodeID) NChannel
	RDFS(n NodeID, traverse TraversalOrder) NChannel
	RBFS(n NodeID) NChannel

	AddEdge(v1 Node, v2 Node)
	// OutEdges(nID NodeID) []Edge
	// InEdges(nID NodeID) []Edge
}

// WGraph is an interface for weighted graphs.
type WGraph interface {
	NodeAccess
	AddEdge(v1 Node, w float64, v2 Node)
}
