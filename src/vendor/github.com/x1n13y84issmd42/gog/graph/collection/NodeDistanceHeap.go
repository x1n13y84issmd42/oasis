package collection

import (
	"github.com/x1n13y84issmd42/gog/graph/contract"
)

// NodeDistance contains distance to a node.
type NodeDistance struct {
	ID contract.NodeID
	D  float64
}

// NewNodeDistance crates a new NodeDistance instance.
func NewNodeDistance(nID contract.NodeID, d float64) NodeDistance {
	return NodeDistance{
		ID: nID,
		D:  d,
	}
}

// NodeDistanceHeap is a heap of node distances.
// It is used in the PQ implementation.
type NodeDistanceHeap []NodeDistance

// Push adds an element to the heap. Used by heap.Push().
func (heap *NodeDistanceHeap) Push(v interface{}) {
	(*heap) = append(*heap, v.(NodeDistance))
}

// Pop removes an element from the tail of the heap. Used by heap.Pop().
func (heap *NodeDistanceHeap) Pop() interface{} {
	l := len(*heap)
	v := (*heap)[l-1]
	(*heap) = (*heap)[:l-1]
	return v
}

// Len returns length of the heap.
func (heap NodeDistanceHeap) Len() int {
	return len(heap)
}

// Less compares elements.
func (heap NodeDistanceHeap) Less(i, j int) bool {
	return heap[i].D < heap[j].D
}

// Swap swaps elements.
func (heap NodeDistanceHeap) Swap(i, j int) {
	heap[i], heap[j] = heap[j], heap[i]
}
