package apikey

import (
	"fmt"

	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/contract"
	"github.com/x1n13y84issmd42/oasis/src/errors"
	"github.com/x1n13y84issmd42/oasis/src/params"
)

// Security implements the security type that requires a unique API key
// to be present in each API request. The exact location of the key
// is implemented corresponding classes.
type Security struct {
	Name      string
	ParamName string
	Value     contract.ParameterAccess
	Log       contract.Logger
}

// New creates a new API Key security scheme.
func New(name string, location string, paramName string, value string, logger contract.Logger) contract.Security {
	switch location {
	case "cookie":
		return &Cookie{
			Security{name, paramName, params.Value(value), logger},
		}

	case "header":
		return &Header{
			Security{name, paramName, params.Value(value), logger},
		}

	case "query":
		return &Query{
			Security{name, paramName, params.Value(value), logger},
		}
	}

	//TODO: return error
	return api.NoSecurity(errors.Oops(fmt.Sprintf("Unknown location \"%s\" for the API Key security parameter \"%s\".", location, paramName), nil), logger)
}

// GetName returns name.
func (sec Security) GetName() string {
	return sec.Name
}

// SetValue sets Value.
func (sec *Security) SetValue(v contract.ParameterAccess) {
	sec.Value = v
}

// SetToken does nothig.
func (sec *Security) SetToken(v contract.ParameterAccess) {
}

// SetUsername does nothig.
func (sec *Security) SetUsername(v contract.ParameterAccess) {
}

// SetPassword does nothig.
func (sec *Security) SetPassword(v contract.ParameterAccess) {
}
