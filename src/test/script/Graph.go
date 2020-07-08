package script

import (
	gog "github.com/x1n13y84issmd42/gog/graph"
	gcontract "github.com/x1n13y84issmd42/gog/graph/contract"
	contract "github.com/x1n13y84issmd42/oasis/src/contract"
)

// ExecutionNode represents a single operation in the script execution graph.
type ExecutionNode struct {
	Operation contract.Operation
	OpID      string
}

// ID returns a uniqe operation node ID.
func (node *ExecutionNode) ID() gcontract.NodeID {
	return gcontract.NodeID(node.OpID)
}

// ExecutionGraph is a graph representing interdependencies between operations.
type ExecutionGraph struct {
	contract.EntityTrait
	*gog.DGraph
}

// NewExecutionGraph creates a new OperationGraph instance.
func NewExecutionGraph(log contract.Logger) *ExecutionGraph {
	return &ExecutionGraph{
		EntityTrait: contract.Entity(log),
		DGraph:      gog.NewDGraph(),
	}
}
