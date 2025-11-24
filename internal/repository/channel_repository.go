package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/marchuknikolay/rss-parser/internal/model"
	"github.com/marchuknikolay/rss-parser/internal/storage"
)

var ErrChannelNotFound = errors.New("channel not found")

type ChannelRepositoryInterface interface {
	Save(ctx context.Context, channel *model.Channel) (int, error)
	GetAll(ctx context.Context) ([]model.Channel, error)
	GetById(ctx context.Context, id int) (model.Channel, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, id int, title, language, description string) (model.Channel, error)
}

type ChannelRepository struct {
	storage.Interface
}

func (r *ChannelRepository) Save(ctx context.Context, channel *model.Channel) (int, error) {
	var channelId int
	query := "INSERT INTO channels (title, language, description) VALUES ($1, $2, $3) RETURNING id"

	executor := r.QueryExecutor()
	err := executor.QueryRow(ctx, query, channel.Title, channel.Language, channel.Description).Scan(&channelId)

	return channelId, err
}

func (r *ChannelRepository) GetAll(ctx context.Context) ([]model.Channel, error) {
	query := `SELECT id, title, language, description FROM channels`

	executor := r.QueryExecutor()
	rows, err := executor.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query channels: %w", err)
	}
	defer rows.Close()

	var channels []model.Channel

	for rows.Next() {
		var channel model.Channel

		if err := rows.Scan(&channel.Id, &channel.Title, &channel.Language, &channel.Description); err != nil {
			return nil, fmt.Errorf("failed to scan channel row: %w", err)
		}

		channels = append(channels, channel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return channels, nil
}

func (r *ChannelRepository) GetById(ctx context.Context, id int) (model.Channel, error) {
	query := `SELECT id, title, language, description FROM channels WHERE id = $1`

	executor := r.QueryExecutor()
	row := executor.QueryRow(ctx, query, id)

	var channel model.Channel
	if err := row.Scan(&channel.Id, &channel.Title, &channel.Language, &channel.Description); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Channel{}, ErrChannelNotFound
		}

		return model.Channel{}, fmt.Errorf("failed to scan channel: %w", err)
	}

	return channel, nil
}

func (r *ChannelRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM channels WHERE id = $1`

	executor := r.ExecExecutor()
	tag, err := executor.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete channel with id=%d: %w", id, err)
	}

	if tag.RowsAffected() == 0 {
		return ErrChannelNotFound
	}

	return nil
}

func (r *ChannelRepository) Update(
	ctx context.Context,
	id int,
	title, language, description string,
) (model.Channel, error) {
	query := `
		UPDATE channels
		SET title = $1, language = $2, description = $3
		WHERE id = $4
		RETURNING id, title, language, description
	`

	executor := r.QueryExecutor()
	row := executor.QueryRow(ctx, query, title, language, description, id)

	var channel model.Channel
	if err := row.Scan(&channel.Id, &channel.Title, &channel.Language, &channel.Description); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Channel{}, ErrChannelNotFound
		}

		return model.Channel{}, fmt.Errorf("failed to update channel with id=%d: %w", id, err)
	}

	return channel, nil
}
