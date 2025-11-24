package mock

import (
	"context"

	"github.com/marchuknikolay/rss-parser/internal/testutils"
)

type MockFetcher struct {
	FetchFunc func(ctx context.Context, url string) ([]byte, error)
}

func (m MockFetcher) Fetch(ctx context.Context, url string) ([]byte, error) {
	if m.FetchFunc != nil {
		return m.FetchFunc(ctx, url)
	}

	return nil, testutils.ErrNotImplemented
}
