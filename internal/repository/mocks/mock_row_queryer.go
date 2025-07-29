package mock

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type MockRowQueryer struct {
	QueryRowFunc func(ctx context.Context, sql string, args ...any) pgx.Row
	QueryFunc    func(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

func (m MockRowQueryer) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	if m.QueryRowFunc != nil {
		return m.QueryRowFunc(ctx, sql, args...)
	}

	return &MockRow{
		ScanFunc: func(dest ...any) error {
			return ErrNotImplemented
		},
	}
}

func (m MockRowQueryer) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	if m.QueryFunc != nil {
		return m.QueryFunc(ctx, sql, args...)
	}

	return nil, ErrNotImplemented
}
