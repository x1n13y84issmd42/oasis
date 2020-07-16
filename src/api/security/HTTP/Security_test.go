package http_test

import (
	gohttp "net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/x1n13y84issmd42/oasis/src/api"
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

func Test_New(T *testing.T) {
	token := "0011223344"

	T.Run("Basic", func(T *testing.T) {
		sec := http.New("test sec", "basic", token, "", "", log.NewPlain(0))
		tsec, ok := sec.(http.Basic)
		assert.True(T, ok)

		if ok {
			assert.Equal(T, token, tsec.Token)
		}
	})

	T.Run("Digest", func(T *testing.T) {
		sec := http.New("test sec", "digest", token, "", "", log.NewPlain(0))
		tsec, ok := sec.(http.Digest)
		assert.True(T, ok)

		if ok {
			assert.Equal(T, token, tsec.Token)
		}
	})

	T.Run("Invalid", func(T *testing.T) {
		sec := http.New("test sec", "INVALID SCHEMA", "", "", "", log.NewPlain(0))
		_, ok := sec.(*api.NullSecurity)
		assert.True(T, ok)
	})
}

func Test_ParseWWWAuthenticate(T *testing.T) {
	sec := http.Security{}
	expected := map[string]string{
		"Digest realm": "Oasis",
		"nonce":        "61b6948856629ad7fd3da9d6179393ec",
		"qop":          "auth",
		"opaque":       "f9a0f11abf3f6710d22c5a2aa65e19036",
	}
	assert.Equal(T, expected, sec.ParseWWWAuthenticate(`Digest realm="Oasis",nonce="61b6948856629ad7fd3da9d6179393ec",qop="auth",opaque="f9a0f11abf3f6710d22c5a2aa65e19036"`))
}
