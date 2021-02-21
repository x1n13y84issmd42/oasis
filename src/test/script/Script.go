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
	Security string           `yaml:"security"`
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
	SpecPaths  map[string]string                   `yaml:"specs"`
	Securities map[string]*contract.ScriptSecurity `yaml:"security"`
	Operations map[string]*OperationRef            `yaml:"operations"`

	Sec map[string]*contract.SecurityAccess `yaml:-`
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

		err = script.SetupDataDependency(graph, &opRef.Use.Path, opNode.Data.URL, opNode, opRef, opRefID)
		if err != nil {
			return NoGraph(err, script.Log)
		}

		err = script.SetupDataDependency(graph, &opRef.Use.Query, opNode.Data.Query, opNode, opRef, opRefID)
		if err != nil {
			return NoGraph(err, script.Log)
		}

		err = script.SetupDataDependency(graph, &opRef.Use.Headers, opNode.Data.Headers, opNode, opRef, opRefID)
		if err != nil {
			return NoGraph(err, script.Log)
		}

		err = script.SetupDataDependency(graph, &opRef.Use.Body, opNode.Data.Body, opNode, opRef, opRefID)
		if err != nil {
			return NoGraph(err, script.Log)
		}

		err = script.SetupDataDependency(graph, &opRef.Expect.Body, opNode.ExpectBody, opNode, opRef, opRefID)
		if err != nil {
			return NoGraph(err, script.Log)
		}

		err = script.SetupAfterDependency(graph, opRef, opNode)
		if err != nil {
			return NoGraph(err, script.Log)
		}

		err = script.SetupSecurityDependency(graph, opRef, opNode)
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

// SetupDependency adds an edge to the execution graph between two spec nodes.
func (script *Script) SetupDependency(
	scriptNodeName string,
	graph *ExecutionGraph,
	opRef *OperationRef,
	opNode *ExecutionNode,
) (contract.Operation, error) {
	opRef2 := script.Operations[scriptNodeName]
	if opRef2 == nil {
		return nil, errors.NotFound("Operation reference", scriptNodeName, nil)
	}

	op2 := script.GetOperation(opRef2.OperationID)
	if opRef2 == nil {
		return nil, errors.NotFound("Spec operation", opRef2.OperationID, nil)
	}

	// Adding an edge to the execution graph.
	graph.AddEdge(opNode.ID(), script.GetNode(graph, scriptNodeName, op2, opRef2).ID())

	return op2, nil
}

// SetupSecurityDependency adds an edge to the execution graph if opRef has an 'after' specified.
func (script *Script) SetupSecurityDependency(graph *ExecutionGraph, opRef *OperationRef, opNode *ExecutionNode) error {
	refdep := func(p *contract.ParameterAccess, v string) error {
		isref, op2RefID, selector := Dereference(v)
		if isref {
			op2, err := script.SetupDependency(op2RefID, graph, opRef, opNode)

			if err != nil {
				return err
			}
			// script.Log.NOMESSAGE("SetSecDep op '%s' depends on '%s' via security.", opNode.ID(), op2RefID)

			(*p) = (params.Reference{
				OpID:     op2.ID(),
				Result:   op2.Result(),
				Selector: selector,
				Log:      script.Log,
			}).Value()
		} else {
			(*p) = params.Value(v)
		}

		return nil
	}

	var err error

	for secName, sec := range script.Securities {
		opNodeSec := opNode.Operation.Resolve().Security(opRef.Use.Security)
		// script.Log.NOMESSAGE("SetSecDep.opRef.Use.Security: '%s'", opRef.Use.Security)
		// script.Log.NOMESSAGE("SetSecDep.opNodeSec.GetName: '%s'", opNodeSec.GetName())

		if opNodeSec.GetName() == secName {
			script.Sec[secName] = &contract.SecurityAccess{}

			err = refdep(&script.Sec[secName].Value, sec.Value)
			if err != nil {
				return err
			}

			err = refdep(&script.Sec[secName].Token, sec.Token)
			if err != nil {
				return err
			}

			err = refdep(&script.Sec[secName].Username, sec.Username)
			if err != nil {
				return err
			}

			err = refdep(&script.Sec[secName].Password, sec.Password)
			if err != nil {
				return err
			}

			// script.Log.NOMESSAGE("SetSecDep.script.Sec[secName]: %#v", script.Sec[secName])
		}

	}

	return nil
}

// SetupAfterDependency adds an edge to the execution graph if opRef has an 'after' specified.
func (script *Script) SetupAfterDependency(graph *ExecutionGraph, opRef *OperationRef, opNode *ExecutionNode) error {
	if opRef.After != "" {
		_, err := script.SetupDependency(opRef.After, graph, opRef, opNode)
		return err
	}

	return nil
}

// SetupDataDependency iterates over the provided map, looks for reference values,
// collects a list of references operations, and adds edges b/w them & opNode.
func (script *Script) SetupDataDependency(
	graph *ExecutionGraph,
	srcParams *OperationDataMap,
	dstParams contract.Set,
	opNode *ExecutionNode,
	opRef *OperationRef,
	opRefID string,
) error {
	refParams := params.NewReferenceSource(script.Log)
	memParams := params.NewMemorySource("script data")

	for pn, pv := range *srcParams {
		isref, op2RefID, selector := Dereference(pv)
		if isref {
			op2, err := script.SetupDependency(op2RefID, graph, opRef, opNode)
			if err != nil {
				return err
			}

			// Adding the value so it's available for op later.
			refParams.AddReference(pn, op2.ID()+" node", op2.Result(), selector)
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

// GetSecurity returns the script's security parameters.
func (script *Script) GetSecurity(name string) *contract.SecurityAccess {
	return script.Sec[name]
}
