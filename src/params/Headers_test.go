package params_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/x1n13y84issmd42/oasis/src/log"
	"github.com/x1n13y84issmd42/oasis/src/params"
)

func Test_Headers(T *testing.T) {
	src1 := params.NewMemorySource()
	src1.Add("X-Oasis-Stuff", "value 1")
	src1.Add("Connection", "close")

	src2 := params.NewMemorySource()
	src2.Add("X-Oasis-Stuff", "value 2")
	src2.Add("Accept", "image/png")

	qry := params.Headers(log.New("plain", 0))
	qry.Load(src1)
	qry.Load(src2)

	defer func() {
		if r := recover(); r != nil {
			T.Error("The Headers.Enrich() panicked when it shouldn't have.")
		}
	}()

	expected := http.Header{
		"Accept": []string{
			"image/png",
		},
		"Connection": []string{
			"close",
		},
		"X-Oasis-Stuff": []string{
			"value 1",
			"value 2",
		},
	}

	req, _ := http.NewRequest("GET", "example.com", nil)

	qry.Enrich(req, log.New("plain", 0))

	assert.Equal(T, expected, req.Header)
}
