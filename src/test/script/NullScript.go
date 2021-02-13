package script

import (
	gcontract "github.com/x1n13y84issmd42/gog/graph/contract"
	"github.com/x1n13y84issmd42/oasis/src/contract"
	"github.com/x1n13y84issmd42/oasis/src/errors"
)

// NullScript is used whenever we can't have a real one.
// Reports the contained error on every method call.
type NullScript struct {
	errors.NullObjectPrototype
}

// NoScript creates a new NullScript instance.
func NoScript(err error, log contract.Logger) contract.Script {
	return &NullScript{
		NullObjectPrototype: errors.NullObject(err, log),
	}
}

// GetExecutionGraph reports an error.
func (s *NullScript) GetExecutionGraph() gcontract.Graph {
	s.Report()
	return nil
}

// GetSecurity reports an error.
func (s *NullScript) GetSecurity(name string) *contract.SecurityAccess {
	s.Report()
	return nil
}
