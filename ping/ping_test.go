package ping_test

import (
	"bytes"
	"encoding/json"
	"github.com/b3ntly/twelvefactor_ping/ping"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var DEFAULT_PING_RESPONSE = "PONG"

func TestService_Endpoint(t *testing.T) {
	service := ping.New(DEFAULT_PING_RESPONSE, log.New(os.Stdout, "logger: ", log.Lshortfile))
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		service.Endpoint(w, r)
	}))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	require.Nil(t, err)

	response, err := ioutil.ReadAll(res.Body)
	require.Nil(t, err)
	res.Body.Close()

	// json.Marshal will surround the string with quotes so we marshal it here to get an identical representation for comparison
	repr, err := json.Marshal(DEFAULT_PING_RESPONSE)
	require.Nil(t, err)
	require.True(t, bytes.Equal(repr, response))
}

// YOU MIGHT NEED TO RAISE YOUR ULIMIT ON MACOS TO RUN THIS
func BenchmarkService_Endpoint(b *testing.B) {
	service := ping.New("DEFAULT_PING_RESPONSE", log.New(os.Stdout, "logger: ", log.Lshortfile))
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		service.Endpoint(w, r)
	}))
	defer ts.Close()

	b.Run("", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = http.Get(ts.URL)
		}
	})
}
