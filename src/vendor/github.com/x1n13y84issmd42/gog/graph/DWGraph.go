package graph

import (
	"github.com/x1n13y84issmd42/gog/graph/collection"
	"github.com/x1n13y84issmd42/gog/graph/contract"
	"github.com/x1n13y84issmd42/gog/graph/iterator"
	"github.com/x1n13y84issmd42/gog/graph/storage"
)

// DWEdge is a pair of nodes making an undirected edge.
type DWEdge contract.WEdge

// Reverse creates a new edge by swapping the receiver's nodes.
func (e DWEdge) Reverse() contract.Edge {
	return DWEdge{
		A: e.B,
		B: e.A,
		W: e.W,
	}
}

// DWGraph is an unweighted directed graph.
type DWGraph struct {
	adjacency contract.IWAdjacency
}

// NewDWGraph creates a new DWGraph instance.
// Provided nodes will be added pairwise as edges.
func NewDWGraph(nodes ...contract.Node) *DWGraph {
	g := &DWGraph{
		adjacency: storage.NewWAdjacencyList(),
	}

	for i := 0; i < len(nodes); i++ {
		g.AddNode(nodes[i])
	}

	return g
}

// AddNode adds a node to the graph.
func (graph *DWGraph) AddNode(n contract.Node) {
	graph.adjacency.AddNode(n)
}

// AddEdge creates an edge between v1 and v2 nodes.
func (graph *DWGraph) AddEdge(v1 contract.NodeID, w float64, v2 contract.NodeID) {
	graph.adjacency.AddEdge(v1, w, v2)
}

// Node returns a node instance by it's ID.
func (graph *DWGraph) Node(nID contract.NodeID) contract.Node {
	return graph.adjacency.Node(nID)
}

// Nodes returns a set of all graph's nodes.
func (graph *DWGraph) Nodes() contract.Nodes {
	return graph.adjacency.Nodes()
}

// Len returns number of nodes in the graph.
func (graph *DWGraph) Len() uint {
	return graph.adjacency.Len()
}

// AdjacentNodes returns a list of adjacent nodes for a node defined by nID.
func (graph *DWGraph) AdjacentNodes(nID contract.NodeID) contract.Nodes {
	return graph.adjacency.AdjacentNodes(nID)
}

// UpstreamNodes returns a list of adjacent nodes for a node defined by nID.
func (graph *DWGraph) UpstreamNodes(nID contract.NodeID) contract.Nodes {
	return graph.adjacency.UpstreamNodes(nID)
}

// DFS returns a DFS node iterator.
func (graph *DWGraph) DFS(nID contract.NodeID, traverse contract.TraversalOrder) contract.NChannel {
	return iterator.DFS(iterator.Forward, traverse).Iterate(graph, collection.NewNodes(graph.Node(nID)))
}

// BFS returns a BFS node iterator.
func (graph *DWGraph) BFS(nID contract.NodeID) contract.NChannel {
	return iterator.BFS(iterator.Forward).Iterate(graph, collection.NewNodes(graph.Node(nID)))
}

// RDFS returns a reversed DFS node iterator.
func (graph *DWGraph) RDFS(nID contract.NodeID, traverse contract.TraversalOrder) contract.NChannel {
	return iterator.DFS(iterator.Backward, traverse).Iterate(graph, collection.NewNodes(graph.Node(nID)))
}

// RBFS returns a RBFS node iterator.
func (graph *DWGraph) RBFS(nID contract.NodeID) contract.NChannel {
	return iterator.BFS(iterator.Backward).Iterate(graph, collection.NewNodes(graph.Node(nID)))
}

// OutEdges returns a list of outbound edges for a node defined by nID.
func (graph *DWGraph) OutEdges(nID contract.NodeID) []contract.WEdge {
	res := []contract.WEdge{}

	for _, e := range graph.adjacency.OutEdges(nID) {
		res = append(res, e)
	}

	return res
}

// InEdges returns a list of inbound edges for a node defined by nID.
func (graph *DWGraph) InEdges(nID contract.NodeID) []DWEdge {
	res := []DWEdge{}
	n := graph.Node(nID)
	if n != nil {
		for na := range graph.UpstreamNodes(nID).Range() {
			res = append(res, DWEdge{
				A: na,
				B: n,
			})
		}
	}

	return res
}
