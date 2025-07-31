package mock

import "github.com/marchuknikolay/rss-parser/internal/testutils"

type MockFetcher struct {
	FetchFunc func(url string) ([]byte, error)
}

func (m MockFetcher) Fetch(url string) ([]byte, error) {
	if m.FetchFunc != nil {
		return m.FetchFunc(url)
	}

	return nil, testutils.ErrNotImplemented
}
