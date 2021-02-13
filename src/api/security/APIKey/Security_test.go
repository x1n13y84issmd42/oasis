package apikey_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/x1n13y84issmd42/oasis/src/api"
	apikey "github.com/x1n13y84issmd42/oasis/src/api/security/APIKey"
	"github.com/x1n13y84issmd42/oasis/src/log"
	"github.com/x1n13y84issmd42/oasis/src/params"
)

func Test_Cookie(T *testing.T) {
	log := log.NewPlain(0)
	name := "api_key"
	key := "00112233"
	sec := &apikey.Cookie{
		Security: apikey.Security{
			Name:      "test sec",
			ParamName: name,
			Value:     params.Value(key),
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
			Value:     params.Value(key),
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
			Value:     params.Value(key),
			Log:       log,
		},
	}

	req, _ := http.NewRequest("GET", "example.com", nil)

	sec.Enrich(req, log)

	assert.Equal(T, key, req.Header.Get(name))
}

func Test_New(T *testing.T) {
	key := "the_security_key"
	value := "0011223344"

	T.Run("Cookie", func(T *testing.T) {
		sec := apikey.New("test sec", "cookie", key, value, log.NewPlain(0))
		tsec, ok := sec.(*apikey.Cookie)
		assert.True(T, ok)

		if ok {
			assert.Equal(T, key, tsec.ParamName)
			assert.Equal(T, value, tsec.Value())
		}
	})

	T.Run("Cookie", func(T *testing.T) {
		sec := apikey.New("test sec", "query", key, value, log.NewPlain(0))
		tsec, ok := sec.(*apikey.Query)
		assert.True(T, ok)

		if ok {
			assert.Equal(T, key, tsec.ParamName)
			assert.Equal(T, value, tsec.Value())
		}
	})

	T.Run("Header", func(T *testing.T) {
		sec := apikey.New("test sec", "header", key, value, log.NewPlain(0))
		tsec, ok := sec.(*apikey.Header)
		assert.True(T, ok)

		if ok {
			assert.Equal(T, key, tsec.ParamName)
			assert.Equal(T, value, tsec.Value())
		}
	})

	T.Run("Invalid", func(T *testing.T) {
		sec := apikey.New("test sec", "INVALID LOCATION", key, value, log.NewPlain(0))
		_, ok := sec.(*api.NullSecurity)
		assert.True(T, ok)
	})
}
