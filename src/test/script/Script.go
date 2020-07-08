package script

import (
	"github.com/x1n13y84issmd42/gog/graph/comp"
	gcontract "github.com/x1n13y84issmd42/gog/graph/contract"
	"github.com/x1n13y84issmd42/oasis/src/contract"
	"github.com/x1n13y84issmd42/oasis/src/errors"
	"github.com/x1n13y84issmd42/oasis/src/params"
	"github.com/x1n13y84issmd42/oasis/src/utility"
)

// OperationRef ...
type OperationRef struct {
	OperationID string           `yaml:"operationId"`
	Use         OperationDataUse `yaml:"use"`
	Expect      OperationDataUse `yaml:"expect"`
}

// OperationDataMap ...
type OperationDataMap map[string]string

// Iterate creates an iterable channel.
func (m OperationDataMap) Iterate() contract.ParameterIterator {
	ch := make(contract.ParameterIterator)

	go func() {
		for pN, pV := range m {
			// fmt.Printf("Params %s = %s\n", pN, pV)
			ch <- contract.ParameterTuple{
				N: pN,
				Parameter: contract.Parameter{
					V:      params.Value(pV),
					Source: "Operation",
				},
			}
		}

		close(ch)
	}()

	return ch
}

// OperationDataUse ...
type OperationDataUse struct {
	Path    OperationDataMap
	Query   OperationDataMap
	Headers OperationDataMap
}

// Script ...
type Script struct {
	contract.EntityTrait
	Spec       string                   `yaml:"spec"`
	Operations map[string]*OperationRef `yaml:"operations"`
}

// GetExecutionGraph builds and returns an operation execution graph.
func (script *Script) GetExecutionGraph() gcontract.Graph {
	spec := utility.Load(script.Spec, script.Log)
	graph := NewExecutionGraph(script.Log)

	for opID, opRef := range script.Operations {
		//TODO: opRef.OperationID may be absent, use the key as op ID then.
		specOp := spec.GetOperation(opRef.OperationID)
		specOp.Data().URL.Load(opRef.Use.Path)

		var err error

		err = script.SetupDependencies(spec, graph, &opRef.Use.Path, specOp.Data().URL, specOp, opID)
		if err != nil {
			return NoGraph(err, script.Log)
		}

		err = script.SetupDependencies(spec, graph, &opRef.Use.Query, specOp.Data().Query, specOp, opID)
		if err != nil {
			return NoGraph(err, script.Log)
		}

		err = script.SetupDependencies(spec, graph, &opRef.Use.Headers, specOp.Data().Headers, specOp, opID)
		if err != nil {
			return NoGraph(err, script.Log)
		}
	}

	// Checking for cycles.
	cycle := comp.Cycle(graph)
	if len(*cycle) > 0 {
		return NoGraph(errors.GraphHasCycles(cycle, nil), script.Log)
	}

	return graph
}

// SetupDependencies iterates over the provided map, looks for reference values,
// collects a list of references operations along with ParameterAccess functions.
func (script *Script) SetupDependencies(
	spec contract.Spec,
	graph *ExecutionGraph,
	srcParams *OperationDataMap,
	dstParams contract.Set,
	specOp contract.Operation,
	opID string,
) error {
	refParams := params.NewReferenceSource("op reference")
	memParams := params.NewMemorySource("op reference")

	for pn, pv := range *srcParams {
		isref, op2Ref, selector := Dereference(pv)
		if isref {
			// fmt.Printf("%s is a ref to %s\n", pn, op2Ref)
			// Retrieving the referenced operation.
			op2 := script.Operations[op2Ref]
			if op2 == nil {
				return errors.NotFound("operation reference", op2Ref, nil)
			}

			specOp2 := spec.GetOperation(op2.OperationID)
			if op2 == nil {
				return errors.NotFound("spec operation", op2.OperationID, nil)
			}

			refParams.AddReference(pn, specOp2, selector)

			// Adding an edge to the execution graph.
			graph.AddEdge(&ExecutionNode{
				Operation: specOp,
				OpID:      opID,
			}, &ExecutionNode{
				Operation: specOp2,
				OpID:      op2Ref,
			})
		} else {
			// fmt.Printf("%s is a value\n", pn)
			memParams.Add(pn, pv)
		}
	}

	dstParams.Load(refParams)
	dstParams.Load(memParams)

	return nil
}
