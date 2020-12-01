package script

import (
	"sync"

	"github.com/x1n13y84issmd42/gog/graph/comp"
	gcontract "github.com/x1n13y84issmd42/gog/graph/contract"
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
func (ex Executor) Execute(graph gcontract.Graph) {
	n0 := comp.MotherNode(graph)
	// TODO: it returns nil sometimes (on script/noosa_test.yaml)
	ex.Log.ScriptExecutionStart(string(n0.ID()))

	results := contract.OperationResults{}

	wg := sync.WaitGroup{}
	wg.Add(1)
	ex.Walk(graph, n0.(*ExecutionNode), &wg, &results)
	wg.Wait()
}

// Walk walks the execution graph and executes operations.
func (ex Executor) Walk(
	graph gcontract.Graph,
	n *ExecutionNode,
	nwg *sync.WaitGroup,
	nresults *contract.OperationResults,
) {
	// Executing child nodes first (post-order).
	anwg := sync.WaitGroup{}
	anwg.Add(graph.AdjacentNodes(n.ID()).Count())
	anresults := contract.OperationResults{}

	for _an := range graph.AdjacentNodes(n.ID()).Range() {
		an := _an.(*ExecutionNode)
		go ex.Walk(graph, an, &anwg, &anresults)
	}

	anwg.Wait()

	//TODO: check for successful outcome of the previous ops.

	n.Lock()

	if n.Result == nil {
		// Executing the current node after it's children.
		n.Operation.Data().Load(&n.Data)
		n.Operation.Data().URL.Load(n.Operation.Resolve().Host(""))

		enrichment := []contract.RequestEnrichment{
			n.Operation.Data().Query,
			n.Operation.Data().Headers,
			n.Operation.Data().Body,

			n.Operation.Resolve().Security(""),
		}

		ex.Log.TestingOperation(n.Operation)

		v := n.Operation.Resolve().Response(0, "")

		n.Result = test.Operation(n.Operation, &enrichment, v, ex.Log)
		*nresults = append(*nresults, n.Result)
	}

	n.Unlock()

	nwg.Done()
}
