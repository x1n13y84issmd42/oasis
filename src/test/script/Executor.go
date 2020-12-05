package script

import (
	"os"
	"sync"

	gcontract "github.com/x1n13y84issmd42/gog/graph/contract"
	"github.com/x1n13y84issmd42/oasis/src/contract"
	"github.com/x1n13y84issmd42/oasis/src/test"
	"github.com/x1n13y84issmd42/oasis/src/test/expect"
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
	success := true
	results := make(contract.OperationResults)

	wg := sync.WaitGroup{}
	for node := range graph.Nodes().Range() {
		wg.Add(1)
		ex.Walk(graph, node.(*ExecutionNode), &wg, &results)
	}

	wg.Wait()

	for nID, nRes := range results {
		if !nRes.Success {
			ex.Log.NOMESSAGE("\rOperation %s has failed.", nID)
			success = false
		}
	}

	if !success {
		os.Exit(255)
	}
}

// Walk walks the execution graph and executes operations.
func (ex Executor) Walk(
	graph gcontract.Graph,
	n *ExecutionNode,
	nwg *sync.WaitGroup,
	nresults *contract.OperationResults,
) {
	// logger := log.NewBuffer(ex.Log)
	logger := ex.Log

	// Executing child nodes first (post-order).
	anwg := sync.WaitGroup{}
	anwg.Add(int(graph.AdjacentNodes(n.ID()).Count()))
	anresults := contract.OperationResults{}

	for _an := range graph.AdjacentNodes(n.ID()).Range() {
		an := _an.(*ExecutionNode)
		go ex.Walk(graph, an, &anwg, &anresults)
	}

	anwg.Wait()

	//TODO: check for successful outcome of the previous ops.

	n.Lock()

	if n.Result == nil {
		// Setting the request enrichment.
		n.Operation.Data().Reload()
		n.Operation.Data().Load(&n.Data)
		n.Operation.Data().URL.Load(n.Operation.Resolve().Host(""))

		enrichment := []contract.RequestEnrichment{
			n.Operation.Data().Query,
			n.Operation.Data().Headers,
			n.Operation.Data().Body,

			n.Operation.Resolve().Security(""),
		}

		logger.TestingOperation(n.Operation)

		// Setting the response validation.
		v := n.Operation.Resolve().Response(n.Expect.Status, "")
		v.Expect(expect.JSONBody(n.ExpectBody, graph, logger))

		n.Result = test.Operation(n.Operation, &enrichment, v, logger)
		(*nresults)[string(n.ID())] = n.Result
	}

	n.Unlock()
	nwg.Done()
	logger.Flush()
}
