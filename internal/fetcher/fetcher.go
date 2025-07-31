package fetcher

import (
	"fmt"
	"io"
	"net/http"
)

type Interface interface {
	Fetch(url string) ([]byte, error)
}

type Fetcher struct{}

func (f Fetcher) Fetch(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed getting data from %v, %w", url, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}
