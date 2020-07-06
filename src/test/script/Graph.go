package script

import (
	gog "github.com/x1n13y84issmd42/gog/graph"
	gcontract "github.com/x1n13y84issmd42/gog/graph/contract"
	contract "github.com/x1n13y84issmd42/oasis/src/contract"
)

// OperationNode represents a single operation in the script execution graph.
type OperationNode struct {
	Operation contract.Operation
	// Parameters api.OperationParameters
}

// ID returns a uniqe operation node ID.
func (node *OperationNode) ID() gcontract.NodeID {
	return gcontract.NodeID(node.Operation.Name())
}

// OperationGraph is a graph representing interdependencies between operations.
type OperationGraph struct {
	contract.EntityTrait
	*gog.DGraph
}

// NewOperationGraph creates a new OperationGraph instance.
func NewOperationGraph(log contract.Logger) *OperationGraph {
	return &OperationGraph{
		EntityTrait: contract.Entity(log),
		DGraph:      gog.NewDGraph(),
	}
}
