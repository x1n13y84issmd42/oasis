package storage

import (
	"github.com/x1n13y84issmd42/gog/graph/collection"
	"github.com/x1n13y84issmd42/gog/graph/contract"
)

type adjWN struct {
	ID contract.NodeID
	W  float64
}

// adjWList is a mapping between node IDs and their instances.
type adjWList map[contract.NodeID][]adjWN

// WAdjacencyList is a list of nodes adjacent to other nodes.
type WAdjacencyList struct {
	nodes contract.Nodes
	list  adjWList
}

// NewWAdjacencyList creates a new WAdjacencyList instance.
func NewWAdjacencyList() *WAdjacencyList {
	return &WAdjacencyList{
		nodes: collection.NewNodes(),
		list:  adjWList{},
	}
}

// AddNode adds a node to the list of nodes.
func (list *WAdjacencyList) AddNode(n contract.Node) {
	list.nodes.Add(n)
}

// AddEdge adds a v2 to the adjacency list of v1.
func (list *WAdjacencyList) AddEdge(v1ID contract.NodeID, w float64, v2ID contract.NodeID) {
	v1 := list.Node(v1ID)
	v2 := list.Node(v2ID)

	if v1 == nil {
		panic("There is no node with ID=" + v1ID)
	}

	if v2 == nil {
		panic("There is no node with ID=" + v2ID)
	}

	if list.list[v1ID] == nil {
		list.list[v1ID] = []adjWN{}
	}

	for _, n := range list.list[v1ID] {
		if n.ID == v2ID {
			return
		}
	}

	list.list[v1ID] = append(list.list[v1ID], adjWN{ID: v2ID, W: w})
}

// Node returns a node instance by it's ID.
func (list *WAdjacencyList) Node(nID contract.NodeID) contract.Node {
	if list.nodes.Has(nID) {
		return list.nodes.Get(nID)
	}

	return nil
}

// Nodes returns a set of all nodes.
func (list *WAdjacencyList) Nodes() contract.Nodes {
	return list.nodes
}

// Len returns number of nodes in the list.
func (list *WAdjacencyList) Len() uint {
	return list.nodes.Count()
}

// AdjacentNodes returns a set of nodes adjacent to n.
func (list *WAdjacencyList) AdjacentNodes(nID contract.NodeID) contract.Nodes {
	if list.list[nID] != nil {
		nodes := collection.NewNodes()
		for _, n := range list.list[nID] {
			nodes.Add(list.nodes.Get(n.ID))
		}

		return nodes
	}

	return collection.NewNodes()
}

// UpstreamNodes returns a set of nodes adjacent to n.
func (list *WAdjacencyList) UpstreamNodes(nID contract.NodeID) contract.Nodes {
	res := collection.NewNodes()

	for upID, adjacent := range list.list {
		for _, nA := range adjacent {
			if nA.ID == nID {
				res.Add(list.nodes.Get(upID))
				break
			}
		}
	}

	return res
}

// OutEdges returns a list of edges outgoing from nID.
func (list *WAdjacencyList) OutEdges(nID contract.NodeID) []contract.WEdge {
	if list.list[nID] != nil {
		n := list.nodes.Get(nID)
		edges := []contract.WEdge{}
		for _, an := range list.list[nID] {
			edges = append(edges, contract.WEdge{A: n, B: list.nodes.Get(an.ID), W: an.W})
		}

		return edges
	}

	return []contract.WEdge{}
}
