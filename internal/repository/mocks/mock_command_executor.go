package mock

import (
	"context"

	"github.com/jackc/pgx/v5/pgconn"
)

type MockCommandExecutor struct {
	ExecFunc func(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

func (m MockCommandExecutor) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	if m.ExecFunc != nil {
		return m.ExecFunc(ctx, sql, args...)
	}

	return pgconn.CommandTag{}, ErrNotImplemented
}
