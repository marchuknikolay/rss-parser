package fetcher

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Fetcher struct {
	client HTTPClient
}

func New(c HTTPClient) Fetcher {
	return Fetcher{client: c}
}

func (f Fetcher) Fetch(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed creating a GET request for %v, %w", url, err)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed getting data from %v, %w", url, err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}
