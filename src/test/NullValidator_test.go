package test_test

import (
	"errors"
	"testing"

	"github.com/x1n13y84issmd42/oasis/src/log"
	"github.com/x1n13y84issmd42/oasis/src/test"
)

func Test_NullValidator(T *testing.T) {
	v := test.NoValidator(errors.New("yolo"), log.NewPlain(0))

	recovery := func() {
		if r := recover(); r != nil {
			//
		} else {
			T.Error("Should have recovered.")
		}
	}

	T.Run("Expect", func(T *testing.T) {
		defer recovery()
		v.Expect(nil)
		T.Error("Should have panicked.")
	})

	T.Run("Validate", func(T *testing.T) {
		defer recovery()
		v.Validate(nil)
		T.Error("Should have panicked.")
	})
}
