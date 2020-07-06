package storage

import (
	"github.com/x1n13y84issmd42/gog/graph/collection"
	"github.com/x1n13y84issmd42/gog/graph/contract"
)

// adjList is a mapping between node IDs and their instances.
type adjList map[contract.NodeID]contract.Nodes

// AdjacencyList is a list of nodes adjacent to other nodes.
type AdjacencyList struct {
	nodes contract.Nodes
	list  adjList
}

// NewAdjacencyList creates a new AdjacencyList instance.
func NewAdjacencyList() *AdjacencyList {
	return &AdjacencyList{
		nodes: collection.NewNodes(),
		list:  adjList{},
	}
}

// AddEdge adds a v2 to the adjacency list of v1.
func (list *AdjacencyList) AddEdge(v1 contract.Node, v2 contract.Node) {
	v1ID := v1.ID()
	v2ID := v2.ID()

	if list.list[v1ID] == nil {
		list.list[v1ID] = collection.NewNodes()
	}

	for n := range list.list[v1ID].Range() {
		if n.ID() == v2ID {
			return
		}
	}

	list.nodes.Add(v1)
	list.nodes.Add(v2)

	list.list[v1ID].Add(v2)
}

// Node returns a node instance by it's ID.
func (list *AdjacencyList) Node(nID contract.NodeID) contract.Node {
	if list.nodes.Has(nID) {
		return list.nodes.Get(nID)
	}

	return nil
}

// Nodes returns a set of all nodes.
func (list *AdjacencyList) Nodes() contract.Nodes {
	return list.nodes
}

// AdjacentNodes returns a set of nodes adjacent to n.
func (list *AdjacencyList) AdjacentNodes(nID contract.NodeID) contract.Nodes {
	if list.list[nID] != nil {
		return list.list[nID]
	}

	return collection.NewNodes()
}

// UpstreamNodes returns a set of nodes adjacent to n.
func (list *AdjacencyList) UpstreamNodes(nID contract.NodeID) contract.Nodes {
	res := collection.NewNodes()

	for upID, adjacent := range list.list {
		for nA := range adjacent.Range() {
			if nA.ID() == nID {
				res.Add(list.nodes.Get(upID))
				break
			}
		}
	}

	return res
}
