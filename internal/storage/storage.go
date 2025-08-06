package storage

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RowQueryer interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

type CommandExecutor interface {
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
}

type Interface interface {
	QueryExecutor() RowQueryer
	ExecExecutor() CommandExecutor

	WithTransaction(ctx context.Context, fn func(txStorage Interface) error) error
	Close()
}

type Storage struct {
	Pool *pgxpool.Pool
	Tx   pgx.Tx
}

func New(connString string) (*Storage, error) {
	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, err
	}

	return &Storage{Pool: pool}, nil
}

func (s *Storage) Close() {
	s.Pool.Close()
}

func (s *Storage) WithTx(tx pgx.Tx) Interface {
	return &Storage{
		Pool: s.Pool,
		Tx:   tx,
	}
}

func (s *Storage) WithTransaction(ctx context.Context, fn func(Interface) error) error {
	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			log.Printf("transaction rollback failed: %v", err)
		}
	}()

	txStorage := s.WithTx(tx)

	if err := fn(txStorage); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *Storage) QueryExecutor() RowQueryer {
	if s.Tx != nil {
		return s.Tx
	}

	return s.Pool
}

func (s *Storage) ExecExecutor() CommandExecutor {
	if s.Tx != nil {
		return s.Tx
	}

	return s.Pool
}
