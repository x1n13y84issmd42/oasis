package params_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/x1n13y84issmd42/oasis/src/log"
	"github.com/x1n13y84issmd42/oasis/src/params"
)

func Test_Query(T *testing.T) {
	src1 := params.NewMemorySource("test")
	src1.Add("bar", "B4R")
	src1.Add("yolo", "Y010")

	src2 := params.NewMemorySource("test")
	src2.Add("bar", "B4444R")
	src2.Add("yolo", "you only live once")

	qry := params.Query(log.New("plain", 0))
	qry.Load(src1)
	qry.Load(src2)

	defer func() {
		if r := recover(); r != nil {
			T.Error("The Query.Enrich() panicked when it shouldn't have.")
		}
	}()

	expected := url.Values{
		"bar": []string{
			"B4R",
			"B4444R",
		},

		"yolo": []string{
			"Y010",
			"you only live once",
		},
	}

	req, _ := http.NewRequest("GET", "example.com", nil)

	qry.Enrich(req, log.New("plain", 0))

	assert.Equal(T, expected, req.URL.Query())
}
