package script

import (
	gog "github.com/x1n13y84issmd42/gog/graph"
	gcontract "github.com/x1n13y84issmd42/gog/graph/contract"
	"github.com/x1n13y84issmd42/oasis/src/contract"
	"github.com/x1n13y84issmd42/oasis/src/params"
)

// ExecutionNode represents a single operation in the script execution graph.
// It contains a pointer to an operation and keeps an OperationData instance
// where operation parameters are stored. Those later get loaded into
// an Operation's own Data() instance.
type ExecutionNode struct {
	Operation contract.Operation
	OpRefID   string
	Data      contract.OperationData
}

// NewExecutionNode creates a new ExecutionNode instance.
func NewExecutionNode(op contract.Operation, opRefID string, log contract.Logger) *ExecutionNode {
	n := &ExecutionNode{
		Operation: op,
		OpRefID:   opRefID,
	}

	n.Data.URL = params.URL("", log)
	n.Data.Query = params.Query(log)
	n.Data.Headers = params.Headers(log)
	n.Data.Body = params.Body(log)

	return n
}

// ID returns a uniqe operation node ID.
func (node *ExecutionNode) ID() gcontract.NodeID {
	return gcontract.NodeID(node.OpRefID)
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
