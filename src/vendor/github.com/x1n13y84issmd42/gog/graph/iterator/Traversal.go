package iterator

import "github.com/x1n13y84issmd42/gog/graph/contract"

// Forward is the default traversal direction for directed graphs,
// i.e. it follows the edges' directions.
func Forward(g contract.NodeAccess, n contract.NodeID) contract.Nodes {
	return g.AdjacentNodes(n)
}

// Backward is a reversed traversal direction for directed graphs,
// i.e. it reverses the edges' directions.
func Backward(g contract.NodeAccess, n contract.NodeID) contract.Nodes {
	return g.UpstreamNodes(n)
}

// PreOrder implements a pre-order traversal, i.e. visits a node first,
// then it's adjacent nodes.
func PreOrder(traverseNode func(), traverseAdjacent func()) {
	traverseNode()
	traverseAdjacent()
}

// PostOrder implements a post-order traversal, i.e. visits a node's
// adjacent nodes first, then the node itself.
func PostOrder(traverseNode func(), traverseAdjacent func()) {
	traverseAdjacent()
	traverseNode()
}
