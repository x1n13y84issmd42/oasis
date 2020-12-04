package script

import (
	gcontract "github.com/x1n13y84issmd42/gog/graph/contract"
	"github.com/x1n13y84issmd42/oasis/src/contract"
	"github.com/x1n13y84issmd42/oasis/src/errors"
)

// NullGraph is used whenever we can't have a real execution graph from a script.
// Reports the contained error on every method call.
type NullGraph struct {
	errors.NullObjectPrototype
}

// NoGraph creates a new NullGraph instance.
func NoGraph(err error, log contract.Logger) gcontract.Graph {
	return &NullGraph{
		NullObjectPrototype: errors.NullObject(err, log),
	}
}

// Nodes reports an error.
func (g NullGraph) Nodes() gcontract.Nodes {
	g.Report()
	return nil
}

// Node reports an error.
func (g NullGraph) Node(nID gcontract.NodeID) gcontract.Node {
	g.Report()
	return nil
}

// AdjacentNodes reports an error.
func (g NullGraph) AdjacentNodes(nID gcontract.NodeID) gcontract.Nodes {
	g.Report()
	return nil
}

// UpstreamNodes reports an error.
func (g NullGraph) UpstreamNodes(nID gcontract.NodeID) gcontract.Nodes {
	g.Report()
	return nil
}

// DFS reports an error.
func (g NullGraph) DFS(n gcontract.NodeID, traverse gcontract.TraversalOrder) gcontract.NChannel {
	g.Report()
	return nil
}

// BFS reports an error.
func (g NullGraph) BFS(n gcontract.NodeID) gcontract.NChannel {
	g.Report()
	return nil
}

// RDFS reports an error.
func (g NullGraph) RDFS(n gcontract.NodeID, traverse gcontract.TraversalOrder) gcontract.NChannel {
	g.Report()
	return nil
}

// RBFS reports an error.
func (g NullGraph) RBFS(n gcontract.NodeID) gcontract.NChannel {
	g.Report()
	return nil
}

// AddNode reports an error.
func (g NullGraph) AddNode(n gcontract.Node) {
	g.Report()
}

// AddEdge reports an error.
func (g NullGraph) AddEdge(v1 gcontract.NodeID, v2 gcontract.NodeID) {
	g.Report()
}

// Len reports an error.
func (g *NullGraph) Len() uint {
	g.Report()
	return 0
}
