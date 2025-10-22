package mock

import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/marchuknikolay/rss-parser/internal/testutils"
)

type MockRows struct {
	ErrFunc  func() error
	NextFunc func() bool
	ScanFunc func(dest ...any) error
}

func (MockRows) Close() {
}

func (m MockRows) Err() error {
	if m.ErrFunc != nil {
		return m.ErrFunc()
	}

	return nil
}

func (m MockRows) CommandTag() pgconn.CommandTag {
	return pgconn.CommandTag{}
}

func (m MockRows) FieldDescriptions() []pgconn.FieldDescription {
	return nil
}

func (m MockRows) Next() bool {
	if m.NextFunc != nil {
		return m.NextFunc()
	}

	return false
}

func (m MockRows) Scan(dest ...any) error {
	if m.ScanFunc != nil {
		return m.ScanFunc(dest...)
	}

	return testutils.ErrNotImplemented
}

func (m MockRows) Values() ([]any, error) {
	return nil, testutils.ErrNotImplemented
}

func (m MockRows) RawValues() [][]byte {
	return nil
}

func (m MockRows) Conn() *pgx.Conn {
	return nil
}
