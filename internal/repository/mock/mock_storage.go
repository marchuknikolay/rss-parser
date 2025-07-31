package mock

import (
	"context"

	"github.com/marchuknikolay/rss-parser/internal/storage"
	"github.com/marchuknikolay/rss-parser/internal/testutils"
)

type MockStorage struct {
	QueryExecutorFunc storage.RowQueryer
	ExecExecutorFunc  storage.CommandExecutor

	WithTransactionFunc func(ctx context.Context, fn func(txStorage storage.Interface) error) error
}

func (m MockStorage) QueryExecutor() storage.RowQueryer {
	return m.QueryExecutorFunc
}

func (m MockStorage) ExecExecutor() storage.CommandExecutor {
	return m.ExecExecutorFunc
}

func (m MockStorage) WithTransaction(ctx context.Context, fn func(txStorage storage.Interface) error) error {
	if m.WithTransactionFunc != nil {
		return m.WithTransactionFunc(ctx, fn)
	}

	return testutils.ErrNotImplemented
}

func (m MockStorage) Close() {
}
