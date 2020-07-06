package script

import (
	gcontract "github.com/x1n13y84issmd42/gog/graph/contract"
	"github.com/x1n13y84issmd42/oasis/src/contract"
)

// NullScript is used whenever we can't have a real one.
// Reports the contained error on every method call.
type NullScript struct {
	contract.NullObjectPrototype
}

// NoScript creates a new NullScript instance.
func NoScript(err error, log contract.Logger) contract.Script {
	return &NullScript{
		NullObjectPrototype: contract.NullObject(err, log),
	}
}

// GetGraph reports an error.
func (s *NullScript) GetGraph() gcontract.Graph {
	s.Report()
	return nil
}
