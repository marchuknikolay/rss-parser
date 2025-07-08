package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/marchuknikolay/rss-parser/internal/model"
	"github.com/marchuknikolay/rss-parser/internal/storage"
)

type ItemRepository struct {
	Storage *storage.Storage
}

func NewItemRepository(storage *storage.Storage) *ItemRepository {
	return &ItemRepository{Storage: storage}
}

func (r *ItemRepository) Save(ctx context.Context, item model.Item, channelId int) error {
	query := "INSERT INTO items (title, description, pub_date, channel_id) VALUES ($1, $2, $3, $4)"

	executor := r.Storage.ExecExecutor()
	_, err := executor.Exec(ctx, query, item.Title, item.Description, item.PubDate, channelId)

	return err
}

func (r *ItemRepository) GetAll(ctx context.Context) ([]model.Item, error) {
	query := `SELECT id, title, description, pub_date FROM items`
	return r.getItems(ctx, query)
}

func (r *ItemRepository) GetByChannelId(ctx context.Context, channelId int) ([]model.Item, error) {
	query := `SELECT id, title, description, pub_date FROM items WHERE channel_id = $1`
	return r.getItems(ctx, query, channelId)
}

func (r *ItemRepository) GetById(ctx context.Context, itemId int) (model.Item, error) {
	query := `SELECT id, title, description, pub_date FROM items WHERE id = $1`

	executor := r.Storage.QueryExecutor()

	var (
		id                 int
		title, description string
		pubDate            time.Time
	)

	if err := executor.QueryRow(ctx, query, itemId).Scan(&id, &title, &description, &pubDate); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Item{}, fmt.Errorf("item with id=%d not found: %w", itemId, err)
		}

		return model.Item{}, fmt.Errorf("failed to scan item: %w", err)
	}

	return model.Item{
		Id:          id,
		Title:       title,
		Description: description,
		PubDate:     model.DateTime(pubDate),
	}, nil
}

func (r *ItemRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM items WHERE id = $1`

	executor := r.Storage.ExecExecutor()
	tag, err := executor.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete item with id=%d: %w", id, err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("no item found with id=%d", id)
	}

	return nil
}

func (r *ItemRepository) Update(ctx context.Context, id int, title, description string, pubTime time.Time) error {
	query := `UPDATE items SET title = $1, description = $2, pub_date = $3 WHERE id = $4`

	executor := r.Storage.ExecExecutor()
	tag, err := executor.Exec(ctx, query, title, description, pubTime, id)
	if err != nil {
		return fmt.Errorf("failed to update item with id=%d: %w", id, err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("no item found with id=%d", id)
	}

	return nil
}

func (r *ItemRepository) getItems(ctx context.Context, query string, args ...any) ([]model.Item, error) {
	executor := r.Storage.QueryExecutor()

	rows, err := executor.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var items []model.Item

	for rows.Next() {
		var (
			id                 int
			title, description string
			pubDate            time.Time
		)

		if err := rows.Scan(&id, &title, &description, &pubDate); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}

		items = append(items, model.Item{
			Id:          id,
			Title:       title,
			Description: description,
			PubDate:     model.DateTime(pubDate),
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration failed: %w", err)
	}

	return items, nil
}
