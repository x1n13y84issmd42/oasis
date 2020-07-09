package script

import (
	"fmt"
	"sync"

	"github.com/x1n13y84issmd42/gog/graph/comp"
	"github.com/x1n13y84issmd42/oasis/src/contract"
	"github.com/x1n13y84issmd42/oasis/src/test"
)

// Executor executes an ExecutionGraph that comes from a script.
type Executor struct {
	contract.EntityTrait
}

// NewExecutor creates a new Executor instance.
func NewExecutor(log contract.Logger) *Executor {
	return &Executor{
		EntityTrait: contract.Entity(log),
	}
}

// Execute executes.
func (ex Executor) Execute(graph *ExecutionGraph) {
	n0 := comp.MotherNode(graph)
	fmt.Printf("Execution starts from the node '%s'\n", n0.ID())
	wg := sync.WaitGroup{}
	results := contract.OperationResults{}
	ex.Walk(graph, n0.(*ExecutionNode), &wg, &results)
	wg.Wait()
}

// Walk ...
func (ex Executor) Walk(
	graph *ExecutionGraph,
	n *ExecutionNode,
	nwg *sync.WaitGroup,
	nresults *contract.OperationResults,
) {
	nwg.Add(1)
	anwg := sync.WaitGroup{}
	anresults := contract.OperationResults{}

	for _an := range graph.AdjacentNodes(n.ID()).Range() {
		an := _an.(*ExecutionNode)
		go ex.Walk(graph, an, &anwg, &anresults)
	}

	anwg.Wait()

	//TODO: use results as parameter source for n.Operation

	*nresults = append(*nresults, test.Operation(n.Operation, nil, nil, ex.Log))
	nwg.Done()
}
