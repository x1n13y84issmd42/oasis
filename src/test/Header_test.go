package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/log"
)

func TestHeader(T *testing.T) {
	//TODO: mock it & assert logging calls.
	logger := log.NewFestive(0)

	cases := []struct {
		HeaderName   string
		Header       *api.Header
		HeaderValues []string
		Expected     bool
	}{
		{
			HeaderName: "X-Index",
			Header: &api.Header{
				Name:        "X-Index",
				Description: "whatever",
				Required:    true,
			},
			HeaderValues: []string{},
			Expected:     false,
		},
		{
			HeaderName: "X-Index",
			Header: &api.Header{
				Name:        "X-Index",
				Description: "whatever",
				Required:    true,
			},
			HeaderValues: []string{
				"33",
				"44",
			},
			Expected: true,
		},
	}

	for _, c := range cases {
		actual := Header(c.HeaderName, c.Header, c.HeaderValues, logger)
		assert.Equal(T, c.Expected, actual)
	}
}
