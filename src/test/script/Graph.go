package script

import (
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
