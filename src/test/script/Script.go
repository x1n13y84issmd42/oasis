package script

import (
	"fmt"

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
			fmt.Printf("Params %s = %s\n", pN, pV)
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
	graph := NewOperationGraph(script.Log)

	for _, opRef := range script.Operations {
		specOp := spec.GetOperation(opRef.OperationID)
		specOp.Data().URL.Load(opRef.Use.Path)

		var err error

		err = script.SetupDependencies(spec, graph, &opRef.Use.Path, specOp.Data().URL, specOp)
		if err != nil {
			return NoGraph(err, script.Log)
		}

		err = script.SetupDependencies(spec, graph, &opRef.Use.Query, specOp.Data().Query, specOp)
		if err != nil {
			return NoGraph(err, script.Log)
		}

		err = script.SetupDependencies(spec, graph, &opRef.Use.Headers, specOp.Data().Headers, specOp)
		if err != nil {
			return NoGraph(err, script.Log)
		}
	}

	return graph
}

// SetupDependencies iterates over the provided map, looks for reference values,
// collects a list of references operations along with ParameterAccess functions.
func (script *Script) SetupDependencies(
	spec contract.Spec,
	graph *OperationGraph,
	srcParams *OperationDataMap,
	dstParams contract.Set,
	specOp contract.Operation,
) error {
	refParams := params.NewReferenceSource("op reference")
	memParams := params.NewMemorySource("op reference")

	for pn, pv := range *srcParams {
		isref, op2Ref, selector := Dereference(pv)
		if isref {
			fmt.Printf("%s is a ref to %s\n", pn, op2Ref)
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
			graph.AddEdge(&OperationNode{
				Operation: specOp,
			}, &OperationNode{
				Operation: specOp2,
			})
		} else {
			fmt.Printf("%s is a value\n", pn)
			memParams.Add(pn, pv)
		}
	}

	dstParams.Load(refParams)
	dstParams.Load(memParams)

	return nil
}
