package apikey_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	apikey "github.com/x1n13y84issmd42/oasis/src/api/security/APIKey"
	"github.com/x1n13y84issmd42/oasis/src/log"
)

func Test_Cookie(T *testing.T) {
	log := log.NewPlain(0)
	name := "api_key"
	key := "00112233"
	sec := &apikey.Cookie{
		Security: apikey.Security{
			Name:      "test sec",
			ParamName: name,
			Value:     key,
			Log:       log,
		},
	}

	req, _ := http.NewRequest("GET", "example.com", nil)

	sec.Enrich(req, log)

	cookie, _ := req.Cookie(name)

	assert.Equal(T, key, cookie.Value)
}

func Test_Query(T *testing.T) {
	log := log.NewPlain(0)
	name := "api_key"
	key := "00112233"
	sec := &apikey.Query{
		Security: apikey.Security{
			Name:      "test sec",
			ParamName: name,
			Value:     key,
			Log:       log,
		},
	}

	req, _ := http.NewRequest("GET", "example.com", nil)

	sec.Enrich(req, log)

	assert.Equal(T, key, req.URL.Query().Get(name))
}

func Test_Header(T *testing.T) {
	log := log.NewPlain(0)
	name := "api_key"
	key := "00112233"
	sec := &apikey.Header{
		Security: apikey.Security{
			Name:      "test sec",
			ParamName: name,
			Value:     key,
			Log:       log,
		},
	}

	req, _ := http.NewRequest("GET", "example.com", nil)

	sec.Enrich(req, log)

	assert.Equal(T, key, req.Header.Get(name))
}
