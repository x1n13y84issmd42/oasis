package api_test

import (
	"errors"
	"testing"

	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/log"
)

func Test_NullSpec(T *testing.T) {
	spec := api.NoSpec(errors.New("yolo"), log.NewPlain(0))

	recovery := func() {
		if r := recover(); r != nil {
			//
		} else {
			T.Error("Should have recovered.")
		}
	}

	T.Run("Operations", func(T *testing.T) {
		defer recovery()
		spec.Operations()
		T.Error("Should have panicked.")
	})

	T.Run("GetOperation", func(T *testing.T) {
		defer recovery()
		spec.GetOperation("")
		T.Error("Should have panicked.")
	})

	T.Run("Title", func(T *testing.T) {
		defer recovery()
		spec.Title()
		T.Error("Should have panicked.")
	})

	T.Run("Description", func(T *testing.T) {
		defer recovery()
		spec.Description()
		T.Error("Should have panicked.")
	})

	T.Run("Version", func(T *testing.T) {
		defer recovery()
		spec.Version()
		T.Error("Should have panicked.")
	})
}
