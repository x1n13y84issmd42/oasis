package contract_test

import (
	"errors"
	"testing"

	"github.com/x1n13y84issmd42/oasis/src/contract"
	"github.com/x1n13y84issmd42/oasis/src/log"
)

func Test_Entity(T *testing.T) {
	e := contract.Entity(log.NewPlain(0))

	defer func() {
		if r := recover(); r != nil {
			//
		} else {
			T.Error("Should have recovered.")
		}
	}()

	e.Error(errors.New("test"))
	T.Error("Should have panicked.")
}
