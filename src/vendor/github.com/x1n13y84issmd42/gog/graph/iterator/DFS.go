package iterator

import (
	"github.com/x1n13y84issmd42/gog/graph/collection"
	"github.com/x1n13y84issmd42/gog/graph/contract"
)

// DFSIterator implements depth-first search to traverse graphs.
type DFSIterator struct {
	ch       contract.NChannel
	traverse contract.TraversalOrder
	next     contract.TraversalDirection
	Visited  collection.NodeVisitMap
}

// Iterate performs DFS starting from each of provided nodes.
func (i *DFSIterator) Iterate(graph contract.NodeAccess, nodes contract.Nodes) contract.NChannel {
	i.ch = make(contract.NChannel)
	go func() {
		for n := range nodes.Range() {
			i.walk(graph, n)
		}
		close(i.ch)
	}()

	return i.ch
}

func (i *DFSIterator) walk(graph contract.NodeAccess, node contract.Node) {
	nID := node.ID()
	if i.Visited.Visited(nID) {
		return
	}

	i.Visited.Visit(nID)
	i.traverse(func() {
		i.ch <- node
	}, func() {
		for n := range i.next(graph, nID).Range() {
			i.walk(graph, n)
		}
	})
}

// DFS creates a depth-first search iterator to traverse the graph.
func DFS(next contract.TraversalDirection, traverse contract.TraversalOrder) *DFSIterator {
	return &DFSIterator{
		next:     next,
		traverse: traverse,
		Visited:  collection.NodeVisitMap{},
	}
}
