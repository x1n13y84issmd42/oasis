package collection

import "github.com/x1n13y84issmd42/gog/graph/contract"

// NodeStack is a slice of nodes.
type NodeStack []contract.Node

// Push adds a node to the end of the stack.
func (stack *NodeStack) Push(node contract.Node) {
	*stack = append(*stack, node)
}

// Append appends a list of nodes to the end of the stack.
func (stack *NodeStack) Append(nodes []contract.Node) {
	*stack = append(*stack, nodes...)
}

// Pop removes the first node from the stack and returns it.
func (stack *NodeStack) Pop() contract.Node {
	s := *stack
	if len(s) >= 1 {
		res := s[len(s)-1]
		*stack = s[:len(s)-1]
		return res
	}

	return nil
}

// PopFront removes the first node from the stack and returns it.
func (stack *NodeStack) PopFront() contract.Node {
	s := *stack
	if len(s) >= 1 {
		res := s[0]
		*stack = s[1:]
		return res
	}

	return nil
}
