package graph

import (
	"github.com/x1n13y84issmd42/gog/graph/collection"
	"github.com/x1n13y84issmd42/gog/graph/contract"
	"github.com/x1n13y84issmd42/gog/graph/iterator"
	"github.com/x1n13y84issmd42/gog/graph/storage"
)

// UEdge is an undirected graph edge.
type UEdge DEdge

// Reverse creates a copy of the receiver edge.
func (e UEdge) Reverse() contract.Edge {
	return UEdge(e)
}

// UGraph is an unweighted undirected graph.
type UGraph struct {
	A contract.IAdjacency
}

// NewUGraph creates a new UGraph instance.
// Provided nodes will be added pairwise as edges.
func NewUGraph(nodes ...contract.Node) *UGraph {
	g := &UGraph{
		A: storage.NewAdjacencyList(),
	}

	for i := 0; i < len(nodes); i += 2 {
		g.AddEdge(nodes[i], nodes[i+1])
	}

	return g
}

// AddEdge creates an edge between v1 and v2 nodes.
func (graph *UGraph) AddEdge(v1 contract.Node, v2 contract.Node) {
	graph.A.AddEdge(v1, v2)
	graph.A.AddEdge(v2, v1)
}

// Node returns a node instance by it's ID.
func (graph *UGraph) Node(nID contract.NodeID) contract.Node {
	return graph.A.Node(nID)
}

// Nodes returns a set of all graph's nodes.
func (graph *UGraph) Nodes() contract.Nodes {
	return graph.A.Nodes()
}

// AdjacentNodes returns a list of adjacent nodes for a node defined by nID.
func (graph *UGraph) AdjacentNodes(nID contract.NodeID) contract.Nodes {
	return graph.A.AdjacentNodes(nID)
}

// UpstreamNodes returns a list of adjacent nodes for a node defined by nID.
func (graph *UGraph) UpstreamNodes(nID contract.NodeID) contract.Nodes {
	return graph.A.UpstreamNodes(nID)
}

// DFS returns a DFS node iterator.
func (graph *UGraph) DFS(nID contract.NodeID, traverse contract.TraversalOrder) contract.NChannel {
	return iterator.DFS(iterator.Forward, traverse).Iterate(graph, collection.NewNodes(graph.Node(nID)))
}

// BFS returns a BFS node iterator.
func (graph *UGraph) BFS(nID contract.NodeID) contract.NChannel {
	return iterator.BFS(iterator.Forward).Iterate(graph, collection.NewNodes(graph.Node(nID)))
}

// RDFS returns a reversed DFS node iterator.
func (graph *UGraph) RDFS(nID contract.NodeID, traverse contract.TraversalOrder) contract.NChannel {
	return iterator.DFS(iterator.Backward, traverse).Iterate(graph, collection.NewNodes(graph.Node(nID)))
}

// RBFS returns a RBFS node iterator.
func (graph *UGraph) RBFS(nID contract.NodeID) contract.NChannel {
	return iterator.BFS(iterator.Backward).Iterate(graph, collection.NewNodes(graph.Node(nID)))
}

// OutEdges returns a list of outbound edges for a node defined by nID.
func (graph *UGraph) OutEdges(nID contract.NodeID) []UEdge {
	res := []UEdge{}
	n := graph.Node(nID)
	if n != nil {
		for na := range graph.AdjacentNodes(nID).Range() {
			res = append(res, UEdge{
				A: n,
				B: na,
			})
		}
	}

	return res
}

// InEdges returns a list of inbound edges for a node defined by nID.
func (graph *UGraph) InEdges(nID contract.NodeID) []UEdge {
	res := []UEdge{}
	n := graph.Node(nID)
	if n != nil {
		for na := range graph.UpstreamNodes(nID).Range() {
			res = append(res, UEdge{
				A: na,
				B: n,
			})
		}
	}

	return res
}
