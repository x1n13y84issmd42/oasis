package iterator

import (
	"github.com/x1n13y84issmd42/gog/graph/collection"
	"github.com/x1n13y84issmd42/gog/graph/contract"
)

// BFSIterator implements depth-first search to traverse graphs.
type BFSIterator struct {
	ch      contract.NChannel
	next    contract.TraversalDirection
	stack   collection.NodeStack
	Visited collection.NodeVisitMap
}

// Iterate performs BFS starting from each of provided nodes.
func (i *BFSIterator) Iterate(graph contract.NodeAccess, nodes contract.Nodes) contract.NChannel {
	i.ch = make(contract.NChannel)
	go func() {
		for n := range nodes.Range() {
			i.stack.Push(n)
			i.walk(graph, n)
		}
		close(i.ch)
	}()

	return i.ch
}

func (i *BFSIterator) walk(graph contract.NodeAccess, node contract.Node) {
	for len(i.stack) > 0 {
		node := i.stack.PopFront()
		nID := node.ID()

		if node != nil && !i.Visited.Visited(nID) {
			i.ch <- node
			i.Visited.Visit(nID)

			next := i.next(graph, nID)
			if next.Count() > 0 {
				i.stack.Append(next.Values())
			}
		}
	}
}

// BFS creates a breadth-first search iterator to traverse nodes.
func BFS(next contract.TraversalDirection) *BFSIterator {
	return &BFSIterator{
		next:    next,
		Visited: collection.NodeVisitMap{},
		stack:   collection.NodeStack{},
	}
}
