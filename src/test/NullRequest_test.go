package test_test

import (
	"errors"
	"testing"

	"github.com/x1n13y84issmd42/oasis/src/log"
	"github.com/x1n13y84issmd42/oasis/src/test"
)

func Test_NullRequest(T *testing.T) {
	req := test.NoRequest(errors.New("yolo"), log.NewPlain(0))

	recovery := func() {
		if r := recover(); r != nil {
			//
		} else {
			T.Error("Should have recovered.")
		}
	}

	T.Run("Enrich", func(T *testing.T) {
		defer recovery()
		req.Enrich(nil)
		T.Error("Should have panicked.")
	})

	T.Run("Execute", func(T *testing.T) {
		defer recovery()
		req.Execute()
		T.Error("Should have panicked.")
	})
}
