package script

import (
	"sync"

	gcontract "github.com/x1n13y84issmd42/gog/graph/contract"
	"github.com/x1n13y84issmd42/oasis/src/contract"
	"github.com/x1n13y84issmd42/oasis/src/test"
	"github.com/x1n13y84issmd42/oasis/src/test/expect"
)

// Executor executes an ExecutionGraph that comes from a script.
type Executor struct {
	contract.EntityTrait

	Script contract.Script
}

// NewExecutor creates a new Executor instance.
func NewExecutor(logger contract.Logger, script contract.Script) *Executor {
	return &Executor{
		EntityTrait: contract.Entity(logger),
		Script:      script,
	}
}

// Execute executes.
func (ex Executor) Execute(graph gcontract.Graph) {
	success := true
	results := make(contract.OperationResults)

	wg := sync.WaitGroup{}
	for node := range graph.Nodes().Range() {
		wg.Add(1)
		go ex.Walk(graph, node.(*ExecutionNode), &wg, &results)
	}

	wg.Wait()

	for nID, nRes := range results {
		if !nRes.Success {
			ex.Log.NOMESSAGE("\rOperation %s has failed.", nID)
			success = false
		}
	}

	if !success {
		// os.Exit(255)
	}
}

// Walk walks the execution graph and executes operations.
func (ex Executor) Walk(
	graph gcontract.Graph,
	n *ExecutionNode,
	nwg *sync.WaitGroup,
	nresults *contract.OperationResults,
) {

	// ex.Log.NOMESSAGE("Walking %s", n.ID())
	// Executing child nodes first (post-order).
	anwg := sync.WaitGroup{}
	anwg.Add(int(graph.AdjacentNodes(n.ID()).Count()))
	anresults := contract.OperationResults{}

	// TODO: consider moving execution of adjacent nodes
	// into the Reference.Value() function in truly lazy fashion.
	for _an := range graph.AdjacentNodes(n.ID()).Range() {
		// ex.Log.NOMESSAGE("Child node %s of %s", _an.ID(), n.ID())
		an := _an.(*ExecutionNode)
		go ex.Walk(graph, an, &anwg, &anresults)
	}

	anwg.Wait()

	//TODO: check for successful outcome of the previous ops.

	n.Lock()

	if n.Result == nil {
		logger := n.Operation.GetLogger()
		logger.Buffer(true)

		// Setting the request enrichment.
		n.Operation.Data().Reload()
		n.Operation.Data().Load(&n.Data)
		n.Operation.Data().URL.Load(n.Operation.Resolve().Host(""))

		opSecurity := n.Operation.Resolve().Security("")
		// ex.Log.NOMESSAGE("security.GetName() = %s", opSecurity.GetName())

		if scriptSec := ex.Script.GetSecurity(opSecurity.GetName()); scriptSec != nil {
			opSecurity.SetValue(scriptSec.Value)
			opSecurity.SetToken(scriptSec.Token)
			opSecurity.SetUsername(scriptSec.Username)
			opSecurity.SetPassword(scriptSec.Password)
		}

		enrichment := []contract.RequestEnrichment{
			n.Operation.Data().Query,
			n.Operation.Data().Headers,
			n.Operation.Data().Body,

			opSecurity,
		}

		logger.TestingOperation(n.Operation)

		// Setting the response validation.
		v := n.Operation.Resolve().Response(n.Expect.Status, "")
		// v.SetLogger(logger)
		v.Expect(expect.JSONBody(n.ExpectBody, graph, logger))

		n.Result = test.Operation(n.Operation, &enrichment, v, logger)
		(*nresults)[string(n.ID())] = n.Result

		logger.Flush()
		logger.Buffer(false)
	}

	n.Unlock()
	nwg.Done()
}
