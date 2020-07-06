package collection

import "github.com/x1n13y84issmd42/gog/graph/contract"

// NodeVisitMap is a map of visited nodes.
// Used in iterators.
type NodeVisitMap map[contract.NodeID]int64

// Visit marks nID as visited and increments the visit counter.
func (m *NodeVisitMap) Visit(nID contract.NodeID) {
	(*m)[nID]++
}

// Visited returns true if nID has been visited before.
func (m *NodeVisitMap) Visited(nID contract.NodeID) bool {
	_, visited := (*m)[nID]
	return visited
}

// Visits tells exactly how many time nID has been visited.
func (m *NodeVisitMap) Visits(nID contract.NodeID) int64 {
	return (*m)[nID]
}
