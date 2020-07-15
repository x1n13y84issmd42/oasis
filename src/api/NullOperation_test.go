package api_test

import (
	"errors"
	"testing"

	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/log"
)

func Test_NullOperation(T *testing.T) {
	op := api.NoOperation(errors.New("yolo"), log.NewPlain(0))

	recovery := func() {
		if r := recover(); r != nil {
			//
		} else {
			T.Error("Should have recovered.")
		}
	}

	T.Run("GetRequest", func(T *testing.T) {
		defer recovery()
		op.GetRequest()
		T.Error("Should have panicked.")
	})

	T.Run("ID", func(T *testing.T) {
		defer recovery()
		op.ID()
		T.Error("Should have panicked.")
	})

	T.Run("Name", func(T *testing.T) {
		defer recovery()
		op.Name()
		T.Error("Should have panicked.")
	})

	T.Run("Description", func(T *testing.T) {
		defer recovery()
		op.Description()
		T.Error("Should have panicked.")
	})

	T.Run("Method", func(T *testing.T) {
		defer recovery()
		op.Method()
		T.Error("Should have panicked.")
	})

	T.Run("Path", func(T *testing.T) {
		defer recovery()
		op.Path()
		T.Error("Should have panicked.")
	})

	T.Run("Data", func(T *testing.T) {
		defer recovery()
		op.Data()
		T.Error("Should have panicked.")
	})

	T.Run("Resolve", func(T *testing.T) {
		defer recovery()
		op.Resolve()
		T.Error("Should have panicked.")
	})

	T.Run("Result", func(T *testing.T) {
		defer recovery()
		op.Result()
		T.Error("Should have panicked.")
	})
}
