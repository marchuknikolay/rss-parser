package mock

import "github.com/marchuknikolay/rss-parser/internal/utils/mock"

type MockFetcher struct {
	FetchFunc func(string) ([]byte, error)
}

func (m MockFetcher) Fetch(url string) ([]byte, error) {
	if m.FetchFunc != nil {
		return m.FetchFunc(url)
	}

	return nil, mock.ErrNotImplemented
}
