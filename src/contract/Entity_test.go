package contract_test

import (
	"testing"

	"github.com/x1n13y84issmd42/oasis/src/contract"
	"github.com/x1n13y84issmd42/oasis/src/log"
)

// This exists just to get coverage.
func Test_Entity(T *testing.T) {
	e := contract.Entity(log.NewPlain(0))
	e.GetLogger().NOMESSAGE("")
}
