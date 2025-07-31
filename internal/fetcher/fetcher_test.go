package fetcher

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFetch(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		const content = "Content"

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, content)
		}))
		defer server.Close()

		bs, err := Fetcher{}.Fetch(server.URL)

		require.NoError(t, err)
		require.Equal(t, content, string(bs))
	})

	t.Run("NetworkError", func(t *testing.T) {
		_, err := Fetcher{}.Fetch("https://invalid.url")

		require.Error(t, err)
	})

	t.Run("StatusNotOK", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			http.Error(w, "InternalServerError", http.StatusInternalServerError)
		}))
		defer server.Close()

		bs, err := Fetcher{}.Fetch(server.URL)

		require.Error(t, err)
		require.Nil(t, bs)
	})
}
