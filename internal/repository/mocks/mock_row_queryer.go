package mock

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type MockRowQueryer struct {
	QueryRowFunc func(ctx context.Context, sql string, args ...any) pgx.Row
	QueryFunc    func(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

func (m *MockRowQueryer) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return m.QueryRowFunc(ctx, sql, args...)
}

func (m *MockRowQueryer) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return m.QueryFunc(ctx, sql, args...)
}
