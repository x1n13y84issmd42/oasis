package comp

import (
	"github.com/x1n13y84issmd42/gog/graph/collection"
	gcontract "github.com/x1n13y84issmd42/gog/graph/contract"
)

// Cycle finds a cycle in a graph if there is at least one.
func Cycle(graph gcontract.Graph) *collection.NodeStack {

	for n := range graph.Nodes().Range() {
		visited := collection.NodeVisitMap{}
		stack := &collection.NodeStack{}

		if dfs(graph, n, stack, &visited) {
			cycle := &collection.NodeStack{}
			cycleStart := false
			cycleNode := stack.Pop()

			for _, n := range *stack {
				if n.ID() == cycleNode.ID() {
					cycleStart = true
				}

				if cycleStart {
					cycle.Push(n)
				}
			}

			cycle.Push(cycleNode)

			return cycle
		}
	}

	return &collection.NodeStack{}
}

func dfs(graph gcontract.Graph, n gcontract.Node, stack *collection.NodeStack, visited *collection.NodeVisitMap) bool {
	nID := n.ID()
	stack.Push(n)

	if visited.Visited(nID) {
		return true
	}

	visited.Visit(nID)

	for n := range graph.AdjacentNodes(nID).Range() {
		if dfs(graph, n, stack, visited) {
			return true
		}
	}

	stack.Pop()

	return false
}
