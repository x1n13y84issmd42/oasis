package comp

import (
	"github.com/x1n13y84issmd42/gog/graph/collection"
	"github.com/x1n13y84issmd42/gog/graph/contract"
	"github.com/x1n13y84issmd42/gog/graph/iterator"
)

// KosarajuSCC finds strongly connected components of a graph using the Sambasiva Rao Kosaraju's algorithm.
func KosarajuSCC(graph contract.Graph) []*collection.Nodes {
	// First DFS pass to collect the nodes into a stack.
	stack := collection.NodeStack{}
	i1 := iterator.DFS(iterator.Forward, iterator.PostOrder)
	for n := range i1.Iterate(graph, graph.Nodes()) {
		stack.Push(n)
	}

	res := []*collection.Nodes{}

	// Second pass over the stack. Each node gets reverse DFS'd, i.e. over inbound edges.
	// I'm reusing the i2's Visited map.
	// This is the reason Iterate() doesn't reset it, by the way.
	i2 := iterator.DFS(iterator.Backward, iterator.PostOrder)
	for n := stack.Pop(); n != nil; n = stack.Pop() {
		if !i2.Visited.Visited(n.ID()) {
			scc := collection.NewNodes()
			for scn := range i2.Iterate(graph, collection.NewNodes(n)) {
				scc.Add(scn)
			}
			res = append(res, scc)
		}
	}

	return res
}
