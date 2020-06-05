package storage

import "github.com/x1n13y84issmd42/gog/graph/contract"

// AdjacencyMatrix is a representation of graphs as a matrix of adjacent nodes.
type AdjacencyMatrix struct {
	Nodes contract.Nodes
	M     [][]contract.Node
}
