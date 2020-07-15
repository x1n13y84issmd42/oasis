package errors_test

import (
	goerrors "errors"
	"testing"

	"github.com/x1n13y84issmd42/oasis/src/errors"
	"github.com/x1n13y84issmd42/oasis/src/log"
)

func Test_NullObjectPrototype(T *testing.T) {

	recovery := func() {
		if r := recover(); r != nil {
			//
		} else {
			T.Error("Should have recovered.")
		}
	}

	T.Run("Report Error", func(T *testing.T) {
		o := errors.NullObject(goerrors.New("yolo"), log.NewPlain(0))
		defer recovery()
		o.Report()
		T.Error("Should have panicked.")
	})

	T.Run("Report NoError", func(T *testing.T) {
		o := errors.NullObject(nil, log.NewPlain(0))
		defer recovery()
		o.Report()
		T.Error("Should have panicked.")
	})
}
