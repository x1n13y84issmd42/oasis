package graph

import (
	"github.com/x1n13y84issmd42/gog/graph/collection"
	"github.com/x1n13y84issmd42/gog/graph/contract"
	"github.com/x1n13y84issmd42/gog/graph/iterator"
	"github.com/x1n13y84issmd42/gog/graph/storage"
)

// DEdge is a pair of nodes making an undirected edge.
type DEdge struct {
	A contract.Node
	B contract.Node
}

// Reverse creates a new edge by swapping the receiver's nodes.
func (e DEdge) Reverse() contract.Edge {
	return DEdge{
		A: e.B,
		B: e.A,
	}
}

// DGraph is an unweighted directed graph.
type DGraph struct {
	adjacency contract.IAdjacency
}

// NewDGraph creates a new DGraph instance.
// Provided nodes will be added pairwise as edges.
func NewDGraph(nodes ...contract.Node) *DGraph {
	g := &DGraph{
		adjacency: storage.NewAdjacencyList(),
	}

	for i := 0; i < len(nodes); i++ {
		g.AddNode(nodes[i])
	}

	return g
}

// AddNode adds a node to the graph.
func (graph *DGraph) AddNode(n contract.Node) {
	graph.adjacency.AddNode(n)
}

// AddEdge creates an edge between v1 and v2 nodes.
func (graph *DGraph) AddEdge(v1 contract.NodeID, v2 contract.NodeID) {
	graph.adjacency.AddEdge(v1, v2)
}

// Node returns a node instance by it's ID.
func (graph *DGraph) Node(nID contract.NodeID) contract.Node {
	return graph.adjacency.Node(nID)
}

// Nodes returns a set of all graph's nodes.
func (graph *DGraph) Nodes() contract.Nodes {
	return graph.adjacency.Nodes()
}

// AdjacentNodes returns a list of adjacent nodes for a node defined by nID.
func (graph *DGraph) AdjacentNodes(nID contract.NodeID) contract.Nodes {
	return graph.adjacency.AdjacentNodes(nID)
}

// UpstreamNodes returns a list of adjacent nodes for a node defined by nID.
func (graph *DGraph) UpstreamNodes(nID contract.NodeID) contract.Nodes {
	return graph.adjacency.UpstreamNodes(nID)
}

// DFS returns a DFS node iterator.
func (graph *DGraph) DFS(nID contract.NodeID, traverse contract.TraversalOrder) contract.NChannel {
	return iterator.DFS(iterator.Forward, traverse).Iterate(graph, collection.NewNodes(graph.Node(nID)))
}

// BFS returns a BFS node iterator.
func (graph *DGraph) BFS(nID contract.NodeID) contract.NChannel {
	return iterator.BFS(iterator.Forward).Iterate(graph, collection.NewNodes(graph.Node(nID)))
}

// RDFS returns a reversed DFS node iterator.
func (graph *DGraph) RDFS(nID contract.NodeID, traverse contract.TraversalOrder) contract.NChannel {
	return iterator.DFS(iterator.Backward, traverse).Iterate(graph, collection.NewNodes(graph.Node(nID)))
}

// RBFS returns a RBFS node iterator.
func (graph *DGraph) RBFS(nID contract.NodeID) contract.NChannel {
	return iterator.BFS(iterator.Backward).Iterate(graph, collection.NewNodes(graph.Node(nID)))
}

// OutEdges returns a list of outbound edges for a node defined by nID.
func (graph *DGraph) OutEdges(nID contract.NodeID) []DEdge {
	res := []DEdge{}
	n := graph.Node(nID)
	if n != nil {
		for na := range graph.AdjacentNodes(nID).Range() {
			res = append(res, DEdge{
				A: n,
				B: na,
			})
		}
	}

	return res
}

// InEdges returns a list of inbound edges for a node defined by nID.
func (graph *DGraph) InEdges(nID contract.NodeID) []DEdge {
	res := []DEdge{}
	n := graph.Node(nID)
	if n != nil {
		for na := range graph.UpstreamNodes(nID).Range() {
			res = append(res, DEdge{
				A: na,
				B: n,
			})
		}
	}

	return res
}
