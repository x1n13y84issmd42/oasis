package script

import (
	"sync"

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
	Mutex     sync.Mutex
	Result    *contract.OperationResult
	Use       *OperationDataUse
	Expect    *OperationDataUse
}

// NewExecutionNode creates a new ExecutionNode instance.
func NewExecutionNode(op contract.Operation, opRefID string, opRef *OperationRef, log contract.Logger) *ExecutionNode {
	n := &ExecutionNode{
		Operation: op,
		OpRefID:   opRefID,
	}

	n.Mutex = sync.Mutex{}

	n.Data.URL = params.URL("", log)
	n.Data.Query = params.Query(log)
	n.Data.Headers = params.Headers(log)
	n.Data.Body = params.Body(log)

	n.Use = &opRef.Use
	n.Expect = &opRef.Expect

	return n
}

// ID returns a uniqe operation node ID.
func (node *ExecutionNode) ID() gcontract.NodeID {
	return gcontract.NodeID(node.OpRefID)
}

// Lock locks the node to prevent parallel executions.
func (node *ExecutionNode) Lock() {
	node.Mutex.Lock()
}

// Unlock unlocks the node.
func (node *ExecutionNode) Unlock() {
	node.Mutex.Unlock()
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
