package fetcher

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/marchuknikolay/rss-parser/internal/fetcher/mock"
	"github.com/stretchr/testify/require"
)

func TestFetch(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		const content = "Content"

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, err := fmt.Fprint(w, content)
			require.NoError(t, err)
		}))
		defer server.Close()

		fetcher := New(http.DefaultClient)
		bs, err := fetcher.Fetch(server.URL)

		require.NoError(t, err)
		require.Equal(t, content, string(bs))
	})

	t.Run("InvalidURL", func(t *testing.T) {
		fetcher := New(&mock.MockHTTPClient{})
		bs, err := fetcher.Fetch("::://invalid-url")

		require.Error(t, err)
		require.Nil(t, bs)
	})

	t.Run("NetworkError", func(t *testing.T) {
		fetcher := New(http.DefaultClient)
		_, err := fetcher.Fetch("https://invalid.url")

		require.Error(t, err)
	})

	t.Run("BodyCloseError", func(t *testing.T) {
		body := mock.MockReadCloser{
			ReadFunc: func(p []byte) (int, error) {
				return 0, io.EOF
			},
			CloseFunc: func() error {
				return fmt.Errorf("Close error")
			},
		}

		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       body,
		}

		fetcher := New(&mock.MockHTTPClient{Resp: resp})
		_, err := fetcher.Fetch("http://example.com")

		require.NoError(t, err)
	})

	t.Run("StatusNotOK", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			http.Error(w, "InternalServerError", http.StatusInternalServerError)
		}))
		defer server.Close()

		fetcher := New(http.DefaultClient)
		bs, err := fetcher.Fetch(server.URL)

		require.Error(t, err)
		require.Nil(t, bs)
	})
}
