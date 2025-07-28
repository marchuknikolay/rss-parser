package mock

import (
	"github.com/marchuknikolay/rss-parser/internal/storage"
)

type MockStorage struct {
	QueryExecutorFunc storage.RowQueryer
	ExecExecutorFunc  storage.CommandExecutor
}

func (m *MockStorage) QueryExecutor() storage.RowQueryer {
	return m.QueryExecutorFunc
}

func (m *MockStorage) ExecExecutor() storage.CommandExecutor {
	return m.ExecExecutorFunc
}
