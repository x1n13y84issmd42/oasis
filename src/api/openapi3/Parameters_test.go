package openapi3_test

import (
	"testing"

	kinopenapi3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
	"github.com/x1n13y84issmd42/oasis/src/api/openapi3"
	"github.com/x1n13y84issmd42/oasis/src/contract"
)

func Test_Parameters(T *testing.T) {
	params := kinopenapi3.Parameters{
		&kinopenapi3.ParameterRef{
			Value: &kinopenapi3.Parameter{
				In:      "query",
				Name:    "param2",
				Example: "P2_VALUE",
			},
		},

		&kinopenapi3.ParameterRef{
			Value: &kinopenapi3.Parameter{
				In:       "path",
				Name:     "param3",
				Required: true,
				Example:  "P3_VALUE",
			},
		},

		&kinopenapi3.ParameterRef{
			Value: &kinopenapi3.Parameter{
				In:       "path",
				Name:     "param1",
				Required: true,
				Example:  "P1_VALUE",
			},
		},

		&kinopenapi3.ParameterRef{
			Value: &kinopenapi3.Parameter{
				In:      "query",
				Name:    "qp",
				Example: "QP_VALUE",
			},
		},

		&kinopenapi3.ParameterRef{
			Value: &kinopenapi3.Parameter{
				In:       "header",
				Name:     "param3",
				Required: true,
				Example:  "P3_VALUE",
			},
		},

		&kinopenapi3.ParameterRef{
			Value: &kinopenapi3.Parameter{
				In:       "header",
				Name:     "abra",
				Required: true,
				Example:  "CADABRA",
			},
		},
	}

	iterate := func(src contract.ParameterSource) string {
		res := ""
		for p := range src.Iterate() {
			res = res + p.N + ":" + p.V + " "
		}
		return res
	}

	T.Run("PathParameterSource", func(T *testing.T) {
		src := openapi3.PathParameterSource(&params)
		expected := "param1:P1_VALUE param3:P3_VALUE "
		assert.Equal(T, expected, iterate(src))
		assert.Equal(T, "P1_VALUE", src.Get("param1"))
		assert.Equal(T, "", src.Get("NONEXISTENT_param11111111"))
	})

	T.Run("QueryParameterSource", func(T *testing.T) {
		src := openapi3.QueryParameterSource(&params)
		expected := "param2:P2_VALUE qp:QP_VALUE "
		assert.Equal(T, expected, iterate(src))
		assert.Equal(T, "QP_VALUE", src.Get("qp"))
	})

	T.Run("HeaderParameterSource", func(T *testing.T) {
		src := openapi3.HeadersParameterSource(&params)
		expected := "abra:CADABRA param3:P3_VALUE "
		assert.Equal(T, expected, iterate(src))
		assert.Equal(T, "CADABRA", src.Get("abra"))
	})
}
