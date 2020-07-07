package comp

import (
	"github.com/x1n13y84issmd42/gog/graph/collection"
	gcontract "github.com/x1n13y84issmd42/gog/graph/contract"
	"github.com/x1n13y84issmd42/gog/graph/iterator"
)

// MotherNode finds a node from which every other node
// in the graph is accessible.
func MotherNode(graph gcontract.Graph) gcontract.Node {

	// The last visited node, a candidate for mother node.
	var lastNode gcontract.Node

	// Iterator iterates and collects a set of visited maps.
	i := iterator.DFS(iterator.Forward, iterator.PostOrder)

	for n := range graph.Nodes().Range() {
		if !i.Visited.Visited(n.ID()) {
			for range i.Iterate(graph, collection.NewNodes(n)) {
				//
			}
			lastNode = n
		}
	}

	// Second DFS starts from the lastNode and collects a set of visited nodes.
	i1 := iterator.DFS(iterator.Forward, iterator.PostOrder)
	for range i1.Iterate(graph, collection.NewNodes(lastNode)) {
		//
	}

	// Checking if i1 has visited all the nodes.
	for n := range graph.Nodes().Range() {
		if !i1.Visited.Visited(n.ID()) {
			return nil
		}
	}

	return lastNode
}
