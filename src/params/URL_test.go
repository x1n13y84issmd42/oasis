package params_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/x1n13y84issmd42/oasis/src/log"
	"github.com/x1n13y84issmd42/oasis/src/params"
)

func Test_URL(T *testing.T) {
	T.Run("KeyHost Validation Failure", func(T *testing.T) {
		src := params.NewMemorySource()
		src.Add("bar", "B4R")
		src.Add("yolo", "y010")

		url := params.URL("/foo/{bar}/qeq/{yolo}", log.New("plain", 0))
		url.Load(src)

		defer func() {
			if r := recover(); r != nil {
			}
		}()

		_ = url.String()
		T.Error("The URL.String() didn't panic when it should have.")
	})

	T.Run("Success", func(T *testing.T) {
		src := params.NewMemorySource()
		src.Add("bar", "B4R")
		src.Add("yolo", "Y010")
		src.Add(params.KeyHost, "example.com")

		url := params.URL("/foo/{bar}/qeq/{yolo}", log.New("plain", 0))
		url.Load(src)

		defer func() {
			if r := recover(); r != nil {
				T.Error("The URL.String() panicked when it shouldn't have.")
			}
		}()

		expected := "example.com/foo/B4R/qeq/Y010"
		assert.Equal(T, expected, url.String())
	})

}
