package storage

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/marchuknikolay/rss-parser/internal/model"
)

type RowQueryer interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

type CommandExecutor interface {
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
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

func (s *Storage) WithTx(tx pgx.Tx) *Storage {
	return &Storage{
		Pool: s.Pool,
		Tx:   tx,
	}
}

func (s *Storage) WithTransaction(ctx context.Context, fn func(*Storage) error) error {
	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

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

func (s *Storage) FetchItemsByChannelId(channelId int) ([]model.Item, error) {
	rows, err := s.Pool.Query(context.Background(),
		"SELECT id, title, description, pub_date FROM items WHERE channel_id = $1", channelId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	items := []model.Item{}

	for rows.Next() {
		var (
			id                 int
			title, description string
			pubDate            time.Time
		)

		if err := rows.Scan(&id, &title, &description, &pubDate); err != nil {
			return nil, err
		}

		items = append(items, model.Item{Id: id, Title: title, Description: description, PubDate: model.DateTime(pubDate)})
	}

	return items, nil
}

func (s *Storage) FetchChannels() ([]model.Channel, error) {
	rows, err := s.Pool.Query(context.Background(), "SELECT id, title, language, description FROM channels")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	channels := []model.Channel{}

	for rows.Next() {
		var (
			id                           int
			title, language, description string
		)

		if err := rows.Scan(&id, &title, &language, &description); err != nil {
			return nil, err
		}

		items, err := s.FetchItemsByChannelId(id)
		if err != nil {
			return nil, err
		}

		channels = append(channels, model.Channel{Id: id, Title: title, Language: language, Description: description, Items: items})
	}

	return channels, nil
}

func (s *Storage) FetchItemById(id int) (model.Item, error) {
	var (
		title, description string
		pubDate            time.Time
	)

	err := s.Pool.QueryRow(context.Background(),
		"SELECT title, description, pub_date FROM items WHERE id = $1", id).
		Scan(&title, &description, &pubDate)

	if err != nil {
		return model.Item{}, err
	}

	return model.Item{Id: id, Title: title, Description: description, PubDate: model.DateTime(pubDate)}, nil
}
