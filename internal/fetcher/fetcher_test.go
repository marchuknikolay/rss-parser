package fetcher

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFetch_Success(t *testing.T) {
	const content = "Content"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, content)
	}))
	defer server.Close()

	bs, err := Fetcher{}.Fetch(server.URL)

	require.NoError(t, err)
	require.Equal(t, content, string(bs))
}

func TestFetch_NetworkError(t *testing.T) {
	_, err := Fetcher{}.Fetch("https://invalid.url")

	require.Error(t, err)
}

func TestFetch_StatusNotOK(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		http.Error(w, "InternalServerError", http.StatusInternalServerError)
	}))
	defer server.Close()

	bs, err := Fetcher{}.Fetch(server.URL)

	require.Nil(t, bs)
	require.Error(t, err)
}
