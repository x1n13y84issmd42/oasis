package script

import (
	gcontract "github.com/x1n13y84issmd42/gog/graph/contract"
	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/contract"
	"github.com/x1n13y84issmd42/oasis/src/errors"
	"github.com/x1n13y84issmd42/oasis/src/params"
)

// OperationRef is a node of execution graph as desfined in the script file.
// It references a spec operation and contains the needed data.
type OperationRef struct {
	OperationID string           `yaml:"operationId"`
	After       string           `yaml:"after"`
	Use         OperationDataUse `yaml:"use"`
	Expect      OperationDataUse `yaml:"expect"`
}

// OperationDataMap is a map of parameters for an OperationRef.
type OperationDataMap map[string]string

// Iterate creates an iterable channel.
func (m OperationDataMap) Iterate() contract.ParameterIterator {
	ch := make(contract.ParameterIterator)

	go func() {
		for pN, pV := range m {
			ch <- contract.ParameterTuple{
				N: pN,
				Parameter: contract.Parameter{
					V:      params.Value(pV),
					Source: "script OperationRef",
				},
			}
		}

		close(ch)
	}()

	return ch
}

// OperationDataUse corresponds to the 'use' block of the OperationRef in a script file.
type OperationDataUse struct {
	Path     OperationDataMap `yaml:"path"`
	Body     OperationDataMap `yaml:"body"`
	Query    OperationDataMap `yaml:"query"`
	Headers  OperationDataMap `yaml:"headers"`
	Security OperationDataMap `yaml:"security"`
	CT       string           `yaml:"CT"`
	Status   int64            `yaml:"status"`
}

// OperationDataExpect corresponds to the 'expect' block of the OperationRef in a script file.
type OperationDataExpect struct {
	Body    OperationDataMap `yaml:"body"`
	Headers OperationDataMap `yaml:"headers"`
	CT      string           `yaml:"CT"`
	Status  int64            `yaml:"status"`
}

// Script is a complex API testing scenario.
// It defines dependencies between various operations
// and order of their execution.
type Script struct {
	api.OperationCache
	contract.EntityTrait
	SpecPaths  map[string]string        `yaml:"specs"`
	Operations map[string]*OperationRef `yaml:"operations"`
}

// GetExecutionGraph builds and returns an operation execution graph.
func (script *Script) GetExecutionGraph() gcontract.Graph {
	if len(script.Operations) == 0 {
		return NoGraph(errors.Oops("The execution graph contains no nodes. Check your script file syntax, YAML format is very sensitive to errors.", nil), script.Log)
	}

	graph := NewExecutionGraph(script.Log)

	for opRefID, opRef := range script.Operations {
		//TODO: opRef.OperationID may be absent, use opRefID then.
		op := script.GetOperation(opRef.OperationID)
		opNode := script.GetNode(graph, opRefID, op, opRef)

		var err error

		err = script.SetupDependencies(graph, &opRef.Use.Path, opNode.Data.URL, opNode, opRefID)
		if err != nil {
			return NoGraph(err, script.Log)
		}

		err = script.SetupDependencies(graph, &opRef.Use.Query, opNode.Data.Query, opNode, opRefID)
		if err != nil {
			return NoGraph(err, script.Log)
		}

		err = script.SetupDependencies(graph, &opRef.Use.Headers, opNode.Data.Headers, opNode, opRefID)
		if err != nil {
			return NoGraph(err, script.Log)
		}

		err = script.SetupDependencies(graph, &opRef.Use.Body, opNode.Data.Body, opNode, opRefID)
		if err != nil {
			return NoGraph(err, script.Log)
		}

		err = script.SetupDependencies(graph, &opRef.Expect.Body, opNode.ExpectBody, opNode, opRefID)
		if err != nil {
			return NoGraph(err, script.Log)
		}

		err = script.SetupAfterDependency(graph, opRef, opNode)
		if err != nil {
			return NoGraph(err, script.Log)
		}
	}

	// Checking for cycles.
	// cycle := comp.Cycle(graph)
	// if len(*cycle) > 0 {
	// 	return NoGraph(errors.GraphHasCycles(cycle, nil), script.Log)
	// }

	return graph
}

// SetupAfterDependency adds an edge to the execution graph if opRef has an 'after' specified.
func (script *Script) SetupAfterDependency(graph *ExecutionGraph, opRef *OperationRef, opNode *ExecutionNode) error {
	if opRef.After != "" {
		opRef2 := script.Operations[opRef.After]
		if opRef2 == nil {
			return errors.NotFound("Operation reference", opRef.After, nil)
		}

		op2 := script.GetOperation(opRef2.OperationID)
		if opRef2 == nil {
			return errors.NotFound("Spec operation", opRef2.OperationID, nil)
		}

		// Adding an edge to the execution graph.
		graph.AddEdge(opNode.ID(), script.GetNode(graph, opRef.After, op2, opRef2).ID())
	}

	return nil
}

// SetupDependencies iterates over the provided map, looks for reference values,
// collects a list of references operations along with ParameterAccess functions.
func (script *Script) SetupDependencies(
	graph *ExecutionGraph,
	srcParams *OperationDataMap,
	dstParams contract.Set,
	opNode *ExecutionNode,
	opRefID string,
) error {
	refParams := params.NewReferenceSource(script.Log)
	memParams := params.NewMemorySource("script data")

	for pn, pv := range *srcParams {
		isref, op2RefID, selector := Dereference(pv)
		if isref {
			// Retrieving the referenced operation.
			opRef2 := script.Operations[op2RefID]
			if opRef2 == nil {
				return errors.NotFound("Operation reference", op2RefID, nil)
			}

			op2 := script.GetOperation(opRef2.OperationID)
			if opRef2 == nil {
				return errors.NotFound("Spec operation", opRef2.OperationID, nil)
			}

			// Adding the value so it's available for op later.
			refParams.AddReference(pn, op2.ID()+" node", op2.Result(), selector)

			// Adding an edge to the execution graph.
			graph.AddEdge(opNode.ID(), script.GetNode(graph, op2RefID, op2, opRef2).ID())
		} else {
			memParams.Add(pn, pv)
		}
	}

	dstParams.Load(refParams)
	dstParams.Load(memParams)

	return nil
}

// GetNode returns an ExecutionNode instance corresponding to the opRefID.
// If such a node exists in the graph, it will be returned, otherwise a new
// node is created.
func (script *Script) GetNode(graph gcontract.Graph, opRefID string, op contract.Operation, opRef *OperationRef) *ExecutionNode {
	var opNode *ExecutionNode
	_opNode := graph.Node(gcontract.NodeID(opRefID))
	if _opNode != nil {
		opNode = _opNode.(*ExecutionNode)
	} else {
		opNode = NewExecutionNode(op, opRefID, opRef, script.Log)
		graph.AddNode(opNode)
	}

	return opNode
}
