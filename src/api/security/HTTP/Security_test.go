package http_test

import (
	gohttp "net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	http "github.com/x1n13y84issmd42/oasis/src/api/security/HTTP"
	"github.com/x1n13y84issmd42/oasis/src/log"
)

func Test_Basic(T *testing.T) {
	log := log.NewPlain(0)
	token := "00112233"
	sec := &http.Basic{
		Security: http.Security{
			Name:  "test sec",
			Token: token,
			Log:   log,
		},
	}

	req, _ := gohttp.NewRequest("GET", "example.com", nil)

	sec.Enrich(req, log)

	assert.Equal(T, token, req.Header.Get("Authorization"))
}

func Test_Digest(T *testing.T) {
	log := log.NewPlain(0)
	token := "00112233"
	sec := &http.Digest{
		Security: http.Security{
			Name:  "test sec",
			Token: token,
			Log:   log,
		},
	}

	req, _ := gohttp.NewRequest("GET", "example.com", nil)

	sec.Enrich(req, log)

	assert.Equal(T, token, req.Header.Get("Authorization"))
}
