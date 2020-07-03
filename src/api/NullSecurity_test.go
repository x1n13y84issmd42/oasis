package api_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/log"
)

func Test_NullSecurity(T *testing.T) {
	log := log.NewPlain(0)
	sec := api.NoSecurity(fmt.Errorf("alas"), log)

	recovery := func() {
		if r := recover(); r != nil {
			//
		} else {
			T.Error("Should have recovered.")
		}
	}

	T.Run("GetName", func(T *testing.T) {
		defer recovery()
		sec.GetName()
		T.Error("Should have panicked")
	})

	T.Run("Enrich", func(T *testing.T) {
		defer recovery()
		req, _ := http.NewRequest("GET", "example.com", nil)
		sec.Enrich(req, log)
		T.Error("Should have panicked")
	})

}
